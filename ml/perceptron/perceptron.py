#!/bin/python

import os, pickle, re, argparse, time
import string, random

parser = argparse.ArgumentParser(description="Prepare for perceptron.")
parser.add_argument("--input",required=True,help="Specify the input data file.")
parser.add_argument("--iterations",help="Number of iterations.",type=int,default=5)
parser.add_argument("--train",help="What file are we training the model from.")
parser.add_argument("--test",help="What file are we testing the trained model from.")
parser.add_argument("--model",help="Where are we dumping the model after generated.")
args = parser.parse_args()

documents = {}
inputs = []
training = {}

BIAS = -1
THRESHOLD = .75 # > = spam
LEARNING_RATE = .2
NEURONS = 2 
INPUTS = 0

def load_model():
	model = pickle.load(open(args.model, "r"))
	weights = model["weights"]
	inputs = model["inputs"]
	return weights, inputs, len(inputs)

def initialize_data():
	data = pickle.load(open(args.input, "r"))
	documents = data["documents"]
	inputs = sorted(data["words"].keys())
	return documents, inputs, len(inputs) + 1

def initialize_random_weights():
	weights = []
	for i in range(INPUTS):
		weight = []
		for j in range(NEURONS):
			if random.random() > .5:
				pn = 1
			else:
				pn = -1
			weight.append(pn * random.random())
		weights.append(weight)
	return weights

def initialize_comparison_data(f):
	data = {}
	for line in f.readlines():
		split = line.split()
		t = []
		if split[0] == "spam":
			t = [1, 0]
		else:
			t = [0, 1]
		f = split[1].split("/")[-1]
		data[f] = t
	return data

def assign_input_values(words):
	tmp = [BIAS]
	for i in range(INPUTS-1):
		_input = inputs[i]
		if _input in words:
			tmp.append(words[_input])
		else:
			tmp.append(0)
	return tmp

def g(doc_inputs):
	neurons = [0 for j in range(NEURONS)] #initialize activations
	for j in range(NEURONS):
		for i in range(INPUTS):
			neurons[j] += (doc_inputs[i] * weights[i][j])
	return neurons

def activation(neurons):
	activated = []
	for i in range(NEURONS):
		if neurons[i] > THRESHOLD:
			activated.append(1)
		else:
			activated.append(0)
	return activated

def adjust_weights(doc_inputs, y, t):
	for j in range(NEURONS):
		for i in range(INPUTS):
			weights[i][j] = weights[i][j] + (LEARNING_RATE * (t[j] - y[j]) * doc_inputs[i])

def run():
	decisions = {}
	for document in documents:
		doc = documents[document]
		doc_inputs = assign_input_values(doc["words"])
		neurons = g(doc_inputs)
		activated_neurons = activation(neurons)
		decisions[document] = activated_neurons
		#print document, training[document], activated_neurons
		if args.train:
			adjust_weights(doc_inputs, activated_neurons, training[document])
	return decisions

def test_neurons(results):
	test_data = initialize_comparison_data(open(args.test, "r"))
	print results
	for result in results:
		_result = results[result]
		_test = test_data[result]
		for j in range(NEURONS):
			if _result[j] <> _test[j]:
				print result, "INCORRECT MATCH:", _result, _test
				break
			else:
				print result, "CORRECT MATCH"

def spam_or_ham(results):
	pass

def save_model(model):
	pickle.dump(model, open(args.model, "w"))




#stop using functions and actually do something!
documents, inputs, INPUTS = initialize_data()
if args.model and not args.train:
	weights, inputs, INPUTS = load_model()
else:
	weights = initialize_random_weights()
if args.train:
	training = initialize_comparison_data(open(args.train, "r"))

results = None
if args.train:
	for i in range(args.iterations):
		start = time.time()
		print "Started iteration:",i
		run()
		end = time.time()
		print "Completed iteration:",i,"and took",(end - start),"s."
else:
	results = run()

model = {
		"weights": weights,
		"inputs": inputs
		}
if args.train and args.model:
	save_model(model)
elif args.train:
	print model
elif args.test:
	test_neurons(results)
else:
	spam_or_ham(results)

#!/bin/python

import os, pickle, re, argparse, time
import string, random

parser = argparse.ArgumentParser(description="Prepare for perceptron.")
parser.add_argument("--input",required=True,help="Specify the input data file.")
parser.add_argument("--iterations",help="Number of iterations.",type=int,default=5)
parser.add_argument("--train",action="store_true")
parser.add_argument("--test",action="store_true")
parser.add_argument("--model",help="Where are we dumping the model after generated.")
args = parser.parse_args()

inputs = []

BIAS = -1
THRESHOLD = .75 # > = spam
LEARNING_RATE = .2
NEURONS = 1
INPUTS = 0

def load_model():
	model = pickle.load(open(args.model, "r"))
	weights = model["weights"]
	inputs = model["inputs"]
	return weights, inputs, len(inputs)

def initialize_data():
	data = pickle.load(open(args.input, "r"))
	vectors = data["vectors"]
	inputs = data["inputs"]
	return vectors, inputs, len(inputs)

def initialize_random_weights():
	weights = []
	for i in range(INPUTS):
		weight = []
		for j in range(NEURONS):
			pn = 1
			if random.random() > .5:
				pn = -1
			weight.append(pn * random.random())
		weights.append(weight)
	return weights

def g(weights, vector):
	neurons = [-1.0 for j in range(NEURONS)] #initialize activations
	for j in range(NEURONS):
		for i in range(INPUTS):
			neurons[j] = neurons[j] + vector[i] * weights[i][j]
	return neurons

def activation(neurons):
	activated = []
	for i in range(NEURONS):
		if neurons[i] > THRESHOLD:
			activated.append(1)
		else:
			activated.append(0)
	return activated

def adjust_weights(weights, vector, y, t):
	for j in range(NEURONS):
		for i in range(INPUTS):
			weights[i][j] = weights[i][j] + (LEARNING_RATE * (t[j] - y[j]) * vector[i])
	return weights

def run(vectors, weights):
	decisions = {}
	for f in vectors:
		doc = vectors[f]
		vector = doc["vector"][0]
		neurons = g(weights, vector)
		#print f, neurons
		activated_neurons = activation(neurons)
		decisions[f] = activated_neurons
		if args.train:
			weights = adjust_weights(weights, vector, activated_neurons, doc["vector"][-1])
	return decisions, weights

def test_neurons(results, vectors):
	print results
	for result in results:
		_result = results[result]
		_test = vectors[result]["vector"][-1]
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
vectors, inputs, INPUTS = initialize_data()
if args.model and not args.train:
	weights, inputs, INPUTS = load_model()
else:
	weights = initialize_random_weights()


results = None
if args.train:
	for i in range(args.iterations):
		start = time.time()
		#print "Started iteration:",i
		results, weights = run(vectors, weights)
		end = time.time()
		#print "Completed iteration:",i,"and took",(end - start),"s."
else:
	results, weights = run(vectors, weights)

model = {
		"weights": weights,
		"inputs": inputs
		}
if args.train and args.model:
	save_model(model)
elif args.train:
	print model
elif args.test:
	test_neurons(results, vectors)
else:
	spam_or_ham(results)

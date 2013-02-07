#!/bin/python

from __future__ import division
import os, pickle, re, argparse
import string, math

parser = argparse.ArgumentParser(description="Prepare for perceptron.")
parser.add_argument("--input",required=True,help="Specify the input data directory.")
parser.add_argument("--limit",help="Where to limit processing files from.",type=int)
parser.add_argument("--comparison",help="File containing spam/ham specifications.")
parser.add_argument("--test",help="File to dump the testing data into.")
parser.add_argument("--train",help="File to dump the training data into.")
args = parser.parse_args()

stop_words = ['a','able','about','across','after','all','almost','also','am','among',
		'an','and','any','are','as','at','be','because','been','but','by','can',
		'cannot','could','dear','did','do','does','either','else','ever','every',
		'for','from','get','got','had','has','have','he','her','hers','him','his',
		'how','however','i','if','in','into','is','it','its','just','least','let',
		'like','likely','may','me','might','most','must','my','neither','no','nor',
		'not','of','off','often','on','only','or','other','our','own','rather','said',
		'say','says','she','should','since','so','some','than','that','the','their',
		'them','then','there','these','they','this','tis','to','too','twas','us',
		'wants','was','we','were','what','when','where','which','while','who',
		'whom','why','will','with','would','yet','you','your']
re_boundary = re.compile("boundary\=\"([^\"]*)\"")
re_newline_boundary = re.compile("^\n")
re_plain_text = re.compile("text\/plain")
d = args.input
PERCENT_TEST = .1

words = {}

def strip_stopwords():
	keys = words.keys()
	for word in stop_words:
		if word in keys:
			del(words[word])

def count_words(lines):
	doc_words = {}
	for line in lines:
		split = string.split(line)
		for s in split:
			if s not in words:
				words[s] = 0
			if s not in doc_words:
				doc_words[s] = 0
			words[s] += 1
			doc_words[s] += 1
	return doc_words

############
#Email Code

def find_boundary(lines):
	for line in lines:
		m = re.search(re_boundary, line)
		if m <> None:
			return m.group(1)

def get_boundary_starts(lines, boundary):
	if boundary <> None:
		return [i for i, x in enumerate(lines) if boundary in x]
	if boundary == None:
		return [i for i, x in enumerate(lines) if re.match(re_newline_boundary, x)]

def split_boundaries(lines, boundary):
	boundaries = []
	boundary_indices = get_boundary_starts(lines, boundary)

	#short circuit since we don't care about the rest
	#if there is no boundary, the entire doc is "valid"
	if boundary == None:
		boundaries.append(lines[boundary_indices[0]:])
		return boundaries

	boundary_pairs = []
	previous = None
	first = None
	for boundary_point in boundary_indices:
		if first == None:
			first = True
			continue
		if previous:
			boundary_pairs.append([previous, boundary_point])
		previous = boundary_point
	for boundary_pair in boundary_pairs:
		boundaries.append(lines[boundary_pair[0]:boundary_pair[1]])
	return boundaries

def pick_boundary(boundaries):
	if len(boundaries) == 1:
		return boundaries[0]

	for boundary in boundaries:
		for line in boundary:
			m = re.search(re_plain_text, line)
			if m <> None:
				start = get_boundary_starts(boundary, None)
				return boundary[start[0]:]

def parse_documents(files, targets):
	documents = {}
	for f in files:
		_f = open(os.path.join(d, f))
		contents = _f.readlines()
		boundary = find_boundary(contents)
		boundaries = split_boundaries(contents, boundary)
		chosen = pick_boundary(boundaries)
		if chosen == None:
			continue
		strip_stopwords()
		doc_words = count_words(chosen)
		doc_words_count = len(doc_words.keys())
		for word in doc_words:
			doc_words[word] = doc_words[word] / doc_words_count
		documents[f] = {
				"target": targets[f],
				"words": doc_words
				}
	return documents

###########
# ML Stuff

def initialize_target_data(f):
	data = {}
	for line in open(f, "r").readlines():
		split = line.split()
		if split[0] == "spam":
			t = [1]
		else:
			t = [0]
		f = split[1].split("/")[-1]
		data[f] = t
	return data

def setup_inputs(doc_words):
	tmp = []
	for word in words:
		if word in doc_words:
			tmp.append(doc_words[word])
		else:
			tmp.append(0)
	return tmp

def process_inputs(doc, f):
	doc_words = doc["words"]
	target = doc["target"]
	vector = None
	doc_vector = setup_inputs(doc_words)
	if target <> None:
		vector = (doc_vector, target)
	else:
		vector = (doc_vector)
	return {
			"vector": vector
			}

def convert_to_vectors(docs):
	for doc in docs:
		docs[doc] = process_inputs(docs[doc], doc)
	return docs

############
# Functional Code

files = os.listdir(d)
files = files[:(args.limit or len(files))]
targets = initialize_target_data(args.comparison)
max_test = int(math.floor(len(files) * PERCENT_TEST))
test = parse_documents(files[:max_test], targets)
train = parse_documents(files[(max_test+1):], targets)
all_inputs = sorted(words.keys())
test = convert_to_vectors(test)
train = convert_to_vectors(train)

pickle.dump({
	"inputs": all_inputs,
	"vectors": test
	}, open(args.test, "w"))
pickle.dump({
	"inputs": all_inputs,
	"vectors": train
	}, open(args.train, "w"))

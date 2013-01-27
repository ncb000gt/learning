#!/bin/python

import os, pickle, re, argparse
import string

parser = argparse.ArgumentParser(description="Prepare for perceptron.")
parser.add_argument("--input",required=True,help="Specify the input data directory.")
parser.add_argument("--output",required=True,help="Specify the output file.")
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
d = args.input #"/home/ncampbell/Downloads/mldata/ceas08-1/data"

words = {}
documents = {}

def strip_stopwords():
	keys = words.keys()
	for word in stop_words:
		if word in keys:
			del(words[word])

def count_words(lines):
	for line in lines:
		split = string.split(line)
		for s in split:
			if s not in words:
				words[s] = 0
			words[s] += 1

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

files = os.listdir(d)[:1000]

#print "Files:",files

for f in files:
	_f = open(os.path.join(d, f))
	contents = _f.readlines()
	boundary = find_boundary(contents)
	boundaries = split_boundaries(contents, boundary)
	chosen = pick_boundary(boundaries)
	documents[f] = chosen

for document in documents:
	count_words(documents[document])

strip_stopwords()


tmp = {
		"words": words,
		"documents": documents
		}
pickle.dump(tmp, open(args.output, "w"))

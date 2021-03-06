package perceptron

import (
	"fmt"
	"os"
	"io"
	"path"
	"strconv"
	"bufio"
	"regexp"
	"net/mail"
	"path/filepath"
	"encoding/csv"
	"strings"
)
var re_boundary, _ = regexp.Compile("boundary=\"([^\"]*)\"")
var re_selected_body, _ = regexp.Compile("text\\/plain|text\\/html")

func getFiles(D *string, offset *int, limit *int) []string {
	files := make([]string, 0)
	start := *offset
	stop := start + *limit
	count := 0
	filepath.Walk(*D, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		if strings.HasSuffix(name, "swp") || name == filepath.Base(*D) {
			return nil
		} else if count < start {
			count++
			return nil
		} else if count > stop {
			return nil
		}

		files = append(files, p)
		count++
		return nil
	})

	return files
}

func getMessage(P string) *mail.Message {
	file, _ := os.Open(P)
	defer file.Close()

	r := bufio.NewReader(file)
	message, _ := mail.ReadMessage(r)
	return message
}

func findPart(header mail.Header) string {
	var part string 
	for _, v := range header {
		if matches := re_boundary.FindStringSubmatch(strings.Join(v, " ")); len(matches) > 0 {
			part = matches[1]
			break
		}
	}

	return part
}

func getLines(r *io.Reader) []string {
	lines := make([]string, 0) 
	mr := bufio.NewReader(*r)

	for {
		line, rerr := mr.ReadString('\n')
		if rerr != nil {
			break
		}

		lines = append(lines, line)
	}

	return lines
}

func splitBodies(lines *[]string, boundary *string) [][]string {
	re_part, _ := regexp.Compile(strings.Replace(*boundary, "+", "\\+", -1))
	m := make([][]string, 0)

	count := 0
	nlines := make([]string, 0)
	for _, line := range *lines {
		if re_part.MatchString(line) {
			if count > 0 {
				m = append(m, nlines)
				nlines = make([]string, 0)
			}
			count++
			continue //skip the part line...not needed
		}
		
		nlines = append(nlines, line)
	}

	return m
}

func splitBodyAndHeader(body *[]string) map[string][]string {
	nbody := make(map[string][]string)

	idx:= First(*body, func(item string) bool {
		item = strings.TrimSpace(item)
		return (item == "\n" || item == "")
	})
	if idx > 0 {
		nbody["headers"] = (*body)[0:idx]
		nbody["body"] = (*body)[idx+1:]
	}

	return nbody
}

func selectBody(bodies *[][]string) []string {
	for _, body := range *bodies {
		bnh := splitBodyAndHeader(&body)

		header := bnh["headers"]
		test_string := strings.Join(header, " ")
		if re_selected_body.MatchString(test_string) {
			return bnh["body"]
		}
	}

	return make([]string, 0)
}

func getText(P string) string {
	/*fmt.Println("file: ", P)*/
	message := getMessage(P)
	lines := getLines(&message.Body)

	part := findPart(message.Header)
	var body []string
	if part != "" {
		bodies := splitBodies(&lines, &part)
		body = selectBody(&bodies)
	} else {
		body = lines
	}

	str := strings.Join(body, " ")
	return str
}

var re_junk, _ = regexp.Compile("(?ism)\n|<[^>]*>")
func removeJunk(junked *string) string {
	dejunked := re_junk.ReplaceAllString(*junked, " ")

	return dejunked
}

type VectorizeHandler func([]string, float64, int, float64)

func vectorize(inputs *[]string, docs *map[string]map[string]float64, tests *map[string][]int, NEURONS *int, TestNum *float64, Inputs *csv.Writer, f VectorizeHandler) {
	//dump inputs to disk
	if Inputs != nil {
		if err := (*Inputs).Write(*inputs); err != nil {
			panic("Error writing inputs.")
		}
		(*Inputs).Flush()
	}

	count := 0.0
	ldocs := len(*docs)
	linputs := len(*inputs)
	for d, doc := range *docs {
		vector := make([]string, linputs + *NEURONS + 1)
		vector[0] = d

		for idx, word := range *inputs {
			val := 0.0
			if _, ok := doc[word]; ok {
				val = doc[word]
			}
			vector[idx+1] = strconv.FormatFloat(val, 'f', -1, 64)
		}
		
		if tests != nil {
			test := (*tests)[d]
			for i := 1; i <= *NEURONS; i++ {
				vector[linputs+i] = strconv.FormatInt(int64(test[i-1]), 10)
			}
		}

		f(vector, count, ldocs, *TestNum)
		count++
	}
}

var re_non_ascii, _ = regexp.Compile("(?i)[^a-z0-9]")
func splitAndGatherCounts(docs *map[string]string) ([]string, map[string]map[string]float64) {
	words := make([]string, 0)
	ndocs := make(map[string]map[string]float64, len(*docs))

	for k, doc := range *docs {
		ndoc := make(map[string]float64)

		for _, word := range strings.Split(doc, " ") {
			word = strings.Replace(strings.TrimSpace(re_non_ascii.ReplaceAllString(word, " ")), " ", "", -1)

			if word == "" || word == " " {
				continue
			}

			if First(words, func(_word string) bool {
				return (_word == word)
			}) < 0 {
				words = append(words, word)
			}
			ndoc[word] += 1.0
		}

		l := len(ndoc)
		for word, v := range ndoc {
			ndoc[word] = v/float64(l)
		}
		ndocs[k] = ndoc
	}

	return words, ndocs
}

type VectorHandler func([]string)

func Preprocess(config *Config, file string, f VectorHandler) {
	Inputs, _ := ReadInput((*config).Inputs)

	name := path.Base(file)
	text := strings.TrimSpace(getText(file))
	if text == "" {
		panic("File has no good data.")
	}
	docs := make(map[string]string)

	text = removeJunk(&text)
	docs[name] = text
	_, ndocs := splitAndGatherCounts(&docs)
	
	vectorize(&Inputs, &ndocs, nil, &(*config).Neurons, &(*config).Preprocessor.TestCount, nil, func(vector []string, _ float64, _ int, _ float64) {
		f(vector)
	})
}

func RunPreprocessor(config *Config) {
	fmt.Println("Preprocess Starting")
	fmt.Println(*config)
	paths := getFiles(&config.Preprocessor.Files, &config.Preprocessor.Offset, &config.Preprocessor.Limit)
	docs := make(map[string]string)

	for _, file := range paths {
		docs[path.Base(file)] = getText(file)
	}

	docs = Filter(docs, func(item string) bool {
		if strings.TrimSpace(item) != "" { return true }
		return false
	})
	fmt.Println("# messages:", len(docs))

	for k, v := range docs {
		docs[k] = removeJunk(&v)
	}

	train, _ := os.Create((*config).Train)
	defer train.Close()
	Train := csv.NewWriter(train)
	test, _ := os.Create((*config).Test)
	defer test.Close()
	Test := csv.NewWriter(test)
	inputs, _ := os.Create((*config).Inputs)
	defer inputs.Close()
	Inputs := csv.NewWriter(inputs)

	words, ndocs := splitAndGatherCounts(&docs)
	m := GetTestValues((*config).Preprocessor.Test)
	vectorize(&words, &ndocs, &m, &(*config).Neurons, &(*config).Preprocessor.TestCount, Inputs, func(vector []string, count float64, ldocs int, testNum float64) {
		if count/float64(ldocs) >= testNum {
			(*Train).Write(vector)
			(*Train).Flush()
		} else {
			(*Test).Write(vector)
			(*Test).Flush()
		}
	})

	fmt.Println("Preprocess Finished")
}

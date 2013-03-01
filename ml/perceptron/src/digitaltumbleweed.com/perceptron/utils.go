package perceptron

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"os"
	"bufio"
	"io"
	"strings"
)
/*
 * Configuration
 */
type Perceptron struct {
	Iterations int
	Model string
}

type Preprocessor struct {
	Limit int
	Offset int
	TestCount float64
	Files string
	Test string
}

type Config struct {
	Train string
	Test string
	Inputs string
	Neurons int
	Charset string
	Perceptron Perceptron
	Preprocessor Preprocessor
}

type Doc struct {
	File string
	Vector []float64
	Test []int64
}

type Words struct {
	Inputs []string
}

func ReadConfig(C string) Config {
	var config Config
	err := json.Unmarshal(ReadFile(&C), &config)
	if err != nil {
		fmt.Println("err: ", err)
	}

	return config
}

func ReadFile(C *string) []byte {
	bytes, rerr := ioutil.ReadFile(*C)
	if rerr != nil {
		fmt.Println("err: ", rerr)
	}
	return bytes
}

type InStrOutBoolFunc func(string) (bool)

func First(a []string, f InStrOutBoolFunc) int {
	for idx, item := range a {
		if f(item) {
			return idx
		}
	}

	return -1
}

func Map(a []string, f InStrOutBoolFunc) []bool {
	reta := make([]bool, 0)
	for _, item := range a {
		reta = append(reta, f(item))
	}

	return reta
}

func Filter(a map[string]string, f InStrOutBoolFunc) map[string]string {
	for k, v := range a {
		if !f(v) {
			delete(a, k)
		}
	}

	return a
}

func FloatToStringSlice(a []float64) []string {
	s := make([]string, len(a))
	for idx, f := range a {
		s[idx] = strconv.FormatFloat(f, 'f', -1, 64)
	}

	return s
}

func IntToStringSlice(a []int) []string {
	s := make([]string, len(a))
	for idx, i := range a {
		s[idx] = strconv.FormatInt(int64(i), 10)
	}

	return s
}

func StringToFloatSlice(a []string) []float64 {
	s := make([]float64, len(a))
	for idx, f := range a {
		s[idx], _ = strconv.ParseFloat(f, 64)
	}

	return s
}

func GetTestValues(file string) map[string][]int {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bf := bufio.NewReader(f)

	m := make(map[string][]int)

	for {
		line, err := bf.ReadString('\n')
		if err == io.EOF {
			break
		}
		
		s := strings.Split(line, " ")
		v := make([]int, 1)
		if s[0] == "spam" {
			v[0] = 1
		} else {
			v[0] = 0
		}
		s1 := strings.Split(s[1], "/")
		k := strings.TrimSpace(s1[len(s1) - 1])

		m[k] = v //we're only doing this because we have a predefined notion of spam and ham which has 1 neuron
	}

	return m
}

func ReadLine(reader *csv.Reader) ([]string, error) {
	return (*reader).Read()
}

func ReadInput(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return ReadLine(csv.NewReader(f))
}

package perceptron

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
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

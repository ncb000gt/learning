package perceptron

import (
	"math/rand"
	"time"
	"encoding/csv"
	"os"
	"fmt"
	"io"
	"strconv"
)

func dotProduct(X []float64, Y []float64) float64 {
	N := len(X)
	if N != len(Y) {
		panic("N != Y")
	}
	V := 0.0
	for i := 0; i < N; i++ {
		V += X[i] * Y[i]
	}

	return V
}

var BIAS = -1.0
func g(W [][]float64, V []float64) []float64 {
	N := len(W)
	NEURONS := make([]float64, N)

	for j := 0; j < N; j++ {
		NEURONS[j] = BIAS + dotProduct(W[j], V)
	}

	return NEURONS
}

func activation(NEURONS []float64) []int64 {
	N := len(NEURONS)
	ACTIVATED := make([]int64, N)

	for j := 0; j < N; j++ {
		active := 0

		if NEURONS[j] > .75 {
			active = 1
		}

		ACTIVATED[j] = int64(active)
	}

	return ACTIVATED
}

var LEARNING_RATE = 0.1
func adjustWeights(W [][]float64, V []float64, Y []int64, T []int64) [][]float64 {
	N := len(W)
	I := len(V)
	NW := make([][]float64, N)

	for j := 0; j < N; j++ {
		NWn := make([]float64, I)
		for i := 0; i < I; i++ {
			NWn[i] = W[j][i] + (LEARNING_RATE * float64(T[j] - Y[j]) * V[i])
		}
		NW[j] = NWn
	}

	return NW
}

/*
 * Setup Application
 */
func randomWeights(INPUTS int, NEURONS int) [][]float64 {
	W := make([][]float64, NEURONS)

	for j := 0; j < NEURONS; j++ {
		I := make([]float64, INPUTS)

		for i := 0; i < INPUTS; i++ {
			m := 1.0
			if rand.Float64() > .5 {
				m = -1.0
			}
			I[i] = m * rand.Float64()
		}

		W[j] = I
	}

	return W
}

type Model struct {
	Weights [][]float64
}

type VectorFunc func(Doc)

func readLine(r *csv.Reader, NEURONS int) (Doc, error) {
	line, err := ReadLine(r)
	if err == io.EOF {
		var doc Doc
		return doc, err
	}
	if err != nil {
		panic(err)
	}

	return lineToDoc(line, NEURONS), nil
}

func lineToDoc(line []string, NEURONS int) Doc {
	var doc Doc
	doc.File = line[0]
	for i := 1; i < (len(line) - NEURONS); i++ {
		fl, _ := strconv.ParseFloat(line[i], 64)
		doc.Vector = append(doc.Vector, fl)
	}

	for i := len(line) - NEURONS; i < len(line); i++ {
		il, _ := strconv.ParseInt(line[i], 10, 64)
		doc.Test = append(doc.Test, il)
	}

	return doc
}

func ProcessVectors(file string, NEURONS int, model Model) Model {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	W := model.Weights

	for {
		doc, err := readLine(r, NEURONS)
		if err == io.EOF {
			break
		}

		W = adjustWeights(W, doc.Vector, activation(g(W, doc.Vector)), doc.Test)
	}

	model.Weights = W
	return model
}

func TestVectors(file string, NEURONS int, model Model) map[string][][]int64 {
	m := make(map[string][][]int64)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	W := model.Weights

	for {
		doc, err := readLine(r, NEURONS)
		if err == io.EOF {
			break
		}

		m[doc.File] = testVector(doc, W)
	}

	return m
}

func testVector(doc Doc, weights [][]float64) [][]int64 {
	activations := activation(g(weights, doc.Vector))
	a := make([][]int64, 2)
	a[0] = activations
	a[1] = doc.Test

	return a
}

func spamOrHam(results []int64) string {
	if results[0] == 1 {
		return "spam"
	}
	return "ham"
}

func writeModel(model Model, file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	N := len(model.Weights)
	for j := 0; j < N; j++ {
		w.Write(FloatToStringSlice(model.Weights[j]))
	}
	w.Flush()
}

func loadModel(file string, I int, N int) [][]float64 {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	fl := make([][]float64, N)
	for j := 0; j < N; j++ {
		line, _ := r.Read()
		fl[j] = StringToFloatSlice(line)
	}

	return fl
}

func TrainPerceptron(config *Config) {
	rand.Seed(time.Now().UTC().UnixNano())
	inputs, _ := ReadInput((*config).Inputs)

	var model Model
	if _, err := os.Stat((*config).Perceptron.Model); os.IsNotExist(err) {
		fmt.Println("Random new weights")
		model.Weights = randomWeights(len(inputs), (*config).Neurons)
	} else {
		fmt.Println("Load existing weights")
		model.Weights = loadModel((*config).Perceptron.Model, len(inputs), (*config).Neurons)
	}

	for i := 0; i < (*config).Perceptron.Iterations; i++ {
		model = ProcessVectors((*config).Train, (*config).Neurons, model)
	}

	writeModel(model, (*config).Perceptron.Model)
}

func evaluate(m map[string][][]int64, NEURONS int) {
	total := 0
	matched := 0
	positive := make([]int, 2)
	/*[0, 0]*/
	negative := make([]int, 2)
	/*[0, 0]*/

	for f, v := range m {
		total++

		_result := v[0]
		_test := v[1]
		for j := 0; j < NEURONS; j++ {
			if _result[j] != _test[j] {
				fmt.Println(f, "NOT MATCH:", _result, _test)
				if _result[j] == 1 {
					positive[1]++
				} else {
					negative[1]++
				}
				break
			} else {
				fmt.Println(f, "MATCH")
				if _result[j] == 1 {
					positive[0]++
				} else {
					negative[0]++
				}
				matched++
			}
		}
	}

	fmt.Println("-------------------")
	fmt.Println("    Totals")
	fmt.Println("-------------------")
	fmt.Println("Total Results:",total)
	fmt.Println("Total Matched:",matched)
	fmt.Println("Positive is spam, Negative is ham")
	fmt.Println("Positive [T, F]:", positive)
	fmt.Println("Negative [T, F]:", negative)
	fmt.Println("Error:", int((float64(total - matched) / float64(total)) * 100.0),"%")
}

func TestPerceptron(config *Config) {
	inputs, _ := ReadInput((*config).Inputs)
	var model Model
	model.Weights = loadModel((*config).Perceptron.Model, len(inputs), (*config).Neurons)
	m := TestVectors((*config).Test, (*config).Neurons, model)
	evaluate(m, (*config).Neurons)
}

func RunPerceptron(config *Config, vector []string) string {
	inputs, _ := ReadInput((*config).Inputs)
	var model Model
	model.Weights = loadModel((*config).Perceptron.Model, len(inputs), (*config).Neurons)
	doc := lineToDoc(vector, (*config).Neurons)
	return spamOrHam(testVector(doc, model.Weights)[0])
}

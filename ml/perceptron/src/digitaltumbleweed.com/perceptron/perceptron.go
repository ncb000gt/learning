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

func activation(NEURONS []float64) []int {
	N := len(NEURONS)
	ACTIVATED := make([]int, N)

	for j := 0; j < N; j++ {
		active := 0

		if NEURONS[j] > .75 {
			active = 1
		}

		ACTIVATED[j] = active
	}

	return ACTIVATED
}

var LEARNING_RATE = 0.1
func adjustWeights(W [][]float64, V []float64, Y []int, T []int64) [][]float64 {
	N := len(W)
	I := len(V)
	NW := make([][]float64, N)

	for j := 0; j < N; j++ {
		NWn := make([]float64, I)
		for i := 0; i < I; i++ {
			NWn[i] = W[j][i] + (LEARNING_RATE * float64(T[j] - int64(Y[j])) * V[i])
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

type Model struct {
	Weights [][]float64
}

type VectorFunc func(Doc)

func ProcessVectors(file string, NEURONS int, vf VectorFunc) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	for {
		line, err := ReadLine(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

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

		vf(doc)
	}
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
		ProcessVectors((*config).Train, (*config).Neurons, func(doc Doc) {
			model.Weights = adjustWeights(model.Weights, doc.Vector, activation(g(model.Weights, doc.Vector)), doc.Test)
		})
	}

	writeModel(model, (*config).Perceptron.Model)
	fmt.Println(model)
}

func TestPerceptron(config *Config) {
}

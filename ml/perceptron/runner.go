package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	perceptron "digitaltumbleweed.com/perceptron"
)

func main() {
	configFile := goopt.String([]string{"--config"}, "./config.json", "Configuration File")
	var action = goopt.String([]string{"--action"}, "", "Action to run")
	goopt.Description = func() string {
		return "Perceptron 2.0"
	}
	goopt.Version = "2.0"
	goopt.Summary = "Perceptron"
	goopt.Parse(nil)

	fmt.Println("WUUUT? ", *configFile)
	fmt.Println("Action: ", *action)
	json := perceptron.ReadConfig(*configFile)
	if *action == "preprocess" {
		perceptron.RunPreprocessor(&json)
	} else if *action == "train" {
		perceptron.TrainPerceptron(&json)
	} else if *action == "test" {
		perceptron.TestPerceptron(&json)
	}
}

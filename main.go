package main

import (
	"fmt"
	"io"
	"os"
	"gopkg.in/yaml.v2"
)

// Summary is a struct that contains the summary of the YAML file
type YAMLData struct {
	Spec ControlsSummary `yaml:"spec"`
}

type ControlsSummary struct {
	Severities Severities `yaml:"severities"`
}

type Severities struct {
	Critical map[string]SeverityControls `yaml:"critical"`
	High     map[string]SeverityControls `yaml:"high"`
	Medium   map[string]SeverityControls `yaml:"medium"`
	Low      map[string]SeverityControls `yaml:"low"`
	Other    map[string]SeverityControls `yaml:"other"`
}

type SeverityControls struct {
	All      int `yaml:"all"`
	Relevant int `yaml:"relevant"`
}

func (sc *SeverityControls) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values map[string]int
	err := unmarshal(&values)
	if err == nil {
		sc.All = values["all"]
		sc.Relevant = values["relevant"]
		return nil
	}

	return unmarshal(&sc.All)
}

func parseYAML(fileContent []byte) (*ControlsSummary, error) {
	var data YAMLData
	err := yaml.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}
	return &data.Spec, nil
}

func main() {
	// Read the sample YAML file
	file, err := os.Open("test.yaml")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Parse the YAML file
	summary, err := parseYAML(yamlFile)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	// Print the controls for each category
	fmt.Println("Summary:")
	printControls(summary.Severities.Critical, "Critical")
	printControls(summary.Severities.High, "High")
	printControls(summary.Severities.Medium, "Medium")
	printControls(summary.Severities.Low, "Low")
	printControls(summary.Severities.Other, "Other")
}

func printControls(controlsMap map[string]SeverityControls, category string) {
	fmt.Printf("Category: %s\n", category)
	for controlName, control := range controlsMap {
		if control.Relevant == 0 {
		fmt.Printf("%s: %d,\n", controlName, control.All)
		} else {
		fmt.Printf("%s: %d,\n", controlName, control.Relevant)
		}
	}
	fmt.Println()
}

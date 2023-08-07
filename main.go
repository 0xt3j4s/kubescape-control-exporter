package main

import (
	"fmt"
	"io"
	"os"
	"gopkg.in/yaml.v2"
)

type YAMLData struct {
	Spec ControlsSummary `yaml:"spec"`
}

type ControlsSummary struct {
	Severities Severities `yaml:"severities"`
}

type Severities struct {
	Critical map[string]interface{} `yaml:"critical"`
	High     map[string]interface{} `yaml:"high"`
	Medium   map[string]interface{} `yaml:"medium"`
	Low      map[string]interface{} `yaml:"low"`
	Other    map[string]interface{} `yaml:"other"`
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

	fmt.Printf("Summary: %v\n", summary)

}



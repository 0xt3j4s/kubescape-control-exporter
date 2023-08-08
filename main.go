package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
	"gopkg.in/yaml.v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	Unknown   map[string]SeverityControls `yaml:"unknown"`
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

// Define the metrics
var (
	controlsClusterCritical = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "controls_cluster_critical",
			Help: "Number of critical severity controls in the cluster",
		},
		[]string{"scope", "severity"},)
	
	controlsClusterHigh = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "controls_cluster_high",
			Help: "Number of high severity controls in the cluster",
		},
		[]string{"scope", "severity"},)
	
	controlsClusterMedium = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "controls_cluster_medium",
			Help: "Number of medium severity controls in the cluster",
		},
		[]string{"scope", "severity"},)
	
	controlsClusterLow = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "controls_cluster_low",
			Help: "Number of low severity controls in the cluster",
		},
		[]string{"scope", "severity"},)
	
	controlsClusterUnknown = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "controls_cluster_unknown",
			Help: "Number of unknown severity controls in the cluster",
		},
		[]string{"scope", "severity"},)

)

func updateMetrics(summary ControlsSummary) {
	for severity, controls := range summary.Severities.Critical {
		controlsClusterCritical.WithLabelValues("Critical", severity).Add(float64(controls.All))
	}
	for severity, controls := range summary.Severities.High {
		controlsClusterHigh.WithLabelValues("High", severity).Add(float64(controls.All))
	}
	for severity, controls := range summary.Severities.Medium {
		controlsClusterMedium.WithLabelValues("Medium", severity).Add(float64(controls.All))
	}
	for severity, controls := range summary.Severities.Low {
		controlsClusterLow.WithLabelValues("Low", severity).Add(float64(controls.All))
	}
	for severity, controls := range summary.Severities.Unknown {
		controlsClusterUnknown.WithLabelValues("Unknown", severity).Add(float64(controls.All))
	}
}

func main() {
	prometheus.MustRegister(controlsClusterCritical)
	prometheus.MustRegister(controlsClusterHigh)
	prometheus.MustRegister(controlsClusterMedium)
	prometheus.MustRegister(controlsClusterLow)
	prometheus.MustRegister(controlsClusterUnknown)

	// Read the sample YAML file
	file, err := os.Open("VulnerabilityManifestSummary.yaml")
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
	printControls(summary.Severities.Unknown, "Unknown")

	updateMetrics(*summary)

	// Expose metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("\nExporter is running at :8070/metrics")
	log.Fatal(http.ListenAndServe(":8070", nil))
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
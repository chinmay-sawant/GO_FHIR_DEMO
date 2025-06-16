package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chinmay/junit-html-generator/internal/generator"
	"github.com/chinmay/junit-html-generator/internal/parser"
)

func main() {
	var (
		inputFile  = flag.String("input", "junit-report.xml", "Path to JUnit XML file")
		outputDir  = flag.String("output", ".", "Output directory for generated HTML report")
		title      = flag.String("title", "JUnit Test Report", "Title for the HTML report")
		standalone = flag.Bool("standalone", false, "Generate standalone HTML with embedded CSS/JS")
	)
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Input file is required")
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", *inputFile)
	}

	// Parse the JUnit XML
	testSuites, err := parser.ParseJUnitXML(*inputFile)
	if err != nil {
		log.Fatalf("Failed to parse JUnit XML: %v", err)
	}

	// Generate the HTML report
	reportGenerator := generator.New(*title, *standalone)
	if err := reportGenerator.Generate(testSuites, *outputDir); err != nil {
		log.Fatalf("Failed to generate HTML report: %v", err)
	}

	fmt.Printf("HTML report generated successfully in: %s\n", *outputDir)
	fmt.Printf("Open %s in your browser to view the report\n", filepath.Join(*outputDir, "index.html"))
}

package parser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/chinmay/junit-html-generator/internal/models"
)

// forbiddenDirs is a set of directory names to skip (lowercase)
var forbiddenDirs = func() map[string]struct{} {
	dirs := map[string]struct{}{
		"config":     {},
		"docs":       {},
		"mocks":      {},
		"routes":     {},
		"domain":     {},
		"middleware": {},
		"database":   {},
		"fhirclient": {},
		"logger":     {},
		"utils":      {},
	}
	if name, err := GetCurrentFolderName(); err == nil {
		dirs[lower(name)] = struct{}{}
	}
	fmt.Println("Forbidden directories:")
	for dir := range dirs {
		fmt.Printf("  - %s\n", dir)
	}
	return dirs
}()

// GetCurrentFolderName returns the name of the current working directory
func GetCurrentFolderName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	foldername := filepath.Base(dir)
	foldername = lower(foldername)
	foldername = string([]rune(foldername))
	foldername = replaceUnderscoreWithDash(foldername)
	return foldername, nil
}

// ParseJUnitXML parses a JUnit XML file and returns the test suites
func ParseJUnitXML(filePath string) (*models.TestSuites, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var testSuites models.TestSuites
	if err := xml.Unmarshal(data, &testSuites); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &testSuites, nil
}

// CalculateAnalytics calculates analytics data from test suites
func CalculateAnalytics(testSuites *models.TestSuites) *models.Analytics {
	analytics := &models.Analytics{
		TotalTests:    0,
		TotalFailures: 0,
		TotalErrors:   0,
		TotalTime:     0,
	}

	var testedSuites, untestedSuites []models.TestSuite
	var allTestTimes []float64
	failurePatterns := make(map[string]int)

	for _, suite := range testSuites.TestSuites {
		if shouldSkipSuite(suite.Name) {
			continue
		}
		analytics.TotalTests += suite.Tests
		analytics.TotalFailures += suite.Failures
		analytics.TotalErrors += suite.Errors
		analytics.TotalTime += suite.Time

		if suite.Tests > 0 {
			testedSuites = append(testedSuites, suite)

			// Collect test times and failure patterns
			for _, testCase := range suite.TestCases {
				allTestTimes = append(allTestTimes, testCase.Time)

				if testCase.Failure != nil {
					pattern := extractPattern(testCase.Failure.Message)
					failurePatterns[pattern]++
				}
				if testCase.Error != nil {
					pattern := extractPattern(testCase.Error.Message)
					failurePatterns[pattern]++
				}
			}
		} else {
			untestedSuites = append(untestedSuites, suite)
		}
	}

	// Sort suites by test count
	sort.Slice(testedSuites, func(i, j int) bool {
		return testedSuites[i].Tests > testedSuites[j].Tests
	})

	sort.Slice(untestedSuites, func(i, j int) bool {
		return untestedSuites[i].Name < untestedSuites[j].Name
	})

	analytics.TestedSuites = testedSuites
	analytics.UntestedSuites = untestedSuites
	analytics.FailurePatterns = failurePatterns

	// Calculate reliability score
	if analytics.TotalTests > 0 {
		passed := analytics.TotalTests - analytics.TotalFailures - analytics.TotalErrors
		analytics.ReliabilityScore = float64(passed) / float64(analytics.TotalTests) * 100
	}

	// Calculate average test time
	if analytics.TotalTests > 0 {
		analytics.AvgTestTime = (analytics.TotalTime / float64(analytics.TotalTests)) * 1000 // in milliseconds
	}

	// Calculate coverage percentage
	totalModules := len(testedSuites) + len(untestedSuites)
	testedModules := len(testedSuites)
	if totalModules > 0 {
		analytics.CoveragePercent = float64(testedModules) / float64(totalModules) * 100
	}

	// Calculate performance data
	analytics.PerformanceData = calculatePerformanceData(allTestTimes, testedSuites)

	// Generate recommendations
	analytics.Recommendations = generateRecommendations(analytics)

	return analytics
}

func extractPattern(message string) string {
	if message == "" {
		return "Unknown"
	}

	// Extract first part before colon or take first 50 characters
	for i, char := range message {
		if char == ':' {
			return message[:i]
		}
	}

	if len(message) > 50 {
		return message[:50]
	}

	return message
}

func calculatePerformanceData(allTestTimes []float64, testedSuites []models.TestSuite) models.PerformanceData {
	perfData := models.PerformanceData{}

	// Categorize tests by execution time
	for _, time := range allTestTimes {
		if time < 0.1 {
			perfData.FastTests++
		} else if time < 1.0 {
			perfData.MediumTests++
		} else {
			perfData.SlowTests++
		}
	}

	// Calculate test velocity (tests per second)
	totalTime := 0.0
	totalTests := 0
	for _, suite := range testedSuites {
		totalTime += suite.Time
		totalTests += suite.Tests
	}

	if totalTime > 0 {
		perfData.TestVelocity = float64(totalTests) / totalTime
	}

	// Find bottlenecks (slowest suites)
	sort.Slice(testedSuites, func(i, j int) bool {
		return testedSuites[i].Time > testedSuites[j].Time
	})

	maxBottlenecks := 5
	if len(testedSuites) < maxBottlenecks {
		maxBottlenecks = len(testedSuites)
	}

	for i := 0; i < maxBottlenecks; i++ {
		if testedSuites[i].Time > 0 {
			perfData.Bottlenecks = append(perfData.Bottlenecks, models.Bottleneck{
				Name: testedSuites[i].Name,
				Time: testedSuites[i].Time,
			})
		}
	}

	return perfData
}

func generateRecommendations(analytics *models.Analytics) []string {
	var recommendations []string

	if analytics.AvgTestTime > 500 {
		recommendations = append(recommendations, "⚡ Consider optimizing test setup/teardown - average test time is high")
	}

	if len(analytics.PerformanceData.Bottlenecks) > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🎯 Focus on optimizing %s - it's the slowest module",
				analytics.PerformanceData.Bottlenecks[0].Name))
	}

	if analytics.ReliabilityScore < 90 {
		recommendations = append(recommendations, "🔧 Improve test stability - reliability score below 90%")
	}

	if len(analytics.UntestedSuites) > len(analytics.TestedSuites) {
		recommendations = append(recommendations, "📈 Increase test coverage - more modules without tests than with tests")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ All metrics look good!")
	}

	return recommendations
}

// shouldSkipSuite returns true if the suite name contains any forbidden directory as a path segment
func shouldSkipSuite(suiteName string) bool {

	// Normalize path separators and lowercase
	suiteName = lower(suiteName)
	suiteName = replaceUnderscoreWithDash(suiteName)

	// Extract the last path segment (directory or file name)
	var last string
	if idx := len(suiteName) - 1; idx >= 0 && (suiteName[idx] == '/' || suiteName[idx] == '\\') {
		suiteName = suiteName[:idx]
	}
	fmt.Println("Normalized suite name:", suiteName)
	last = filepath.Base(suiteName)
	fmt.Println("Last part of suite name:", last)
	if _, forbidden := forbiddenDirs[last]; forbidden {
		return true
	}

	println("Last part of suite name:", last)
	if _, forbidden := forbiddenDirs[last]; forbidden {
		return true
	}
	return false
}

func lower(s string) string {
	// Fast path for ASCII
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func replaceUnderscoreWithDash(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r == '_' {
			result[i] = '-'
		} else {
			result[i] = r
		}
	}
	return string(result)
}

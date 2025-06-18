package models

import (
	"encoding/xml"
	"time"
)

// TestSuites represents the root element of JUnit XML
type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Name       string      `xml:"name,attr,omitempty"`
	Tests      int         `xml:"tests,attr"`
	Failures   int         `xml:"failures,attr"`
	Errors     int         `xml:"errors,attr"`
	Time       float64     `xml:"time,attr"`
	Timestamp  string      `xml:"timestamp,attr,omitempty"`
	TestSuites []TestSuite `xml:"testsuite"`
}

// TestSuite represents a test suite
type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr"`
	Tests      int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Errors     int        `xml:"errors,attr"`
	Time       float64    `xml:"time,attr"`
	Timestamp  string     `xml:"timestamp,attr,omitempty"`
	TestCases  []TestCase `xml:"testcase"`
	Properties []Property `xml:"properties>property"`
}

// TestCase represents a single test case
type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Name      string   `xml:"name,attr"`
	ClassName string   `xml:"classname,attr,omitempty"`
	Time      float64  `xml:"time,attr"`
	Failure   *Failure `xml:"failure,omitempty"`
	Error     *Error   `xml:"error,omitempty"`
	Skipped   *Skipped `xml:"skipped,omitempty"`
}

// Failure represents a test failure
type Failure struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// Error represents a test error
type Error struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// Skipped represents a skipped test
type Skipped struct {
	Message string `xml:"message,attr,omitempty"`
}

// Property represents a property in the test suite
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// Analytics represents calculated analytics data
type Analytics struct {
	TotalTests       int             `json:"totalTests"`
	TotalFailures    int             `json:"totalFailures"`
	TotalErrors      int             `json:"totalErrors"`
	TotalTime        float64         `json:"totalTime"`
	ReliabilityScore float64         `json:"reliabilityScore"`
	AvgTestTime      float64         `json:"avgTestTime"`
	SlowTests        int             `json:"slowTests"`
	CoveragePercent  float64         `json:"coveragePercent"`
	TestedSuites     []TestSuite     `json:"testedSuites"`
	UntestedSuites   []TestSuite     `json:"untestedSuites"`
	PerformanceData  PerformanceData `json:"performanceData"`
	FailurePatterns  map[string]int  `json:"failurePatterns"`
	Recommendations  []string        `json:"recommendations"`
}

// PerformanceData represents performance analysis data
type PerformanceData struct {
	FastTests    int          `json:"fastTests"`
	MediumTests  int          `json:"mediumTests"`
	SlowTests    int          `json:"slowTests"`
	TestVelocity float64      `json:"testVelocity"`
	Bottlenecks  []Bottleneck `json:"bottlenecks"`
}

// Bottleneck represents a performance bottleneck
type Bottleneck struct {
	Name string  `json:"name"`
	Time float64 `json:"time"`
}

// Status returns the status of a test case
func (tc *TestCase) Status() string {
	if tc.Failure != nil {
		return "failed"
	}
	if tc.Error != nil {
		return "error"
	}
	if tc.Skipped != nil {
		return "skipped"
	}
	return "passed"
}

// IsHealthy returns true if the test suite has no failures or errors
func (ts *TestSuite) IsHealthy() bool {
	return ts.Failures == 0 && ts.Errors == 0
}

// PassRate returns the pass rate as a percentage
func (ts *TestSuite) PassRate() float64 {
	if ts.Tests == 0 {
		return 0
	}
	passed := ts.Tests - ts.Failures - ts.Errors
	return float64(passed) / float64(ts.Tests) * 100
}

// ParseTimestamp parses the timestamp string to time.Time
func (ts *TestSuite) ParseTimestamp() (time.Time, error) {
	if ts.Timestamp == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, ts.Timestamp)
}

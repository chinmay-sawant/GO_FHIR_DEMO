package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/chinmay/junit-html-generator/internal/models"
	"github.com/chinmay/junit-html-generator/internal/parser"
)

// Generator handles HTML report generation
type Generator struct {
	title      string
	standalone bool
}

// New creates a new generator instance
func New(title string, standalone bool) *Generator {
	return &Generator{
		title:      title,
		standalone: standalone,
	}
}

// Generate creates the HTML report
func (g *Generator) Generate(testSuites *models.TestSuites, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Calculate analytics
	analytics := parser.CalculateAnalytics(testSuites)

	// Generate HTML file
	if err := g.generateHTML(analytics, outputDir); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	return nil
}

func (g *Generator) generateHTML(analytics *models.Analytics, outputDir string) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
<link rel="stylesheet" href="styles.css" />
<script src="script.js"></script>

<link rel="stylesheet" href="styles.css" />
<script src="script.js"></script>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>

</head>
<body data-theme="dark">
    <div class="container">
        <div class="header">
            <button class="theme-toggle" onclick="toggleTheme()">
                <span id="theme-icon">🌙</span>
            </button>
            <h1>{{.Title}}</h1>
            <p>Comprehensive test results and coverage analysis</p>
        </div>

        <div id="content">
            <div class="summary">
                <div class="summary-card">
                    <div class="number tests">{{.Analytics.TotalTests}}</div>
                    <div class="label">Total Tests</div>
                </div>
                <div class="summary-card">
                    <div class="number failures">{{.Analytics.TotalFailures}}</div>
                    <div class="label">Failures</div>
                </div>
                <div class="summary-card">
                    <div class="number errors">{{.Analytics.TotalErrors}}</div>
                    <div class="label">Errors</div>
                </div>
                <div class="summary-card">
                    <div class="number time">{{printf "%.3f" .Analytics.TotalTime}}s</div>
                    <div class="label">Total Time</div>
                </div>
                <div class="summary-card">
                    <div class="number coverage">{{printf "%.0f" .Analytics.CoveragePercent}}%</div>
                    <div class="label">Coverage</div>
                </div>
                <div class="summary-card">
                    <div class="number reliability">{{printf "%.0f" .Analytics.ReliabilityScore}}%</div>
                    <div class="label">Reliability Score</div>
                </div>
                <div class="summary-card">
                    <div class="number performance">{{printf "%.0f" .Analytics.AvgTestTime}}ms</div>
                    <div class="label">Avg Test Time</div>
                </div>
                <div class="summary-card">
                    <div class="number flaky">{{.Analytics.PerformanceData.SlowTests}}</div>
                    <div class="label">Slow Tests</div>
                </div>
            </div>

            <div class="tabs">
                <button class="tab active" onclick="showTab('all')">All Modules</button>
                <button class="tab" onclick="showTab('tested')">Tested Modules</button>
                <button class="tab" onclick="showTab('untested')">Untested Modules</button>
                <button class="tab" onclick="showTab('analytics')">Analytics</button>
            </div>

            <div class="tab-content active" id="all">
                <div id="all-suites">{{template "testsuites" .AllSuites}}</div>
            </div>

            <div class="tab-content" id="tested">
                <div id="tested-suites">{{template "testsuites" .Analytics.TestedSuites}}</div>
            </div>

            <div class="tab-content" id="untested">
                <div id="untested-suites">{{template "testsuites" .Analytics.UntestedSuites}}</div>
            </div>

            <div class="tab-content" id="analytics">
                <div class="analytics-grid">
                    <div class="analytics-section">
                        <h3>Test Distribution</h3>
                        <canvas id="testDistributionChart" width="400" height="200"></canvas>
                    </div>
                    <div class="analytics-section">
                        <h3>Performance Analysis</h3>
                        <canvas id="performanceChart" width="400" height="200"></canvas>
                    </div>
                    <div class="analytics-section">
                        <h3>Failure Analysis</h3>
                        <div id="failureAnalysis">{{template "failureanalysis" .Analytics}}</div>
                    </div>
                    <div class="analytics-section">
                        <h3>Performance Insights</h3>
                        <div id="performanceInsights">{{template "insights" .Analytics}}</div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        const analyticsData = {{.AnalyticsJSON}};
        // Chart.js Test Distribution Chart
        document.addEventListener("DOMContentLoaded", function() {
            if (window.Chart && analyticsData && analyticsData.TestedSuites && analyticsData.UntestedSuites) {
                const ctx = document.getElementById('testDistributionChart').getContext('2d');
                const tested = analyticsData.TestedSuites.length;
                const untested = analyticsData.UntestedSuites.length;
                new Chart(ctx, {
                    type: 'doughnut',
                    data: {
                        labels: ['Tested', 'Untested'],
                        datasets: [{
                            data: [tested, untested],
                            backgroundColor: ['#4caf50', '#f44336'],
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: {
                            legend: { position: 'bottom' }
                        }
                    }
                });
            }
        });
        {{if not .Standalone}}
        {{.JS}}
        {{else}}
        {{.InlineJS}}
        {{end}}
    </script>
</body>
</html>

{{define "testsuites"}}
{{range .}}
<div class="testsuite">
    <div class="testsuite-header {{if eq .Tests 0}}coverage-0{{else if .IsHealthy}}coverage-100{{else}}coverage-partial{{end}}" onclick="toggleTestsuite(this)">
        <div class="testsuite-title">
            <h3>{{.Name}}</h3>
            <span class="coverage-indicator">
                {{if eq .Tests 0}}No Tests{{else if .IsHealthy}}100% Pass{{else}}{{printf "%.0f" .PassRate}}% Pass{{end}}
            </span>
        </div>
        <div class="testsuite-stats">
            <span class="stat">Tests: {{.Tests}}</span>
            {{if gt .Failures 0}}<span class="stat" style="color: var(--accent-red);">Failures: {{.Failures}}</span>{{end}}
            {{if gt .Errors 0}}<span class="stat" style="color: var(--accent-orange);">Errors: {{.Errors}}</span>{{end}}
            <span class="stat">Time: {{printf "%.3f" .Time}}s</span>
        </div>
    </div>
    <div class="testcases">
        {{range .TestCases}}
        <div class="testcase {{.Status}}">
            <div class="testcase-name">
                <span class="status-icon {{.Status}}">
                    {{if eq .Status "passed"}}✓{{else if eq .Status "failed"}}✗{{else}}!{{end}}
                </span>
                {{.Name}}
            </div>
            <div class="testcase-time">{{printf "%.3f" .Time}}s</div>
        </div>
        {{end}}
        {{if eq (len .TestCases) 0}}
        <div class="no-coverage">No tests found in this module</div>
        {{end}}
    </div>
</div>
{{end}}
{{if eq (len .) 0}}
<div class="no-coverage">No modules found in this category</div>
{{end}}
{{end}}

{{define "failureanalysis"}}
{{if gt (len .FailurePatterns) 0}}
<div class="failure-patterns">
    <h4>Common Failure Patterns</h4>
    {{range $pattern, $count := .FailurePatterns}}
    <div class="failure-pattern">
        <span class="pattern-text">{{$pattern}}</span>
        <span class="pattern-count">{{$count}} occurrences</span>
    </div>
    {{end}}
</div>
{{else}}
<div class="metric-card success">🎉 No failing tests detected!</div>
{{end}}
{{end}}

{{define "insights"}}
<div class="insights-section">
    <h4>Performance Bottlenecks</h4>
    {{if gt (len .PerformanceData.Bottlenecks) 0}}
    <div class="bottleneck-list">
        {{range .PerformanceData.Bottlenecks}}
        <div class="bottleneck-item">
            <span class="bottleneck-name">{{.Name}}</span>
            <span class="bottleneck-time">{{printf "%.2f" .Time}}s</span>
        </div>
        {{end}}
    </div>
    {{else}}
    <p>No significant performance bottlenecks detected.</p>
    {{end}}
</div>
<div class="insights-section">
    <h4>Recommendations</h4>
    <div class="recommendations">
        {{range .Recommendations}}
        <div class="recommendation{{if eq . "✅ All metrics look good!"}} success{{end}}">{{.}}</div>
        {{end}}
    </div>
</div>
{{end}}`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	allSuites := append(analytics.TestedSuites, analytics.UntestedSuites...)
	analyticsJSON, _ := json.Marshal(analytics)

	data := struct {
		Title         string
		Standalone    bool
		Analytics     *models.Analytics
		AllSuites     []models.TestSuite
		AnalyticsJSON template.JS
		CSS           template.CSS
		JS            template.JS
		InlineJS      template.JS
	}{
		Title:         g.title,
		Standalone:    g.standalone,
		Analytics:     analytics,
		AllSuites:     allSuites,
		AnalyticsJSON: template.JS(analyticsJSON),
	}

	// Create HTML file
	htmlFile := filepath.Join(outputDir, "index.html")
	file, err := os.Create(htmlFile)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	return t.Execute(file, data)
}

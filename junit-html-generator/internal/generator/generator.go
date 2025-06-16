package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
	if err := g.generateHTML(testSuites, analytics, outputDir); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Copy/embed assets based on standalone flag
	if g.standalone {
		return g.generateStandaloneHTML(testSuites, analytics, outputDir)
	} else {
		return g.copyAssets(outputDir)
	}
}

func (g *Generator) generateHTML(testSuites *models.TestSuites, analytics *models.Analytics, outputDir string) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    {{if not .Standalone}}
    <link rel="stylesheet" href="styles.css">
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    {{else}}
    <style>{{.CSS}}</style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    {{end}}
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

func (g *Generator) generateStandaloneHTML(testSuites *models.TestSuites, analytics *models.Analytics, outputDir string) error {
	// Read CSS and JS content and embed them
	// This would require reading the actual CSS and JS files and embedding them
	return nil
}

func (g *Generator) copyAssets(outputDir string) error {
	// Copy CSS file
	cssContent := getCSSContent()
	cssFile := filepath.Join(outputDir, "styles.css")
	if err := ioutil.WriteFile(cssFile, []byte(cssContent), 0644); err != nil {
		return fmt.Errorf("failed to write CSS file: %w", err)
	}

	// Copy JS file
	jsContent := getJSContent()
	jsFile := filepath.Join(outputDir, "script.js")
	if err := ioutil.WriteFile(jsFile, []byte(jsContent), 0644); err != nil {
		return fmt.Errorf("failed to write JS file: %w", err)
	}

	return nil
}

func getCSSContent() string {
	// Return the CSS content from your styles.css file
	return `:root {
  --bg-primary: #0d1117;
  --bg-secondary: #161b22;
  --bg-tertiary: #21262d;
  --text-primary: #f0f6fc;
  --text-secondary: #8b949e;
  --text-muted: #656d76;
  --border-color: #30363d;
  --accent-blue: #58a6ff;
  --accent-green: #3fb950;
  --accent-red: #f85149;
  --accent-orange: #d29922;
  --accent-purple: #a5a5ff;
  --shadow: rgba(0, 0, 0, 0.3);
  --gradient-primary: linear-gradient(135deg, #161b22 0%, #0d1117 100%);
  --gradient-accent: linear-gradient(135deg, #58a6ff 0%, #a5a5ff 100%);
  --hover-bg: rgba(255, 255, 255, 0.05);
}

/* Add more CSS content from your styles.css file here */`
}

func getJSContent() string {
	// Return basic JS functionality
	return `function toggleTheme() {
  const body = document.body;
  const themeIcon = document.getElementById("theme-icon");
  
  if (body.getAttribute("data-theme") === "dark") {
    body.setAttribute("data-theme", "light");
    themeIcon.textContent = "☀️";
    localStorage.setItem("theme", "light");
  } else {
    body.setAttribute("data-theme", "dark");
    themeIcon.textContent = "🌙";
    localStorage.setItem("theme", "dark");
  }
}

function showTab(tabName) {
  document.querySelectorAll(".tab").forEach(tab => tab.classList.remove("active"));
  document.querySelector('[onclick="showTab(\'' + tabName + '\')"]').classList.add("active");
  
  document.querySelectorAll(".tab-content").forEach(content => content.classList.remove("active"));
  document.getElementById(tabName).classList.add("active");
}

function toggleTestsuite(header) {
  const testcases = header.nextElementSibling;
  testcases.classList.toggle("show");
}

document.addEventListener("DOMContentLoaded", () => {
  const savedTheme = localStorage.getItem("theme") || "dark";
  const themeIcon = document.getElementById("theme-icon");
  
  document.body.setAttribute("data-theme", savedTheme);
  themeIcon.textContent = savedTheme === "dark" ? "🌙" : "☀️";
  
  // Initialize charts if Chart.js is available
  if (typeof Chart !== 'undefined') {
    initializeCharts();
  }
});

function initializeCharts() {
  // Initialize charts with analyticsData
  // Add chart initialization code here
}`
}

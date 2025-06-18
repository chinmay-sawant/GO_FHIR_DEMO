# JUnit HTML Report Generator

A Go-based command-line tool that generates beautiful, interactive HTML reports from JUnit XML test results.

## Features

- 📊 **Advanced Analytics**: Reliability scores, performance metrics, failure analysis
- 📈 **Interactive Charts**: Test distribution, performance analysis with Chart.js
- 🎨 **Modern UI**: Dark/light theme support, responsive design
- 🚀 **Fast Processing**: Efficient XML parsing and report generation
- 📱 **Mobile Friendly**: Responsive design works on all devices
- 🔍 **Detailed Insights**: Test coverage heatmaps, performance bottlenecks
- 💡 **Smart Recommendations**: Automated suggestions for test improvements

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd junit-html-generator

# Build the application
go build -o junit-html-generator main.go
```

## Usage

### Basic Usage

```bash
# Generate report from junit-report.xml in current directory
./junit-html-generator

# Specify input file and output directory
./junit-html-generator -input /path/to/junit-report.xml -output ./reports

# Generate standalone HTML (embeds CSS/JS)
./junit-html-generator -input junit-report.xml -output ./reports -standalone

# Custom report title
./junit-html-generator -input junit-report.xml -title "My Project Test Report"
```

### Command Line Options

- `-input`: Path to JUnit XML file (default: "junit-report.xml")
- `-output`: Output directory for generated HTML report (default: ".")
- `-title`: Title for the HTML report (default: "JUnit Test Report")
- `-standalone`: Generate standalone HTML with embedded CSS/JS (default: false)

## Generated Report Features

### Summary Dashboard
- Total tests, failures, errors, and execution time
- Test coverage percentage
- Reliability score and average test time
- Slow test identification

### Tabbed Views
- **All Modules**: Complete overview of all test suites
- **Tested Modules**: Modules with test cases
- **Untested Modules**: Modules without tests
- **Analytics**: Advanced analytics and visualizations

### Analytics Dashboard
- **Test Distribution Chart**: Visual breakdown of test results
- **Performance Analysis**: Execution time analysis by module
- **Failure Analysis**: Common failure patterns identification
- **Test Health Metrics**: Performance distribution and test velocity
- **Coverage Heatmap**: Visual representation of module coverage
- **Performance Insights**: Bottleneck identification and recommendations

## Project Structure

```
junit-html-generator/
├── main.go                    # Application entry point
├── internal/
│   ├── models/
│   │   └── junit.go          # Data models for JUnit XML
│   ├── parser/
│   │   └── parser.go         # XML parsing and analytics
│   └── generator/
│       └── generator.go      # HTML report generation
├── go.mod                    # Go module file
└── README.md                # This file
```

## Example Output

The generated HTML report includes:

1. **Interactive Summary Cards** showing key metrics
2. **Expandable Test Suites** with individual test case details
3. **Visual Charts** for test distribution and performance
4. **Smart Recommendations** based on test results analysis
5. **Theme Toggle** for dark/light mode switching

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]
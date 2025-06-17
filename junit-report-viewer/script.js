
function showFileInput() {
  const loadingElem = document.getElementById("loading");
  if (loadingElem) {
    loadingElem.style.display = "none";
  }
  const fileInputSection = document.getElementById("file-input-section");
  if (fileInputSection) {
    fileInputSection.style.display = "block";
  }
  setupFileInput();
}

function setupFileInput() {
  const fileInput = document.getElementById("fileInput");
  const dropZone = document.getElementById("dropZone");

  // Only add event listeners if the elements exist
  if (fileInput) {
    fileInput.addEventListener("change", handleFileSelect);
  }
  if (dropZone) {
    dropZone.addEventListener("dragover", (e) => {
      e.preventDefault();
      dropZone.classList.add("dragover");
    });

    dropZone.addEventListener("dragleave", () => {
      dropZone.classList.remove("dragover");
    });

    dropZone.addEventListener("drop", (e) => {
      e.preventDefault();
      dropZone.classList.remove("dragover");
      const files = e.dataTransfer.files;
      if (files.length > 0) {
        handleFile(files[0]);
      }
    });
  }
}

function handleFileSelect(event) {
  const file = event.target.files[0];
  if (file) {
    handleFile(file);
  }
}

function handleFile(file) {
  if (!file.name.endsWith(".xml")) {
    showError("Please select an XML file");
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    try {
      parseXMLAndDisplay(e.target.result);
    } catch (error) {
      showError(`Failed to parse XML: ${error.message}`);
    }
  };
  reader.onerror = () => {
    showError("Failed to read file");
  };
  reader.readAsText(file);
}

function parseXMLAndDisplay(xmlText) {
  const parser = new DOMParser();
  const xmlDoc = parser.parseFromString(xmlText, "text/xml");

  if (xmlDoc.querySelector("parsererror")) {
    throw new Error("Invalid XML format");
  }

  parseAndDisplayReport(xmlDoc);
}

function toggleTheme() {
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

function initTheme() {
  const savedTheme = localStorage.getItem("theme") || "dark";
  const themeIcon = document.getElementById("theme-icon");

  document.body.setAttribute("data-theme", savedTheme);
  themeIcon.textContent = savedTheme === "dark" ? "🌙" : "☀️";
}

function showTab(tabName) {
  // Update tabs
  document
    .querySelectorAll(".tab")
    .forEach((tab) => tab.classList.remove("active"));
  document
    .querySelector(`[onclick="showTab('${tabName}')"]`)
    .classList.add("active");

  // Update content
  document
    .querySelectorAll(".tab-content")
    .forEach((content) => content.classList.remove("active"));
  document.getElementById(tabName).classList.add("active");
}

function parseAndDisplayReport(xmlDoc) {
  const testsuites = xmlDoc.querySelector("testsuites");
  if (!testsuites) {
    showError("Invalid JUnit XML format: No testsuites found");
    return;
  }

  // Hide file input section
  document.getElementById("file-input-section").style.display = "none";

  const totalTests = parseInt(testsuites.getAttribute("tests") || "0");
  const totalFailures = parseInt(
    testsuites.getAttribute("failures") || "0"
  );
  const totalErrors = parseInt(testsuites.getAttribute("errors") || "0");
  const totalTime = parseFloat(testsuites.getAttribute("time") || "0");

  // Update summary
  document.getElementById("total-tests").textContent = totalTests;
  document.getElementById("total-failures").textContent = totalFailures;
  document.getElementById("total-errors").textContent = totalErrors;
  document.getElementById("total-time").textContent =
    totalTime.toFixed(3) + "s";

  // Parse test suites
  const testsuiteElements = xmlDoc.querySelectorAll("testsuite");
  const testedSuites = [];
  const untestedSuites = [];

  testsuiteElements.forEach((testsuite) => {
    const name = testsuite.getAttribute("name");
    const tests = parseInt(testsuite.getAttribute("tests") || "0");
    const failures = parseInt(testsuite.getAttribute("failures") || "0");
    const errors = parseInt(testsuite.getAttribute("errors") || "0");
    const time = parseFloat(testsuite.getAttribute("time") || "0");
    const timestamp = testsuite.getAttribute("timestamp");

    const suiteData = {
      name,
      tests,
      failures,
      errors,
      time,
      timestamp,
      element: testsuite,
    };

    if (tests > 0) {
      testedSuites.push(suiteData);
    } else {
      untestedSuites.push(suiteData);
    }
  });

  // Calculate coverage
  const totalModules = testsuiteElements.length;
  const testedModules = testedSuites.length;
  const coveragePercentage =
    totalModules > 0
      ? Math.round((testedModules / totalModules) * 100)
      : 0;
  document.getElementById("coverage-percentage").textContent =
    coveragePercentage + "%";

  // Sort by test count (descending)
  testedSuites.sort((a, b) => b.tests - a.tests);
  untestedSuites.sort((a, b) => a.name.localeCompare(b.name));

  // Calculate advanced analytics
  calculateAdvancedAnalytics(
    testedSuites,
    untestedSuites,
    testsuiteElements
  );

  // Render suites
  renderTestSuites([...testedSuites, ...untestedSuites], "all-suites");
  renderTestSuites(testedSuites, "tested-suites");
  renderTestSuites(untestedSuites, "untested-suites");

  // Show content
  document.getElementById("content").style.display = "block";
}

function calculateAdvancedAnalytics(
  testedSuites,
  untestedSuites,
  testsuiteElements
) {
  // Calculate reliability score
  const totalTests = testedSuites.reduce(
    (sum, suite) => sum + suite.tests,
    0
  );
  const totalFailures = testedSuites.reduce(
    (sum, suite) => sum + suite.failures + suite.errors,
    0
  );
  const reliabilityScore =
    totalTests > 0
      ? Math.round(((totalTests - totalFailures) / totalTests) * 100)
      : 0;

  // Calculate average test time
  const totalTime = testedSuites.reduce(
    (sum, suite) => sum + suite.time,
    0
  );
  const avgTestTime =
    totalTests > 0 ? Math.round((totalTime / totalTests) * 1000) : 0;

  // Find slow tests (top 10% by time)
  const allTestTimes = [];
  testedSuites.forEach((suite) => {
    const testcases = suite.element.querySelectorAll("testcase");
    testcases.forEach((testcase) => {
      const time = parseFloat(testcase.getAttribute("time") || "0");
      allTestTimes.push(time);
    });
  });
  allTestTimes.sort((a, b) => b - a);
  const slowTestThreshold =
    allTestTimes[Math.floor(allTestTimes.length * 0.1)] || 0;
  const slowTests = allTestTimes.filter(
    (time) => time >= slowTestThreshold
  ).length;

  // Update summary cards
  document.getElementById("reliability-score").textContent =
    reliabilityScore + "%";
  document.getElementById("avg-test-time").textContent =
    avgTestTime + "ms";
  document.getElementById("flaky-tests").textContent = slowTests;

  // Store analytics data
  analyticsData = {
    testedSuites,
    untestedSuites,
    totalTests,
    totalFailures,
    reliabilityScore,
    avgTestTime,
    slowTests,
    slowTestThreshold,
    allTestTimes,
  };

  // Render analytics
  renderAnalytics();
}

function renderAnalytics() {
  renderTestDistributionChart();
  renderPerformanceChart();
  renderFailureAnalysis();
  renderHealthMetrics();
  renderCoverageHeatmap();
  renderPerformanceInsights();
}

function renderTestDistributionChart() {
  const ctx = document
    .getElementById("testDistributionChart")
    .getContext("2d");
  const data = {
    labels: ["Passing", "Failing", "Error", "No Tests"],
    datasets: [
      {
        data: [
          analyticsData.totalTests - analyticsData.totalFailures,
          analyticsData.testedSuites.reduce(
            (sum, suite) => sum + suite.failures,
            0
          ),
          analyticsData.testedSuites.reduce(
            (sum, suite) => sum + suite.errors,
            0
          ),
          analyticsData.untestedSuites.length,
        ],
        backgroundColor: ["#4ade80", "#ef4444", "#f97316", "#6b7280"],
        borderWidth: 0,
      },
    ],
  };

  new Chart(ctx, {
    type: "doughnut",
    data: data,
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: "bottom",
          labels: { color: "var(--text-color)" },
        },
      },
    },
  });
}

function renderPerformanceChart() {
  const ctx = document
    .getElementById("performanceChart")
    .getContext("2d");
  const suiteData = analyticsData.testedSuites
    .slice(0, 10)
    .map((suite) => ({
      label:
        suite.name.length > 20
          ? suite.name.substring(0, 20) + "..."
          : suite.name,
      time: suite.time,
      tests: suite.tests,
    }));

  new Chart(ctx, {
    type: "bar",
    data: {
      labels: suiteData.map((s) => s.label),
      datasets: [
        {
          label: "Execution Time (s)",
          data: suiteData.map((s) => s.time),
          backgroundColor: "#3b82f6",
          borderRadius: 4,
        },
      ],
    },
    options: {
      responsive: true,
      scales: {
        y: {
          beginAtZero: true,
          ticks: { color: "var(--text-color)" },
        },
        x: {
          ticks: {
            color: "var(--text-color)",
            maxRotation: 45,
          },
        },
      },
      plugins: {
        legend: {
          labels: { color: "var(--text-color)" },
        },
      },
    },
  });
}

function renderFailureAnalysis() {
  const container = document.getElementById("failureAnalysis");
  const failingSuites = analyticsData.testedSuites.filter(
    (suite) => suite.failures > 0 || suite.errors > 0
  );

  if (failingSuites.length === 0) {
    container.innerHTML =
      '<div class="metric-card success">🎉 No failing tests detected!</div>';
    return;
  }

  const failurePatterns = {};
  failingSuites.forEach((suite) => {
    const testcases = suite.element.querySelectorAll("testcase");
    testcases.forEach((testcase) => {
      const failure = testcase.querySelector("failure");
      const error = testcase.querySelector("error");
      if (failure || error) {
        const message =
          (failure || error).getAttribute("message") || "Unknown";
        const pattern = message.split(":")[0] || message.substring(0, 50);
        failurePatterns[pattern] = (failurePatterns[pattern] || 0) + 1;
      }
    });
  });

  const sortedPatterns = Object.entries(failurePatterns)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 5);

  container.innerHTML = `
    <div class="failure-patterns">
      <h4>Common Failure Patterns</h4>
      ${sortedPatterns
        .map(
          ([pattern, count]) => `
        <div class="failure-pattern">
          <span class="pattern-text">${pattern}</span>
          <span class="pattern-count">${count} occurrences</span>
        </div>
      `
        )
        .join("")}
    </div>
    <div class="failure-summary">
      <div class="metric-card error">
        <div class="metric-value">${failingSuites.length}</div>
        <div class="metric-label">Failing Modules</div>
      </div>
    </div>
  `;
}

function renderHealthMetrics() {
  const container = document.getElementById("healthMetrics");
  const { testedSuites, reliabilityScore, avgTestTime } = analyticsData;

  const fastTests = analyticsData.allTestTimes.filter(
    (time) => time < 0.1
  ).length;
  const mediumTests = analyticsData.allTestTimes.filter(
    (time) => time >= 0.1 && time < 1
  ).length;
  const slowTests = analyticsData.allTestTimes.filter(
    (time) => time >= 1
  ).length;

  const totalExecutionTimeForTestedSuites = testedSuites.reduce(
    (sum, suite) => sum + suite.time,
    0
  );
  const testVelocity =
    analyticsData.totalTests > 0 && totalExecutionTimeForTestedSuites > 0
      ? Math.round(
          analyticsData.totalTests / totalExecutionTimeForTestedSuites
        )
      : 0;

  const totalTrackedTimes = analyticsData.allTestTimes.length;
  const fastPercent =
    totalTrackedTimes > 0 ? (fastTests / totalTrackedTimes) * 100 : 0;
  const mediumPercent =
    totalTrackedTimes > 0 ? (mediumTests / totalTrackedTimes) * 100 : 0;
  const slowPercent =
    totalTrackedTimes > 0 ? (slowTests / totalTrackedTimes) * 100 : 0;

  container.innerHTML = `
    <div class="health-grid">
      <div class="metric-card ${
        reliabilityScore >= 90
          ? "success"
          : reliabilityScore >= 70
          ? "warning"
          : "error"
      }">
        <div class="metric-value">${reliabilityScore}%</div>
        <div class="metric-label">Test Reliability</div>
      </div>
      <div class="metric-card ${
        avgTestTime < 100
          ? "success"
          : avgTestTime < 500
          ? "warning"
          : "error"
      }">
        <div class="metric-value">${avgTestTime}ms</div>
        <div class="metric-label">Avg Test Duration</div>
      </div>
      <div class="metric-card info">
        <div class="metric-value">${testVelocity}</div>
        <div class="metric-label">Tests/Second</div>
      </div>
    </div>
    <div class="performance-breakdown">
      <h4>Performance Distribution</h4>
      <div class="perf-bar">
        <div class="perf-segment fast" style="width: ${fastPercent.toFixed(
          2
        )}%">
          Fast (&lt;100ms): ${fastTests}
        </div>
        <div class="perf-segment medium" style="width: ${mediumPercent.toFixed(
          2
        )}%">
          Medium (100ms-1s): ${mediumTests}
        </div>
        <div class="perf-segment slow" style="width: ${slowPercent.toFixed(
          2
        )}%">
          Slow (&gt;1s): ${slowTests}
        </div>
      </div>
    </div>
  `;
}

function renderCoverageHeatmap() {
  const container = document.getElementById("coverageHeatmap");
  const allSuites = [
    ...analyticsData.testedSuites,
    ...analyticsData.untestedSuites,
  ];

  container.innerHTML = `
    <div class="heatmap-grid">
      ${allSuites
        .slice(0, 20)
        .map((suite) => {
          let coverage = 0;
          let cssClass = "no-coverage";

          if (suite.tests > 0) {
            coverage = Math.round(
              ((suite.tests - suite.failures - suite.errors) /
                suite.tests) *
                100
            );
            if (coverage === 100) cssClass = "full-coverage";
            else if (coverage >= 75) cssClass = "high-coverage";
            else if (coverage >= 50) cssClass = "medium-coverage";
            else cssClass = "low-coverage";
          }

          return `
          <div class="heatmap-cell ${cssClass}" title="${
            suite.name
          }: ${coverage}% success rate">
            <div class="cell-label">${suite.name.substring(0, 8)}...</div>
            <div class="cell-value">${coverage}%</div>
          </div>
        `;
        })
        .join("")}
    </div>
  `;
}

function renderPerformanceInsights() {
  const container = document.getElementById("performanceInsights");
  const { testedSuites, allTestTimes, slowTestThreshold } = analyticsData;

  const bottlenecks = testedSuites
    .filter((suite) => suite.time > slowTestThreshold)
    .sort((a, b) => b.time - a.time)
    .slice(0, 5);

  const recommendations = [];

  if (analyticsData.avgTestTime > 500) {
    recommendations.push(
      "⚡ Consider optimizing test setup/teardown - average test time is high"
    );
  }

  if (bottlenecks.length > 0) {
    recommendations.push(
      `🎯 Focus on optimizing ${bottlenecks[0].name} - it's the slowest module`
    );
  }

  if (analyticsData.reliabilityScore < 90) {
    recommendations.push(
      "🔧 Improve test stability - reliability score below 90%"
    );
  }

  if (
    analyticsData.untestedSuites.length >
    analyticsData.testedSuites.length
  ) {
    recommendations.push(
      "📈 Increase test coverage - more modules without tests than with tests"
    );
  }

  container.innerHTML = `
    <div class="insights-section">
      <h4>Performance Bottlenecks</h4>
      ${
        bottlenecks.length > 0
          ? `
        <div class="bottleneck-list">
          ${bottlenecks
            .map(
              (suite) => `
            <div class="bottleneck-item">
              <span class="bottleneck-name">${suite.name}</span>
              <span class="bottleneck-time">${suite.time.toFixed(
                2
              )}s</span>
            </div>
          `
            )
            .join("")}
        </div>
      `
          : "<p>No significant performance bottlenecks detected.</p>"
      }
    </div>

    <div class="insights-section">
      <h4>Recommendations</h4>
      <div class="recommendations">
        ${
          recommendations.length > 0
            ? recommendations
                .map((rec) => `<div class="recommendation">${rec}</div>`)
                .join("")
            : '<div class="recommendation success">✅ All metrics look good!</div>'
        }
      </div>
    </div>
  `;
}

function renderTestSuites(suites, containerId) {
  const container = document.getElementById(containerId);
  container.innerHTML = "";

  if (suites.length === 0) {
    container.innerHTML =
      '<div class="no-coverage">No modules found in this category</div>';
    return;
  }

  suites.forEach((suite) => {
    const testsuiteDiv = document.createElement("div");
    testsuiteDiv.className = "testsuite";

    // Determine coverage class
    let coverageClass = "coverage-0";
    let coverageText = "No Tests";

    if (suite.tests > 0) {
      if (suite.failures === 0 && suite.errors === 0) {
        coverageClass = "coverage-100";
        coverageText = "100% Pass";
      } else {
        coverageClass = "coverage-partial";
        const passRate = Math.round(
          ((suite.tests - suite.failures - suite.errors) / suite.tests) *
            100
        );
        coverageText = `${passRate}% Pass`;
      }
    }

    let testcasesHtml = "";

    if (suite.tests > 0) {
      const testcases = suite.element.querySelectorAll("testcase");
      testcasesHtml = Array.from(testcases)
        .map((testcase) => {
          const testName = testcase.getAttribute("name");
          const testTime = parseFloat(
            testcase.getAttribute("time") || "0"
          );
          const hasFailure = testcase.querySelector("failure") !== null;
          const hasError = testcase.querySelector("error") !== null;

          let status = "passed";
          let statusIcon = "✓";

          if (hasFailure) {
            status = "failed";
            statusIcon = "✗";
          } else if (hasError) {
            status = "error";
            statusIcon = "!";
          }

          return `
                      <div class="testcase ${status}">
                          <div class="testcase-name">
                              <span class="status-icon ${status}">${statusIcon}</span>
                              ${testName}
                          </div>
                          <div class="testcase-time">${testTime.toFixed(
                            3
                          )}s</div>
                      </div>
                  `;
        })
        .join("");
    } else {
      testcasesHtml =
        '<div class="no-coverage">No tests found in this module</div>';
    }

    testsuiteDiv.innerHTML = `
              <div class="testsuite-header ${coverageClass}" onclick="toggleTestsuite(this)">
                  <div class="testsuite-title">
                      <h3>${suite.name}</h3>
                      <span class="coverage-indicator">${coverageText}</span>
                  </div>
                  <div class="testsuite-stats">
                      <span class="stat">Tests: ${suite.tests}</span>
                      ${
                        suite.failures > 0
                          ? `<span class="stat" style="color: var(--accent-red);">Failures: ${suite.failures}</span>`
                          : ""
                      }
                      ${
                        suite.errors > 0
                          ? `<span class="stat" style="color: var(--accent-orange);">Errors: ${suite.errors}</span>`
                          : ""
                      }
                      <span class="stat">Time: ${suite.time.toFixed(
                        3
                      )}s</span>
                      ${
                        suite.timestamp
                          ? `<span class="stat">${new Date(
                              suite.timestamp
                            ).toLocaleString()}</span>`
                          : ""
                      }
                  </div>
              </div>
              <div class="testcases">
                  ${testcasesHtml}
              </div>
          `;

    container.appendChild(testsuiteDiv);
  });
}

function toggleTestsuite(header) {
  const testcases = header.nextElementSibling;
  testcases.classList.toggle("show");
}

function showError(message) {
  document.getElementById("loading").style.display = "none";
  document.getElementById("file-input-section").style.display = "none";
  document.getElementById("error-message").textContent = message;
  document.getElementById("error-container").style.display = "block";
}

// Initialize theme on page load
document.addEventListener("DOMContentLoaded", () => {
  initTheme();
});

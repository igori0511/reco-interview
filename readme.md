# 🐙 Asana Project and User Data Extractor

This application is a robust Go client designed to extract user-to-project associations from the Asana API, handle potential network transient failures using exponential backoff, and save the data locally. The extraction runs on a timed loop.

## ⚙️ Build and Run

### Prerequisites

You must have the following installed:

1. **Go (1.18+):** The application is written in Go.

2. **Personal Access Token (PAT):** An Asana Personal Access Token configured in the source code (`main.go`).

### Compilation

Navigate to the project root directory and compile the application into an executable named `app`:

go build -o app .

### Execution

Run the compiled executable and specify the desired running interval using the `-interval` flag.

# Run every 5 seconds (for testing/development)
./app -interval=5

# Run every 30 minutes (for production/scheduled use)
./app -interval=1800

## 🛠️ Command-Line Arguments

The application requires a single flag:

| **Flag** | **Description** | **Required Format** | **Allowed Values** |
| :--- | :--- | :--- | :--- |
| `-interval` | The frequency at which the data extraction will run. | Integer (seconds) | **`5`** (seconds) or **`1800`** (seconds, equivalent to 30 minutes) |

### Invalid Argument Handling

If you run the application with an invalid interval, it will immediately exit:

./app -interval=60
# Output: invalid -interval=60 (allowed: 5 or 1800)

## 📂 Output

The application creates a local directory and saves the extracted data as a JSON file.

| **Component** | **Path** | **Description** |
| :--- | :--- | :--- |
| **Directory** | `extracted-asana-data/` | Automatically created on first run. |
| **File** | `extracted-asana-data/asana_data.json` | Contains the map of `Project Name` to `[List of Member Names]`. |

## 🛡️ Reliability (Rate Limiting)

The client uses the `github.com/cenkalti/backoff/v4` library to ensure reliable communication with the Asana API.

* **Transient Errors (500, 503, 429):** The application will retry the request using **exponential backoff** to avoid overwhelming the server.

* **Permanent Errors (401, 404):** The application will fail immediately, as retrying will not resolve the issue.

*Press `Ctrl+C` at any time to gracefully shut down the application.*
# SHELA - SHELA Helps Evaluate Linux Access

![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.18+-blue)
![License](https://img.shields.io/badge/license-MIT-lightgrey)

**SHELA** is a lightweight command-line tool written in Go, designed to help evaluate basic security aspects of a Linux system. It performs various audits such as detecting SUID files, checking user login shells, identifying open ports, locating world-writable directories, and listing crontab entries. SHELA aims to help system administrators quickly spot potential security risks.

## ✨ Features

- Detects files with the SUID bit set
- Audits user accounts with shell access
- Scans for open ports (using `ss`)
- Lists world-writable directories
- Enumerates crontab entries per user
- Interactive progress bar for each step
- Saves a complete report to file

## 🛠️ Building with Makefile

This project includes a `Makefile` to simplify common development tasks.

### Available Makefile Commands

| Command       | Description                                                              |
|---------------|--------------------------------------------------------------------------|
| `make`        | Alias for `make build`                                                   |
| `make build`  | Builds the SHELA binary to the `bin/` directory                          |
| `make run`    | Builds and runs SHELA immediately                                        |
| `make fmt`    | Formats Go source files using `go fmt`                                   |
| `make vet`    | Performs static analysis using `go vet`                                  |
| `make lint`   | Runs `golint` (if installed)                                             |
| `make deps`   | Ensures all dependencies are downloaded via `go mod tidy`                |
| `make clean`  | Removes the `bin/` directory and all build artifacts                     |

## 🚀 How to Use

### 1. Build & Run

```bash
sudo make run
```

To only build the binary:

```bash
make build
```

The binary will be created in the `bin/` directory as `shela`.

### Specify Output Directory (Optional)

By default, the report is saved in the current working directory. To specify a custom output directory:

```bash
sudo ./bin/shela /path/to/output
```

## 📦 Alternative Installation

You can also install SHELA directly using Go (requires Go 1.18+):

```bash
go install github.com/wprimadi/shela@latest
```

## 🔒 Important Notes

- Some checks (e.g., world-writable directories) may require `sudo` privileges.
- Run SHELA as a user with sufficient permissions for best results.

## 📦 Requirements

- Go 1.18 or newer
- Standard Linux utilities: `find`, `ss`, `ls`, etc.

## 📄 License

This project is licensed under the MIT License.
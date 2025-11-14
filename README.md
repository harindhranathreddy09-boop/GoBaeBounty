```markdown
# 🚀 GoBaeBounty - Installation, Setup & Usage Guide

Welcome to **GoBaeBounty**, your comprehensive all-in-one bug bounty automation tool! This step-by-step guide will help you install, configure, and run the tool effectively for authorized security testing.

---

## 🛠️ Installation

### 1. Prerequisites

- **Go 1.21+** installed on your system  
- Basic command line experience  
- Internet connection to download dependencies  
- Optional but recommended tools (for enhanced recon):  
  - `waybackurls`  
  - `gau`  
  - `ffuf`  
  - `nuclei`  
  - `sqlmap`

---

### 2. Clone the Repository

```
git clone https://github.com/harindhranathreddy09-boop/GoBaeBounty.git
cd GoBaeBounty
```

---

### 3. Download Dependencies & Build

```
go mod tidy
make build
```

This will build the `gobaebounty` executable in the `./bin` directory.

---

## ⚙️ Setup

### 4. Authorization File (Mandatory)

Due to legal and ethical considerations, *GoBaeBounty requires an authorization file* to run.

- Locate the example at `configs/auth.example.json`  
- Request or generate your authorization file following your bug bounty program policy  
- Ensure the file matches your target domain and contains a valid signature  

*Note:* The tool enforces validity and expiry checks on this file.

---

### 5. Configuration

Customize scanning parameters in `configs/config.example.yml` or provide your own config path with `--config`.

**Key config options:**

| Option            | Default            | Description                                         |
|-------------------|--------------------|-----------------------------------------------------|
| `target`          | `example.com`      | Target domain to scan                               |
| `workers`         | 50                 | Concurrent workers for crawling and fuzzing        |
| `max_rate`        | 100                | Maximum HTTP requests per second                     |
| `crawl_depth`     | 3                  | Maximum recursion depth in crawling                 |
| `ignore_robots`   | false              | Ignore robots.txt restrictions                       |
| `intrusive`       | false              | Enable intrusive vulnerability checks (use with care) |
| `wordlists`       | preset files       | Wordlists for directory fuzzing and admin discovery |
| `fuzzing`         | enabled            | Parameter fuzzing settings                           |
| `external_tools`  | partial enabled    | Toggle integration of external tools (waybackurls, gau, etc.) |

---

## 🏃 Running the Tool

Basic example command:

```
./bin/gobaebounty \
--target example.com \
--auth-file ./configs/auth.example.json \
--o ./results \
--workers 50 \
--max-rate 100 \
--depth 3 \
--v
```

### Command Flags:

| Flag           | Description                                  | Required | Default        |
|----------------|----------------------------------------------|----------|----------------|
| `--target`     | Domain or IP to scan                          | Yes      | N/A            |
| `--auth-file`  | Authorization JSON file path                  | Yes      | N/A            |
| `--o`          | Output directory for reports and data        | No       | `./results`    |
| `--workers`    | Number of concurrent goroutines               | No       | 50             |
| `--max-rate`   | Max HTTP requests per second                   | No       | 100            |
| `--depth`      | Crawl recursion depth                          | No       | 3              |
| `--ignore-robots` | Ignore robots.txt rules                      | No       | false          |
| `--intrusive`  | Enable intrusive / destructive checks (care) | No       | false          |
| `--v`          | Verbose logging                               | No       | false          |

---

## 📂 Output

By default, output files save under the directory specified by `--o` (e.g., `./results`):

- `report.md` — Human-friendly HackerOne-style Markdown report  
- `report.json` — Machine-readable findings JSON  
- `endpoints.json` — Extracted and scored endpoints  
- `poc/` — Proof-of-concept request and response snippets  
- `screenshots/` (optional) — Captured screenshots if enabled

---

## ⚠️ Legal & Ethical Notice

🔒 **Only scan targets with explicit, written permission!** Unauthorized scanning is illegal and unethical.  
The tool **enforces an authorization file check** to help ensure compliance.

---

## 📚 Troubleshooting & Tips

- If the tool fails due to missing dependencies, install external tools like `waybackurls`, `gau`, or `ffuf` for enhanced recon.  
- Use `--dry-run` flag to check authorization without scanning.  
- Adjust concurrency (`--workers`) and rate limits (`--max-rate`) to avoid overwhelming the target or your network.  
- Review logs with verbose mode (`--v`) for debugging.  
- Customize wordlists for targeted fuzzing to improve scan relevance and speed.

---

## 🙌 Getting Help & Contribution

- For bug reports, feature requests, or code contributions, visit the project GitHub repository:  
  https://github.com/harindhranathreddy09-boop/GoBaeBounty  
- Contributions are welcome! Please follow the Developer Guide.

---

🎯 **Happy and responsible bounty hunting with GoBaeBounty!** 🎯
```

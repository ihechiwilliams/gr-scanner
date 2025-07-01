# gr-scanner

## Introduction
`gr-scanner` is a robust command-line tool written in Go that scans GitHub repositories to identify and report files exceeding a specified size threshold.

## Setup Instructions
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/ihechiwilliams/gr-scanner.git
   cd gr-scanner
   ```
2. **Set Up Environment Variables**:
   Create a `.env` file in the root directory of the project and add your GitHub token:
   ```bash
   cp .env.example .env
   ```
   Then edit the `.env` file to include your GitHub token:
   ```
3.  **Install Dependencies** (for local development):
   ```bash
   go mod tidy
   ```
4. **Build the Application** (optional, for local runs):
   ```bash
   go build -o gr-scanner ./cmd/gr-scanner
   ```
5. **Build Docker Image** (for containerized runs):
   ```bash
   docker build -t gr-scanner .
   ```

### Running Locally
Scan a repository for files larger than 1MB:
```bash
./gr-scanner scan '{"clone_url":"https://github.com/owner/repo.git","size":1.0}'
```
**Output** (JSON):
```json
{
  "total": 2,
  "files": [
    {"name": "code.zip", "size": 1500000},
    {"name": "sub/dir/big_image.png", "size": 1200000}
  ]
}
```

### Running with Docker
```bash
docker run --env-file .env gr-scanner scan '{"clone_url":"https://github.com/owner/repo.git","size":1.0}'
```
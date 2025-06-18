# Project Overview

TRMNL-Markdown-Kanban is a minimal Go service that reads a local Markdown file, parses ATX headings as columns and list items as cards, and exposes a Liquid-templated Kanban board via HTTP. It integrates with TRMNL’s public plugin model using OAuth2 and persists tokens across restarts. README files like this serve as the first point of contact for potential contributors and users  ￼.

# Features
	•	Universal Markdown Parsing: Uses Goldmark to convert any Markdown into a Kanban Board of columns and cards.
	•	TRMNL OAuth2 Integration: Implements the full OAuth2 installation flow and persistent token storage.
	•	Token Persistence: Securely saves access_token and refresh_token in JSON on disk for automatic refresh.
	•	Distroless Deployment: Multi-stage Docker build produces a ~30 MB static binary image.
	•	Clean Liquid UI: Renders the board with a simple, e-ink-friendly Liquid/HTML template.
	•	Self-Hosted: You provide the service; TRMNL polls the /markup endpoint for pre-rendered Kanban boards  ￼.

# Installation
	1.	Clone the repository
``` bash
git clone https://github.com/Derek-Roberts/trmnl-markdown-kanban.git
cd trmnl-markdown-kanban
```
	2.	Build the Docker image
``` bash
docker build -t yourrepo/trmnl-markdown-kanban:latest .
```
	3.	Prepare your data volume
	•	Place kanban.md (your task board) in a host folder, e.g. /volume1/kanban/md.
	•	Create a data directory for tokens, e.g. /volume1/kanban/data.
	4.	Deploy with Docker Compose
``` yaml
version: "3.8"
services:
  kanban:
    image: yourrepo/trmnl-markdown-kanban:latest
    ports:
      - "8080:8080"
    environment:
      TRMNL_CLIENT_ID:     <your-client-id>
      TRMNL_CLIENT_SECRET: <your-client-secret>
      REDIRECT_URL:        https://yourdomain.com/install
      DATA_DIR:            /data
    volumes:
      - /volume1/kanban/md:/data/kanban.md
      - /volume1/kanban/data:/data
    restart: always
```
``` bash
docker-compose up -d
```

# Configuration

Set the following environment variables in your container or host’s Docker Compose:
	•	TRMNL_CLIENT_ID & TRMNL_CLIENT_SECRET
	•	REDIRECT_URL (matches your plugin manifest)
	•	DATA_DIR (defaults to /data)
	•	PORT (optional; defaults to 8080)

# Usage
	1.	Register the plugin on TRMNL’s Developer Dashboard with your manifest, pointing to /install, /markup, and /uninstall URLs.
	2.	Install the plugin in your TRMNL workspace—this triggers /install, exchanges the code, and persists tokens.
	3.	Update your kanban.md file; TRMNL will poll /markup to render the latest board.

# Development
	•	Go Module:
``` bash
go mod tidy
go build ./cmd/server
```

	•	Local Test:
``` bash
DATA_DIR=./test-data \
  TRMNL_CLIENT_ID=foo \
  TRMNL_CLIENT_SECRET=bar \
  REDIRECT_URL=https://example.com/install \
  go run cmd/server/main.go
```

	•	Smoke Test:
``` bash
curl -X POST http://localhost:8080/install \
  -H "Content-Type: application/json" \
  -d '{"code":"XYZ","installation_callback_url":"http://localhost:8080/install/callback"}'
curl http://localhost:8080/markup
```

# Contributing
	1.	Fork the repository and create a feature branch (git checkout -b feature/awesome).
	2.	Commit your changes (git commit -m "Add awesome feature").
	3.	Push to your fork (git push origin feature/awesome).
	4.	Open a Pull Request against main.

Please adhere to Go formatting (go fmt) and include tests where applicable.

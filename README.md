# House of Bounce Website (Go)

Business web portal for **House of Bounce**, a bounce house rental company in Maine.

## Features

- Professional, cheerful homepage with primary-color branding
- Informational section about services and service area
- Contact section with a message form
- Scheduling request section for event booking inquiries
- Go standard library server (`net/http`) with no external dependencies

## Run locally

1. Open a terminal in this folder.
2. Run:

```bash
go run ./cmd/web
```

3. Visit:

- http://localhost:8080

## Project structure

- `cmd/web/main.go` - HTTP server and routes
- `templates/index.html` - Portal layout and sections
- `static/css/styles.css` - Branding and responsive UI styles
- `static/js/main.js` - Frontend helpers for active nav and smooth scrolling

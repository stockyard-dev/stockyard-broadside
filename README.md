# Stockyard Broadside

**OG image and screenshot API — POST a template, get back a PNG for social cards**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9080:9080 -v broadside_data:/data ghcr.io/stockyard-dev/stockyard-broadside
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9080` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9080` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `BROADSIDE_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 3 templates, 100 renders/mo | Unlimited templates and renders |
| Price | Free | $4.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Developer Tools

## License

Apache 2.0

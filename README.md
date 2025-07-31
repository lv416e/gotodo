# GoTODO

A modern TODO application built with Go, featuring categories, priorities, and SQLite persistence.

## Features

- ✅ Create, edit, toggle, and delete TODOs
- 📁 Category management with color coding
- 🔴🟡⚪ Priority levels (High, Medium, Low)
- 🔍 Filter by category and priority
- 💾 SQLite database persistence
- 📱 Responsive design
- 🐳 Docker support

## Quick Start

### Using Go

```bash
go run main.go
```

### Using Docker Compose

```bash
# Production mode
docker-compose up -d

# Development mode with hot reload
docker-compose -f docker-compose.dev.yml up
```

Open http://localhost:8080

## Docker Usage

### Build and Run

```bash
# Build the image
docker-compose build

# Start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

### Data Persistence

The SQLite database is stored in the `./data` directory and persisted between container restarts.

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: ./data/todos.db)

## Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (optional)

### Local Development

```bash
# Install dependencies
go mod download

# Run the server
go run main.go
```

### Docker Development

The development setup includes hot reload using Air:

```bash
docker-compose -f docker-compose.dev.yml up
```

Any changes to Go files will automatically rebuild and restart the server.

## Project Structure

```
gotodo/
├── main.go              # Application entry point
├── handlers/            # HTTP request handlers
├── models/              # Data models and database operations
├── database/            # Database connection and migrations
├── static/              # CSS and JavaScript files
├── templates/           # HTML templates
├── data/                # SQLite database (auto-created)
├── docker-compose.yml   # Production Docker configuration
├── docker-compose.dev.yml # Development Docker configuration
├── Dockerfile           # Production Docker image
└── Dockerfile.dev       # Development Docker image
```

## License

MIT
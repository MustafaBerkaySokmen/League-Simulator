# League Simulator

A professional football league simulation and prediction platform built with Go, PostgreSQL, and modern web technologies. This project allows you to simulate artificial leagues, analyze real-world football data, and predict championship probabilities using advanced statistical and machine learning models.

---

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Simulation Models](#simulation-models)
- [API Endpoints](#api-endpoints)
- [Web Interface](#web-interface)
- [Database Schema](#database-schema)
- [Development](#development)
- [License](#license)

---

## Features
- **Artificial League Simulation:** Create and simulate custom leagues with configurable team strengths.
- **Real League Analysis:** Fetch live data from [football-data.org](https://www.football-data.org/) to analyze and simulate real competitions.
- **Multiple Prediction Models:** Includes Poisson, Elo, Bradley-Terry, Logistic Regression, Neural Networks, and more.
- **Interactive Web UI:** Modern Bootstrap 5 interface with charts, tables, and simulation controls.
- **RESTful API:** Easily integrate with other tools or automate simulations.
- **PostgreSQL Persistence:** All league data is stored in a robust relational database.

---

## Architecture
- **Backend:** Go (Golang), REST API, PostgreSQL
- **Frontend:** HTML, Bootstrap 5, Chart.js, Vanilla JS
- **Data Source:** [football-data.org](https://www.football-data.org/) for real league data
- **Containerization:** Docker Compose for easy setup

---

## Getting Started

### Prerequisites
- Go 1.20+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Quick Start (with Docker)
1. Clone the repository:
   ```sh
   git clone https://github.com/MustafaBerkaySokmen/League-Simulator.git
   cd League-Simulator
   ```
2. Start PostgreSQL:
   ```sh
   docker-compose up -d
   ```
3. Initialize the database:
   ```sh
   psql -h localhost -U insider -d insider_league -f sql/schema.sql
   # (Optional) Seed data:
   # psql -h localhost -U insider -d insider_league -f sql/seeds.sql
   ```
4. Build and run the server:
   ```sh
   cd cmd/league
   go build -o league.exe
   ./league.exe
   ```
5. Open your browser and go to [http://localhost:8080](http://localhost:8080)

### Manual Setup
- Edit the Postgres DSN in `cmd/league/main.go` if needed.
- Use `go run cmd/league/main.go` for development.

---

## Project Structure
```
cmd/league/         # Main server entrypoint
internal/
  interfaces/       # Interfaces for repository, generator, updater, etc.
  models/           # Data models (Team, Match, Table, etc.)
  predictor/        # Prediction models (Poisson, Elo, ML, etc.)
  repo/             # Postgres repository implementation
  service/          # Simulation and table update logic
sql/                # Database schema and seed data
web/                # Static web UI (HTML, JS, CSS)
docker-compose.yml  # Docker config for Postgres
```

---

## Simulation Models
- **Poisson:** Classic Poisson goal model
- **Elo:** Elo rating-based simulation
- **Bradley-Terry:** Pairwise comparison model
- **Logistic Regression:** Pre-trained logistic model
- **Bivariate Poisson:** Correlated goal model
- **Zero-Inflated Poisson:** Poisson with extra 0-0 draws
- **MLP Neural Net:** Simple neural network (demo)

---

## API Endpoints
- `/table` - Get current league table
- `/simulate/week` - Simulate the next week
- `/simulate/weeks?weeks=N` - Simulate N weeks
- `/simulate/all` - Simulate all remaining matches
- `/predict?model=poisson|elo|bt|logistic|mlp|bivariate|zip&sims=5000` - Predict championship probabilities
- `/reset?teams=N&type=random|homogeneous` - Reset league with N teams
- `/edit_match` - Edit a match result (POST JSON)
- `/real/leagues` - List real leagues (from football-data.org)
- `/real/standings?league=CODE` - Get real league standings
- `/real/fixtures?league=CODE` - Get real league fixtures
- `/real/simulate?league=CODE&week=N` - Simulate up to week N for a real league
- `/real/predict?league=CODE&model=...&sims=...&week=N` - Predict real league outcomes

---

## Web Interface
- **Artificial League Tab:** Create, simulate, and analyze custom leagues
- **Real League Tab:** Analyze and simulate real-world leagues
- **Charts:** Visualize championship probabilities
- **Standings Table:** Sortable, searchable, downloadable as CSV
- **Simulation Controls:** Simulate weeks, edit results, reset league

---

## Database Schema
- See `sql/schema.sql` for table definitions:
  - `teams` (id, name, strength)
  - `matches` (id, week, home_team_id, away_team_id, home_goals, away_goals)

---

## Development
- Use `go mod tidy` to manage dependencies
- All business logic is in `internal/`
- Add new models in `internal/predictor/` as needed
- Web UI is static and can be extended in `web/`

---

## License
MIT License. See [LICENSE](LICENSE) for details.

---

## Credits
- [football-data.org](https://www.football-data.org/) for real league data
- Bootstrap, Chart.js, Go, PostgreSQL

---

## Contact
For questions or contributions, open an issue or contact [Mustafa Berkay Sokmen](mailto:info@yourcompany.com).

## Badges

![Go Version](https://img.shields.io/badge/go-1.20%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

# Architecture Overview

This document describes the architecture and design of the League Simulator project.

---

## Overview
The League Simulator is a modular, interface-driven Go application for simulating and predicting football league outcomes. It uses PostgreSQL for persistence and provides a REST API and a modern web UI.

---

## Main Components

### 1. Backend (Go)
- **Entry Point:** `cmd/league/main.go`
- **Business Logic:** `internal/`
  - `interfaces/`: Defines interfaces for repository, match generator, table updater, etc.
  - `models/`: Data models for teams, matches, standings, etc.
  - `predictor/`: Implements various prediction models (Poisson, Elo, ML, etc.).
  - `repo/`: PostgreSQL repository implementation.
  - `service/`: Simulation and table update logic.

### 2. Database (PostgreSQL)
- **Schema:** Defined in `sql/schema.sql`
- **Tables:**
  - `teams`: Team info and strength
  - `matches`: Match results (week, teams, goals)

### 3. Web UI
- **Static files:** `web/index.html` (Bootstrap 5, Chart.js, JS)
- **Served by:** Go HTTP server
- **Features:** League table, simulation controls, probability charts, real league data

---

## Data Flow
1. **User Action:** User interacts with the web UI or sends API requests.
2. **API Handler:** Go HTTP handler processes the request.
3. **Business Logic:** Calls into `internal/` modules for simulation, prediction, or data retrieval.
4. **Database:** Reads/writes to PostgreSQL via repository interface.
5. **Response:** Returns JSON data or serves updated web UI.

---

## Extensibility
- **Add new prediction models:** Implement the `Predictor` interface in `internal/predictor/`.
- **Change database:** Implement a new repository in `internal/repo/`.
- **Add endpoints:** Extend handlers in `cmd/league/main.go`.

---

## Diagram

```
+-------------------+
|   Web Browser     |
+-------------------+
          |
          v
+-------------------+
|   Go HTTP Server  |
| (cmd/league)      |
+-------------------+
          |
          v
+-------------------+
|   Business Logic  |
|  (internal/*)     |
+-------------------+
          |
          v
+-------------------+
|   PostgreSQL DB   |
+-------------------+
```

---

## Technologies Used
- Go 1.20+
- PostgreSQL
- Docker, Docker Compose
- Bootstrap 5, Chart.js

---

## References
- [football-data.org](https://www.football-data.org/)
- [Go Documentation](https://go.dev/doc/)
- [Bootstrap](https://getbootstrap.com/)
- [Chart.js](https://www.chartjs.org/)

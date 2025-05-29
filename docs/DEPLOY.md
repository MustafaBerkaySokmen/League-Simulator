# Deployment Guide

This document explains how to set up and deploy the League Simulator project on your local machine or a server.

---

## Prerequisites
- **Go 1.20+** (https://go.dev/dl/)
- **Docker & Docker Compose** (https://docs.docker.com/get-docker/)
- **PostgreSQL** (optional, if not using Docker)

---

## Quick Start (Recommended: Docker Compose)

1. **Clone the repository:**
   ```sh
   git clone https://github.com/MustafaBerkaySokmen/League-Simulator.git
   cd League-Simulator
   ```

2. **Start PostgreSQL with Docker Compose:**
   ```sh
   docker-compose up -d
   ```
   This will start a PostgreSQL instance with the correct user, password, and database.

3. **Initialize the database schema:**
   ```sh
   psql -h localhost -U insider -d insider_league -f sql/schema.sql
   # (Optional) Seed data:
   # psql -h localhost -U insider -d insider_league -f sql/seeds.sql
   ```
   Default credentials (see `docker-compose.yml`):
   - User: insider
   - Password: insider
   - Database: insider_league

4. **Build and run the Go server:**
   ```sh
   cd cmd/league
   go build -o league.exe
   ./league.exe
   ```
   The server will listen on `http://localhost:8080` by default.

5. **Access the web interface:**
   Open your browser and go to [http://localhost:8080](http://localhost:8080)

---

## Manual Setup (Without Docker)

1. **Install PostgreSQL** and create a database/user matching the credentials in `cmd/league/main.go`.
2. **Run the schema and (optionally) seed scripts** in `sql/`.
3. **Build and run the Go server** as above.

---

## Environment Variables
- You can change the Postgres DSN in `cmd/league/main.go` if needed.
- For production, consider using environment variables for DB credentials.

---

## Deployment to a Server
- Use a Linux VM or cloud service (AWS, GCP, Azure, DigitalOcean, etc.).
- Install Docker and follow the same steps as above.
- Use a process manager (e.g., systemd, supervisor) to keep the Go server running.
- Configure firewall to allow access to port 8080.

---

## Updating the Application
- Pull the latest changes:
   ```sh
   git pull origin main
   ```
- Rebuild and restart the server as above.

---

## Troubleshooting
- Check logs for errors.
- Ensure PostgreSQL is running and accessible.
- Make sure ports 5432 (Postgres) and 8080 (app) are open.

---

## Contact
For help, open an issue on GitHub or contact the maintainer.

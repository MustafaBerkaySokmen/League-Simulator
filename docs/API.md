# API Documentation

This document describes the REST API endpoints for the League Simulator project. All endpoints are accessible via HTTP and return JSON responses.

---

## Table of Contents
- [League Table](#league-table)
- [Simulate Matches](#simulate-matches)
- [Predictions](#predictions)
- [Reset League](#reset-league)
- [Edit Match Result](#edit-match-result)
- [Real League Data](#real-league-data)

---

## League Table

### Get Current League Table
- **Endpoint:** `/table`
- **Method:** GET
- **Description:** Returns the current league table.
- **Response Example:**
```json
[
  {
    "team_id": 1,
    "team_name": "Team 1",
    "played": 4,
    "won": 2,
    "drawn": 1,
    "lost": 1,
    "goals_for": 7,
    "goals_against": 5,
    "goal_diff": 2,
    "points": 7
  },
  ...
]
```

---

## Simulate Matches

### Simulate Next Week
- **Endpoint:** `/simulate/week`
- **Method:** POST
- **Description:** Simulates all matches for the next unplayed week.
- **Response:** Updated league table (see above).

### Simulate N Weeks
- **Endpoint:** `/simulate/weeks?weeks=N`
- **Method:** POST
- **Description:** Simulates the next N weeks.
- **Response Example:**
```json
{
  "simulated": [
    {
      "week": 5,
      "fixtures": [
        { "week": 5, "home_team_id": 1, "away_team_id": 2, "home_goals": 2, "away_goals": 1 },
        ...
      ]
    },
    ...
  ],
  "table": [ ... ]
}
```

### Simulate All Remaining Matches
- **Endpoint:** `/simulate/all`
- **Method:** POST
- **Description:** Simulates all remaining matches in the league.
- **Response:** Same as above, with all weeks simulated.

---

## Predictions

### Predict Championship Probabilities
- **Endpoint:** `/predict?model=poisson|elo|bt|logistic|mlp|bivariate|zip&sims=5000`
- **Method:** GET
- **Description:** Runs Monte Carlo simulations to estimate each team's probability of winning the league.
- **Query Parameters:**
  - `model`: Prediction model to use (default: poisson)
  - `sims`: Number of simulations (default: 5000)
- **Response Example:**
```json
{
  "model": "Poisson",
  "sims": 5000,
  "probs": {
    "1": 0.32,
    "2": 0.18,
    "3": 0.25,
    "4": 0.25
  }
}
```

---

## Reset League

### Reset League and Teams
- **Endpoint:** `/reset?teams=N&type=random|homogeneous`
- **Method:** POST
- **Description:** Resets the league with N teams. `type` can be `random` (random strengths) or `homogeneous` (all teams equal strength).
- **Response:**
```json
{
  "teams": [ ... ],
  "table": [ ... ]
}
```

---

## Edit Match Result

### Edit a Match Result
- **Endpoint:** `/edit_match`
- **Method:** POST
- **Content-Type:** application/json
- **Body Example:**
```json
{
  "week": 3,
  "home_team_id": 1,
  "away_team_id": 2,
  "home_goals": 2,
  "away_goals": 2
}
```
- **Response:**
```json
{ "status": "ok" }
```

---

## Real League Data

### List Real Leagues
- **Endpoint:** `/real/leagues`
- **Method:** GET
- **Description:** Returns a list of real football leagues from football-data.org.

### Get Real League Standings
- **Endpoint:** `/real/standings?league=CODE`
- **Method:** GET
- **Description:** Returns the current standings for a real league.

### Get Real League Fixtures
- **Endpoint:** `/real/fixtures?league=CODE`
- **Method:** GET
- **Description:** Returns upcoming fixtures for a real league.

### Simulate Real League Weeks
- **Endpoint:** `/real/simulate?league=CODE&week=N`
- **Method:** GET
- **Description:** Simulates up to week N for a real league, returning real results and standings.

### Predict Real League Outcomes
- **Endpoint:** `/real/predict?league=CODE&model=poisson|elo|bt|logistic|mlp|bivariate|zip&sims=5000&week=N`
- **Method:** GET
- **Description:** Predicts championship probabilities for a real league using the selected model.

---

## Error Handling
- All errors are returned as JSON with an appropriate HTTP status code and message.

---

## Authentication
- No authentication is required for local use. If deploying publicly, consider adding authentication.

---

## Example Usage (curl)
```sh
# Get league table
curl http://localhost:8080/table

# Simulate next week
curl -X POST http://localhost:8080/simulate/week

# Predict with Elo model
curl "http://localhost:8080/predict?model=elo&sims=10000"

# Reset league with 4 teams, random strengths
curl -X POST "http://localhost:8080/reset?teams=4&type=random"
```

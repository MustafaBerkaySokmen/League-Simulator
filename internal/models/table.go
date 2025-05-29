package models

// TeamStanding tracks a single team's cumulative stats.
type TeamStanding struct {
    TeamID       int    `json:"team_id"`
    TeamName     string `json:"team_name"`
    Played       int    `json:"played"`
    Won          int    `json:"won"`
    Drawn        int    `json:"drawn"`
    Lost         int    `json:"lost"`
    GoalsFor     int    `json:"goals_for"`
    GoalsAgainst int    `json:"goals_against"`
    GoalDiff     int    `json:"goal_diff"`
    Points       int    `json:"points"`
}

// LeagueTable is the full standings for all teams.
type LeagueTable []TeamStanding

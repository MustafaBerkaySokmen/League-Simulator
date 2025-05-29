package interfaces

import "github.com/musta/insider-league/internal/models"

// Repository abstracts data persistence.
type Repository interface {
    ListTeams() ([]models.Team, error)
    SaveMatch(m models.MatchResult) error
    GetTable() (models.LeagueTable, error)
    GetAllMatches() ([]models.MatchResult, error)
    ResetMatches() error
    ResetTeams() error
    SaveTeam(team models.Team) error


    // ListRemainingMatches returns all future fixtures.
    ListRemainingMatches() ([]Matchup, error)
}

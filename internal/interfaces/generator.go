package interfaces

import "github.com/musta/insider-league/internal/models"

// MatchGenerator defines how to simulate a fixture result.
type MatchGenerator interface {
    // Generate takes two teams and returns the MatchResult.
    Generate(home, away models.Team) (models.MatchResult, error)
}

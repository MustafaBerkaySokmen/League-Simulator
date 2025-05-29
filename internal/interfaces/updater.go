package interfaces

import "github.com/musta/insider-league/internal/models"

// TableUpdater applies a MatchResult to a LeagueTable.
type TableUpdater interface {
    // Update takes the current table and a match result,
    // and returns the new updated table.
    Update(current models.LeagueTable, result models.MatchResult) (models.LeagueTable, error)
}

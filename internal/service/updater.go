package service

import (
    "fmt"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

// TableUpdater implements interfaces.TableUpdater
type TableUpdater struct{}

// NewTableUpdater returns a ready-to-use TableUpdater.
func NewTableUpdater() interfaces.TableUpdater {
    return &TableUpdater{}
}

// Update applies a single MatchResult to the current LeagueTable.
func (u *TableUpdater) Update(current models.LeagueTable, result models.MatchResult) (models.LeagueTable, error) {
    // We expect each team to already have a standing entry in current.
    // Find indices for home and away.
    var homeIdx, awayIdx = -1, -1
    for i, row := range current {
        if row.TeamID == result.HomeTeamID {
            homeIdx = i
        }
        if row.TeamID == result.AwayTeamID {
            awayIdx = i
        }
    }
    if homeIdx < 0 || awayIdx < 0 {
        return nil, fmt.Errorf("team not found in table: homeIdx=%d awayIdx=%d", homeIdx, awayIdx)
    }

    // Helper to update a single side
    apply := func(idx int, goalsFor, goalsAgainst int, winPts, drawPts int) {
        row := &current[idx]
        row.Played++
        row.GoalsFor += goalsFor
        row.GoalsAgainst += goalsAgainst
        row.GoalDiff = row.GoalsFor - row.GoalsAgainst
        if goalsFor > goalsAgainst {
            row.Won++
            row.Points += winPts
        } else if goalsFor == goalsAgainst {
            row.Drawn++
            row.Points += drawPts
        } else {
            row.Lost++
        }
    }

    // Apply home
    apply(homeIdx, result.HomeGoals, result.AwayGoals, 3, 1)
    // Apply away
    apply(awayIdx, result.AwayGoals, result.HomeGoals, 3, 1)

    return current, nil
}

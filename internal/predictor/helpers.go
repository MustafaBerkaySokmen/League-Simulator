package predictor

import (
    "math"
    "math/rand"

    "github.com/musta/insider-league/internal/models"
)


// findTableRow returns the TeamStanding and its index for a given teamID.
func findTableRow(table models.LeagueTable, teamID int) (models.TeamStanding, int) {
    for i, row := range table {
        if row.TeamID == teamID {
            return row, i
        }
    }
    return models.TeamStanding{}, -1
}

// applyResult applies a single fixture result to a TeamStanding.
func applyResult(ps *models.TeamStanding, goalsFor, goalsAgainst int) {
    ps.Played++
    ps.GoalsFor += goalsFor
    ps.GoalsAgainst += goalsAgainst
    ps.GoalDiff = ps.GoalsFor - ps.GoalsAgainst
    if goalsFor > goalsAgainst {
        ps.Won++
        ps.Points += 3
    } else if goalsFor == goalsAgainst {
        ps.Drawn++
        ps.Points++
    } else {
        ps.Lost++
    }
}

// applyResultByID applies a result to both home and away entries in the table.
func applyResultByID(table *models.LeagueTable, homeID, awayID, homeGoals, awayGoals int) {
    for i := range *table {
        if (*table)[i].TeamID == homeID {
            applyResult(&(*table)[i], homeGoals, awayGoals)
        }
        if (*table)[i].TeamID == awayID {
            applyResult(&(*table)[i], awayGoals, homeGoals)
        }
    }
}

// Exported version for use in main.go
func ApplyResultByID(table *models.LeagueTable, homeID, awayID, homeGoals, awayGoals int) {
	applyResultByID(table, homeID, awayID, homeGoals, awayGoals)
}

// findChampion returns the teamID with highest points, breaking ties by goal diff.
func findChampion(table models.LeagueTable) int {
    champ := table[0]
    champID := champ.TeamID
    for _, row := range table[1:] {
        if row.Points > champ.Points || (row.Points == champ.Points && row.GoalDiff > champ.GoalDiff) {
            champ = row
            champID = row.TeamID
        }
    }
    return champID
}

// Exported version for use in main.go
func FindChampion(table models.LeagueTable) int {
	return findChampion(table)
}

// samplePoisson generates a Poisson-distributed int with mean lambda.
func samplePoisson(lambda float64, rng *rand.Rand) int {
    L := math.Exp(-lambda)
    k := 0
    p := 1.0
    for p > L {
        k++
        p *= rng.Float64()
    }
    return k - 1
}

// computeLambda scales a team's strength into an expected-goals lambda.
func computeLambda(strength, totalStrength, maxGoals float64) float64 {
    return strength / totalStrength * maxGoals
}

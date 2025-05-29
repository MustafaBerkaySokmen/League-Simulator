package predictor

import (
    "math"
    "math/rand"
    "time"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

const (
    KFactor    = 20.0
    BaseRating = 1500.0
)

// EloMC uses Elo rating simulations.
type EloMC struct{}

// NewEloMonteCarlo creates an Elo predictor.
func NewEloMonteCarlo() Predictor {
    return &EloMC{}
}

func (e *EloMC) Name() string { return "Elo Monte Carlo" }

func (e *EloMC) Predict(
    teams []models.Team,
    table models.LeagueTable,
    remaining []interfaces.Matchup,
    sims int,
) (map[int]float64, error) {
    wins := make(map[int]int)

    for i := 0; i < sims; i++ {
        // initialize ratings
        ratings := make(map[int]float64)
        for _, t := range teams {
            ratings[t.ID] = BaseRating
        }
        simTable := make(models.LeagueTable, len(table))
        copy(simTable, table)
        rng := rand.New(rand.NewSource(time.Now().UnixNano()))

        for _, m := range remaining {
            Ra := ratings[m.HomeTeamID]
            Rb := ratings[m.AwayTeamID]
            Ea := 1.0 / (1.0 + math.Pow(10, (Rb-Ra)/400))
            Eb := 1.0 - Ea

            // simulate win/loss (no draws here for simplicity)
            var hg, ag int
            if rng.Float64() < Ea {
                hg, ag = 1, 0
            } else {
                hg, ag = 0, 1
            }

            // update ratings
            ratings[m.HomeTeamID] += KFactor * (float64(hg) - Ea)
            ratings[m.AwayTeamID] += KFactor * (float64(ag) - Eb)

            applyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
        }

        champ := findChampion(simTable)
        wins[champ]++
    }

    probs := make(map[int]float64)
    for id, cnt := range wins {
        probs[id] = float64(cnt) / float64(sims)
    }
    return probs, nil
}

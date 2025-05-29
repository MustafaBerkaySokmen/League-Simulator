package predictor

import (
    "math/rand"
    "time"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

// PoissonMC uses Poisson sampling for the rest of the season.
type PoissonMC struct{}

// NewPoissonMonteCarlo creates a Poisson Monte Carlo predictor.
func NewPoissonMonteCarlo() Predictor {
    return &PoissonMC{}
}

func (p *PoissonMC) Name() string { return "Poisson Monte Carlo" }

func (p *PoissonMC) Predict(
    teams []models.Team,
    table models.LeagueTable,
    remaining []interfaces.Matchup,
    sims int,
) (map[int]float64, error) {
    wins := make(map[int]int)
    totalStr := 0.0
    for _, t := range teams {
        totalStr += float64(t.Strength)
    }

    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < sims; i++ {
        simTable := make(models.LeagueTable, len(table))
        copy(simTable, table)

        for _, m := range remaining {
            // look up teams
            var home, away models.Team
            for _, t := range teams {
                if t.ID == m.HomeTeamID {
                    home = t
                }
                if t.ID == m.AwayTeamID {
                    away = t
                }
            }
            // compute lambdas
            hl := computeLambda(float64(home.Strength), totalStr, 3.0)
            al := computeLambda(float64(away.Strength), totalStr, 2.5)
            // sample goals
            hg := samplePoisson(hl, rng)
            ag := samplePoisson(al, rng)
            // apply to table
            applyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
        }

        champ := findChampion(simTable)
        wins[champ]++
    }

    probs := make(map[int]float64, len(wins))
    for id, cnt := range wins {
        probs[id] = float64(cnt) / float64(sims)
    }
    return probs, nil
}

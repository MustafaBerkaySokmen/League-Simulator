package predictor

import (
    "math/rand"
    "time"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

// btMC uses a Bradley–Terry simulation.
type btMC struct{}

// NewBradleyTerryMonteCarlo creates a BT predictor.
func NewBradleyTerryMonteCarlo() Predictor {
    return &btMC{}
}

func (b *btMC) Name() string {
    return "Bradley-Terry Monte Carlo"
}

func (b *btMC) Predict(
    teams []models.Team,
    table models.LeagueTable,
    remaining []interfaces.Matchup,
    sims int,
) (map[int]float64, error) {
    wins := make(map[int]int)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    // use Strength as BT parameter
    strengths := make(map[int]float64)
    for _, t := range teams {
        strengths[t.ID] = float64(t.Strength)
    }

    for i := 0; i < sims; i++ {
        simTable := make(models.LeagueTable, len(table))
        copy(simTable, table)

        for _, m := range remaining {
            sa := strengths[m.HomeTeamID]
            sb := strengths[m.AwayTeamID]
            pa := sa / (sa + sb)

            if rng.Float64() < pa {
                applyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, 1, 0)
            } else {
                applyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, 0, 1)
            }
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

package predictor

import (
    "math"
    "math/rand"
    "time"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

const (
    b0 = 0.1
    b1 = 0.05
    b2 = -0.03
)

// LogisticMC uses a pre-trained logistic model.
type LogisticMC struct{}

func NewLogisticMonteCarlo() Predictor { return &LogisticMC{} }
func (l *LogisticMC) Name() string    { return "Logistic Regression MC" }

func (l *LogisticMC) Predict(
    teams []models.Team,
    table models.LeagueTable,
    remaining []interfaces.Matchup,
    sims int,
) (map[int]float64, error) {
    wins := make(map[int]int)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    for i := 0; i < sims; i++ {
        simTable := make(models.LeagueTable, len(table))
        copy(simTable, table)

        for _, m := range remaining {
            homeRow, _ := findTableRow(simTable, m.HomeTeamID)
            awayRow, _ := findTableRow(simTable, m.AwayTeamID)

            λh := float64(homeRow.GoalsFor) / float64(homeRow.Played+1)
            λa := float64(awayRow.GoalsFor) / float64(awayRow.Played+1)

            x := b0 + b1*λh + b2*λa
            ph := 1.0 / (1.0 + math.Exp(-x))

            if rng.Float64() < ph {
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

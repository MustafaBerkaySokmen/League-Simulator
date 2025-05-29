package predictor

import (
    "errors"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

// MonteCarlo is a stub for the Monte Carlo model.
type MonteCarlo struct{}

// NewMonteCarlo creates a Monte Carlo predictor.
func NewMonteCarlo() Predictor {
    return &MonteCarlo{}
}

func (m *MonteCarlo) Name() string {
    return "Monte Carlo"
}

func (m *MonteCarlo) Predict(
    teams []models.Team,
    table models.LeagueTable,
    remaining []interfaces.Matchup,
    sims int,
) (map[int]float64, error) {
    return nil, errors.New("Monte Carlo model not implemented")
}

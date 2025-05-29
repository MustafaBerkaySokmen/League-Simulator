package service

import (
    "math"
    "math/rand"
    "time"

    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
)

// PoissonGenerator implements interfaces.MatchGenerator
// using a Poisson distribution for goal scoring.
type PoissonGenerator struct {
    rnd *rand.Rand
}

// NewPoissonGenerator returns a ready-to-use MatchGenerator.
func NewPoissonGenerator() interfaces.MatchGenerator {
    return &PoissonGenerator{
        rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// Generate simulates a match between home and away teams.
func (g *PoissonGenerator) Generate(home, away models.Team) (models.MatchResult, error) {
    // Calculate expected goals (λ) proportional to strengths
    total := float64(home.Strength + away.Strength)
    homeLambda := float64(home.Strength) / total * 3.0
    awayLambda := float64(away.Strength) / total * 2.5

    homeGoals := samplePoisson(homeLambda, g.rnd)
    awayGoals := samplePoisson(awayLambda, g.rnd)

    return models.MatchResult{
        Week:       0, // set this when you schedule the match
        HomeTeamID: home.ID,
        AwayTeamID: away.ID,
        HomeGoals:  homeGoals,
        AwayGoals:  awayGoals,
    }, nil
}

// samplePoisson draws a random integer from a Poisson distribution.
func samplePoisson(lambda float64, rnd *rand.Rand) int {
    L := math.Exp(-lambda)
    k, p := 0, 1.0
    for p > L {
        k++
        p *= rnd.Float64()
    }
    return k - 1
}

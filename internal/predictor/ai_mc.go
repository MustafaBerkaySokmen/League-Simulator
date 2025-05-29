package predictor

import (
	"math"
	"math/rand"
	"time"

	"github.com/musta/insider-league/internal/interfaces"
	"github.com/musta/insider-league/internal/models"
)

// AIMC uses a simple neural network for match outcome prediction.
type AIMC struct{}

func NewAIMonteCarlo() Predictor {
	return &AIMC{}
}

// MLP Neural Net Monte Carlo (renamed from AI)
func NewMLPNeuralNetMonteCarlo() Predictor {
	return &AIMC{}
}

func (a *AIMC) Name() string { return "AI Monte Carlo" }

// For demo: a tiny neural net with hardcoded weights (normally you'd train this!)
func predictOutcome(homeStrength, awayStrength float64) (ph, pd, pa float64) {
	// Simple 1-hidden-layer MLP with made-up weights
	// Inputs: [homeStrength, awayStrength] scaled to [0,1]
	hs := homeStrength / 100.0
	as := awayStrength / 100.0

	// Hidden layer (2 neurons)
	h1 := 0.8*hs - 0.5*as + 0.1
	h2 := -0.3*hs + 0.9*as - 0.2
	// ReLU
	if h1 < 0 {
		h1 = 0
	}
	if h2 < 0 {
		h2 = 0
	}

	// Output layer (3 neurons: home win, draw, away win)
	outH := 1.2*h1 - 0.7*h2 + 0.2
	outD := 0.5*h1 + 0.5*h2 + 0.1
	outA := -0.6*h1 + 1.3*h2 + 0.2

	// Softmax
	expH := exp(outH)
	expD := exp(outD)
	expA := exp(outA)
	sum := expH + expD + expA
	return expH / sum, expD / sum, expA / sum
}

func exp(x float64) float64 {
	return math.Exp(x)
}

// Exported version for use in main.go
func PredictOutcome(homeStrength, awayStrength float64) (ph, pd, pa float64) {
	return predictOutcome(homeStrength, awayStrength)
}

func (a *AIMC) Predict(
	teams []models.Team,
	table models.LeagueTable,
	remaining []interfaces.Matchup,
	sims int,
) (map[int]float64, error) {
	wins := make(map[int]int)
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < sims; i++ {
		simTable := make(models.LeagueTable, len(table))
		copy(simTable, table)

		for _, m := range remaining {
			home := teamMap[m.HomeTeamID]
			away := teamMap[m.AwayTeamID]
			ph, pd, _ := predictOutcome(float64(home.Strength), float64(away.Strength))
			r := rng.Float64()
			var hg, ag int
			if r < ph {
				hg, ag = 1, 0 // home win
			} else if r < ph+pd {
				hg, ag = 1, 1 // draw
			} else {
				hg, ag = 0, 1 // away win
			}
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

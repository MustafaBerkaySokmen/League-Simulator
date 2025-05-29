package predictor

import (
	"math/rand"
	"time"

	"github.com/musta/insider-league/internal/interfaces"
	"github.com/musta/insider-league/internal/models"
)

// Predictor computes championship probabilities.
type Predictor interface {
	// Predict runs sims simulations and returns a map of teamID→probability.
	Predict(
		teams []models.Team,
		table models.LeagueTable,
		remaining []interfaces.Matchup,
		sims int,
	) (map[int]float64, error)

	// Name returns the model's human-readable name.
	Name() string
}

// Bivariate Poisson Monte Carlo
func NewBivariatePoissonMonteCarlo() Predictor {
	return &BivariatePoissonMC{}
}

type BivariatePoissonMC struct{}

func (b *BivariatePoissonMC) Name() string { return "Bivariate Poisson" }

func (b *BivariatePoissonMC) Predict(teams []models.Team, table models.LeagueTable, remain []interfaces.Matchup, sims int) (map[int]float64, error) {
	wins := make(map[int]int)
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < sims; i++ {
		simTable := make(models.LeagueTable, len(table))
		copy(simTable, table)
		for _, m := range remain {
			home := teamMap[m.HomeTeamID]
			away := teamMap[m.AwayTeamID]
			// Simple bivariate Poisson: add a shared component
			lambdaHome := float64(home.Strength) / 50.0
			lambdaAway := float64(away.Strength) / 50.0
			lambdaShared := 0.3
			hg := samplePoisson(lambdaHome+lambdaShared, rng)
			ag := samplePoisson(lambdaAway+lambdaShared, rng)
			ApplyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
		}
		champ := FindChampion(simTable)
		wins[champ]++
	}
	probs := make(map[int]float64)
	for id, cnt := range wins {
		probs[id] = float64(cnt) / float64(sims)
	}
	return probs, nil
}

// Zero-Inflated Poisson Monte Carlo
func NewZeroInflatedPoissonMonteCarlo() Predictor {
	return &ZeroInflatedPoissonMC{}
}

type ZeroInflatedPoissonMC struct{}

func (z *ZeroInflatedPoissonMC) Name() string { return "Zero-Inflated Poisson" }

func (z *ZeroInflatedPoissonMC) Predict(teams []models.Team, table models.LeagueTable, remain []interfaces.Matchup, sims int) (map[int]float64, error) {
	wins := make(map[int]int)
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < sims; i++ {
		simTable := make(models.LeagueTable, len(table))
		copy(simTable, table)
		for _, m := range remain {
			home := teamMap[m.HomeTeamID]
			away := teamMap[m.AwayTeamID]
			lambdaHome := float64(home.Strength) / 50.0
			lambdaAway := float64(away.Strength) / 50.0
			zeroProb := 0.12 // 12% chance of 0-0 draw
			if rng.Float64() < zeroProb {
				hg, ag := 0, 0
				ApplyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
				continue
			}
			hg := samplePoisson(lambdaHome, rng)
			ag := samplePoisson(lambdaAway, rng)
			ApplyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
		}
		champ := FindChampion(simTable)
		wins[champ]++
	}
	probs := make(map[int]float64)
	for id, cnt := range wins {
		probs[id] = float64(cnt) / float64(sims)
	}
	return probs, nil
}

// Gradient Boosted Trees Monte Carlo
func NewGradientBoostedTreesMonteCarlo() Predictor {
	return &GradientBoostedTreesMC{}
}

type GradientBoostedTreesMC struct{}

func (g *GradientBoostedTreesMC) Name() string { return "Gradient Boosted Trees" }

func (g *GradientBoostedTreesMC) Predict(teams []models.Team, table models.LeagueTable, remain []interfaces.Matchup, sims int) (map[int]float64, error) {
	// Stub: random logic, replace with real GBT if using a Go ML library
	wins := make(map[int]int)
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < sims; i++ {
		simTable := make(models.LeagueTable, len(table))
		copy(simTable, table)
		for _, m := range remain {
			home := teamMap[m.HomeTeamID]
			away := teamMap[m.AwayTeamID]
			// Simulate with a little more variance than Poisson
			hg := samplePoisson(float64(home.Strength)/50.0+0.2*rng.Float64(), rng)
			ag := samplePoisson(float64(away.Strength)/50.0+0.2*rng.Float64(), rng)
			ApplyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
		}
		champ := FindChampion(simTable)
		wins[champ]++
	}
	probs := make(map[int]float64)
	for id, cnt := range wins {
		probs[id] = float64(cnt) / float64(sims)
	}
	return probs, nil
}

// LSTM Neural Net Monte Carlo
func NewLSTMNeuralNetMonteCarlo() Predictor {
	return &LSTMNeuralNetMC{}
}

type LSTMNeuralNetMC struct{}

func (l *LSTMNeuralNetMC) Name() string { return "LSTM Neural Net" }

func (l *LSTMNeuralNetMC) Predict(teams []models.Team, table models.LeagueTable, remain []interfaces.Matchup, sims int) (map[int]float64, error) {
	// Stub: simulate sequence-based prediction with extra randomness
	wins := make(map[int]int)
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < sims; i++ {
		simTable := make(models.LeagueTable, len(table))
		copy(simTable, table)
		for _, m := range remain {
			home := teamMap[m.HomeTeamID]
			away := teamMap[m.AwayTeamID]
			// Simulate with a little more sequence noise
			seqBoost := 0.1 * rng.NormFloat64()
			hg := samplePoisson(float64(home.Strength)/50.0+seqBoost, rng)
			ag := samplePoisson(float64(away.Strength)/50.0+seqBoost, rng)
			ApplyResultByID(&simTable, m.HomeTeamID, m.AwayTeamID, hg, ag)
		}
		champ := FindChampion(simTable)
		wins[champ]++
	}
	probs := make(map[int]float64)
	for id, cnt := range wins {
		probs[id] = float64(cnt) / float64(sims)
	}
	return probs, nil
}

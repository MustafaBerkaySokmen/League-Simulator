package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/musta/insider-league/internal/interfaces"
	"github.com/musta/insider-league/internal/models"
	"github.com/musta/insider-league/internal/predictor"
	"github.com/musta/insider-league/internal/repo"
	"github.com/musta/insider-league/internal/service"
)

/* -------------------------------------------------------------------------- */
/*                                App helpers                                 */
/* -------------------------------------------------------------------------- */

type App struct {
	repo interfaces.Repository
	gen  interfaces.MatchGenerator
}

// nextFixtures returns *all* fixtures that belong to the next un-played week.
func (a *App) nextFixtures() ([]interfaces.Matchup, error) {
	remain, err := a.repo.ListRemainingMatches()
	if err != nil || len(remain) == 0 {
		return nil, err
	}
	nextWeek := remain[0].Week
	var out []interfaces.Matchup
	for _, f := range remain {
		if f.Week != nextWeek {
			break
		}
		out = append(out, f)
	}
	return out, nil
}

/* -------------------------------------------------------------------------- */
/*                                  Handlers                                  */
/* -------------------------------------------------------------------------- */

func (a *App) handleGetTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	table, err := a.repo.GetTable()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(table)
}

/* ------------ 1 week ------------------------------------------------------ */

func (a *App) handleSimulateWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	fixtures, err := a.nextFixtures()
	if err != nil {
		http.Error(w, "no remaining fixtures", http.StatusInternalServerError)
		return
	}

	teams, _ := a.repo.ListTeams()
	for _, f := range fixtures {
		home := findTeam(teams, f.HomeTeamID)
		away := findTeam(teams, f.AwayTeamID)
		match, _ := a.gen.Generate(home, away)
		match.Week = f.Week
		_ = a.repo.SaveMatch(match)
	}

	table, _ := a.repo.GetTable()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(table)
}

/* ------------ N weeks ----------------------------------------------------- */

func (a *App) handleSimulateWeeks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	weeks := 1
	if v := r.URL.Query().Get("weeks"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			weeks = n
		}
	}

	type WeekResult struct {
		Week     int                  `json:"week"`
		Fixtures []models.MatchResult `json:"fixtures"`
	}
	var simulated []WeekResult

	for i := 0; i < weeks; i++ {
		fixtures, _ := a.nextFixtures()
		if len(fixtures) == 0 {
			break // season finished
		}

		teams, _ := a.repo.ListTeams()
		var played []models.MatchResult
		for _, f := range fixtures {
			home := findTeam(teams, f.HomeTeamID)
			away := findTeam(teams, f.AwayTeamID)
			m, _ := a.gen.Generate(home, away)
			m.Week = f.Week
			_ = a.repo.SaveMatch(m)
			played = append(played, m)
		}
		simulated = append(simulated, WeekResult{Week: fixtures[0].Week, Fixtures: played})
	}

	table, _ := a.repo.GetTable()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"simulated": simulated,
		"table":     table,
	})
}

/* ------------ Predictions ------------------------------------------------- */

func (a *App) handlePredict(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	teams, _ := a.repo.ListTeams()
	table, _ := a.repo.GetTable()
	remain, _ := a.repo.ListRemainingMatches()

	// query params
	modelName := r.URL.Query().Get("model")
	sims := 5000
	if q := r.URL.Query().Get("sims"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			sims = n
		}
	}

	var model predictor.Predictor
	switch modelName {
	case "elo":
		model = predictor.NewEloMonteCarlo()
	case "bt":
		model = predictor.NewBradleyTerryMonteCarlo()
	case "logistic":
		model = predictor.NewLogisticMonteCarlo()
	case "mlp":
		model = predictor.NewMLPNeuralNetMonteCarlo()
	case "bivariate":
		model = predictor.NewBivariatePoissonMonteCarlo()
	case "zip":
		model = predictor.NewZeroInflatedPoissonMonteCarlo()
	default:
		model = predictor.NewPoissonMonteCarlo()
	}

	probs, err := model.Predict(teams, table, remain, sims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"model": model.Name(),
		"sims":  sims,
		"probs": probs,
	})
}

/* ------------ Hard-reset -------------------------------------------------- */

func (a *App) handleReset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[RESET] Resetting league...")
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	n, _ := strconv.Atoi(r.URL.Query().Get("teams"))
	fmt.Printf("[RESET] Requested team count: %d\n", n)
	if n < 2 || n%2 != 0 {
		http.Error(w, "teams must be an even integer ≥2", http.StatusBadRequest)
		return
	}

	typeParam := r.URL.Query().Get("type")
	initType := "random"
	if typeParam == "homogeneous" {
		initType = "homogeneous"
	}
	fmt.Printf("[RESET] Initialization type: %s\n", initType)

	if err := a.repo.ResetMatches(); err != nil {
		fmt.Printf("[RESET] Error truncating matches: %v\n", err)
	} else {
		fmt.Println("[RESET] Matches truncated.")
	}
	if err := a.repo.ResetTeams(); err != nil {
		fmt.Printf("[RESET] Error truncating teams: %v\n", err)
	} else {
		fmt.Println("[RESET] Teams truncated.")
	}

	for i := 1; i <= n; i++ {
		team := models.Team{
			Name: fmt.Sprintf("Team %d", i),
		}
		if initType == "homogeneous" {
			team.Strength = 75 // All teams same strength
		} else {
			team.Strength = 50 + rand.Intn(51) // 50-100
		}
		fmt.Printf("[RESET] Inserting team: %s, strength: %d\n", team.Name, team.Strength)
		_ = a.repo.SaveTeam(team)
	}

	teams, _ := a.repo.ListTeams()
	fmt.Printf("[RESET] Number of teams after insert: %d\n", len(teams))
	for _, t := range teams {
		fmt.Printf("[RESET] Team: id=%d, name=%s, strength=%d\n", t.ID, t.Name, t.Strength)
	}
	table, _ := a.repo.GetTable()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"teams": teams,
		"table": table,
	})
}

/* ------------ Real Leagues ----------------------------------------------- */

// Handler to fetch real leagues from football-data.org
func (a *App) handleRealLeagues(w http.ResponseWriter, r *http.Request) {
	apiKey := "fb0aff78982d4d7c9835ada840fdabda" // Hardcoded for demo
	req, _ := http.NewRequest("GET", "https://api.football-data.org/v4/competitions", nil)
	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Handler to fetch real standings from football-data.org
func (a *App) handleRealStandings(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	if league == "" {
		http.Error(w, "league param required", http.StatusBadRequest)
		return
	}
	apiKey := "fb0aff78982d4d7c9835ada840fdabda"
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings", league)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Handler to fetch real fixtures from football-data.org
func (a *App) handleRealFixtures(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	if league == "" {
		http.Error(w, "league param required", http.StatusBadRequest)
		return
	}
	apiKey := "fb0aff78982d4d7c9835ada840fdabda"
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/matches?status=SCHEDULED", league)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Handler to simulate N weeks for a real league (homogeneous strengths)
func (a *App) handleRealSimulate(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	weekParam := r.URL.Query().Get("week")
	var weekLimit int
	if weekParam != "" {
		weekLimit, _ = strconv.Atoi(weekParam)
	}
	apiKey := "fb0aff78982d4d7c9835ada840fdabda"
	// Fetch teams (standings)
	standingsUrl := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings", league)
	req, _ := http.NewRequest("GET", standingsUrl, nil)
	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	var standingsResp struct {
		Standings []struct {
			Table []struct {
				Team struct{ ID int; Name string }
			} `json:"table"`
		} `json:"standings"`
	}
	json.NewDecoder(resp.Body).Decode(&standingsResp)
	if len(standingsResp.Standings) == 0 {
		http.Error(w, "No standings data available for this league/season", http.StatusNotFound)
		return
	}
	var teams []models.Team
	for _, row := range standingsResp.Standings[0].Table {
		teams = append(teams, models.Team{ID: row.Team.ID, Name: row.Team.Name, Strength: 75})
	}
	// Fetch all matches (not just scheduled)
	matchesUrl := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/matches", league)
	req2, _ := http.NewRequest("GET", matchesUrl, nil)
	req2.Header.Set("X-Auth-Token", apiKey)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp2.Body.Close()
	var matchesResp struct {
		Matches []struct {
			Matchday int `json:"matchday"`
			Status   string `json:"status"`
			HomeTeam struct{ ID int } `json:"homeTeam"`
			AwayTeam struct{ ID int } `json:"awayTeam"`
			Score struct {
				FullTime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"fullTime"`
			} `json:"score"`
		} `json:"matches"`
	}
	json.NewDecoder(resp2.Body).Decode(&matchesResp)
	// Build table from all finished matches up to weekLimit (if set)
	table := make(models.LeagueTable, len(teams))
	for i, t := range teams {
		table[i] = models.TeamStanding{TeamID: t.ID, TeamName: t.Name}
	}
	var maxPlayedMatchday int
	for _, m := range matchesResp.Matches {
		if m.Status == "FINISHED" && (weekLimit == 0 || m.Matchday <= weekLimit) {
			predictor.ApplyResultByID(&table, m.HomeTeam.ID, m.AwayTeam.ID, m.Score.FullTime.Home, m.Score.FullTime.Away)
			if m.Matchday > maxPlayedMatchday {
				maxPlayedMatchday = m.Matchday
			}
		}
	}
	// Collect real results up to weekLimit
	var realResults []map[string]interface{}
	for w := 1; w <= maxPlayedMatchday && (weekLimit == 0 || w <= weekLimit); w++ {
		var weekFixtures []map[string]interface{}
		for _, m := range matchesResp.Matches {
			if m.Matchday == w && m.Status == "FINISHED" {
				weekFixtures = append(weekFixtures, map[string]interface{}{
					"week": w,
					"home_team_id": m.HomeTeam.ID,
					"away_team_id": m.AwayTeam.ID,
					"home_goals": m.Score.FullTime.Home,
					"away_goals": m.Score.FullTime.Away,
				})
			}
		}
		realResults = append(realResults, map[string]interface{}{
			"week": w,
			"fixtures": weekFixtures,
		})
	}
	// Find actual champion from final table
	finalTable := make(models.LeagueTable, len(teams))
	for i, t := range teams {
		finalTable[i] = models.TeamStanding{TeamID: t.ID, TeamName: t.Name}
	}
	for _, m := range matchesResp.Matches {
		if m.Status == "FINISHED" {
			predictor.ApplyResultByID(&finalTable, m.HomeTeam.ID, m.AwayTeam.ID, m.Score.FullTime.Home, m.Score.FullTime.Away)
		}
	}
	actualChampionID := predictor.FindChampion(finalTable)
	var actualChampionName string
	for _, t := range teams {
		if t.ID == actualChampionID {
			actualChampionName = t.Name
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"realResults": realResults,
		"table":     table,
		"teams":     teams,
		"actualChampion": actualChampionName,
	})
}

// Handler to predict championship probabilities for a real league
func (a *App) handleRealPredict(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	modelName := r.URL.Query().Get("model")
	sims, _ := strconv.Atoi(r.URL.Query().Get("sims"))
	weekParam := r.URL.Query().Get("week")
	var weekLimit int
	if weekParam != "" {
		weekLimit, _ = strconv.Atoi(weekParam)
	}
	if league == "" || modelName == "" || sims < 1 {
		http.Error(w, "league, model, sims param required", http.StatusBadRequest)
		return
	}
	apiKey := "fb0aff78982d4d7c9835ada840fdabda"
	// Fetch teams (standings)
	standingsUrl := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings", league)
	req, _ := http.NewRequest("GET", standingsUrl, nil)
	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	var standingsResp struct {
		Standings []struct {
			Table []struct {
				Team struct{ ID int; Name string }
			} `json:"table"`
		} `json:"standings"`
	}
	json.NewDecoder(resp.Body).Decode(&standingsResp)
	if len(standingsResp.Standings) == 0 {
		http.Error(w, "No standings data available for this league/season", http.StatusNotFound)
		return
	}
	var teams []models.Team
	for _, row := range standingsResp.Standings[0].Table {
		teams = append(teams, models.Team{ID: row.Team.ID, Name: row.Team.Name, Strength: 75})
	}
	// Fetch all matches
	matchesUrl := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/matches", league)
	req2, _ := http.NewRequest("GET", matchesUrl, nil)
	req2.Header.Set("X-Auth-Token", apiKey)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp2.Body.Close()
	var matchesResp struct {
		Matches []struct {
			Matchday int `json:"matchday"`
			Status   string `json:"status"`
			HomeTeam struct{ ID int } `json:"homeTeam"`
			AwayTeam struct{ ID int } `json:"awayTeam"`
			Score struct {
				FullTime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"fullTime"`
			} `json:"score"`
		} `json:"matches"`
	}
	json.NewDecoder(resp2.Body).Decode(&matchesResp)
	// Build table from all finished matches up to weekLimit (if set)
	table := make(models.LeagueTable, len(teams))
	for i, t := range teams {
		table[i] = models.TeamStanding{TeamID: t.ID, TeamName: t.Name}
	}
	for _, m := range matchesResp.Matches {
		if m.Status == "FINISHED" && (weekLimit == 0 || m.Matchday <= weekLimit) {
			predictor.ApplyResultByID(&table, m.HomeTeam.ID, m.AwayTeam.ID, m.Score.FullTime.Home, m.Score.FullTime.Away)
		}
	}
	// Use all future fixtures after weekLimit
	var matchups []interfaces.Matchup
	for _, m := range matchesResp.Matches {
		if m.Matchday > weekLimit && m.Status != "FINISHED" {
			matchups = append(matchups, interfaces.Matchup{
				HomeTeamID: m.HomeTeam.ID,
				AwayTeamID: m.AwayTeam.ID,
			})
		}
	}
	var model predictor.Predictor
	switch modelName {
	case "elo":
		model = predictor.NewEloMonteCarlo()
	case "bt":
		model = predictor.NewBradleyTerryMonteCarlo()
	case "logistic":
		model = predictor.NewLogisticMonteCarlo()
	case "mlp":
		model = predictor.NewMLPNeuralNetMonteCarlo()
	case "bivariate":
		model = predictor.NewBivariatePoissonMonteCarlo()
	case "zip":
		model = predictor.NewZeroInflatedPoissonMonteCarlo()
	default:
		model = predictor.NewPoissonMonteCarlo()
	}
	probs, err := model.Predict(teams, table, matchups, sims)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	teamNames := make(map[int]string)
	for _, t := range teams {
		teamNames[t.ID] = t.Name
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"model": model.Name(),
		"sims":  sims,
		"probs": probs,
		"teamNames": teamNames,
	})
}

/* ------------ Edit match result ------------------------------------------- */

// Edit a match result by week, home_team_id, away_team_id
func (a *App) handleEditMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Week       int `json:"week"`
		HomeTeamID int `json:"home_team_id"`
		AwayTeamID int `json:"away_team_id"`
		HomeGoals  int `json:"home_goals"`
		AwayGoals  int `json:"away_goals"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	match := models.MatchResult{
		Week: req.Week,
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		HomeGoals: req.HomeGoals,
		AwayGoals: req.AwayGoals,
	}
	repoImpl, ok := a.repo.(*repo.PostgresRepo)
	if !ok {
		http.Error(w, "edit not supported for this repo type", http.StatusInternalServerError)
		return
	}
	if err := repoImpl.UpdateMatch(match); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/* ------------ Simulate all remaining matches ------------------------------- */

// Simulate all remaining matches in one go
func (a *App) handleSimulateAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	type WeekResult struct {
		Week     int                  `json:"week"`
		Fixtures []models.MatchResult `json:"fixtures"`
	}
	var simulated []WeekResult
	for {
		fixtures, _ := a.nextFixtures()
		if len(fixtures) == 0 {
			break
		}
		teams, _ := a.repo.ListTeams()
		var played []models.MatchResult
		for _, f := range fixtures {
			home := findTeam(teams, f.HomeTeamID)
			away := findTeam(teams, f.AwayTeamID)
			m, _ := a.gen.Generate(home, away)
			m.Week = f.Week
			_ = a.repo.SaveMatch(m)
			played = append(played, m)
		}
		simulated = append(simulated, WeekResult{Week: fixtures[0].Week, Fixtures: played})
	}
	table, _ := a.repo.GetTable()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"simulated": simulated,
		"table":     table,
	})
}

/* -------------------------------------------------------------------------- */
/*                               Helper utils                                */
/* -------------------------------------------------------------------------- */

func findTeam(teams []models.Team, id int) models.Team {
	for _, t := range teams {
		if t.ID == id {
			return t
		}
	}
	return models.Team{ID: id, Name: fmt.Sprintf("Team %d", id), Strength: 50}
}

/* -------------------------------------------------------------------------- */
/*                                   main                                    */
/* -------------------------------------------------------------------------- */

func main() {
	rand.Seed(time.Now().UnixNano())

	dsn := "postgres://insider:insider@localhost:5432/insider_league?sslmode=disable"
	repository, err := repo.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	app := &App{
		repo: repository,
		gen:  service.NewPoissonGenerator(),
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/table", app.handleGetTable)
	http.HandleFunc("/simulate/week", app.handleSimulateWeek)
	http.HandleFunc("/simulate/weeks", app.handleSimulateWeeks)
	http.HandleFunc("/predict", app.handlePredict)
	http.HandleFunc("/reset", app.handleReset)
	http.HandleFunc("/real/leagues", app.handleRealLeagues)
	http.HandleFunc("/real/standings", app.handleRealStandings)
	http.HandleFunc("/real/fixtures", app.handleRealFixtures)
	http.HandleFunc("/real/simulate", app.handleRealSimulate)
	http.HandleFunc("/real/predict", app.handleRealPredict)
	http.HandleFunc("/edit_match", app.handleEditMatch)
	http.HandleFunc("/simulate/all", app.handleSimulateAll)

	log.Println("Listening on :8080 …")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

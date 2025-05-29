package repo

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
    "github.com/musta/insider-league/internal/interfaces"
    "github.com/musta/insider-league/internal/models"
    "github.com/musta/insider-league/internal/service"
)

// PostgresRepo implements interfaces.Repository.
type PostgresRepo struct {
    db      *sql.DB
    updater interfaces.TableUpdater
}

// NewPostgresRepo connects to Postgres and returns a Repository.
func NewPostgresRepo(dsn string) (interfaces.Repository, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return &PostgresRepo{
        db:      db,
        updater: service.NewTableUpdater(),
    }, nil
}

func (r *PostgresRepo) ListTeams() ([]models.Team, error) {
    rows, err := r.db.Query(`
        SELECT id, name, strength
        FROM teams
        ORDER BY id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teams []models.Team
    for rows.Next() {
        var t models.Team
        if err := rows.Scan(&t.ID, &t.Name, &t.Strength); err != nil {
            return nil, err
        }
        fmt.Printf("[REPO] ListTeams: id=%d, name=%s, strength=%d\n", t.ID, t.Name, t.Strength)
        teams = append(teams, t)
    }
    return teams, nil
}

func (r *PostgresRepo) SaveTeam(team models.Team) error {
    fmt.Printf("[REPO] SaveTeam: name=%s, strength=%d\n", team.Name, team.Strength)
    // Always insert with default serial id, ignore provided id
    _, err := r.db.Exec(`
        INSERT INTO teams (name, strength)
        VALUES ($1, $2)
    `, team.Name, team.Strength)
    return err
}

func (r *PostgresRepo) SaveMatch(m models.MatchResult) error {
    _, err := r.db.Exec(`
        INSERT INTO matches (week, home_team_id, away_team_id, home_goals, away_goals)
        VALUES ($1, $2, $3, $4, $5)
    `, m.Week, m.HomeTeamID, m.AwayTeamID, m.HomeGoals, m.AwayGoals)
    return err
}

// UpdateMatch updates the result of a match by week, home_team_id, and away_team_id.
func (r *PostgresRepo) UpdateMatch(m models.MatchResult) error {
    res, err := r.db.Exec(`
        UPDATE matches
        SET home_goals = $1, away_goals = $2
        WHERE week = $3 AND home_team_id = $4 AND away_team_id = $5
    `, m.HomeGoals, m.AwayGoals, m.Week, m.HomeTeamID, m.AwayTeamID)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return fmt.Errorf("no such match found to update")
    }
    return nil
}

func (r *PostgresRepo) GetTable() (models.LeagueTable, error) {
    teams, err := r.ListTeams()
    if err != nil {
        return nil, err
    }
    table := make(models.LeagueTable, len(teams))
    for i, t := range teams {
        table[i] = models.TeamStanding{
            TeamID:   t.ID,
            TeamName: t.Name,
        }
    }

    rows, err := r.db.Query(`
        SELECT week, home_team_id, away_team_id, home_goals, away_goals
        FROM matches
        ORDER BY week, id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var m models.MatchResult
        if err := rows.Scan(
            &m.Week, &m.HomeTeamID, &m.AwayTeamID,
            &m.HomeGoals, &m.AwayGoals,
        ); err != nil {
            return nil, err
        }
        table, err = r.updater.Update(table, m)
        if err != nil {
            return nil, fmt.Errorf("updating table: %w", err)
        }
    }
    return table, nil
}

func (r *PostgresRepo) GetAllMatches() ([]models.MatchResult, error) {
    rows, err := r.db.Query(`
        SELECT week, home_team_id, away_team_id, home_goals, away_goals
        FROM matches
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.MatchResult
    for rows.Next() {
        var m models.MatchResult
        if err := rows.Scan(
            &m.Week, &m.HomeTeamID, &m.AwayTeamID,
            &m.HomeGoals, &m.AwayGoals,
        ); err != nil {
            return nil, err
        }
        all = append(all, m)
    }
    return all, nil
}

func (r *PostgresRepo) ListRemainingMatches() ([]interfaces.Matchup, error) {
    teams, err := r.ListTeams()
    if err != nil {
        return nil, err
    }
    played, err := r.GetAllMatches()
    if err != nil {
        return nil, err
    }

    // build set of played fixtures
    playedSet := make(map[string]bool)
    for _, m := range played {
        key := fmt.Sprintf("%d-%d-%d", m.Week, m.HomeTeamID, m.AwayTeamID)
        playedSet[key] = true
    }

    // generate full schedule
    schedule := generateDoubleRoundRobin(teams)

    // filter out already played
    var remaining []interfaces.Matchup
    for _, f := range schedule {
        key := fmt.Sprintf("%d-%d-%d", f.Week, f.HomeTeamID, f.AwayTeamID)
        if !playedSet[key] {
            remaining = append(remaining, f)
        }
    }
    return remaining, nil
}

func generateDoubleRoundRobin(teams []models.Team) []interfaces.Matchup {
    // extract IDs
    ids := make([]int, len(teams))
    for i, t := range teams {
        ids[i] = t.ID
    }
    // if odd, add a dummy bye
    if len(ids)%2 == 1 {
        ids = append(ids, 0)
    }
    n := len(ids)
    rounds := n - 1

    build := func(offset int) []interfaces.Matchup {
        slice := append([]int(nil), ids...)
        var fixtures []interfaces.Matchup
        for week := 1; week <= rounds; week++ {
            for i := 0; i < n/2; i++ {
                home, away := slice[i], slice[n-1-i]
                if home != 0 && away != 0 {
                    fixtures = append(fixtures, interfaces.Matchup{
                        Week:       week + offset,
                        HomeTeamID: home,
                        AwayTeamID: away,
                    })
                }
            }
            // rotate but keep first fixed
            slice = append(slice[:1], append(slice[n-1:], slice[1:n-1]...)...)
        }
        return fixtures
    }

    first := build(0)
    second := make([]interfaces.Matchup, len(first))
    for i, f := range first {
        second[i] = interfaces.Matchup{
            Week:       f.Week + rounds,
            HomeTeamID: f.AwayTeamID,
            AwayTeamID: f.HomeTeamID,
        }
    }
    return append(first, second...)
}

func (r *PostgresRepo) ResetMatches() error {
    _, err := r.db.Exec(`TRUNCATE matches RESTART IDENTITY`)
    return err
}

func (r *PostgresRepo) ResetTeams() error {
    fmt.Println("[REPO] TRUNCATE teams RESTART IDENTITY CASCADE")
    _, err := r.db.Exec(`TRUNCATE teams RESTART IDENTITY CASCADE`)
    return err
}

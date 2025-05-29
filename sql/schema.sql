-- sql/schema.sql

-- Teams table
CREATE TABLE teams (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  strength INT NOT NULL
);

-- Matches table
CREATE TABLE matches (
  id SERIAL PRIMARY KEY,
  week INT NOT NULL,
  home_team_id INT NOT NULL REFERENCES teams(id),
  away_team_id INT NOT NULL REFERENCES teams(id),
  home_goals INT NOT NULL,
  away_goals INT NOT NULL
);

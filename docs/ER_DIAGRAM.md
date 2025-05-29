# Database ER Diagram (Markdown)

```
+---------+        +---------+
| teams   |        | matches |
+---------+        +---------+
| id (PK) |<----+  | id (PK) |
| name    |     |  | week    |
| strength|     +--| home_team_id (FK -> teams.id)
+---------+        | away_team_id (FK -> teams.id)
                   | home_goals
                   | away_goals
                   +---------+
```

- **teams**: Stores team information and strength.
- **matches**: Stores match results, week, and references to home/away teams.
- **PK**: Primary Key
- **FK**: Foreign Key

This diagram shows the relationship between teams and matches. Each match references two teams (home and away).

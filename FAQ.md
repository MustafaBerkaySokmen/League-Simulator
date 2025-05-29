# FAQ & Troubleshooting

## 1. The server won't start or can't connect to the database
- Make sure PostgreSQL is running and accessible at the DSN specified in `main.go`.
- If using Docker, run `docker-compose up -d` first.
- Check that the user, password, and database match those in your Go code.

## 2. How do I reset the league or change the number of teams?
- Use the `/reset?teams=N&type=random|homogeneous` endpoint with a POST request.
- Example: `curl -X POST "http://localhost:8080/reset?teams=4&type=random"`

## 3. How do I simulate all matches at once?
- Use the `/simulate/all` endpoint with a POST request.
- Example: `curl -X POST http://localhost:8080/simulate/all`

## 4. How do I edit a match result?
- Use the `/edit_match` endpoint with a POST request and a JSON body specifying the week, teams, and new result.
- Example:
  ```json
  {
    "week": 3,
    "home_team_id": 1,
    "away_team_id": 2,
    "home_goals": 2,
    "away_goals": 2
  }
  ```

## 5. How do I get real league data?
- Use `/real/leagues`, `/real/standings?league=CODE`, and `/real/fixtures?league=CODE` endpoints.

## 6. How do I run the tests?
- Run `go test ./...` from the project root to execute all tests.

## 7. I get a CORS error when calling the API from a browser
- The default server does not set CORS headers. For browser-based API calls, you may need to add CORS support.

## 8. Where can I find example API requests?
- See the included Postman collection in `docs/postman_collection.json`.

## 9. How do I update dependencies?
- Run `go mod tidy` to clean up and update Go module dependencies.

## 10. Who do I contact for help?
- Open an issue on GitHub or email the maintainer (see README Contact section).

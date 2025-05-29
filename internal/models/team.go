package models

// Team represents a club in the mini-league.
type Team struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Strength int    `json:"strength"` // higher = stronger
}

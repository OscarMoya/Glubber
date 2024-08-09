package model

// Passenger represents a passenger in the system
type Passenger struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

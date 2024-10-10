package model

type Movie struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Director       string `json:"director"`
	Duration       int    `json:"duration"`
	AvailableSeats int    `json:"available_seats"`
}

type Booking struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
	Seats   int    `json:"seats"`
	Status  string `json:"status"`
	Payment string `json:"payment"`
}

type Users struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

package pkg

import "time"

// auth
type ScoreServerAuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// query
type QueryRequest struct {
	Query string `json:"query"`
}

// get penalties response
type ProblemBody struct {
	Title string `json:"title"`
}
type Problem struct {
	Code string      `json:"code"`
	Body ProblemBody `json:"body"`
}

type Team struct {
	Number int `json:"number"`
}

type Penalty struct {
	ID        string    `json:"id"`
	Problem   Problem   `json:"problem"`
	Team      Team      `json:"team"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetPenaltiesResponse struct {
	Data struct {
		Penalties []Penalty `json:"penalties"`
	} `json:"data"`
}

type ResponseErrorExtensions struct {
	Code      string `json:"code"`
	RequestID string `json:"requestId"`
}

type ResponseError struct {
	Message    string                  `json:"message"`
	Extensions ResponseErrorExtensions `json:"extensions"`
}

type ErrorsResponse struct {
	Errors []ResponseError `json:"errors"`
}

package models

type Pagination struct {
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type ResponsePages struct {
	Pagination Pagination `json:"pagination"`
}

type ResponseEvents struct {
	Events []Event `json:"events"`
}

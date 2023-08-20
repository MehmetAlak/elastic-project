package model

import "time"

type CreateResponse struct {
	ID string `json:"id"`
}

type FindResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Job        string     `json:"job"`
	ChildNames []string   `json:"childNames"`
	Comment    string     `json:"comment"`
	CreatedAt  *time.Time `json:"created_at"`
}

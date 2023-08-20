package model

type CreateRequest struct {
	Name       string   `json:"name"`
	Job        string   `json:"job"`
	ChildNames []string `json:"childNames"`
	Comment    string   `json:"comment"`
}

type UpdateRequest struct {
	Name       string   `json:"name"`
	Job        string   `json:"job"`
	ChildNames []string `json:"childNames"`
	Comment    string   `json:"comment"`
}

type DeleteRequest struct {
	ID string
}

type FindRequest struct {
	ID string
}

type FindByRequest struct {
	QueryType string `json:"queryType"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

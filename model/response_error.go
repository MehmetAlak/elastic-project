package model

type ResponseError struct {
	StatusCode int
	Err        error
}

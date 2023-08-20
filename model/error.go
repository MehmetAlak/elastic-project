package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

const (
	StatusOK                  int = 200
	StatusNoContent           int = 204
	StatusBadRequest          int = 400
	StatusUnauthorized        int = 401
	StatusNotFound            int = 404
	StatusInternalServerError int = 500
)

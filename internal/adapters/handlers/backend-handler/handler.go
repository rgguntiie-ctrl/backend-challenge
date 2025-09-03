package handlers

import "github.com/kanta/backend-challenge/internal/core/ports"

type BackEndHandler interface {
}

type backEndHandler struct {
	service ports.Service
}

func NewBackEndHandler(
	service ports.Service,
) BackEndHandler {
	return &backEndHandler{
		service,
	}
}

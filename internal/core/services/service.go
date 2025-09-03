package services

import "github.com/kanta/backend-challenge/internal/core/ports"

type service struct {
}

func NewBackEndService() ports.Service {
	return &service{}
}

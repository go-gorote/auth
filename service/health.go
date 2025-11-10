package service

import "github.com/go-gorote/gorote"

func (s *AppService) Health() (*gorote.Health, error) {
	return gorote.HealthGorm(s.DB)
}

package service

import (
	"github.com/38888/nunu-layout-advanced/internal/repository"
	"github.com/38888/nunu-layout-advanced/pkg/jwt"
	"github.com/38888/nunu-layout-advanced/pkg/log"
	"github.com/38888/nunu-layout-advanced/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}

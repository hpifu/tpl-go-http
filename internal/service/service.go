package service

import (
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	secure bool
	domain string
}

func NewService(
	secure bool,
	domain string,
) *Service {
	return &Service{
		secure: secure,
		domain: domain,
	}
}

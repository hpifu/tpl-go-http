package service

import "github.com/sirupsen/logrus"

type Service struct {
	secure    bool
	domain    string
	infoLog   *logrus.Logger
	warnLog   *logrus.Logger
	accessLog *logrus.Logger
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func NewService(
	secure bool,
	domain string,
) *Service {
	return &Service{
		secure:    secure,
		domain:    domain,
		infoLog:   logrus.New(),
		warnLog:   logrus.New(),
		accessLog: logrus.New(),
	}
}

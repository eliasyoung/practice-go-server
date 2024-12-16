package service

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func SetupLogger(l *zap.SugaredLogger) {
	logger = l
}

type Services struct {
	User *UserService
	Post *PostService
}

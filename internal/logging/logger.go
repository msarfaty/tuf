package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

const (
	TUF_LOG_MODE_KEY               = "TUF_LOG_MODE"
	TUF_LOG_MODE_VALUE_CI          = "CI"
	TUF_LOG_MODE_VALUE_DEVELOPMENT = "DEVELOPMENT"
)

func GetLogger() (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error
	switch os.Getenv(TUF_LOG_MODE_KEY) {
	case TUF_LOG_MODE_VALUE_CI:
		logger, err = zap.NewProduction()
	case TUF_LOG_MODE_VALUE_DEVELOPMENT:
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, fmt.Errorf("unexpected failure to retrieve logger: %w", err)
	}

	defer logger.Sync()

	return logger.Sugar(), nil
}

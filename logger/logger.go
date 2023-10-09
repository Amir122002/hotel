package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

func InitLogger() (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetReportCaller(true)

	logger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("./logger/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(file)

	return logger, nil

}

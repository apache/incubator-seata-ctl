package tool

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// Logger Global logger instance
var Logger *logrus.Logger
var once sync.Once

// InitLogger initializes the global logger instance
func InitLogger() *logrus.Logger {
	once.Do(func() { // Ensures initialization is only done once
		Logger = logrus.New()
		Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		Logger.SetLevel(logrus.InfoLevel)
	})
	return Logger
}

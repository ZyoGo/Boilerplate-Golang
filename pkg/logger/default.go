package logger

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Define a global zerolog.Logger instance and a sync.Once for ensuring singleton initialization.
var (
	log  zerolog.Logger
	once sync.Once
)

// Get returns the initialized zerolog.Logger instance.
func Get() zerolog.Logger {
	// sync.Once to ensure the logger is initialized only once.
	once.Do(func() {
		// Configure lumberjack.Logger to handle log file rotation.
		fileLumberjack := &lumberjack.Logger{
			Filename:   "logs/noodle-svc.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}

		// Create a multi-writer to write logs to both stdout and the log file.
		multiWriter := zerolog.MultiLevelWriter(os.Stdout, fileLumberjack)

		// Set up a channel to listen for SIGHUP signals (usually sent to reload configuration).
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)

		// Start a goroutine that rotates the log file whenever a SIGHUP signal is received.
		go func() {
			for {
				<-c
				fileLumberjack.Rotate()
			}
		}()

		// Initialize the zerolog.Logger with the multi-writer and include timestamps in the log entries.
		log = zerolog.New(multiWriter).
			With().
			Timestamp().
			Logger()
	})

	// Return the logger instance.
	return log
}

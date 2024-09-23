package logging

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
)

type LogManager struct {
	logger     *slog.Logger
	child      *slog.Logger
	attributes *sync.Map
}

var (
	mng    LogManager
	Logger *slog.Logger
)

func SetupLogger(logLevel string, logFormat string, logDest string) error {
	var file *os.File

	switch logDest {
	case "stdout":
		file = os.Stdout
	default:
		var err error
		file, err = os.OpenFile(logDest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("cannot open log file: %w", err)
		}
	}

	var level slog.Level
	var addSource = false

	//attributes := make(map[string]any)
	var attributes sync.Map
	switch logLevel {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "debug":
		level = slog.LevelDebug
		addSource = true

		buildInfo, _ := debug.ReadBuildInfo()
		if buildInfo == nil {
			buildInfo = &debug.BuildInfo{}
		}
		attributes.Store("pid", fmt.Sprintf("%d", os.Getpid()))
		attributes.Store("go_version", buildInfo.GoVersion)
	default:
		return errors.New("unrecognized log level")
	}

	var handler slog.Handler

	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	case "text":
		handler = slog.NewTextHandler(file, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	default:
		return errors.New("unrecognized log format")
	}

	mng = LogManager{
		logger:     slog.New(handler),
		attributes: &attributes,
	}
	childRebuild()

	return nil
}

func AddAttribute(key string, value string) {
	mng.attributes.Store(key, value)
	childRebuild()
}

func RemoveAttribute(key string) {
	mng.attributes.Delete(key)
	childRebuild()
}

func childRebuild() {
	attributes := make([]any, 0)

	mng.attributes.Range(func(k, v interface{}) bool {
		attributes = append(attributes, k, v)
		return true
	})

	mng.child = mng.logger.With(
		slog.Group("extra", attributes...),
	)

	Logger = mng.child
}

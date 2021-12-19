package helper

import (
	"github.com/sirupsen/logrus"
	"os"
)

type ConsoleHook struct {
	writer *os.File
}

func (hook ConsoleHook) Fire(entry *logrus.Entry) error {
	formatter := logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
	}
	bytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(bytes)

	return err
}

func (hook ConsoleHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
	}
}

func NewConsoleHook(logInErr bool) *ConsoleHook {
	if logInErr {
		return &ConsoleHook{writer: os.Stderr}
	} else {
		return &ConsoleHook{writer: os.Stdout}
	}
}

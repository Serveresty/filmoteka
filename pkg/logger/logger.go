package logger

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

func Init(logFilePath string) error {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fileWriter := io.Writer(file)
	InfoLogger = log.New(fileWriter, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(fileWriter, "ERROR: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(fileWriter, "Warning: ", log.Ldate|log.Ltime)

	return nil
}

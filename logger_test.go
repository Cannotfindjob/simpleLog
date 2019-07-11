package simpleLogger

import (
	"log"
	"testing"
)

func TestLevel(t *testing.T) {

	slog := NewSimpleLog()
	slog.SetLevel("WARN")

	logger := log.New(slog, "", 0)
	logger.Print("[WARN] 1")
	logger.Println("[ERROR] 2")
	logger.Println("[DEBUG] 0")
	logger.Println("[WARN] 1")
}

func TestSetLevel(t *testing.T) {

	slog := NewSimpleLog()
	slog.SetLevels("DEBUG", "WARN", "ERROR")
    slog.SetLevel("ERROR")

	logger := log.New(slog, "", 0)
	logger.Print("[WARN] 1")
	logger.Println("[ERROR] 2")
	logger.Println("[DEBUG] 0")
	logger.Println("[WARN] 1")
}
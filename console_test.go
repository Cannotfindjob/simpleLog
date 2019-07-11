package simpleLogger

import (
	"log"
	"testing"
)

func TestNewConsoleObject(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 100000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}


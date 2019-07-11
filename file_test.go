package simpleLogger

import (
	"encoding/json"
	"log"
	"testing"
)

func TestNewFileObject(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	var config map[string]interface{} = map[string]interface{}{"filepath":"./test.log"}

	stringConfig, _ := json.Marshal(config)


	_:logger.SetOutPut(AdapterFile, string(stringConfig))

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 100000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}

func TestNewFileObjectMaxLines1000(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	var config map[string]interface{} = map[string]interface{}{"filepath":"./test.log","max_lines":1000}

	stringConfig, _ := json.Marshal(config)


   _:logger.SetOutPut(AdapterFile, string(stringConfig))

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 10000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}

func TestNewFileObjectMaxSize10M(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	var config map[string]interface{} = map[string]interface{}{"filepath":"./test.log","max_size":10 << 20}

	stringConfig, _ := json.Marshal(config)


    _:logger.SetOutPut(AdapterFile, string(stringConfig))

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 1000000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}

func TestNewFileObjectCloseRotate(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	var config map[string]interface{} = map[string]interface{}{"filepath":"./test.log", "max_lines":1000, "rotate": false}

	stringConfig, _ := json.Marshal(config)


	_:logger.SetOutPut(AdapterFile, string(stringConfig))

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 10000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}

func TestNewFileObjectCloseCompress(t *testing.T) {
	logger := NewSimpleLog()
	logger.SetLevel("DEBUG")

	var config map[string]interface{} = map[string]interface{}{"filepath":"./test.log", "max_lines":1000, "compress": false}

	stringConfig, _ := json.Marshal(config)


_:logger.SetOutPut(AdapterFile, string(stringConfig))

	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
	for i := 0; i < 10000; i++ {
		log.Println("[DEBUG] Do Debugging")
	}

	logger.Close()
}
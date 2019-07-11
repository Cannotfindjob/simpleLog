package simpleLogger

import (
	"encoding/json"
	"os"
)

type ConsoleObject struct {
	file *os.File

}

func NewConsoleObject() Logger {
	console := &ConsoleObject {
		file: os.Stdout,

	}
	return console
}

func (c *ConsoleObject) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(jsonConfig), c)

	return err
}

func (c *ConsoleObject) Write(p []byte) (n int, err error) {
	return c.file.Write(p)
}

func (c *ConsoleObject) Flush() {

}

func (c *ConsoleObject) Close() {

}

func init() {
	Register(AdapterConsole, NewConsoleObject)
}

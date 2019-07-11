## install

	go get github.com/Cannotfindjob/simpleLog


## What adapters are supported?

default console, file


## use it

First you must import it
import (
	"github.com/Cannotfindjob/simpleLog"
	"log"
)

Then init a Log (default console adapter)

```golang
logger := NewSimpleLog()
logger.SetLevel("DEBUG")
log.SetOutput(logger)
log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
```

Use it like this:

```golang
log.Println("[DEBUG] Do Debugging")
```

## File adapter

Configure file adapter like this:

```golang
logger := NewSimpleLog()
logger.SetLevel("DEBUG")
log.SetOutput(logger)
log.SetFlags(log.LstdFlags| log.Lshortfile| log.Lmicroseconds)
log.SetLogger("file", `{"filename":"test.log"}`)
```
Config Json you can set:
```golang
config := `{
	"filepath": "./test.log",   // file path
	"perm": "0660",             // default 0660
	"rotate": true,             // default true
	"compress": true,           // default true
	"max_lines": 10000,         // default 10000
	"max_size": 500 << 20,      // default 500M
	"max_keep_days" : 7,        // default 7 day
}`
log.SetLogger("file", config)
```

## Default Levels
```golang
const(
	LevelFatal         = "FATAL"
	LevelError         = "ERROR"
	LevelWarning       = "WARN"
	LevelInformational = "INFO"
	LevelDebug         = "DEBUG"
)
```
you can set levels like this(must be in order ):
```golang
slog := NewSimpleLog()
slog.SetLevels("DEBUG", "WARN", "ERROR")
slog.SetLevel("ERROR")
```

## Main

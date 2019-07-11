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

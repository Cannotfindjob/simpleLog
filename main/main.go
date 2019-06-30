package main

import (
	"log"
	"simpleLogger"
)

func main() {

    //fmt.Println(0&(log.Lshortfile|log.Llongfile))
	//fmt.Println(0&(log.Lshortfile|log.Llongfile))
    //fmt.Println(1 << 0)
	//logger := log.New(os.Stdout, "deb", 11113421412421 )
	//logger.Print("Hello, log file!")


	//filter := &logutils.LevelFilter{
	//	Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
	//	MinLevel: logutils.LogLevel("WARN"),
	//	Writer: os.Stderr,
	//}
	//log.SetOutput(filter)

	//log.Print("DEBUG Debugging") // this will not print
	//log.Print("[WARN] Warning") // this will
	//log.Print("[EInfo] Erring") // and so will this
	//log.Print("Message I haven't updated") // and so will this




    res := simpleLogger.NewSimpleLog()
    res.SetOutPut("console",``)
    res.SetLevel("DEBUG")
	log.SetOutput(res)
	log.Print("[DEBUG] Debugging")
	//log.Print("[WARN] Warning")
	//log.Print("[INFO] Erring")
	//res.SetLevel("ssssssss")
	//log.Print("[INFO] Debugging")

	//
	//log := logs.NewLogger(10000)
	//log.SetLevel(2)
	//log.Warning("Warning")
}
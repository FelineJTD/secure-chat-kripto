package logger

import (
	"fmt"
	"log"
	"runtime"
)

var (
	verbosity int = 1 // Default at info
	// 3: Minutia, 2: Debug, 1: Info, 0: Error
)

func SetVerbosity(v int) {
	verbosity = v
}

func HandleError(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)

		if !ok {
			log.Println("Error retrieving details of error") // God Help Us
			return
		}

		errorMsg := fmt.Sprintf("[ERROR => %s:%d ]: %s", file, line, err)
		log.Println(errorMsg)
	}
}

func HandleFatal(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)

		if !ok {
			log.Println("Error retrieving details of fatal error") // God Help Us
			return
		}

		errorMsg := fmt.Sprintf("[FATAL => %s:%d ]: %s", file, line, err)
		log.Fatalln(errorMsg)
	}
}

func Info(info string) {
	if verbosity < 1 {
		return
	}

	_, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Println("Error retrieving details of info") // God Help Us, But Kinda Less
		return
	}

	infoMsg := fmt.Sprintf("[INFO => %s:%d ]: %s", file, line, info)
	log.Println(infoMsg)
}

func Debug(info string) {
	if verbosity < 2 {
		return
	}

	_, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Println("Error retrieving details of debug") // God Help Us, But Kinda Less
		return
	}

	infoMsg := fmt.Sprintf("[DEBUG => %s:%d ]: %s", file, line, info)
	log.Println(infoMsg)
}

func Minute(info string) {
	if verbosity < 3 {
		return
	}

	_, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Println("Error retrieving details of minutia") // God Help Us, But Kinda Less
		return
	}

	infoMsg := fmt.Sprintf("[MINUTIA => %s:%d ]: %s", file, line, info)
	log.Println(infoMsg)
}
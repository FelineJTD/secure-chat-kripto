package util

import (
	"fmt"
	"log"
	"runtime"
)

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
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		log.Println("Error retrieving details of info") // God Help Us, But Kinda Less
		return
	}

	infoMsg := fmt.Sprintf("[INFO => %s:%d ]: %s", file, line, info)
	log.Println(infoMsg)
}
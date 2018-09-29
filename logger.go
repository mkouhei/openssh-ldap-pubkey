package main

import (
	"log"
	"os"
)

func logging(err error) {
	if err == errVersion {
		os.Exit(0)
	} else if err != nil {
		log.Fatal(err)
	}
}

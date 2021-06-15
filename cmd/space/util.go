package main

import "log"

func handleErr(msg string, err error) {
	if err != nil {
		log.Panicf("[ERROR] %s: %v\n", msg, err)
		return
	}
	log.Printf("[INFO] %s\n", msg)
}

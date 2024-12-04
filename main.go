package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	C chan Message

	Users         []User
	GlobalEmotes  []string
	ChannelEmotes []string
)

func main() {
	// Keep open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	Page("Set Up", func() bool {
		return true
	})

	readConfig()

	getEmotes(true)
	go updateEmotes()
	go getLiveStatuses()

	C := make(chan Message)
	go Start(C)
	go Mimic(C)

	Page("Started", func() bool {
		return true
	})

	<-sc
	Page("Aborted", func() bool {
		return true
	})
}

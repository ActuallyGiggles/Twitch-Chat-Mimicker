package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	C chan Message

	Users  []User
	Emotes []string
)

func main() {
	// Keep open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// Read config file, if no config file found, start setup process
	readConfig()

	// Emote and live status processes
	configInit()
	getEmotes(true)
	go updateEmotes()
	go getLiveStatuses()

	// Start twitch IRC and Mimic processes
	C := make(chan Message)
	go Start(C)
	go Mimic(C)

	// Show started screen.
	Page("Started", func() bool {
		return true
	})

	<-sc
	// Show aborted screen
	Page("Aborted", func() bool {
		return true
	})
}

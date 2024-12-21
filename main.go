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

	// Show setup screen. Read config file, if no config file found, start setup process
	Page("Setup", func() bool {
		return true
	})
	readConfig()

	// Show initialization screen. Emote and live status processes
	Page("Initialization", func() bool {
		return true
	})
	configInit()
	getEmotes(true)
	go updateEmotes()
	go getLiveStatuses()

	// Show started screen. Start twitch IRC and Mimic processes
	Page("Started", func() bool {
		return true
	})
	C := make(chan Message)
	go Start(C)
	go Mimic(C)

	<-sc
	// Show aborted screen
	Page("Aborted", func() bool {
		return true
	})
}

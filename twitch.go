package main

import (
	"github.com/gempir/go-twitch-irc/v3"
)

var client *twitch.Client

// Start creates a twitch client and connects it.
func Start(in chan Message) {
	client = &twitch.Client{}
	client = twitch.NewClient(Config.Name, "oauth:"+Config.OAuth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		m := Message{
			Channel: message.Channel,
			Message: message.Message,
		}

		in <- m
	})

	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		Join(user.Name)
	}

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

// Say sends a message to a specific twitch chatroom.
func Say(channel string, message string) {
	client.Say(channel, message)
}

// Join joins a twitch chatroom.
func Join(channel string) {
	client.Join(channel)
}

// Depart departs a twitch chatroom.
func Depart(channel string) {
	client.Depart(channel)
}

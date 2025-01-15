package main

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/pterm/pterm"
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

func Respond(u *User, message string) {
	u.Busy = true
	u.LastSentEmote = message

	t := time.Now()
	var waitTime int

	if Config.IntervalMin == Config.IntervalMax {
		waitTime = Config.IntervalMin
	} else {
		waitTime = RandomNumber(Config.IntervalMin, Config.IntervalMax)
	}

	delay := RandomNumber(0, 5)

	pterm.Success.Printf("%-20s     %s\n%s", strings.ToUpper(u.Name), message, pterm.Sprintf(pterm.Gray("%d:%02d | Delay %ds | Cooldown: %s"), t.Hour(), t.Minute(), delay, secondsToMinutes(waitTime)))
	pterm.Println()
	pterm.Println()

	if len(Config.Channels) == 1 {
		go countdown(waitTime + delay)
	}

	time.Sleep(time.Duration(delay) * time.Second)
	Say(u.Name, message)
	time.Sleep(time.Duration(waitTime+1) * time.Second)

	u.Busy = false
}

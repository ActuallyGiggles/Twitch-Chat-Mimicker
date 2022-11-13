package main

import (
	"fmt"
	"strings"
	"time"
)

func Mimic(C chan Message) {
messageRange:
	for c := range C {
		u := c.Channel
		m := c.Message

		if updatingEmotes {
			continue
		}

		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			if user.Name == u {
				if user.Busy || !user.IsLive {
					continue messageRange
				}

				e := ParseEmote(m)

				if e == "" {
					continue messageRange
				}

				if len(user.Emotes) == 0 {
					if e == "" {
						continue messageRange
					} else {
						user.Emotes = append(user.Emotes, e)
					}
					continue messageRange
				}

				if len(user.Emotes) == 1 {
					user.Emotes = append(user.Emotes, e)
					continue messageRange
				}

				if len(user.Emotes) == 2 {
					user.Emotes = append(user.Emotes, e)
					if user.Emotes[0] == user.Emotes[1] && (user.Emotes[0] != "" || user.Emotes[1] != "") {
						go Respond(u, user.Emotes[0])
					} else if user.Emotes[1] == user.Emotes[2] && (user.Emotes[1] != "" || user.Emotes[2] != "") {
						go Respond(u, user.Emotes[1])
					} else if user.Emotes[0] == user.Emotes[2] && (user.Emotes[1] != "" || user.Emotes[2] != "") {
						go Respond(u, user.Emotes[2])
					}

					user.Emotes = nil
				}
			}
		}
	}
}

func ParseEmote(message string) (eJ string) {
	s := strings.Split(message, " ")

	var eS []string
loop1:
	for _, w := range s {
		for _, emote := range Emotes {
			if w == emote {
				for _, blacked := range Config.BlacklistEmotes {
					if strings.ToLower(blacked) == strings.ToLower(emote) {
						return ""
					}
				}
				eS = append(eS, emote)
				continue loop1
			}
		}
		for _, emote := range GlobalEmotes {
			if w == emote {
				for _, blacked := range Config.BlacklistEmotes {
					if strings.ToLower(blacked) == strings.ToLower(emote) {
						return ""
					}
				}
				eS = append(eS, emote)
				continue loop1
			}
		}
	}

	eJ = strings.Join(eS, " ")

	return eJ
}

func Respond(channel string, message string) {
	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		if user.Name == channel {
			user.Busy = true
			go func() {
				if Config.IntervalMin == Config.IntervalMax {
					time.Sleep(time.Duration(Config.IntervalMin) * time.Minute)
					user.Busy = false
				} else {
					time.Sleep(time.Duration(RandomNumber(Config.IntervalMin, Config.IntervalMax)) * time.Minute)
					user.Busy = false
				}
			}()

			fmt.Printf("[Said In %s]: %s\n", channel, message)

			time.Sleep(time.Duration(RandomNumber(2, 10)) * time.Second)
			Say(channel, message)
		}
	}
}

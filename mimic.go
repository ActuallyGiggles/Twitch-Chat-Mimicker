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

		//fmt.Println("message recieved")

		if updatingEmotes {
			//fmt.Println("updating emotes")
			continue
		}

		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			if user.Name == u {
				//fmt.Println("found user")

				if user.Busy || !user.IsLive {
					//fmt.Printf("%s busy: %t, %s live: %t\n", user.Name, user.Busy, user.Name, user.IsLive)
					continue messageRange
				}

				e := ParseEmote(m)

				//fmt.Printf("parsed %s into -> \n\t%s\n", m, e)

				exists := false
				for i := 0; i < len(user.DetectedEmotes); i++ {
					emote := &user.DetectedEmotes[i]
					if e == emote.Name {
						//fmt.Println("found emote")
						exists = true
						emote.Value++
					}

					if emote.Value >= Config.MessageThreshold {
						//fmt.Println("responding")
						go Respond(user, emote.Name)
						user.Messages = 0
						user.DetectedEmotes = nil
						continue messageRange
					}
				}

				if !exists && e != "" {
					//fmt.Println("emote doesn't exist, adding")
					entry := Emote{
						Name:  e,
						Value: 1,
					}
					user.DetectedEmotes = append(user.DetectedEmotes, entry)
				}

				user.Messages++

				if user.Messages > Config.MessageSample {
					//fmt.Println("starting new sample")
					user.Messages = 0
					user.DetectedEmotes = nil
				}
			}
		}
	}
}

func printDetectedEmotes(user *User) {
	fmt.Println()
	for _, emoticon := range user.DetectedEmotes {
		fmt.Printf("\t[%s] %s: %d/%d\n", user.Name, emoticon.Name, emoticon.Value, Config.MessageThreshold)
	}
}

func ParseEmote(message string) (eJ string) {
	s := strings.Split(message, " ")
	var eS []string

loop1:
	for _, w := range s {
		for _, emote := range GlobalEmotes {
			if w == emote {
				for _, blacked := range Config.BlacklistEmotes {
					if strings.ToLower(blacked) == strings.ToLower(emote) {
						//fmt.Println("emote isn't allowed via globalemotes")
						return ""
					}
				}
				eS = append(eS, emote)
				continue loop1
			}
		}

		for _, emote := range ChannelEmotes {
			if w == emote {
				for _, blacked := range Config.BlacklistEmotes {
					if strings.ToLower(blacked) == strings.ToLower(emote) {
						//fmt.Println("emote isn't allowed via channelemotes")
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

func Respond(u *User, message string) {
	u.Busy = true

	rS := RandomNumber(2, 10)
	//fmt.Printf("Saying %s in %s's chat in %d seconds.\n", message, u.Name, rS)
	time.Sleep(time.Duration(rS) * time.Second)
	Say(u.Name, message)
	//clearTerminal()
	fmt.Printf("[%s] <- %s ... ", u.Name, message)

	if Config.IntervalMin == Config.IntervalMax {
		fmt.Println("waiting", Config.IntervalMin, "minutes")
		time.Sleep(time.Duration(Config.IntervalMin) * time.Minute)
	} else {
		r := RandomNumber(Config.IntervalMin, Config.IntervalMax)
		fmt.Println("waiting", r, "minutes")
		time.Sleep(time.Duration(r) * time.Minute)
	}

	u.Busy = false
}

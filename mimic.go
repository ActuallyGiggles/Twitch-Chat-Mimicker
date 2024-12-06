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
				if e == user.LastSentEmote && !Config.AllowConsecutiveDuplicates {
					continue messageRange
				}

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
					//fmt.Println("the emote that could have been sent:")

					var maxValue int
					var maxName string
					for i := 0; i < len(user.DetectedEmotes); i++ {
						emote := &user.DetectedEmotes[i]
						if emote.Value > maxValue {
							maxValue = emote.Value
							maxName = emote.Name
						}
					}

					t := time.Now()

					if maxValue > 1 {
						PrintWarning(Instructions{
							Channel: user.Name,
							Emote:   maxName,
							Note:    fmt.Sprintf("Times Used: %d/%d | Sample Size: %d\n%d:%0d", maxValue, Config.MessageThreshold, Config.MessageSample, t.Hour(), t.Minute()),
						})
					}

					//fmt.Println("starting new sample")
					user.Messages = 0
					user.DetectedEmotes = nil
				}
			}
		}
	}
}

// func printDetectedEmotes(user *User) {
// 	for _, emoticon := range user.DetectedEmotes {
// 		fmt.Printf("\t[%s] %s: %d/%d\n", user.Name, emoticon.Name, emoticon.Value, Config.MessageThreshold)
// 	}
// }

func ParseEmote(message string) string {
	// Parse for word letter combo
	for _, combo := range Config.WordLetterCombos {
		if strings.Contains(strings.ToLower(message), combo) {
			return combo
		}
	}

	sentenceSliced := strings.Split(message, " ")
	var emotesSliced []string

	// Parse for emote or emoji
loop1:
	for _, word := range sentenceSliced {
		for _, emote := range Emotes {
			if word == emote {
				for _, blacked := range Config.BlacklistEmotes {
					// Ignore messages with blacklisted emotes
					if strings.EqualFold(blacked, emote) {
						return ""
					}
				}
				emotesSliced = append(emotesSliced, emote)
				continue loop1
			}
		}
	}

	return strings.Join(emotesSliced, " ")
}

func Respond(u *User, message string) {
	u.Busy = true

	rS := RandomNumber(2, 5)
	time.Sleep(time.Duration(rS) * time.Second)
	Say(u.Name, message)
	u.LastSentEmote = message

	var waitTime int

	if Config.IntervalMin == Config.IntervalMax {
		waitTime = Config.IntervalMin
	} else {
		waitTime = RandomNumber(Config.IntervalMin, Config.IntervalMax)
	}

	t := time.Now()

	PrintSuccess(Instructions{
		Channel: u.Name,
		Emote:   message,
		Note:    fmt.Sprintf("Delay: %ds | Cooldown: %dm\n%d:%0d", rS, waitTime, t.Hour(), t.Minute()),
	})

	time.Sleep(time.Duration(waitTime) * time.Minute)

	u.Busy = false
}

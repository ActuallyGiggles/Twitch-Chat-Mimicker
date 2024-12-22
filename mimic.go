package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pterm/pterm"
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
						// fmt.Println("responding")
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
						// fmt.Println("warning")
						pterm.Warning.Printf("%-20s     %s\n%s", strings.ToUpper(user.Name), maxName, pterm.Gray(fmt.Sprintf("%d:%02d | Times Used: %d/%d | Sample Size: %d", t.Hour(), t.Minute(), maxValue, Config.MessageThreshold, Config.MessageSample)))
						pterm.Println()
						pterm.Println()
					}

					//fmt.Println("starting new sample")
					user.Messages = 0
					user.DetectedEmotes = nil
				}
			}
		}
	}
}

func ParseEmote(message string) string {
	// Parse for word letter combo
	for _, combo := range Config.WordLetterCombos {
		match, err := regexp.MatchString(`(?i)^`+combo+`$`, message)
		if err != nil {
			panic(err)
		}

		if match {
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

	t := time.Now()
	var waitTime int

	if Config.IntervalMin == Config.IntervalMax {
		waitTime = Config.IntervalMin
	} else {
		waitTime = RandomNumber(Config.IntervalMin, Config.IntervalMax)
	}

	delay := RandomNumber(0, 5)

	pterm.Success.Printf("%-20s     %s\n%s", strings.ToUpper(u.Name), message, pterm.Sprintf(pterm.Gray("%d:%02d | Delay %ds | Cooldown: %s"), delay, t.Hour(), t.Minute(), secondsToMinutes(waitTime)))
	pterm.Println()
	pterm.Println()

	time.Sleep(time.Duration(delay))
	Say(u.Name, message)
	u.LastSentEmote = message

	// countdown, _ := pterm.DefaultArea.WithRemoveWhenDone().Start(pterm.Gray("Waiting for " + secondsToMinutes(waitTime) + " seconds..."))
	// for i := waitTime; i >= 0; i-- {
	// 	countdown.Update(pterm.Gray("Waiting for " + secondsToMinutes(i) + " seconds..."))
	// 	time.Sleep(time.Second)
	// }
	// countdown.Stop()

	time.Sleep(time.Duration(waitTime) * time.Second)

	u.Busy = false
}

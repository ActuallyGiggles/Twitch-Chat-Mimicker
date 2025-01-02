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
		channel := c.Channel
		message := c.Message

		//fmt.Println("message recieved")

		if updatingEmotes {
			//fmt.Println("updating emotes")
			continue
		}

		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			if user.Name == channel {
				//fmt.Println("found user")

				if user.Busy || !user.IsLive {
					//fmt.Printf("%s busy: %t, %s live: %t\n", user.Name, user.Busy, user.Name, user.IsLive)
					continue messageRange
				}

				var thingToSend string
				emoteFound := ParseEmote(message)

				// Low priority (not unique)
				if emoteFound != "" {
					thingToSend = emoteFound
				}
				// Medium priority (more unique than regular emote)
				if onlyWordCombo := ParseOnlyWordCombo(message); onlyWordCombo != "" {
					thingToSend = onlyWordCombo
				}
				// Top priority (very unique)
				if emoteWordCombo := ParseEmoteWordCombo(emoteFound, message); emoteWordCombo != "" {
					thingToSend = emoteWordCombo
				}

				if (thingToSend == user.LastSentEmote && !Config.AllowConsecutiveDuplicates) || thingToSend == "" {
					continue messageRange
				}

				//fmt.Printf("parsed %s into -> \n\t%s\n", m, e)

				exists := false
				for i := 0; i < len(user.DetectedEmotes); i++ {
					potentialThingToSend := &user.DetectedEmotes[i]

					if thingToSend == potentialThingToSend.Name {
						//fmt.Println("found emote")
						exists = true
						potentialThingToSend.Value++
					}

					if potentialThingToSend.Value >= Config.MessageThreshold {
						// fmt.Println("responding")
						go Respond(user, potentialThingToSend.Name)
						user.Messages = 0
						user.DetectedEmotes = nil
						continue messageRange
					}
				}

				if !exists && emoteFound != "" {
					//fmt.Println("emote doesn't exist, adding")
					entry := Emote{
						Name:  emoteFound,
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

// Returns one emote, or an only emote combo
func ParseEmote(message string) string {
	sentenceSliced := strings.Split(message, " ")
	var emotesSliced []string
loop1:
	// Find an emote
	for _, word := range sentenceSliced {
		for _, emote := range Emotes {
			if word == emote {
				for _, blacklisted := range Config.BlacklistEmotes {
					// Ignore it if it's a blacklisted emote
					if strings.EqualFold(blacklisted, emote) {
						return ""
					}
				}

				// Append the emote to an emotes to send list
				emotesSliced = append(emotesSliced, emote)
				continue loop1
			} else {
				if len(emotesSliced) > 0 {
					return strings.Join(emotesSliced, " ")
				}
			}
		}
	}

	// Send the emotes to use
	if len(emotesSliced) > 0 {
		return strings.Join(emotesSliced, " ")
	}

	// If no success, return nothing
	return ""
}

// Returns an emote word combo (should account for words before, after, and both before and after emote)
func ParseEmoteWordCombo(emote, message string) string {
	if emote == "" || emote == message {
		return ""
	}

	allWordsAroundEmote := strings.Split(message, emote)
	var emoteWordComboToSend []string

	// Add words before emote
	if len(allWordsAroundEmote) > 0 {
		for _, emoteWordCombo := range Config.EmoteWordCombos {
			if emoteWordCombo == strings.TrimSpace(allWordsAroundEmote[0]) {
				emoteWordComboToSend = append(emoteWordComboToSend, emoteWordCombo)
			}
		}
	}

	// Add emote
	emoteWordComboToSend = append(emoteWordComboToSend, emote)

	// Add words after emote
	if len(allWordsAroundEmote) > 1 {
		for _, emoteWordCombo := range Config.EmoteWordCombos {
			if emoteWordCombo == strings.TrimSpace(allWordsAroundEmote[1]) {
				emoteWordComboToSend = append(emoteWordComboToSend, emoteWordCombo)
			}
		}
	}

	return strings.Join(emoteWordComboToSend, " ")
}

// Returns an only word combo to send
func ParseOnlyWordCombo(message string) string {
	for _, combo := range Config.OnlyWordCombos {
		match, err := regexp.MatchString(`(?i)^`+combo+`$`, message)
		if err != nil {
			panic(err)
		}

		if match {
			return combo
		}
	}

	return ""
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

	time.Sleep(time.Duration(delay) * time.Second)
	Say(u.Name, message)

	// countdown, _ := pterm.DefaultArea.WithRemoveWhenDone().Start(pterm.Gray("Waiting for " + secondsToMinutes(waitTime) + " seconds..."))
	// for i := waitTime; i >= 0; i-- {
	// 	countdown.Update(pterm.Gray("Waiting for " + secondsToMinutes(i) + " seconds..."))
	// 	time.Sleep(time.Second)
	// }
	// countdown.Stop()

	time.Sleep(time.Duration(waitTime) * time.Second)

	u.Busy = false
}

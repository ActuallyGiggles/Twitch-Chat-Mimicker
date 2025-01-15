package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

func Mimic(C chan Message) {
messageRange:
	for c := range C {
		channel := c.Channel
		message := c.Message

		if updatingEmotes {
			continue
		}

		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			if user.Name == channel {

				if user.Busy || !user.IsLive {
					continue messageRange
				}

				var thingToSend string

				// Low priority (not unique)
				emoteFound := ParseEmote(message)
				if emoteFound != "" {
					thingToSend = emoteFound
				}
				// Medium low priority (more unique than regular emote)
				if onlyWordCombo := ParseOnlyWordCombo(message); onlyWordCombo != "" {
					thingToSend = onlyWordCombo
				}
				//  Medium high priority (very unique)
				if emoteWordCombo := ParseEmoteWordCombo(emoteFound, message); emoteWordCombo != "" {
					thingToSend = emoteWordCombo
				}

				if (thingToSend == user.LastSentEmote && !Config.AllowConsecutiveDuplicates) || thingToSend == "" {
					continue messageRange
				}

				user.Responses[thingToSend]++

				if user.Responses[thingToSend] >= Config.MessageThreshold {
					go Respond(user, thingToSend)
					user.Messages = 0
					user.Responses = make(map[string]int)
					user.WordsUsed = nil
					continue messageRange
				}

				user.Messages++

				if user.Messages > Config.MessageSample {
					var maxValue int
					var maxName string
					for name, value := range user.Responses {
						if value > maxValue {
							maxValue = value
							maxName = name
						}
					}

					t := time.Now()

					if maxValue > 1 {
						pterm.Warning.Printf("%-20s     %s\n%s", strings.ToUpper(user.Name), maxName, pterm.Gray(fmt.Sprintf("%d:%02d | Times Used: %d/%d | Sample Size: %d", t.Hour(), t.Minute(), maxValue, Config.MessageThreshold, Config.MessageSample)))
						pterm.Println()
						pterm.Println()
					}

					user.Messages = 0
					user.Responses = make(map[string]int)
					user.WordsUsed = nil
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
			if strings.EqualFold(emoteWordCombo, strings.TrimSpace(allWordsAroundEmote[0])) {
				emoteWordComboToSend = append(emoteWordComboToSend, emoteWordCombo)
			}
		}
	}

	// Add emote
	emoteWordComboToSend = append(emoteWordComboToSend, emote)

	// Add words after emote
	if len(allWordsAroundEmote) > 1 {
		for _, emoteWordCombo := range Config.EmoteWordCombos {
			if strings.EqualFold(emoteWordCombo, strings.TrimSpace(allWordsAroundEmote[1])) {
				emoteWordComboToSend = append(emoteWordComboToSend, emoteWordCombo)
			}
		}
	}

	return strings.Join(emoteWordComboToSend, " ")
}

// Returns an only word combo to send
func ParseOnlyWordCombo(message string) string {
	for _, combo := range Config.OnlyWordCombos {
		if combo == strings.ToUpper(message) {
			return combo
		}

	}

	return ""
}

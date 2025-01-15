package main

import (
	"fmt"
	"strings"
)

func MimicNew(C chan Message) {
messageRange:
	for c := range C {
		channel := c.Channel
		message := c.Message

		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			if user.Name == channel {
				user.Messages++

				// If the twitch channel is not live or is currently waiting on a cooldown, ignore the message
				if user.Busy || !user.IsLive {
					continue messageRange
				}

				// ????
				for phrase := range user.Responses {
					if strings.Contains(message, phrase) {
						user.Responses[phrase]++
					}
				}

				// Break message into a slice
				messageFields := strings.FieldsFunc(message, func(r rune) bool {
					return r == ' ' || r == ',' || r == '.' || r == '!' || r == '?'
				})

				var phraseSeparated []string
			messageFieldsLoop:
				for _, word := range messageFields {
					var wordAlreadyAddedToList bool
					for _, wordUsed := range user.WordsUsed {
						if word == wordUsed {
							phraseSeparated = append(phraseSeparated, word)
							wordAlreadyAddedToList = true
							continue messageFieldsLoop
						}
					}
					if !wordAlreadyAddedToList {
						user.WordsUsed = append(user.WordsUsed, word)
					}
				}

				phraseToSend := strings.Join(phraseSeparated, " ")

				if (phraseToSend == user.LastSentPhrase && !Config.AllowConsecutiveDuplicates) || phraseToSend == "" {
					continue messageRange
				}

				user.Responses[phraseToSend]++

				if user.Responses[phraseToSend] >= Config.MessageThreshold {
					for word, v := range user.Responses {
						fmt.Println("Word:", word, "Value:", v)
					}

					fmt.Println("RESPONDING: ", phraseToSend)
					fmt.Println()
					//go Respond(user, phraseToSend)
					user.LastSentPhrase = phraseToSend
					user.Messages = 0
					user.Responses = make(map[string]int)
					user.WordsUsed = nil
				}

				if user.Messages > Config.MessageSample {
					user.Messages = 0
					user.Responses = make(map[string]int)
					user.WordsUsed = nil
				}
			}
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

func PrintSuccess(p Instructions) {
	pterm.Success.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray(p.Note))
	pterm.Println()
	pterm.Println()
}

func PrintWarning(p Instructions) {
	pterm.Warning.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray(p.Note))
	pterm.Println()
	pterm.Println()
}

// func PrintInfo(p Instructions) {
// 	pterm.Warning.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray(p.Note))
// 	pterm.Println()
// 	pterm.Println()
// }

func Page(title string, content func() bool) {
doAgain:
	print("\033[H\033[2J")
	if title == "Aborted" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
	} else if title == "Set Up" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
	} else if title == "Initialization" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
	} else if title == "Started" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println("Twitch Chat Mimicker " + title)

		pterm.Println()
		pterm.Info.Printf("%-20s     %s\n%s\n\n%s", "TotalEmotes:",
			strconv.Itoa(EmoteAmounts.TwitchGlobal+
				EmoteAmounts.TwitchChannel+
				EmoteAmounts.SevenTVGlobal+
				EmoteAmounts.SevenTVChannel+
				EmoteAmounts.BetterTTVGlobal+
				EmoteAmounts.BetterTTVChannel+
				EmoteAmounts.FFZGlobal+
				EmoteAmounts.FFZChannel),
			pterm.Gray(fmt.Sprintf("Twitch Global: %d\nTwitch Channel: %d\nSevenTV Global: %d\nSevenTV Channel: %d\nBetterTTV Global: %d\nBetterTTV Channel: %d\nFFZ Global: %d\nFFZ Channel: %d",
				EmoteAmounts.TwitchGlobal,
				EmoteAmounts.TwitchChannel,
				EmoteAmounts.SevenTVGlobal,
				EmoteAmounts.SevenTVChannel,
				EmoteAmounts.BetterTTVGlobal,
				EmoteAmounts.BetterTTVChannel,
				EmoteAmounts.FFZGlobal,
				EmoteAmounts.FFZChannel)),
			pterm.Gray(fmt.Sprintf("Emojis: %d", EmoteAmounts.Emojis)))

	}
	pterm.Println()
	pterm.Println()
	if !content() {
		pterm.Error.Println("Press Enter to try again.")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		goto doAgain
	}
}

func Clear() {
	print("\033[H\033[2J")
}

type Instructions struct {
	Channel string
	Emote   string

	Note     string
	NoteOnly bool

	Error bool
}

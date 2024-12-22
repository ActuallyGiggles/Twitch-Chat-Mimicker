package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

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
		pterm.Info.Printf("%-20s     %s\n%s\n%s", "Channels:", strings.Join(Config.Channels, ", "),
			pterm.Gray(fmt.Sprintf("TotalEmotes: %d\nTwitch Global: %d\nTwitch Channel: %d\nSevenTV Global: %d\nSevenTV Channel: %d\nBetterTTV Global: %d\nBetterTTV Channel: %d\nFFZ Global: %d\nFFZ Channel: %d",
				EmoteAmounts.TwitchGlobal+
					EmoteAmounts.TwitchChannel+
					EmoteAmounts.SevenTVGlobal+
					EmoteAmounts.SevenTVChannel+
					EmoteAmounts.BetterTTVGlobal+
					EmoteAmounts.BetterTTVChannel+
					EmoteAmounts.FFZGlobal+
					EmoteAmounts.FFZChannel,
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

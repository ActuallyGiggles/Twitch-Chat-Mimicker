package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func PrintSuccess(p Instructions) {
	pterm.Success.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray(p.Note))
	pterm.Println()
	pterm.Println()
}

func PrintFail(p Instructions) {
	pterm.Info.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray(p.Note))
	pterm.Println()
	pterm.Println()
}

func Page(title string, content func() bool) {
doAgain:
	print("\033[H\033[2J")
	if title == "Aborted" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
	} else if title == "Started" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
	} else if title == "Set Up" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Chat Mimicker " + title)
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

package main

import (
	"strings"

	"github.com/pterm/pterm"
)

func Print(p Instructions) {
	pterm.Success.Printf("%-20s     %s\n%s", strings.ToUpper(p.Channel), p.Emote, pterm.Gray("Note: "+p.Note))
	pterm.Println()
	pterm.Println()
}

func Page(title string, content func()) {
	print("\033[H\033[2J")
	if title == "Exited" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println("Twitch Chat Mimicker by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println(title)
	} else if title == "Started" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println("Twitch Chat Mimicker by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println(title)
	} else if title == "Set Up" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Chat Mimicker by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println(title)
	}
	pterm.Println()
	pterm.Println()
	content()
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

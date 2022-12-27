package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

func readConfig() {
	f, err := os.Open("./config.json")
	defer f.Close()
	if err != nil {
		configSetup()
		return
	}

	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		panic(err)
	}

	for _, channel := range Config.Channels {
		addUser(channel)
	}
}

func writeConfig() {
	f, err := os.OpenFile("./config.json", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}

	readConfig()
}

func configSetup() {
	// Name
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Enter the Twitch account name you will be using.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account Name: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := strings.ToLower(scanner.Text())
		Config.Name = name
	})

	// Client ID
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Obtaining your ClientID is necessary to gather Twitch emotes.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://dev.twitch.tv/console\n"), "Steps:\n1. Give your application a name.\n2. Set the redirect URL to (https://twitchapps.com/tokengen/).\n3. Choose the chatbot category.\n4. Copy and paste the Client ID here.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account Client ID: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		clientID := strings.ToLower(scanner.Text())
		Config.ClientID = clientID
	})

	// OAuth
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Obtaining your OAuth is necessary to connect to Twitch chatrooms as yourself.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://twitchapps.com/tokengen/\n"), "\n\nSteps:\n1. Paste in the Client ID\n2. For scopes, type in: 'chat:read chat:edit'.\n3. Click connect and copy and paste the OAuth Token here.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account OAuth: "), pterm.White("oauth:"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		oauth := strings.ToLower(scanner.Text())
		Config.OAuth = oauth
	})

	// Channels
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the channels in which the program should act in.\nSeparate channel names with spaces.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Channels To Join: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		channels := strings.Split(strings.ToLower(scanner.Text()), " ")
		for _, channel := range channels {
			Config.Channels = append(Config.Channels, channel)
		}
	})

	// Blacklisted Emotes
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the emotes that you want blacklisted.\nSeparate them with spaces.\nExample: 'TriHard KEKW ResidentSleeper'"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Blacklisted Emotes: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		blacklist := scanner.Text()
		for _, blacked := range strings.Split(blacklist, " ") {
			Config.BlacklistEmotes = append(Config.BlacklistEmotes, blacked)
		}
	})

	// Message Sample
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the amount of messages to sample?\nRecommended: 10\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Sample: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		sample := scanner.Text()
		s, err := strconv.Atoi(sample)
		if err != nil {
			pterm.Error.Println(sample, "is not a number!")
			os.Exit(3)
		}
		Config.MessageSample = s
	})

	// Message Threshold
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Out of that sample size, specify the amount of times that an emote has to repeat itself to send it?\nRecommended: 3\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Threshold: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		threshold := scanner.Text()
		t, err := strconv.Atoi(threshold)
		if err != nil {
			pterm.Error.Println(threshold, "is not a number!")
			os.Exit(3)
		}
		Config.MessageThreshold = t
	})

	// Messaging Interval
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the range of minutes for the bot to wait in between message sends.\nRecommended: '1 5'\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Range: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		interval := scanner.Text()
		tS := strings.Split(interval, " ")
		min, err := strconv.Atoi(tS[0])
		if err != nil {
			pterm.Error.Println(tS[0], "is not a number!")
			os.Exit(3)
		}
		max, err := strconv.Atoi(tS[1])
		if err != nil {
			pterm.Error.Println(tS[1], "is not a number!")
			os.Exit(3)
		}
		if min < 0 || max < 0 {
			pterm.Error.Println("Cannot be less than zero!")
			os.Exit(3)
		}
		Config.IntervalMin = min
		Config.IntervalMax = max
	})

	// Consecutive Duplicates
	Page("Set Up", func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify if you allow consecutive duplicate emotes to be sent?\nExample: sending OMEGALUL two times in a row because the chat won't stop spamming OMEGALUL\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--True(t)/False(f): "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		decision, err := strconv.ParseBool(scanner.Text())
		if err != nil {
			pterm.Error.Println(`Please answer "t" for true or "f" for false.`)
			os.Exit(3)
		}
		Config.AllowConsecutiveDuplicates = decision
	})

	writeConfig()
	clearTerminal()
}

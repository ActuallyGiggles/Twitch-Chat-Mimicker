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
	if err != nil {
		configSetup()
		f.Close()
		return
	}

	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		panic(err)
	}

	f.Close()
}

func writeConfig() {
	f, err := os.OpenFile("./config.json", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	c, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}

	f.Close()

	readConfig()
}

func configSetup() {
	ObtainName()
	ObtainClientID()
	ObtainOAuth()

	// Blacklisted Emotes
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the emotes that you want blacklisted.\n\nExample: 'TriHard, KEKW, ResidentSleeper'"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Blacklisted Emotes: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		blacklist := strings.Split(strings.ToLower(scanner.Text()), ", ")
		if strings.Trim(blacklist[0], " ") == "" {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("No emotes entered!")
			return false
		}
		Config.BlacklistEmotes = append(Config.BlacklistEmotes, blacklist...)
		return true
	})

	// Custom Emote Word Combos
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Are there any words or letters that pair with emotes? Separate them with commas.\nExample: 'L OMEGALUL L, W H OMEGALUL, CLM'"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Emote Word Combos: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		EmoteWordCombos := strings.Split(scanner.Text(), ", ")
		Config.EmoteWordCombos = append(Config.EmoteWordCombos, EmoteWordCombos...)
		return true
	})

	// Custom Only Word Combos
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Are there any words or letters that are used by themselves? Separate them with commas.\nExample: 'F, +1, -1, L, W, CLM'"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Only Word Combos: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		OnlyWordCombos := strings.Split(scanner.Text(), ", ")
		Config.OnlyWordCombos = append(Config.OnlyWordCombos, OnlyWordCombos...)
		return true
	})

	// Message Sample
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the amount of messages to sample?\nRecommended: 20\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Sample: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		sample := scanner.Text()
		s, err := strconv.Atoi(sample)
		if err != nil {
			pterm.Error.Println(sample, "is not a number!")
			return false
		}
		Config.MessageSample = s
		return true
	})

	// Message Threshold
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Out of that sample size, specify the amount of times that an emote has to repeat itself to send it?\nRecommended: 5\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Threshold: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		threshold := scanner.Text()
		t, err := strconv.Atoi(threshold)
		if err != nil {
			pterm.Error.Println(threshold, "is not a number!")
			return false
		}
		Config.MessageThreshold = t
		return true
	})

	// Messaging Interval
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the range of seconds for the bot to wait in between message sends.\n\nSeparate the minimum and maximum with spaces.\nRecommended: '60, 300'\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Range: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		interval := scanner.Text()
		tS := strings.Split(interval, ", ")
		min, err := strconv.Atoi(tS[0])
		if err != nil {
			pterm.Error.Println(tS[0], "is not a number!")
			return false
		}
		max, err := strconv.Atoi(tS[1])
		if err != nil {
			pterm.Error.Println(tS[1], "is not a number!")
			return false
		}
		if min < 0 || max < 0 {
			pterm.Error.Println("Cannot be less than zero!")
			return false
		}
		Config.IntervalMin = min
		Config.IntervalMax = max
		return true
	})

	// Consecutive Duplicates
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify if you allow consecutive duplicate emotes to be sent?\nExample: sending OMEGALUL two times in a row because the chat won't stop spamming OMEGALUL\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--True(t)/False(f): "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		decision, err := strconv.ParseBool(scanner.Text())
		if err != nil {
			pterm.Error.Println(`Please answer "t" for true or "f" for false.`)
			return false
		}
		Config.AllowConsecutiveDuplicates = decision
		return true
	})

	writeConfig()
	clearTerminal()
}

func configInit() {
	// Channels to join
	Page("Initialization", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the channels in which the program should act in.\nExample: 'nmplol, xqc, sodapoppin'\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Channels To Join: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		channels := strings.Split(strings.ToLower(scanner.Text()), ", ")
		if strings.Trim(channels[0], " ") == "" {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("No channels entered!")
			return false
		}
		Config.Channels = channels

		for _, channel := range Config.Channels {
			u := User{
				Name: channel,
			}
			Users = append(Users, u)
		}

		return true
	})

	writeConfig()
	clearTerminal()

	Config.EmoteWordCombos = append(Config.EmoteWordCombos, Config.OnlyWordCombos...)

	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		user.Responses = make(map[string]int)
	}
}

func VerifyCredentials() {
	user := User{
		Name: Config.Name,
	}
	success := GetBroadcaster(&user)

	if success {
		return
	}

	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Looks like either your ClientID and/or OAuth is invalid or expired. Let's do that again!"))
		countdown(5)
		return true
	})

	ObtainName()
	ObtainClientID()
	ObtainOAuth()

	writeConfig()
	clearTerminal()
	VerifyCredentials()
}

func ObtainName() {
	// Name
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Enter the Twitch account name you will be using.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account Name: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := strings.ToLower(scanner.Text())
		Config.Name = name
		return true
	})
}

func ObtainOAuth() {
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Obtaining your OAuth is necessary to connect to Twitch chatrooms as yourself.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://twitchapps.com/tokengen/\n"), "\n\nSteps:\n1. Paste in the Client ID\n2. For scopes, type in: 'chat:read chat:edit'.\n3. Click connect and copy and paste the OAuth Token here.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account OAuth: "), pterm.White("oauth:"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		oauth := strings.ToLower(scanner.Text())
		Config.OAuth = oauth
		return true
	})
}

func ObtainClientID() {
	Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Obtaining your ClientID is necessary to gather Twitch emotes.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://dev.twitch.tv/console\n"), "\n\nSteps:\n1. Give your application a name.\n2. Set the redirect URL to (https://twitchapps.com/tokengen/).\n3. Choose the chatbot category.\n4. Copy and paste the Client ID here.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account Client ID: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		clientID := strings.ToLower(scanner.Text())
		Config.ClientID = clientID
		return true
	})
}

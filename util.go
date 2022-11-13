package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	Config = Configuration{}
)

func addUser(user string) {
	u := User{
		Name: user,
	}
	Users = append(Users, u)
}

func RandomNumber(min int, max int) (num int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	num = r.Intn(max-min) + min
	return num
}

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
	if err == nil {
		fmt.Println("Config successfully created.")
	}

	readConfig()
}

func configSetup() {
	// Intro
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("First time? Let's setup your bot.")
	fmt.Println()
	fmt.Println("Press Enter...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Name
	fmt.Println("First, what is the login name of the account you are going to use?")
	fmt.Println()
	fmt.Print("Name: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := strings.ToLower(scanner.Text())
	Config.Name = name
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// OAuth
	fmt.Println("Let's generate an OAuth. Go to this website (https://twitchapps.com/tmi/), and paste the result here.")
	fmt.Println()
	fmt.Print("OAuth: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	oauth := scanner.Text()
	Config.OAuth = oauth
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Client ID
	fmt.Println("Now let's get your Client ID, it can be found here: (https://dev.twitch.tv/console). \n\nSteps:\n1. Give your application a name.\n2. Set the redirect URL to (https://localhost).\n3. Choose the chatbot category.\n4. Copy and paste the Client ID here.")
	fmt.Println()
	fmt.Print("Client ID: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clientID := scanner.Text()
	Config.ClientID = clientID
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Access Token
	fmt.Println("It's time to get your Access Token, it can be found here (https://twitchtokengenerator.com/). \n\nSteps:\n1. Select 'Bot Chat Token'.\n2. Click 'Authorize'.\n3. Copy and paste the Access Token here.")
	fmt.Println()
	fmt.Print("Access Token: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	accessToken := scanner.Text()
	Config.AccessToken = accessToken
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Channels
	fmt.Println(`Now list the channels that you want your bot to be active in. (Separate them with spaces.) (Example: "39daph nmplol sodapoppin veibae")`)
	fmt.Println()
	fmt.Print("Channels: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	channels := scanner.Text()
	for _, channel := range strings.Split(channels, " ") {
		Config.Channels = append(Config.Channels, channel)
	}
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Blacklisted Emotes
	fmt.Println(`Enter the emotes that you want to be blacklisted. (Separate them with spaces.) (Example: "TriHard KEKW ResidentSleeper")`)
	fmt.Println()
	fmt.Print("Blacklist emotes: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	blacklist := scanner.Text()
	for _, blacked := range strings.Split(blacklist, " ") {
		Config.BlacklistEmotes = append(Config.BlacklistEmotes, blacked)
	}
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Messaging Interval
	fmt.Println(`Finally, please specify the range of minutes for the bot to wait in between message sends. (Example: "5 10")`)
	fmt.Println()
	fmt.Print("Range: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	interval := scanner.Text()
	tS := strings.Split(interval, " ")

	min, err := strconv.Atoi(tS[0])
	if err != nil {
		fmt.Println(tS[0], "is not a number!")
		os.Exit(3)
	}

	max, err := strconv.Atoi(tS[1])
	if err != nil {
		fmt.Println(tS[1], "is not a number!")
		os.Exit(3)
	}

	Config.IntervalMin = min
	Config.IntervalMax = max
	fmt.Println()

	writeConfig()
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	updatingEmotes bool
	tempEmotes     []string
)

func getEmotes() {
	fmt.Println()
	fmt.Println("Gathering emotes...")
	updatingEmotes = true

	// Get broadcaster IDs
	GetBroadcasters()

	// Get global emotes
	fmt.Println("[Global]:")
	TGlobal := getTwitchGlobalEmotes()
	fmt.Println("\tTwitch:", TGlobal)
	SevenGlobal := get7tvGlobalEmotes()
	fmt.Println("\t7tv:", SevenGlobal)
	BGlobal := getBttvGlobalEmotes()
	fmt.Println("\tBTTV:", BGlobal)
	FGlobal := getFfzGlobalEmotes()
	fmt.Println("\tFFZ:", FGlobal)

	globalsTotal := TGlobal + SevenGlobal + BGlobal + FGlobal
	fmt.Println("\tTotal:", globalsTotal)

	// Get channel emotes
	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		fmt.Println("[" + user.Name + "]")

		TChannel := getTwitchChannelEmotes(user)
		fmt.Println("\tTwitch:", len(TChannel))
		ChannelEmotes = append(ChannelEmotes, TChannel...)

		SChannel := get7tvChannelEmotes(user)
		fmt.Println("\t7tv:", len(SChannel))
		ChannelEmotes = append(ChannelEmotes, SChannel...)

		BChannel := getBttvChannelEmotes(user)
		fmt.Println("\tBTTV:", len(BChannel))
		ChannelEmotes = append(ChannelEmotes, BChannel...)

		FChannel := getFfzChannelEmotes(user)
		fmt.Println("\tFFZ:", len(FChannel))
		ChannelEmotes = append(ChannelEmotes, FChannel...)

		fmt.Printf("\tTotal: %d\n", len(TChannel)+len(SChannel)+len(BChannel)+len(SChannel))
	}

	updatingEmotes = false

	go updateEmotes()
}

func updateEmotes() {
	for range time.Tick(1 * time.Hour) {
		fmt.Println("Updating emotes...")
		updatingEmotes = true
		ChannelEmotes = nil

		// Get channel emotes
		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			fmt.Println("[" + user.Name + "]")

			TChannel := getTwitchChannelEmotes(user)
			fmt.Println("\tTwitch:", len(TChannel))
			ChannelEmotes = append(ChannelEmotes, TChannel...)

			SChannel := get7tvChannelEmotes(user)
			fmt.Println("\tTwitch:", len(SChannel))
			ChannelEmotes = append(ChannelEmotes, SChannel...)

			BChannel := getBttvChannelEmotes(user)
			fmt.Println("\tTwitch:", len(BChannel))
			ChannelEmotes = append(ChannelEmotes, BChannel...)

			FChannel := getFfzChannelEmotes(user)
			fmt.Println("\tTwitch:", len(FChannel))
			ChannelEmotes = append(ChannelEmotes, FChannel...)

			fmt.Printf("Total: %d\n", len(TChannel)+len(SChannel)+len(BChannel)+len(SChannel))
		}

		fmt.Println("Emotes updated.")
		updatingEmotes = false
	}
}

func GetBroadcasters() {
	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		url := "https://api.twitch.tv/helix/users?login=" + user.Name

		var jsonStr = []byte(`{"":""}`)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+Config.OAuth)
		req.Header.Set("Client-Id", Config.ClientID)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			log.Println("GetBroadcasterID failed\n", err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("GetBroadcasterID failed\n", err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			log.Printf("GetBroadcaters(%s) is not OK\n\t%s", user.Name, string(body))
		}
		broadcaster := Broadcaster[Data]{}
		if err := json.Unmarshal(body, &broadcaster); err != nil {
			log.Println("GetBroadcasterID failed\n", err.Error())
		}
		for _, v := range broadcaster.Data {
			user.ID = v.ID
		}
	}
}

func getTwitchGlobalEmotes() (number int) {
	url := "https://api.twitch.tv/helix/chat/emotes/global"

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+Config.OAuth)
	req.Header.Set("Client-Id", Config.ClientID)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := TwitchEmoteAPIResponse[TwitchGlobalEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes.Data {
		GlobalEmotes = append(GlobalEmotes, emote.Name)
		number++
	}

	return number
}

func getTwitchChannelEmotes(user *User) (c []string) {
	url := "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=" + user.ID

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+Config.OAuth)
	req.Header.Set("Client-Id", Config.ClientID)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Printf("\t getTwitchChannelEmotes failed\n")
		log.Printf("\t For channel %s\n1", user.Name)
		log.Println(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("getTwitchChannelEmotes failed\n", err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := TwitchEmoteAPIResponse[TwitchChannelEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		log.Println("getTwitchChannelEmotes failed\n", err.Error())
	}

	for _, emote := range emotes.Data {
		c = append(c, emote.Name)
	}

	return c
}

func get7tvGlobalEmotes() (number int) {
	url := "https://api.7tv.app/v2/emotes/global"

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := []SevenTVEmote{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		log.Fatal(err)
	}

	for _, emote := range emotes {
		GlobalEmotes = append(GlobalEmotes, emote.Name)
		number++
	}

	return number
}

func get7tvChannelEmotes(user *User) (c []string) {
	url := "https://api.7tv.app/v2/users/" + user.Name + "/emotes"

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return c
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := []SevenTVEmote{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes {
		c = append(c, emote.Name)
	}

	return c
}

func getBttvGlobalEmotes() (number int) {
	url := "https://api.betterttv.net/3/cached/emotes/global"

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := []BttvEmote{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes {
		GlobalEmotes = append(GlobalEmotes, emote.Name)
		number++
	}
	return number
}

func getBttvChannelEmotes(user *User) (c []string) {
	url := "https://api.betterttv.net/3/cached/users/twitch/" + user.ID

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return c
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	emotes := BttvChannelEmotes[BttvEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes.ChannelEmotes {
		c = append(c, emote.Name)
	}
	for _, emote := range emotes.SharedEmotes {
		c = append(c, emote.Name)
	}

	return c
}

func getFfzGlobalEmotes() (number int) {
	url := "https://api.frankerfacez.com/v1/set/global"

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	set := FfzSets{}
	if err := json.Unmarshal(body, &set); err != nil {
		panic(err)
	}

	for _, emotes := range set.Sets {
		for _, emote := range emotes.Emoticons {
			GlobalEmotes = append(GlobalEmotes, emote.Name)
			number++
		}
	}

	return number
}

func getFfzChannelEmotes(user *User) (c []string) {
	url := "https://api.frankerfacez.com/v1/room/id/" + user.ID

	var jsonStr = []byte(`{"":""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return c
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	set := FfzSets{}
	if err := json.Unmarshal(body, &set); err != nil {
		panic(err)
	}

	for _, emotes := range set.Sets {
		for _, emote := range emotes.Emoticons {
			c = append(c, emote.Name)
		}
	}

	return c
}

func getLiveStatuses() {
	getLiveStatus()
	for range time.Tick(30 * time.Second) {
		getLiveStatus()
	}
}

func getLiveStatus() {
	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		url := "https://api.twitch.tv/helix/streams?user_login=" + user.Name
		var jsonStr = []byte(`{"":""}`)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+Config.OAuth)
		req.Header.Set("Client-Id", Config.ClientID)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			log.Println(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("LIVE BODY RESPONSE:", string(body))
		var stream StreamStatusData
		if err := json.Unmarshal(body, &stream); err != nil {
			log.Println(err.Error())
		}
		if len(stream.Data) == 0 {
			user.IsLive = false
		} else {
			user.IsLive = true
		}
	}
}

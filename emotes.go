package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

var (
	updatingEmotes bool
	tempEmotes     []string

	pb *pterm.ProgressbarPrinter
)

func getEmotes(isInit bool) {
	updatingEmotes = true

	if isInit {
		pb, _ = pterm.DefaultProgressbar.WithTotal(4 + len(Config.Channels)*5).WithTitle("Gathering Twitch API information...").WithRemoveWhenDone(true).Start()
		// Get broadcaster IDs
		pb.UpdateTitle("Gathering broadcaster information...")
		for i := 0; i < len(Users); i++ {
			user := &Users[i]
			GetBroadcasters(user)
			pb.Increment()
		}

		// Get global emotes
		getTwitchGlobalEmotes()
		pb.Increment()
		get7tvGlobalEmotes()
		pb.Increment()
		getBttvGlobalEmotes()
		pb.Increment()
		getFfzGlobalEmotes()
		pb.Increment()
	}

	// Get channel emotes
	if isInit {
		pb.UpdateTitle("Gathering channel emotes...")
	}

	ChannelEmotes = nil

	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		TChannel := getTwitchChannelEmotes(user)
		ChannelEmotes = append(ChannelEmotes, TChannel...)
		if isInit {
			pb.Increment()
		}
		SChannel := get7tvChannelEmotes(user)
		ChannelEmotes = append(ChannelEmotes, SChannel...)
		if isInit {
			pb.Increment()
		}
		BChannel := getBttvChannelEmotes(user)
		ChannelEmotes = append(ChannelEmotes, BChannel...)
		if isInit {
			pb.Increment()
		}
		FChannel := getFfzChannelEmotes(user)
		ChannelEmotes = append(ChannelEmotes, FChannel...)
		if isInit {
			pb.Increment()
		}
	}

	updatingEmotes = false
}

func updateEmotes() {
	for range time.Tick(1 * time.Hour) {
		getEmotes(false)
	}
}

func GetBroadcasters(user *User) {
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

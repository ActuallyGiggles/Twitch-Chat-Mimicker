package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

var (
	EmoteAmounts   EmoteAmountsStruct
	updatingEmotes bool

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
			GetBroadcaster(user)
			pb.Increment()
		}

		EmoteAmounts.Emojis += getEmojis()

		// Get global emotes
		EmoteAmounts.TwitchGlobal += getTwitchGlobalEmotes()
		pb.Increment()
		EmoteAmounts.SevenTVGlobal += get7tvGlobalEmotes()
		pb.Increment()
		EmoteAmounts.BetterTTVGlobal += getBttvGlobalEmotes()
		pb.Increment()
		EmoteAmounts.FFZGlobal += getFfzGlobalEmotes()
		pb.Increment()
	}

	// Get channel emotes
	if isInit {
		pb.UpdateTitle("Gathering channel emotes...")
	}

	for i := 0; i < len(Users); i++ {
		user := &Users[i]
		EmoteAmounts.TwitchChannel += getTwitchChannelEmotes(user)
		if isInit {
			pb.Increment()
		}

		EmoteAmounts.SevenTVChannel += get7tvChannelEmotes(user)
		if isInit {
			pb.Increment()
		}
		EmoteAmounts.BetterTTVChannel += getBttvChannelEmotes(user)
		if isInit {
			pb.Increment()
		}

		EmoteAmounts.FFZChannel += getFfzChannelEmotes(user)
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

func GetBroadcaster(user *User) {
	url := "https://api.twitch.tv/helix/users?login=" + user.Name

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	req.Header.Set("Authorization", "Bearer "+Config.OAuth)
	req.Header.Set("Client-Id", Config.ClientID)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Println("GetBroadcaster failed\n", err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("GetBroadcaster failed\n", err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		println()
		log.Printf("GetBroadcater(%s) is not OK\n%s\n", user.Name, string(body))
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

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	req.Header.Set("Authorization", "Bearer "+Config.OAuth)
	req.Header.Set("Client-Id", Config.ClientID)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	emotes := TwitchEmoteAPIResponse[TwitchGlobalEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes.Data {
		Emotes = append(Emotes, emote.Name)
		number++
	}

	return number
}

func getTwitchChannelEmotes(user *User) (number int) {
	url := "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=" + user.ID

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
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
	body, _ := io.ReadAll(resp.Body)
	emotes := TwitchEmoteAPIResponse[TwitchChannelEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		log.Println("getTwitchChannelEmotes failed\n", err.Error())
	}

	for _, emote := range emotes.Data {
		Emotes = append(Emotes, emote.Name)
		number++
	}

	return number
}

func get7tvGlobalEmotes() (number int) {
	url := "https://7tv.io/v3/emote-sets/global"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	emoteSet := SevenTVGlobalEmoteSet{}
	if err := json.Unmarshal(body, &emoteSet); err != nil {
		log.Fatal(err)
	}

	for _, emote := range emoteSet.Emotes {
		Emotes = append(Emotes, emote.Name)
		number++
	}

	return number
}

func get7tvChannelEmotes(user *User) (number int) {
	url := "https://7tv.io/v3/users/twitch/" + user.ID

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	channel := SevenTVChannel{}
	if err := json.Unmarshal(body, &channel); err != nil {
		panic(err)
	}

	for _, emote := range channel.EmoteSet.Emotes {
		Emotes = append(Emotes, emote.Name)
		number++
	}

	return number
}

func getBttvGlobalEmotes() (number int) {
	url := "https://api.betterttv.net/3/cached/emotes/global"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	emotes := []BttvEmote{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes {
		Emotes = append(Emotes, emote.Name)
		number++
	}
	return number
}

func getBttvChannelEmotes(user *User) (number int) {
	url := "https://api.betterttv.net/3/cached/users/twitch/" + user.ID

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	emotes := BttvChannelEmotes[BttvEmote]{}
	if err := json.Unmarshal(body, &emotes); err != nil {
		panic(err)
	}

	for _, emote := range emotes.ChannelEmotes {
		Emotes = append(Emotes, emote.Name)
		number++
	}
	for _, emote := range emotes.SharedEmotes {
		Emotes = append(Emotes, emote.Name)
		number++
	}

	return number
}

func getFfzGlobalEmotes() (number int) {
	url := "https://api.frankerfacez.com/v1/set/global"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	set := FfzSets{}
	if err := json.Unmarshal(body, &set); err != nil {
		panic(err)
	}

	for _, emotes := range set.Sets {
		for _, emote := range emotes.Emoticons {
			Emotes = append(Emotes, emote.Name)
			number++
		}
	}

	return number
}

func getFfzChannelEmotes(user *User) (number int) {
	url := "https://api.frankerfacez.com/v1/room/id/" + user.ID

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	set := FfzSets{}
	if err := json.Unmarshal(body, &set); err != nil {
		panic(err)
	}

	for _, emotes := range set.Sets {
		for _, emote := range emotes.Emoticons {
			Emotes = append(Emotes, emote.Name)
			number++
		}
	}

	return number
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

		req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
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
		body, _ := io.ReadAll(resp.Body)
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

func getEmojis() (number int) {
	url := "https://raw.githubusercontent.com/chalda-pnuzig/emojis.json/refs/heads/master/src/list.json"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	allEmojis := EmojisStruct{}
	if err := json.Unmarshal(body, &allEmojis); err != nil {
		panic(err)
	}

	for _, emoji := range allEmojis.Emojis {
		Emotes = append(Emotes, emoji.Name)
		number++
	}

	return number
}

package main

import "time"

type Configuration struct {
	Name     string
	OAuth    string
	ClientID string

	IntervalMin int
	IntervalMax int

	MessageThreshold int
	MessageSample    int

	Channels        []string
	BlacklistEmotes []string
	EmoteWordCombos []string
	OnlyWordCombos  []string

	AllowConsecutiveDuplicates bool
}

type User struct {
	Name          string
	ID            string
	IsLive        bool
	Busy          bool
	Messages      int
	LastSentEmote string

	WordsUsed      []string
	Responses      map[string]int
	LastSentPhrase string
}

type Response struct {
	Name  string
	Value int
}

type Broadcaster[T any] struct {
	Data []T
}

type Data struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

type TwitchEmoteAPIResponse[T any] struct {
	Data     []T    `json:"data"`
	Template string `json:"template"`
}

type EmoteAmountsStruct struct {
	TwitchGlobal  int
	TwitchChannel int

	SevenTVGlobal  int
	SevenTVChannel int

	BetterTTVGlobal  int
	BetterTTVChannel int

	FFZGlobal  int
	FFZChannel int

	Emojis int
}

type TwitchGlobalEmote struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Images    map[string]string `json:"images"`
	Format    []string          `json:"format"`
	Scale     []string          `json:"scale"`
	ThemeMode []string          `json:"theme_mode"`
}

type TwitchChannelEmote struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Images     map[string]string `json:"images"`
	Tier       string            `json:"tier"`
	EmoteType  string            `json:"emote_type"`
	EmoteSetID string            `json:"emote_set_id"`
	Format     []string          `json:"format"`
	Scale      []string          `json:"scale"`
	ThemeMode  []string          `json:"theme_mode"`
}

type SevenTVGlobalEmoteSet struct {
	Emotes []SevenTVEmote `json:"emotes"`
}

type SevenTVChannel struct {
	EmoteSet SevenTVChannelEmoteSet `json:"emote_set"`
}

type SevenTVChannelEmoteSet struct {
	Emotes []SevenTVEmote `json:"emotes"`
}

type SevenTVEmote struct {
	Name string `json:"name"`
}

type BttvChannelEmotes[T any] struct {
	ChannelEmotes []T `json:"channelEmotes"`
	SharedEmotes  []T `json:"sharedEmotes"`
}

type BttvEmote struct {
	Name string `json:"code"`
	ID   string `json:"id"`
}

type FfzSets struct {
	Sets map[string]FfzSet `json:"sets"`
}

type FfzSet struct {
	Emoticons []FfzEmotes `json:"emoticons"`
}

type FfzEmotes struct {
	Name string            `json:"name"`
	Urls map[string]string `json:"urls"`
}

type EmojisStruct struct {
	Emojis []Emoji `json:"emojis"`
}

type Emoji struct {
	Name string `json:"emoji"`
}

type StreamStatusData struct {
	Data []StreamStatusActual `json:"data"`
}

type StreamStatusActual struct {
	Name string `json:"user_login"`
	Type string `json:"type"`
}

// Twitch message struct
type Message struct {
	Channel string
	Message string
}

type Instructions struct {
	Channel string
	Emote   string

	Note      string
	Delay     int
	Cooldown  int
	TimeStamp time.Time

	Error bool
}

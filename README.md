# Twitch Chat Mimicker
Twitch Chat Mimicker will automatically send emotes in your stead if it detects that a Twitch chatroom sends a certain emote multiple times.


## Instructions

Download [`Twitch Chat Mimicker.exe`](https://github.com/ActuallyGiggles/Twitch-Chat-Mimicker/releases/download/1.0/Twitch.Chat.Mimicker.exe) launch it. It will run you through the first time setup and create a config file in the folder that you run it in.

## Additional Information

1. You can specify which channels to be active in.
2. You can specify which emotes should be blacklisted.
3. You can specify how much of a sample size to use and how much of that sample size should contain the same emote to trigger sending.
4. You can specify how long to wait before detecting emotes again (to seem more human).
5. You can specify if the same emote is allowed to be sent multiple times in a row.

Responses are randomly delayed between 2-5 seconds of finding an emote to mimic.
If you would like to adjust any of the setup settings, edit the config JSON file.
Emotes will only be sent if the channel is live.

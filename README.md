# Twitch Chat Mimicker
Twitch Chat Mimicker will automatically send emotes in your stead if it detects that a Twitch chatroom sends a certain emote multiple times.


## Instructions

Download [`twitch-chat-mimicker.exe`](https://github.com/ActuallyGiggles/Twitch-Chat-Mimicker/releases/tag/1.0.0). Now you can launch `twitch-chat-mimicker.exe` and the program will run you through the first time setup. It will create a config.JSON file that you can edit.

## Additional Information

1. You can specify which channels to answer in.
2. You can specify which emotes should be blacklisted.
3. You can specify how much of a sample size to use and how much of that sample size should contain the same emote to trigger sending.
4. You can specify how long to wait before detecting emotes again (to seem more human).
5. You can specify if the same emote is allowed to be sent multiple times in a row.

Responses are randomly delayed between 2-10 seconds of finding an emote to mimic.
If you would like to add channels to monitor, edit the config JSON file.
Emotes will only be sent if the channel is live.

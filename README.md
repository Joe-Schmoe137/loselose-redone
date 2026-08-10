LoseLose-Redone is a personal project I created to learn Go. This is a TUI recreation of the [(in)famous game](https://en.wikipedia.org/wiki/Lose/Lose) made by Zach Gage for MacOS, although mine is for Linux. In this game you take the place of a space pilot that is presented with the dilemma of incoming enemies. Although if you decide to shoot one of these enemies rather than let them pass by then a random file on your computer will be deleted. As such **RUN AT YOUR OWN RISK**.

Additionally, one feature that I thought would be cool to add since I wanted to learn it was the ability for remote process injection (essentially you force someone to play it within their own terminal by taking over their process). This made my work take a little bit longer but I learned quite a bit from it! :)
## Usage
**If you desire to simply play the game**: 
1. Run the command: `go build ./cmd/loselose-game/`
2. Now you can run the executable `loselose-game`.
Do note that the files you have access to are tied to your current user, so for the true experience why not run it at root. >:)

**If you desire to utilize process injection:**
1. Run the setup script `bash setup.sh`, it will walk you through the process and afterwards you'll have the file `loselose-injection`
2. Using sudo privileges run one of the following :
	1. `sudo ./loselose-injection --pid [PID of the target process]`
	2. `sudo ./loselose-injection --targetUser [Username of the user you want to target]` This one will find the first available shell process to inject into that belongs to the target user. 
## Additional information

When I had originally started this project Zach Gage had hosted a website to view his project simply loselose\[.\]net. However, during the development of this personal project of mine some gremlin has squatted on the domain after it expired. So please don't view that site. Rather to view the original project, it has been open sourced and can be found [here](https://github.com/stfj/LoseLose).

Originally, I wanted the project to work on Windows, Mac, and Linux but as feature creep took place and my desire to work on the project waned, I've decided to leave it here. One day I might come back and flesh out both the game and injection aspect of this a bit more.

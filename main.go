package main

import (
	"os"
	"flag"
	"fmt"

    "github.com/muesli/mango"
    "github.com/muesli/mango/mflag"
    "github.com/muesli/roff"
)

var (
	//help = flag.String("help", "", "If specified it will run this help page then exit.")
	targetUser = flag.String("targetUser", "", "The user to target, if not specififed the running user will be targeted.")
	terminalTarget = flag.String("terminalTarget", "", "The pts/tty to target, if not specified the first one found will be used.")
)
func testFile() {
	os.Create("./success")
}

func main() {
	flag.Parse()

	manPage := mango.NewManPage(1, "loselose-redone", "A TUI based game that uses your files as collateral").
		WithLongDescription("loselose-redone is a reimagining of the original loselose game designed Zach Gage with the goal of being multiplatform and TUI based.\n" +
			"NOTE: This game WILL DESTORY FILES ON YOUR SYSTEM and can be classified as MALWARE. You have been warened\n" +
			"To find out more about the orignial game you can go to the loselose site here: https://loselose.net/")

	flag.VisitAll(mflag.FlagVisitor(manPage))
	
	// Display man page if --help is specified
	var exitBool = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "help" {
			fmt.Println(manPage.Build(roff.NewDocument()))
			exitBool = true
		}
	})
	if (exitBool) {return}

	testFile()
	setup()
}
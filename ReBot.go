package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"log"

	"strconv"

	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	RunAuto()
}

var returnedbot chan string
var dg *discordgo.Session

var owner string
var OC *discordgo.Channel
var OCID string
var botsfolder string

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	OC, err := dg.UserChannelCreate(owner)
	if err != nil {
		log.Fatal(err)
		return
	}
	OCID = OC.ID
	dg.ChannelMessageSend(OC.ID, "`Started ReBot!`")
}

func CM(Ses *discordgo.Session, MesC *discordgo.MessageCreate) {
	Mes := MesC.Message
	if Mes.Content != "" {
		if Mes.Content[0] == '!' && Mes.Author.ID == owner && len(Mes.Content) > 1 {
			var CMD = Mes.Content[1:4]
			switch strings.ToLower(CMD) {
			case "run":
				Runvals := strings.Split(Mes.Content[1:], " ")
				fmt.Println(Runvals)
				if len(Runvals) == 1 {
					dg.ChannelMessageSend(OCID, "`Nil val`")
				}
				if len(Runvals) == 2 {
					go Launch(Runvals[1], "")
					dg.ChannelMessageSend(OCID, "`Bot "+Runvals[1]+" launched `")
				}
				if len(Runvals) == 3 {
					go Launch(Runvals[1], Runvals[2])
					dg.ChannelMessageSend(OCID, "`Bot "+Runvals[1]+" launched with "+Runvals[2]+"`")
				}
			}
		}
	}
}

func RunAuto() {
	BotsB, err := ioutil.ReadFile("AutoStart")

	owner = "132583718291243008" // change this to your own ID if you use this bot for yourself

	botsfolder = "C:/Bots/ReBot/bots/" // change this to the absoluet location of the bots folder here

	if err != nil {
		panic(err)
	}
	Token, err := ioutil.ReadFile("token")
	if err != nil {
		panic(err)
	}
	token := strings.TrimSpace(string(Token))
	dg, err = discordgo.New(token)
	if err != nil {
		panic(err)
	}
	dg.AddHandler(Ready)
	dg.AddHandler(CM)
	BOTS := strings.Split(string(BotsB), "+")
	if len(BOTS) == 0 {
		return
	}
	var Bots []string
	for _, B := range BOTS {
		Bots = append(Bots, strings.TrimSpace(B))
	}
	returnedbot = make(chan string)
	for _, Bot := range Bots {
		go Launch(Bot, "")
	}
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	var RB string
	for {
		RB = <-returnedbot
		fmt.Println(RB)
		CFG, err := ioutil.ReadFile("bots/" + RB + ".bot")
		if err == nil {
			HandleCase(RB, CFG)
		}
	}
}

func Launch(Bot string, subcmd string) {
	var CMD *exec.Cmd
	if subcmd == "" {
		CMD = exec.Command(botsfolder + Bot + ".bat")
	} else {
		CMD = exec.Command(botsfolder + Bot + "-" + subcmd + ".bat")
	}
	cmdReader, err := CMD.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf(Bot+" | %s\n", scanner.Text())
		}
	}()
	err = CMD.Run()
	if err != nil {
		panic(err)
	}
	returnedbot <- Bot
}

func HandleCase(RB string, CFG []byte) {
	cfg := strings.Split(string(CFG), ",")
	var BotCFG = &BotCFG{}
	fmt.Println(cfg)
	for _, a := range cfg {
		vals := strings.Split(a, ":")
		if vals[0] == "Bot" {
			BotCFG.Path = strings.Join(vals[1:], ":")

		} else {
			num, err := strconv.ParseInt(vals[0], 10, 0)
			if err == nil {
				A := &action{
					num: int(num),
					run: vals[1],
				}
				BotCFG.Action = append(BotCFG.Action, A)
			}
		}
	}
	botboot, err := ioutil.ReadFile(BotCFG.Path + "/retcmd.botboot")
	if err == nil {
		handleme := strings.Split(string(botboot), " ")
		for i, val := range handleme {
			for _, A := range BotCFG.Action {
				if A.num-1 == i {
					yus, err := strconv.ParseBool(val)
					if err == nil {
						if yus {
							go Launch(RB, A.run)
						}
					}
				}
			}
		}
	}
}

type BotCFG struct {
	Path   string
	Action []*action
}

type action struct {
	num int
	run string
}

package main

import (
	"io/ioutil"
	"os/exec"
	"strings"

	"log"

	"strconv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	RunAuto()
}

var returnedbot chan string
var dg *discordgo.Session

var owner string
var OC *discordgo.Guild

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	OC, err := dg.UserChannelCreate(owner)
	if err != nil {
		log.Fatal(err)
		return
	}
	dg.ChannelMessageSend(OC.ID, "`Started ReBot!`")
}

func CM(Ses *discordgo.Session, MesC *discordgo.MessageCreate) {
	Mes := MesC.Message
	if Mes.Content != "" {
		if Mes.Content[0] == '!' && Mes.Author.ID == owner && len(Mes.Content) > 1 {
			var CMD = Mes.Content[1:]
			switch strings.ToLower(CMD) {
			case "run":
				Runvals := strings.Split(CMD, " ")
				if len(Runvals) == 1 {
					dg.ChannelMessageSend(OC.ID, "`Nil val`")
				} else if len(Runvals) == 2 {
					go Launch(Runvals[1], "")
					dg.ChannelMessageSend(OC.ID, "`Bot "+Runvals[1]+" launched `")
				} else if len(Runvals) == 3 {
					go Launch(Runvals[1], Runvals[2])
					dg.ChannelMessageSend(OC.ID, "`Bot "+Runvals[1]+" launched with "+Runvals[2]+"`")
				}
			}
		}
	}
}

func RunAuto() {
	BotsB, err := ioutil.ReadFile("AutoStart")

	owner = "132583718291243008" // change this to your own ID if you use this bot for yourself

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
	BOTS := (strings.Split(string(BotsB), "+"))
	if len(BOTS) == 0 {
		return
	}
	var Bots []string
	for _, B := range BOTS {
		Bots = append(Bots, strings.TrimSpace(B))
	}
	returnedbot := make(chan string, 10)
	for _, Bot := range Bots {
		go Launch(Bot, "")
	}
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	for {
		RB := <-returnedbot
		CFG, err := ioutil.ReadFile("bots/" + RB + ".bot")
		if err == nil {
			HandleCase(RB, CFG)
		}
	}
}

func Launch(Bot string, subcmd string) {
	var CMD *exec.Cmd
	if subcmd == "" {
		CMD = exec.Command("cmd", "/c", Bot+".bat")
	} else {
		CMD = exec.Command("cmd", "/c", Bot+"-"+subcmd+".bat")
	}
	CMD.Path = "../bots"
	CMD.Run()
	returnedbot <- Bot
}

func HandleCase(RB string, CFG []byte) {
	cfg := strings.Split(string(CFG), ",")
	var BotCFG *BotCFG
	for _, a := range cfg {
		vals := strings.Split(a, ":")
		if vals[0] == "Bot" {
			BotCFG.Path = vals[1]

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
							Launch(RB, A.run)
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

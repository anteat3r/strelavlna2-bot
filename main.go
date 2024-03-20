package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var DB_URL string

func endpoint(path string) string {
  return "http://" + DB_URL + "/" + path
}

func check(err error) {
  if err != nil { log.Println(err) }
}

func main() {
  err := godotenv.Load("../.env")  
  if err != nil { log.Fatal(err) }

  DB_URL = os.Getenv("DATABASE_URL")

  discord, err := discordgo.New(
    "Bot " + os.Getenv("DISCORD_BOT_TOKEN"),
  )
  if err != nil { log.Fatal(err) }

  discord.AddHandler(onMessage)

  discord.Open()
  defer discord.Close()

  fmt.Println("Bot running...")
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  <-c
}

func onMessage(
  dis *discordgo.Session,
  msg *discordgo.MessageCreate,
) {
  if msg.Author.ID == dis.State.User.ID { return }

  rawcmds := strings.Split(msg.Content, " ")
  if len(rawcmds) < 2 { return }
  if rawcmds[0] != dis.State.User.Mention() { return }
  cmd := rawcmds[1]
  cmds := []string{}
  if len(rawcmds) > 2 { cmds = rawcmds[2:] }

  switch cmd {

  case "ping":
    res, err := http.Get(endpoint("ping"))
    check(err)
    rawbody, err := io.ReadAll(res.Body)
    body := string(rawbody)
    dis.ChannelMessageSend(msg.ChannelID, body)
  case "config":
    if len(cmds) < 1 {
      _, err := dis.ChannelMessageSendReply(
        msg.ChannelID,
        "invalid request",
        msg.Reference(),
      )
      check(err)
      return
    }
    res, err := http.Get(endpoint("config/"+cmds[0]))
    if err != nil { log.Fatal(err) }
    rawbody, err := io.ReadAll(res.Body)
    body := string(rawbody)
    dis.ChannelMessageSend(msg.ChannelID, body)
  }
  // case "delthread":
}

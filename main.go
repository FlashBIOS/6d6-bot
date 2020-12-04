package main

import (
	"6d6/dice"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	check(err)
	defer func() {
		_ = discord.Close()
	}()

	discord.AddHandler(messageCreate)

	err = discord.Open()
	check(err)

	fmt.Printf("Bot 6d6 is now running. Press CTRL-C to exit.\n")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalChannel
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all of the bot's own messages.
	if message.Author.ID == session.State.User.ID {
		return
	}

	if strings.HasPrefix(message.Content, "!6d6") {
		pool, err := dice.Parse(message.Content)
		if err != nil {
			logMessageError(message.Message, err)
			message, err := session.ChannelMessageSend(message.ChannelID, message.Author.Mention()+"\n"+err.Error())
			if err != nil {
				logMessageError(message, err)
			}
			logMessage(message)
			return
		}

		results := dice.Roll(pool)

		message, err := session.ChannelMessageSend(message.ChannelID, message.Author.Mention()+"\n"+results.String())
		if err != nil {
			logMessageError(message, err)
		}
		logMessage(message)
	}
}

func logMessage(message *discordgo.Message) {
	log(fmt.Sprintf("[%v] %s\n", message.Timestamp, strings.ReplaceAll(message.Content, "\n", " ")))
}

func logMessageError(message *discordgo.Message, err error) {
	log(fmt.Sprintf("[%v] ERROR %s\n", message.Timestamp, err))
}

func log(str string) {
	fmt.Print(str)
}

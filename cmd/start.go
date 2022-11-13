package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Killayt/LoadMusic-bot/bot"
)

func main() {
	token := os.Getenv("MUSICAL_LOAD_BOT")
	b, err := bot.NewTelegramBot(
		token,
		6000,           // max download time
		120,            // max video duration
		"Musical Load", // bot name
	)
	if err != nil {
		log.Fatal(err)
	}

	sigint := make(chan os.Signal)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)

	go func() {
		<-sigint
		fmt.Println("Gracefully stopping app")
		b.Stop()
		os.Exit(1)
	}()

	if err := b.Run(true); err != nil {
		fmt.Println(err.Error())
	}
}

package main

import (
	"log"
	"os"

	"github.com/kelr/gundyr"
)

func auth() {
	cfg := &gundyr.HelixConfig{
		ClientID:     os.Getenv("TWITCH_ID"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
	}

	c, err := gundyr.NewHelix(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userID, err := c.UserToID("ohgustie")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userID)

	clips, err := c.GetAllClips(userID, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(clips)
}

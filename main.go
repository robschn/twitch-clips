package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kelr/gundyr"
)

func main() {

	// authenticate and grab session
	session := auth(os.Getenv("TWITCH_ID"), os.Getenv("TWITCH_SECRET"))

	// grab video clips URL
	vidURLs := getVidClips(session, "ohgustie")

	fmt.Printf("vidURLs: %v\n", vidURLs)
}

func auth(clientID string, clientSecret string) *gundyr.Helix {
	// Authenticate
	cfg := &gundyr.HelixConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	c, err := gundyr.NewHelix(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func getVidClips(c *gundyr.Helix, username string) []string {

	// convert username to user ID
	userID, err := c.UserToID(username)
	if err != nil {
		log.Fatal(err)
	}

	// grab raw clip data from user ID
	clips, err := c.GetAllClips(userID, "")
	if err != nil {
		log.Fatal(err)
	}

	// initialize video URL slice
	vidURLs := []string{}

	for _, v := range clips {

		// direct video links are not available, but we can extract it from ThumbnailURL
		if strings.Contains(v.ThumbnailURL, "AT-cm%") {
			splitURL := strings.Split(v.ThumbnailURL, "-preview-")

			// base URL will be the first element, add .mp4
			vidURLs = append(vidURLs, splitURL[0]+".mp4")
		}
	}
	return vidURLs
}

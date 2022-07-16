package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kelr/gundyr"
)

func main() {

	// grab secrets from env and auth
	// export TWITCH_ID=xxxxx
	// export TWITCH_SECRET=xxxxx
	session := auth(os.Getenv("TWITCH_ID"), os.Getenv("TWITCH_SECRET"))

	// grab userIDs from env
	// export TWITCH_USER="user1,user2,user3"
	userIDs := strings.Split(os.Getenv("TWITCH_USERS"), ",")

	for _, userID := range userIDs {

		// grab video clips URL
		vidURLs, err := getVidClips(session, userID)

		// skip if no user found
		if err != nil {
			continue
		}

		// print to m3u file
		printFileM3U(vidURLs, userID)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func auth(clientID string, clientSecret string) *gundyr.Helix {
	// Authenticate
	cfg := &gundyr.HelixConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	c, err := gundyr.NewHelix(cfg)
	check(err)

	return c
}

func getVidClips(c *gundyr.Helix, username string) ([]string, error) {

	// convert username to user ID
	userID, err := c.UserToID(username)
	if err != nil {
		return nil, err
	}

	// grab raw clip data from user ID
	clips, err := c.GetAllClips(userID, "")
	check(err)

	// shuffle vidURLs
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(clips), func(i, j int) { clips[i], clips[j] = clips[j], clips[i] })

	// initialize video URL slice
	vidURLs := []string{}

	// only need the first 25 links
	for _, v := range clips[:25] {
		// direct video links are not available, but we can extract it from ThumbnailURL
		if strings.Contains(v.ThumbnailURL, "-preview-") {
			splitURL := strings.Split(v.ThumbnailURL, "-preview-")

			// base URL will be the first element, add .mp4
			vidURLs = append(vidURLs, splitURL[0]+".mp4")
		}
	}

	return vidURLs, nil
}

func printFileM3U(vidArray []string, userID string) {

	// create clips dir
	os.Mkdir("clips/", 0770)

	// create clips filepath
	filePath := fmt.Sprintf("clips/%s.m3u", userID)
	f, err := os.Create(filePath)

	check(err)

	// remember to close the file
	defer f.Close()

	// write URLs to the file
	for _, line := range vidArray {
		_, err := f.WriteString(line + "\n")
		check(err)
	}
}

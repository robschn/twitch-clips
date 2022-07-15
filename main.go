package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kelr/gundyr"
)

func main() {

	// authenticate and grab session
	session := auth(os.Getenv("TWITCH_ID"), os.Getenv("TWITCH_SECRET"))

	// grab video clips URL
	vidURLs := getVidClips(session, "ohgustie")

	// print to m3u file
	printFileM3U(vidURLs)
}

func printFileM3U(vidArray []string) {

	// create clips dir
	os.Mkdir("clips/", 0770)
	f, err := os.Create("clips/twitch.m3u")

	check(err)

	// remember to close the file
	defer f.Close()

	for _, line := range vidArray {
		_, err := f.WriteString(line + "\n")
		check(err)
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

func getVidClips(c *gundyr.Helix, username string) []string {

	// convert username to user ID
	userID, err := c.UserToID(username)
	check(err)

	// grab raw clip data from user ID
	clips, err := c.GetAllClips(userID, "")
	check(err)

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

	// shuffle vidURLs
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vidURLs), func(i, j int) { vidURLs[i], vidURLs[j] = vidURLs[j], vidURLs[i] })

	return vidURLs
}

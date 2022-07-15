package main

import (
	"bufio"
	"fmt"
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

	printFileM3U(vidURLs)
}

func printFileM3U(vidArray []string) {

	const timeLayout = "01-02-2006"

	timeStamp := time.Now().Format(timeLayout)
	fileName := fmt.Sprintf("twitch-clips-%s.m3u", timeStamp)

	file, err := os.Create(fileName)

	check(err)

	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range vidArray {
		fmt.Fprintln(w, line)
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
	return vidURLs
}

// vlc https://clips-media-assets2.twitch.tv/AT-cm%7C1239136642.mp4 --live-caching=10 --sout '#transcode{vcodec=mp2v,vb=256,acodec=ne}:std{access=udp{caching=10},mux=raw,dst=localhost:8081}'

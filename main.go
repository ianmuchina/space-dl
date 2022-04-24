package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const bearer = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs=1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

func main() {
	spaceID := flag.String("id", "1eaKbNYzMzjKX", "Space ID")
	flag.Parse()
	Guest, err := newGuest()
	if err != nil {
		log.Fatal(err)
	}

	Space, err := Guest.getAudioSpace(*spaceID)
	if err != nil {
		log.Fatal(err)
	}

	Source, err := Guest.getAudioSpaceSource(Space)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Source.Location)
}

func newGuest() (*guest, error) {
	url := "https://api.twitter.com/1.1/guest/activate.json"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	// Add Headers
	req.Header.Set("Authorization", "Bearer "+bearer)
	// Make request
	response, err := new(http.Transport).RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// Make new g
	g := new(guest)
	// Unmarshal json
	if err := json.NewDecoder(response.Body).Decode(g); err != nil {
		return nil, err
	}
	return g, nil
}

// Get Space Data
func (g guest) getAudioSpace(id string) (*audioSpace, error) {
	var space struct {
		Data struct {
			AudioSpace audioSpace
		}
	}
	url_ := "https://twitter.com/i/api/graphql/lFpix9BgFDhAMjn9CrW6jQ/AudioSpaceById"

	// Client
	req, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Authorization", "Bearer "+bearer)
	req.Header.Add("X-Guest-Token", g.Guest_Token)

	// Query String
	buf, err := json.Marshal(spaceRequest{ID: id})
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = "variables=" + url.QueryEscape(string(buf))

	// Make Request
	res, err := new(http.Transport).RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Ensure response is 200
	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	// Decode JSON
	if err := json.NewDecoder(res.Body).Decode(&space); err != nil {
		return nil, err
	}

	// Return Data
	return &space.Data.AudioSpace, nil
}

// Get m3u hls link to space vod
func (g guest) getAudioSpaceSource(space *audioSpace) (*source, error) {
	var video struct {
		Source source
	}
	url := fmt.Sprintf(
		"https://twitter.com/i/api/1.1/live_video_stream/status/%s",
		space.Metadata.Media_Key,
	)

	// Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Authorization", "Bearer "+bearer)
	req.Header.Add("X-Guest-Token", g.Guest_Token)

	// Make Request
	res, err := new(http.Transport).RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Ensure Response is 200
	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	// Decode Json
	if err := json.NewDecoder(res.Body).Decode(&video); err != nil {
		return nil, err
	}

	// Return Data
	return &video.Source, nil
}

type source struct {
	Location string // Segment
}

type guest struct {
	Guest_Token string
}

type variant struct {
	Bitrate      int64
	Content_Type string
	URL          string
}

type media struct {
	Media_URL     string
	Original_Info struct {
		Width  int64
		Height int64
	}
	Video_Info struct {
		Variants []variant
	}
}

type audioSpace struct {
	Metadata struct {
		Media_Key  string
		Title      string
		State      string
		Started_At int64
		Ended_At   int64 `json:"ended_at,string"`
	}
	Participants struct {
		Admins []struct {
			Display_Name string
		}
	}
}

type spaceRequest struct {
	ID                          string `json:"id"`
	IsMetatagsQuery             bool   `json:"isMetatagsQuery"`
	WithBirdwatchPivots         bool   `json:"withBirdwatchPivots"`
	WithDownvotePerspective     bool   `json:"withDownvotePerspective"`
	WithReactionsMetadata       bool   `json:"withReactionsMetadata"`
	WithReactionsPerspective    bool   `json:"withReactionsPerspective"`
	WithReplays                 bool   `json:"withReplays"`
	WithScheduledSpaces         bool   `json:"withScheduledSpaces"`
	WithSuperFollowsTweetFields bool   `json:"withSuperFollowsTweetFields"`
	WithSuperFollowsUserFields  bool   `json:"withSuperFollowsUserFields"`
}

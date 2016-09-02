package google

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/model"
)

// WFromTakeout parses a history file (JSON) as inside a Google Takeout.
func WFromTakeout(historyFile io.Reader) (err error) {
	var watchhistory WatchHistory
	jsoncontent, err := ioutil.ReadAll(historyFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsoncontent, &watchhistory.Videos)
	if err != nil {
		return
	}
	return insertWatchHistory(watchhistory)
}

// Insert the watch history into the datatrack.
func insertWatchHistory(history WatchHistory) (err error) {

	wg := new(sync.WaitGroup)
	logChan := make(chan Video, len(history.Videos))
	parsedChan := make(chan ParsedVideo, len(history.Videos))

	// start workers, one per CPU
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go watchWorker(logChan, parsedChan, wg)
	}

	// feed workers
	for _, video := range history.Videos {
		logChan <- video
	}

	disclosures := make([]model.Disclosure, 0, 1*len(history.Videos))
	disclosed := make([]model.Disclosed, 0, 1*len(history.Videos))
	attributes := make([]model.Attribute, 0, 13*len(history.Videos))

	// wait for workers to finish
	close(logChan)
	wg.Wait()

	// assemble all parsed locations
	close(parsedChan)
	for p := range parsedChan {
		disclosures = append(disclosures, p.Disclosure)
		disclosed = append(disclosed, p.Disclosed)
		attributes = append(attributes, p.Attributes...)

	}

	wg = new(sync.WaitGroup)
	errChan := make(chan error, 4)
	wg.Add(3)

	// send into database
	go database.AddAttributes(attributes, wg, errChan)
	go database.AddDisclosures(disclosures, wg, errChan)
	go database.AddDiscloseds(disclosed, wg, errChan)

	wg.Wait()

	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func watchWorker(locChan chan Video, parsedChan chan ParsedVideo, wg *sync.WaitGroup) {
	defer wg.Done()
	for video := range locChan {
		parsed, err := insertVideo(video)
		if err != nil {
			panic(err)
		}
		parsedChan <- parsed
	}
}

// Converts the timestamp of when the video was watched and returns a parsed video.
func insertVideo(video Video) (parsed ParsedVideo, err error) {

	time, err := time.Parse(time.RFC3339Nano, video.Snippet.PublishedAt)
	if err != nil {
		panic(err)
	}

	timeEpochMs := time.Unix() * 1000

	parsed, err = createParsedVideo(timeEpochMs, video)
	if err != nil {
		return
	}

	return
}

// Creates a video according to the data model by adding attributes and assigning disclosures.
func createParsedVideo(timestamp int64, video Video) (parsedVideo ParsedVideo, err error) {

	parsedVideo.Disclosure, err = model.MakeDisclosure(database.Self, org.ID,
		strconv.Itoa(int(timestamp)), "", "", "", "")
	if err != nil {
		return
	}

	parsedVideo.Attributes = make([]model.Attribute, 13)

	parsedVideo.Attributes[0], err = model.MakeAttribute("Video id", "camera",
		fmt.Sprintf("%s", video.ContentDetails.VideoID))
	if err != nil {
		return
	}

	parsedVideo.Attributes[1], err = model.MakeAttribute("Etag", "tag",
		fmt.Sprintf("%s", video.ETag))
	if err != nil {
		return
	}

	parsedVideo.Attributes[2], err = model.MakeAttribute("Id", "barcode",
		fmt.Sprintf("%s", video.ID))
	if err != nil {
		return
	}

	parsedVideo.Attributes[3], err = model.MakeAttribute("Kind", "filter",
		fmt.Sprintf("%s", video.Kind))
	if err != nil {
		return
	}

	parsedVideo.Attributes[4], err = model.MakeAttribute("Channel id", "barcode",
		fmt.Sprintf("%s", video.Snippet.ChannelID))
	if err != nil {
		return
	}

	parsedVideo.Attributes[5], err = model.MakeAttribute("Channel Title", "info",
		fmt.Sprintf("%s", video.Snippet.ChannelTitle))
	if err != nil {
		return
	}

	parsedVideo.Attributes[6], err = model.MakeAttribute("Description", "info",
		fmt.Sprintf("%s", video.Snippet.Description))
	if err != nil {
		return
	}

	parsedVideo.Attributes[7], err = model.MakeAttribute("Playlist id", "barcode",
		fmt.Sprintf("%s", video.Snippet.PlaylistID))
	if err != nil {
		return
	}

	parsedVideo.Attributes[8], err = model.MakeAttribute("Position", "sort-numeric-asc",
		fmt.Sprintf("%d", video.Snippet.Position))
	if err != nil {
		return
	}

	parsedVideo.Attributes[9], err = model.MakeAttribute("Published at", "clock-o",
		fmt.Sprintf("%s", video.Snippet.PublishedAt))
	if err != nil {
		return
	}

	parsedVideo.Attributes[10], err = model.MakeAttribute("Kind", "filter",
		fmt.Sprintf("%s", video.Snippet.ResourceID.Kind))
	if err != nil {
		return
	}

	parsedVideo.Attributes[11], err = model.MakeAttribute("Title", "info",
		fmt.Sprintf("%s", video.Snippet.Title))
	if err != nil {
		return
	}

	parsedVideo.Attributes[12], err = model.MakeAttribute("Status", "lock",
		fmt.Sprintf("%s", video.Status))
	if err != nil {
		return
	}

	parsedVideo.Disclosed = model.Disclosed{
		Disclosure: parsedVideo.Disclosure.ID,
		Attribute: []string{parsedVideo.Attributes[0].ID, parsedVideo.Attributes[1].ID, parsedVideo.Attributes[2].ID,
			parsedVideo.Attributes[3].ID, parsedVideo.Attributes[4].ID, parsedVideo.Attributes[5].ID,
			parsedVideo.Attributes[6].ID, parsedVideo.Attributes[7].ID, parsedVideo.Attributes[8].ID,
			parsedVideo.Attributes[9].ID, parsedVideo.Attributes[10].ID, parsedVideo.Attributes[11].ID,
			parsedVideo.Attributes[12].ID}}

	return
}

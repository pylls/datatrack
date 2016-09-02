package google

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"runtime"
	"strconv"
	"sync"

	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/model"
)

// LFromTakeout parses a history file (JSON) as inside a Google Takeout.
func LFromTakeout(historyFile io.Reader) (err error) {
	var locationhistory LocationHistory
	jsoncontent, err := ioutil.ReadAll(historyFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsoncontent, &locationhistory)
	if err != nil {
		return
	}
	return insertLocationHistory(locationhistory)
}

func insertLocationHistory(history LocationHistory) (err error) {
	wg := new(sync.WaitGroup)
	logChan := make(chan Location, len(history.Locations))
	parsedChan := make(chan ParsedLocation, len(history.Locations))

	// start workers, one per CPU
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go locWorker(logChan, parsedChan, wg)
	}

	// feed workers
	for _, location := range history.Locations {
		logChan <- location
	}

	disclosures := make([]model.Disclosure, 0, 2*len(history.Locations))
	disclosed := make([]model.Disclosed, 0, 2*len(history.Locations))
	attributes := make([]model.Attribute, 0, 5*len(history.Locations))
	downstream := make([]model.Downstream, 0, 1*len(history.Locations))
	coordinates := make([]model.Coordinate, 0, 1*len(history.Locations))

	// wait for workers to finish
	close(logChan)
	wg.Wait()

	// assemble all parsed locations
	close(parsedChan)
	for p := range parsedChan {
		// position
		disclosures = append(disclosures, p.Position.Disclosure)
		disclosed = append(disclosed, p.Position.Disclosed)
		attributes = append(attributes, p.Position.Attributes...)
		coordinates = append(coordinates, p.Position.Coordinate)

		// activities
		for _, activity := range p.Activities {
			disclosures = append(disclosures, activity.Disclosure)
			disclosed = append(disclosed, activity.Disclosed)
			attributes = append(attributes, activity.Attributes...)
			downstream = append(downstream, activity.Downstream)
		}
	}

	wg = new(sync.WaitGroup)
	errChan := make(chan error, 5)
	wg.Add(5)
	// send into database
	go database.AddAttributes(attributes, wg, errChan)
	go database.AddDisclosures(disclosures, wg, errChan)
	go database.AddDiscloseds(disclosed, wg, errChan)
	go database.AddDownstreams(downstream, wg, errChan)
	go database.AddCoordinates(coordinates, wg, errChan)

	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func locWorker(locChan chan Location, parsedChan chan ParsedLocation, wg *sync.WaitGroup) {
	defer wg.Done()
	for location := range locChan {
		parsed, err := insertLocation(location)
		if err != nil {
			panic(err)
		}
		parsedChan <- parsed
	}
}

func insertLocation(location Location) (parsed ParsedLocation, err error) {
	time, err := strconv.ParseInt(location.TimeStampMs, 10, 64)
	if err != nil {
		return
	}

	parsed.Position, err = createPosition(time, location)
	if err != nil {
		return
	}

	parsed.Activities = make([]Activity, 0, len(location.OuterActivities))
	for _, activity := range location.OuterActivities {
		time, err := strconv.ParseInt(activity.OTimeStampMs, 10, 64)
		if err != nil {
			return parsed, err
		}

		pAct, err := createActivity(parsed.Position, time, activity)
		if err != nil {
			return parsed, err
		}
		parsed.Activities = append(parsed.Activities, pAct)
	}
	return
}

func createPosition(timestamp int64, location Location) (position Position, err error) {
	position.Disclosure, err = model.MakeDisclosure(database.Self, org.ID,
		strconv.Itoa(int(timestamp)), "", "", "", "")
	if err != nil {
		return
	}

	position.Attributes = make([]model.Attribute, 5)
	latitude := big.NewRat(location.LatitudeE7, 10000000).FloatString(7)
	longitude := big.NewRat(location.LongitudeE7, 10000000).FloatString(7)
	// coordinates
	position.Attributes[0], err = model.MakeAttribute("Coordinates", "location-arrow",
		fmt.Sprintf("%s, %s",
			latitude,
			longitude))
	if err != nil {
		return
	}
	// accuracy
	position.Attributes[1], err = model.MakeAttribute("Accuracy", "bullseye",
		fmt.Sprintf("%d", location.Accuracy))
	if err != nil {
		return
	}
	// velocity
	position.Attributes[2], err = model.MakeAttribute("Velocity", "tachometer",
		fmt.Sprintf("%d", location.Velocity))
	if err != nil {
		return
	}
	// heading
	position.Attributes[3], err = model.MakeAttribute("Heading", "arrows-alt",
		fmt.Sprintf("%d", location.Heading))
	if err != nil {
		return
	}
	// altitude
	position.Attributes[4], err = model.MakeAttribute("Altitude", "signal",
		fmt.Sprintf("%d", location.Altitude))
	if err != nil {
		return
	}

	position.Disclosed = model.Disclosed{
		Disclosure: position.Disclosure.ID,
		Attribute: []string{position.Attributes[0].ID, position.Attributes[1].ID,
			position.Attributes[2].ID, position.Attributes[3].ID, position.Attributes[4].ID}}

	position.Coordinate = model.MakeCoordinate(latitude, longitude, position.Disclosure.ID,
		position.Disclosure.Timestamp)

	return
}

func createActivity(position Position, timestamp int64, activity OuterActivity) (pAct Activity, err error) {
	pAct.Disclosure, err = model.MakeDisclosure(org.ID, org.ID,
		strconv.Itoa(int(timestamp)), "", "", "", "")
	if err != nil {
		return
	}

	pAct.Downstream = model.Downstream{
		Origin: position.Disclosure.ID,
		Result: pAct.Disclosure.ID}

	pAct.Attributes = make([]model.Attribute, 0, len(activity.InnerActivities))
	attrIDs := make([]string, 0, len(activity.InnerActivities))
	for _, inner := range activity.InnerActivities {
		enc, err := json.Marshal(inner)
		if err != nil {
			return pAct, err
		}
		attribute, err := model.MakeAttribute("Activity", "child", string(enc))
		if err != nil {
			return pAct, err
		}
		pAct.Attributes = append(pAct.Attributes, attribute)
		attrIDs = append(attrIDs, attribute.ID)
	}

	pAct.Disclosed = model.Disclosed{
		Disclosure: pAct.Disclosure.ID,
		Attribute:  attrIDs}

	return
}

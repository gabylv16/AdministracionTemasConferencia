package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
)

//conference structure
type conference struct {
	topic    string
	duration time.Duration
	hour     time.Time
}

type conferenceFrame struct {
	initialHour          time.Time
	finalHour            time.Time
	scheduledConferences []conference
	currentHour          time.Time
}

func newConferenceFrame(initHour, finHour string) conferenceFrame {
	var a conferenceFrame
	a.initialHour, _ = time.Parse("15:04", initHour)
	a.currentHour, _ = time.Parse("15:04", initHour)
	a.finalHour, _ = time.Parse("15:04", finHour)
	return a
}

type thematic struct {
	morning   conferenceFrame
	afternoon conferenceFrame
}

func main() {
	// file reading
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	//file to string conversion
	s := string(data[:])
	//string separation by lines
	lines := strings.Split(s, "\n")
	//declaration of a slice for conferences
	conferenceList := []conference{}

	totalDuration := 0

	//storage of the conferences in the slice with name and duration
	for _, element := range lines {

		words := strings.Split(element, " ")
		lastWord := words[len(words)-1]
		dur, td, err := determineDuration(lastWord)
		if err != nil {
			fmt.Println(err)
			return
		}
		conf := conference{topic: element, duration: dur}
		conferenceList = append(conferenceList, conf)
		totalDuration = totalDuration + td
	}

	//order conferences from longest to shortest
	sort.SliceStable(conferenceList, func(i, j int) bool {
		return conferenceList[i].duration > conferenceList[j].duration
	})

	var thematicSlice []thematic
	//creates lunch and networking event conference
	lunch := conference{topic: "Lunch"}
	lunch.hour, _ = time.Parse("15:04", lunchHour)
	networkingEvent := conference{topic: "Networking Event"}
	networkingEvent.hour, _ = time.Parse("15:04", eventHout)
	thematicSlice = assingnConferencesTothematic(conferenceList, thematicSlice)

	for i, thematic := range thematicSlice {
		fmt.Println()
		fmt.Print("THEMATIC #")
		fmt.Println(i + 1)
		thematic.morning.scheduledConferences = append(thematic.morning.scheduledConferences, lunch)
		for _, element := range thematic.morning.scheduledConferences {
			fmt.Print(element.hour.Format("15:04"))
			fmt.Print(" ")
			fmt.Println(element.topic)
		}
		thematic.afternoon.scheduledConferences = append(thematic.afternoon.scheduledConferences, networkingEvent)
		for _, element := range thematic.afternoon.scheduledConferences {
			fmt.Print(element.hour.Format("15:04"))
			fmt.Print(" ")
			fmt.Println(element.topic)
		}

	}
}

//determines the duration of the conferece based on the last word (lightning or -min)
func determineDuration(lastWord string) (time.Duration, int, error) {
	dur, _ := time.ParseDuration("0m")
	intDur := 0
	var err error
	if strings.Contains(lastWord, "lightning") {
		dur, _ = time.ParseDuration("5m")
		intDur = 5
	} else if strings.Contains(lastWord, "min") {
		durStr := lastWord[:strings.LastIndex(lastWord, "i")]
		intDurStr := lastWord[:strings.LastIndex(lastWord, "m")]
		dur, _ = time.ParseDuration(durStr)
		intDur, _ = strconv.Atoi(intDurStr)
	} else {
		err = errors.New("Malformed file")
	}
	return dur, intDur, err
}

//time constants are set
const (
	morningInitialHour   = "09:00"
	morningFinalHour     = "12:00"
	afternoonInitialHour = "13:00"
	afternoonFinalHour   = "17:00"
	lunchHour            = "12:00"
	eventHout            = "17:00"
)

//assigns the conference to a frame, increments the current hour + the duration of the conference
func assignConferenceToFrame(conf conference, frame conferenceFrame) (conference, conferenceFrame) {
	conf.hour = frame.currentHour
	frame.scheduledConferences = append(frame.scheduledConferences, conf)
	frame.currentHour = frame.currentHour.Add(conf.duration)
	return conf, frame
}

//assigns morning and afternoon events to a thematic, and the thematic to a final thematic slice
func assingnConferencesTothematic(conferenceList []conference, thematicList []thematic) []thematic {
	var thematic thematic
	thematic.morning = newConferenceFrame(morningInitialHour, morningFinalHour)
	thematic.afternoon = newConferenceFrame(afternoonInitialHour, afternoonFinalHour)

	var nonScheduledConfs []conference
	for _, element := range conferenceList {

		if element.duration <= thematic.morning.finalHour.Sub(thematic.morning.currentHour) {
			element, thematic.morning = assignConferenceToFrame(element, thematic.morning)
		} else if element.duration <= thematic.afternoon.finalHour.Sub(thematic.afternoon.currentHour) {
			element, thematic.afternoon = assignConferenceToFrame(element, thematic.afternoon)
		} else {
			nonScheduledConfs = append(nonScheduledConfs, element)
		}

	}
	thematicList = append(thematicList, thematic)
	if len(nonScheduledConfs) > 0 {
		thematicList = assingnConferencesTothematic(nonScheduledConfs, thematicList)
	}
	return thematicList
}

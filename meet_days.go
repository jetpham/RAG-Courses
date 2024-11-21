package main

import "fmt"

var meetDays = map[string]string{
	"U": "Sunday",
	"M": "Monday",
	"T": "Tuesday",
	"W": "Wednesday",
	"R": "Thursday",
	"F": "Friday",
	"S": "Saturday",
}

type MeetDay struct {
	Code string
	Full string
}

func UnmarshalCSVtoMeetDay(csv string) (MeetDay, error) {

	if full, ok := meetDays[csv]; ok {
		return MeetDay{
			Code: csv,
			Full: full,
		}, nil
	}

	return MeetDay{}, fmt.Errorf("invalid meet day code: %s", csv)
}
func (m MeetDay) String() string {
	return fmt.Sprintf("Code: %s, Full: %s", m.Code, m.Full)
}

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Course struct {
	ID                     uint   `json:"id"`
	SubjectCode            string `json:"subject_code"`
	SubjectName            string `json:"subject_name"`
	CourseNumber           string `json:"course_number"`
	Section                string `json:"section"`
	CRN                    int    `json:"crn"`
	ScheduleTypeCode       string `json:"schedule_type_code"`
	CampusCode             string `json:"campus_code"`
	TitleShortDesc         string `json:"title_short_desc"`
	InstructionModeDesc    string `json:"instruction_mode_desc"`
	MeetingTypeCode        string `json:"meeting_type_code"`
	MeetingTypeName        string `json:"meeting_type_name"`
	MeetDays               string `json:"meet_days"`
	MeetDaysFull           string `json:"meet_days_full"`
	BeginTime              string `json:"begin_time"`
	EndTime                string `json:"end_time"`
	MeetStart              string `json:"meet_start"`
	MeetEnd                string `json:"meet_end"`
	Building               string `json:"building"`
	Room                   string `json:"room"`
	ActualEnrollment       int    `json:"actual_enrollment"`
	PrimaryInstructorFirst string `json:"primary_instructor_first_name"`
	PrimaryInstructorLast  string `json:"primary_instructor_last_name"`
	PrimaryInstructorFull  string `json:"primary_instructor_full_name"`
	PrimaryInstructorEmail string `json:"primary_instructor_email"`
	College                string `json:"college"`
}

func loadCSV(path string) ([]Course, error) {
	defer log.Println("CSV setup complete")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 0

	var courses []Course

	// Read and discard the header line
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		crn, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		actualEnrollment, err := strconv.Atoi(record[16])
		if err != nil {
			return nil, err
		}
		meetDays := ""
		meetDaysFull := ""
		for _, day := range record[9] {
			meetDays += string(day)
			if meetDay, err := UnmarshalCSVtoMeetDay(string(day)); err == nil {
				meetDaysFull += meetDay.Full + " "
			} else {
				meetDaysFull += ""
			}
		}
		meetDaysFull = strings.TrimSpace(meetDaysFull)

		course := Course{
			ID:          uint(len(courses) + 1),
			SubjectCode: record[0],
			SubjectName: func() string {
				if val, ok := courseSubjects[record[0]]; ok {
					return val
				} else {
					return ""
				}
			}(),
			CourseNumber:        record[1],
			Section:             record[2],
			CRN:                 crn,
			ScheduleTypeCode:    record[4],
			CampusCode:          record[5],
			TitleShortDesc:      record[6],
			InstructionModeDesc: record[7],
			MeetingTypeCode:     record[8],
			MeetingTypeName: func() string {
				if val, ok := instructionalMethods[record[8]]; ok {
					return val
				} else {
					return ""
				}
			}(),
			MeetDays:               meetDays,
			MeetDaysFull:           meetDaysFull,
			BeginTime:              record[10],
			EndTime:                record[11],
			MeetStart:              record[12],
			MeetEnd:                record[13],
			Building:               record[14],
			Room:                   record[15],
			ActualEnrollment:       actualEnrollment,
			PrimaryInstructorFirst: record[17],
			PrimaryInstructorLast:  record[18],
			PrimaryInstructorFull:  record[17] + " " + record[18],
			PrimaryInstructorEmail: record[19],
			College:                record[20],
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (c Course) String() string {
	return fmt.Sprintf(
		"%s%s %s Section %s by %s %s, ",
		c.SubjectCode,
		c.CourseNumber,
		c.TitleShortDesc,
		c.Section,
		c.PrimaryInstructorFirst,
		c.PrimaryInstructorLast,
	)
}

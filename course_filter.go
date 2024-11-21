package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/invopop/jsonschema"
)

type CourseFilter struct {
	SubjectCodes            []string `json:"subject_codes" jsonschema_description:"Any subject codes listed explicitly in the prompt. Example: ['CS', 'MATH'], not ['Computer Science', 'Mathematics']"`
	SubjectNames            []string `json:"subject_names" jsonschema_description:"The subject names listed explicitly in the prompt. Example: ['Computer Science', 'Mathematics']"`
	CourseNumbers           []string `json:"course_numbers" jsonschema_description:"The course numbers. Example: ['101', '202']"`
	Sections                []string `json:"sections" jsonschema_description:"The sections of the courses. Example: ['001', '002']"`
	CRNs                    []int    `json:"crns" jsonschema_description:"The course registration numbers. Example: [12345, 67890]"`
	ScheduleTypeCodes       []string `json:"schedule_type_codes" jsonschema_description:"The schedule type codes, if provided as codes. Example: ['LEC', 'LAB']"`
	CampusCodes             []string `json:"campus_codes" jsonschema_description:"The campus codes, if provided as codes. Example: ['MAIN', 'SAT']"`
	Title                   []string `json:"title" jsonschema_description:"The title of the course. Example: ['Intro to CS', 'Calculus I', 'Physics II']"`
	InstructionModeDescs    []string `json:"instruction_mode_descs" jsonschema_description:"The descriptions of the instruction modes. Example: ['In Person', 'Online']"`
	MeetingTypeCodes        []string `json:"meeting_type_codes" jsonschema_description:"The meeting type codes, if provided as codes. Example: ['CLAS', 'LAB']"`
	MeetingTypeNames        []string `json:"meeting_type_names" jsonschema_description:"The meeting type names, if provided as names. Example: ['Class', 'Laboratory']"`
	MeetDays                []string `json:"meet_days" jsonschema_description:"The meeting days. Example: ['M', 'W', 'F']"`
	MeetDaysFull            []string `json:"meet_days_full" jsonschema_description:"The full meeting days. Example: ['Monday', 'Wednesday', 'Friday']"`
	BeginTimes              []string `json:"begin_times" jsonschema_description:"The begin times of the meetings. Example: ['08:00', '10:00']"`
	EndTimes                []string `json:"end_times" jsonschema_description:"The end times of the meetings. Example: ['09:00', '11:00']"`
	MeetStarts              []string `json:"meet_starts" jsonschema_description:"The start dates of the meetings. Example: ['2023-01-10', '2023-01-12']"`
	MeetEnds                []string `json:"meet_ends" jsonschema_description:"The end dates of the meetings. Example: ['2023-05-10', '2023-05-12']"`
	Buildings               []string `json:"buildings" jsonschema_description:"The buildings where the courses are held. Example: ['Engineering Hall', 'Science Building']"`
	Rooms                   []string `json:"rooms" jsonschema_description:"The rooms where the courses are held. Example: ['101', '202']"`
	ActualEnrollments       []int    `json:"actual_enrollments" jsonschema_description:"The actual enrollments in the courses. Example: [30, 25]"`
	PrimaryInstructorNames  []string `json:"primary_instructor_names" jsonschema_description:"The names of the primary instructors. Example: ['John Doe', 'Jane Smith', 'Phil', 'Pham']"`
	PrimaryInstructorEmails []string `json:"primary_instructor_emails" jsonschema_description:"The emails of the primary instructors. Example: ['jdoe@example.com', 'jsmith@example.com']"`
	Colleges                []string `json:"colleges" jsonschema_description:"The colleges offering the courses. Example: ['College of Engineering', 'College of Science']"`
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func (d *chromaDB) correctCourseFilter(collections *collections, filter CourseFilter) CourseFilter {
	combinedSubjectNamesAndTitles := append(filter.SubjectNames, filter.Title...)
	filter.SubjectNames = combinedSubjectNamesAndTitles
	filter.Title = combinedSubjectNamesAndTitles
	filter.SubjectNames = d.correctSubjectNames(collections, filter)
	filter.Title = d.correctTitleShortDescs(collections, filter)
	filter.PrimaryInstructorNames = d.correctInstructorFullNames(collections, filter)
	return filter
}

func (d *chromaDB) correctSubjectNames(collections *collections, filter CourseFilter) []string {
	tempSubjectNames := make([]string, 0)
	for _, code := range filter.SubjectCodes {
		if name, exists := courseSubjects[code]; exists {
			tempSubjectNames = append(tempSubjectNames, name)
		}
	}

	// Query the collections for subject names
	correctedSubjectNames := make([]string, 0, len(filter.SubjectNames))
	for _, subjectName := range filter.SubjectNames {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.SubjectNameList, subjectName) {
			correctedSubjectNames = append(correctedSubjectNames, subjectName)
			continue
		}
		result, err := d.query(collections.SubjectNameCollection, subjectName, 1)
		if err != nil {
			log.Fatalf("Error querying subject name '%s': %s \n", subjectName, err)
		}
		correctedSubjectNames = append(correctedSubjectNames, result...)
	}

	// Add the converted subject names to the corrected ones
	correctedSubjectNames = append(correctedSubjectNames, tempSubjectNames...)

	return correctedSubjectNames
}

func (d *chromaDB) correctTitleShortDescs(collections *collections, filter CourseFilter) []string {
	correctedTitleShortDescs := make([]string, 0, len(filter.Title))
	for _, titleShortDesc := range filter.Title {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.TitleShortDescList, titleShortDesc) {
			correctedTitleShortDescs = append(correctedTitleShortDescs, titleShortDesc)
			continue
		}
		result, err := d.query(collections.TitleShortDescCollection, titleShortDesc, 1)
		if err != nil {
			log.Fatalf("Error querying title short description '%s': %s \n", titleShortDesc, err)
		}
		correctedTitleShortDescs = append(correctedTitleShortDescs, result...)
	}
	return correctedTitleShortDescs
}

func (d *chromaDB) correctInstructorFullNames(collections *collections, filter CourseFilter) []string {
	correctedNames := make([]string, 0, len(filter.PrimaryInstructorNames))
	for _, name := range filter.PrimaryInstructorNames {
		// Shortcut if it's in the collections already
		if slices.Contains(collections.InstructorFullNameList, name) ||
			slices.Contains(collections.InstructorFirstNameList, name) ||
			slices.Contains(collections.InstructorLastNameList, name) {
			correctedNames = append(correctedNames, name)
			continue
		}

		// make a query for the first, last, and full name collections and see their results. Then we'll make a mini collection with those results and see which is the most similar to the original name

		resultFull, err := d.query(collections.InstructorFullNameCollection, name, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor full name '%s': %s \n", name, err)
		}
		log.Printf("Result Full: %v", resultFull)
		resultFirst, err := d.query(collections.instructorFirstNameCollection, name, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor first name '%s': %s \n", name, err)
		}
		log.Printf("Result First: %v", resultFirst)
		resultLast, err := d.query(collections.instructorLastNameCollection, name, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor last name '%s': %s \n", name, err)
		}
		log.Printf("Result Last: %v", resultLast)

		collection, err := d.makeCollectionWithRecords([]string{resultFull[0], resultFirst[0], resultLast[0]})
		if err != nil {
			log.Fatalf("Error creating collection with records for '%s': %s \n", name, err)
		}
		result, err := d.query(collection, name, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor mini collection '%s': %s \n", name, err)
		}
		correctedNames = append(correctedNames, result...)
	}
	return correctedNames
}

func (f CourseFilter) String() string {
	var sb strings.Builder
	sb.WriteString("Filter: ")
	if len(f.SubjectCodes) > 0 {
		sb.WriteString(fmt.Sprintf("SubjectCodes: %v,", f.SubjectCodes))
	}
	if len(f.SubjectNames) > 0 {
		sb.WriteString(fmt.Sprintf("SubjectNames: %v,", f.SubjectNames))
	}
	if len(f.CourseNumbers) > 0 {
		sb.WriteString(fmt.Sprintf("CourseNumbers: %v,", f.CourseNumbers))
	}
	if len(f.Sections) > 0 {
		sb.WriteString(fmt.Sprintf("Sections: %v,", f.Sections))
	}
	if len(f.CRNs) > 0 {
		sb.WriteString(fmt.Sprintf("CRNs: %v,", f.CRNs))
	}
	if len(f.ScheduleTypeCodes) > 0 {
		sb.WriteString(fmt.Sprintf("ScheduleTypeCodes: %v,", f.ScheduleTypeCodes))
	}
	if len(f.CampusCodes) > 0 {
		sb.WriteString(fmt.Sprintf("CampusCodes: %v,", f.CampusCodes))
	}
	if len(f.Title) > 0 {
		sb.WriteString(fmt.Sprintf("TitleShortDescs: %v,", f.Title))
	}
	if len(f.InstructionModeDescs) > 0 {
		sb.WriteString(fmt.Sprintf("InstructionModeDescs: %v,", f.InstructionModeDescs))
	}
	if len(f.MeetingTypeCodes) > 0 {
		sb.WriteString(fmt.Sprintf("MeetingTypeCodes: %v,", f.MeetingTypeCodes))
	}
	if len(f.MeetingTypeNames) > 0 {
		sb.WriteString(fmt.Sprintf("MeetingTypeNames: %v,", f.MeetingTypeNames))
	}
	if len(f.MeetDays) > 0 {
		sb.WriteString(fmt.Sprintf("MeetDays: %v,", f.MeetDays))
	}
	if len(f.MeetDaysFull) > 0 {
		sb.WriteString(fmt.Sprintf("MeetDaysFull: %v,", f.MeetDaysFull))
	}
	if len(f.BeginTimes) > 0 {
		sb.WriteString(fmt.Sprintf("BeginTimes: %v,", f.BeginTimes))
	}
	if len(f.EndTimes) > 0 {
		sb.WriteString(fmt.Sprintf("EndTimes: %v,", f.EndTimes))
	}
	if len(f.MeetStarts) > 0 {
		sb.WriteString(fmt.Sprintf("MeetStarts: %v,", f.MeetStarts))
	}
	if len(f.MeetEnds) > 0 {
		sb.WriteString(fmt.Sprintf("MeetEnds: %v,", f.MeetEnds))
	}
	if len(f.Buildings) > 0 {
		sb.WriteString(fmt.Sprintf("Buildings: %v,", f.Buildings))
	}
	if len(f.Rooms) > 0 {
		sb.WriteString(fmt.Sprintf("Rooms: %v,", f.Rooms))
	}
	if len(f.ActualEnrollments) > 0 {
		sb.WriteString(fmt.Sprintf("ActualEnrollments: %v,", f.ActualEnrollments))
	}
	if len(f.PrimaryInstructorNames) > 0 {
		sb.WriteString(fmt.Sprintf("Names: %v,", f.PrimaryInstructorNames))
	}
	if len(f.PrimaryInstructorEmails) > 0 {
		sb.WriteString(fmt.Sprintf("PrimaryInstructorEmails: %v,", f.PrimaryInstructorEmails))
	}
	if len(f.Colleges) > 0 {
		sb.WriteString(fmt.Sprintf("Colleges: %v,", f.Colleges))
	}
	return sb.String()
}

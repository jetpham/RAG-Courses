package main

import (
	"fmt"
	"log"
	"sync"

	chroma "github.com/amikos-tech/chroma-go"
)

type collections struct {
	SubjectNameCollection         *chroma.Collection
	SubjectNameList               []string
	TitleShortDescCollection      *chroma.Collection
	TitleShortDescList            []string
	instructorFirstNameCollection *chroma.Collection
	InstructorFirstNameList       []string
	instructorLastNameCollection  *chroma.Collection
	InstructorLastNameList        []string
	InstructorFullNameCollection  *chroma.Collection
	InstructorFullNameList        []string
}

func makeCollections(db *chromaDB, courses []Course) (*collections, error) {
	defer log.Println("Collections setup complete")
	subjectNameRecords := make([]string, 0, len(courseSubjects))
	for _, subject := range courseSubjects {
		subjectNameRecords = append(subjectNameRecords, subject)
	}
	titleShortDescRecords := make([]string, 0, len(courses))
	instructorFullNameRecords := make([]string, 0, len(courses))
	instructorFirstNameRecords := make([]string, 0, len(courses))
	instructorLastNameRecords := make([]string, 0, len(courses))

	for _, course := range courses {
		titleShortDescRecords = append(titleShortDescRecords, course.TitleShortDesc)
		instructorFullNameRecords = append(instructorFullNameRecords, course.PrimaryInstructorFull)
		instructorFirstNameRecords = append(instructorFirstNameRecords, course.PrimaryInstructorFirst)
		instructorLastNameRecords = append(instructorLastNameRecords, course.PrimaryInstructorLast)
	}

	var wg sync.WaitGroup
	var err error

	var subjectNameCollection *chroma.Collection
	var titleShortDescCollection *chroma.Collection
	var instructorFirstNameCollection *chroma.Collection
	var instructorLastNameCollection *chroma.Collection
	var instructorFullNameCollection *chroma.Collection

	wg.Add(5)

	go func() {
		defer wg.Done()
		subjectNameCollection, err = db.makeCollectionWithRecords(subjectNameRecords)
		if err != nil {
			err = fmt.Errorf("failed to create SubjectNameCollection: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		titleShortDescCollection, err = db.makeCollectionWithRecords(titleShortDescRecords)
		if err != nil {
			err = fmt.Errorf("failed to create TitleShortDescCollection: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		instructorFirstNameCollection, err = db.makeCollectionWithRecords(instructorFirstNameRecords)
		if err != nil {
			err = fmt.Errorf("failed to create InstructorFirstNameCollection: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		instructorLastNameCollection, err = db.makeCollectionWithRecords(instructorLastNameRecords)
		if err != nil {
			err = fmt.Errorf("failed to create InstructorLastNameCollection: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		instructorFullNameCollection, err = db.makeCollectionWithRecords(instructorFullNameRecords)
		if err != nil {
			err = fmt.Errorf("failed to create InstructorFullNameCollection: %v", err)
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return &collections{
		SubjectNameCollection:         subjectNameCollection,
		SubjectNameList:               subjectNameRecords,
		TitleShortDescCollection:      titleShortDescCollection,
		TitleShortDescList:            titleShortDescRecords,
		instructorFirstNameCollection: instructorFirstNameCollection,
		InstructorFirstNameList:       instructorFirstNameRecords,
		instructorLastNameCollection:  instructorLastNameCollection,
		InstructorLastNameList:        instructorLastNameRecords,
		InstructorFullNameCollection:  instructorFullNameCollection,
		InstructorFullNameList:        instructorFullNameRecords,
	}, nil
}

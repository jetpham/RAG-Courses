package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSqlite(courses []Course) (*gorm.DB, error) {
	defer log.Println("SQLite setup complete")
	sqlDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = sqlDB.AutoMigrate(&Course{})
	if err != nil {
		return nil, err
	}
	for i := range courses {
		courses[i].ID = 0 // Reset ID to zero to avoid conflicts
	}
	sqlDB.CreateInBatches(courses, 100)
	return sqlDB, err
}

func filterCourses(sqlDB *gorm.DB, filter CourseFilter) ([]Course, error) {
	var courses []Course
	query := sqlDB.Model(&Course{})

	if len(filter.SubjectCodes) > 0 || len(filter.SubjectNames) > 0 {
		subQuery := sqlDB
		if len(filter.SubjectCodes) > 0 {
			subQuery = subQuery.Or("subject_code IN ?", filter.SubjectCodes)
		}
		if len(filter.SubjectNames) > 0 {
			subQuery = subQuery.Or("subject_name IN ?", filter.SubjectNames)
		}
		if len(filter.Title) > 0 {
			subQuery = subQuery.Or("title_short_desc IN ?", filter.Title)
		}
		query = query.Where(subQuery)
	}
	if len(filter.CourseNumbers) > 0 {
		subQuery := sqlDB
		for _, courseNumber := range filter.CourseNumbers {
			subQuery = subQuery.Or("course_number LIKE ?", courseNumber+"%")
		}
		query = query.Where(subQuery)
	}
	if len(filter.Sections) > 0 {
		query = query.Where("section IN ?", filter.Sections)
	}
	if len(filter.CRNs) > 0 {
		query = query.Where("crn IN ?", filter.CRNs)
	}
	if len(filter.ScheduleTypeCodes) > 0 {
		query = query.Where("schedule_type_code IN ?", filter.ScheduleTypeCodes)
	}
	if len(filter.CampusCodes) > 0 {
		query = query.Where("campus_code IN ?", filter.CampusCodes)
	}
	if len(filter.InstructionModeDescs) > 0 {
		query = query.Where("instruction_mode_desc IN ?", filter.InstructionModeDescs)
	}
	if len(filter.MeetingTypeCodes) > 0 || len(filter.MeetingTypeNames) > 0 {
		subQuery := sqlDB
		if len(filter.MeetingTypeCodes) > 0 {
			subQuery = subQuery.Or("meeting_type_code IN ?", filter.MeetingTypeCodes)
		}
		if len(filter.MeetingTypeNames) > 0 {
			subQuery = subQuery.Or("meeting_type_name IN ?", filter.MeetingTypeNames)
		}
		query = query.Where(subQuery)
	}
	if len(filter.MeetDays) > 0 || len(filter.MeetDaysFull) > 0 {
		subQuery := sqlDB
		for _, meetDay := range filter.MeetDays {
			subQuery = subQuery.Or("meet_days LIKE ?", "%"+meetDay+"%")
		}
		for _, meetDayFull := range filter.MeetDaysFull {
			subQuery = subQuery.Or("meet_days_full LIKE ?", "%"+meetDayFull+"%")
		}
		query = query.Where(subQuery)
	}
	if len(filter.BeginTimes) > 0 {
		query = query.Where("begin_time IN ?", filter.BeginTimes)
	}
	if len(filter.EndTimes) > 0 {
		query = query.Where("end_time IN ?", filter.EndTimes)
	}
	if len(filter.MeetStarts) > 0 {
		query = query.Where("meet_start IN ?", filter.MeetStarts)
	}
	if len(filter.MeetEnds) > 0 {
		query = query.Where("meet_end IN ?", filter.MeetEnds)
	}
	if len(filter.Buildings) > 0 {
		query = query.Where("building IN ?", filter.Buildings)
	}
	if len(filter.Rooms) > 0 {
		query = query.Where("room IN ?", filter.Rooms)
	}
	if len(filter.ActualEnrollments) > 0 {
		query = query.Where("actual_enrollment IN ?", filter.ActualEnrollments)
	}
	if len(filter.PrimaryInstructorNames) > 0 || len(filter.PrimaryInstructorEmails) > 0 {
		subQuery := sqlDB
		if len(filter.PrimaryInstructorNames) > 0 {
			subQuery = subQuery.Or("primary_instructor_first IN ?", filter.PrimaryInstructorNames).
				Or("primary_instructor_last IN ?", filter.PrimaryInstructorNames).
				Or("primary_instructor_full IN ?", filter.PrimaryInstructorNames)
		}
		if len(filter.PrimaryInstructorEmails) > 0 {
			subQuery = subQuery.Or("primary_instructor_email IN ?", filter.PrimaryInstructorEmails)
		}
		query = query.Where(subQuery)
	}
	if len(filter.Colleges) > 0 {
		query = query.Where("college IN ?", filter.Colleges)
	}

	err := query.Find(&courses).Error
	return courses, err
}

package models

import "database/sql"

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Avatar      []byte `json:"avatar"`
	OwnerID     int    `json:"owner_id"`
}

type CourseModel struct {
	model
}

func NewCourseModel(db *sql.DB) CourseModel {
	return CourseModel{model{db}}
}

func (m CourseModel) GetList(limit int, offset int) []Course {
	courses := make([]Course, 0)

	rows, err := m.db.Query(`SELECT * FROM courses LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		//
		return courses
	}

	for rows.Next() {
		var course Course

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
		if err != nil {
			//
			return courses
		}

		courses = append(courses, course)
	}

	return courses
}

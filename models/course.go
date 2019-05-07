package models

import (
	"database/sql"
)

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Avatar      []byte `json:"avatar"`
	OwnerID     int    `json:"owner_id"`
}

type CourseCreateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
}

type CourseUpdateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ModelCourse struct {
	model
}

type ICourseLister interface {
	GetList(limit int, offset int) []Course
}

type ICourseGetter interface {
	Get(id int) Course
}

type ICourseCreator interface {
	Create(input CourseCreateInput) int
	ICourseGetter
}

type ICourseDeleter interface {
	Delete(id int)
}

type ICourseUpdater interface {
	Update(updatedCourse Course)
	ICourseGetter
}

func NewCourseModel(db *sql.DB) ModelCourse {
	return ModelCourse{model{db}}
}

func (m ModelCourse) GetList(limit int, offset int) []Course {
	courses := make([]Course, 0)

	rows, err := m.db.Query(`SELECT * FROM courses LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return courses
	}
	defer rows.Close()

	for rows.Next() {
		var course Course

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
		if err != nil {
			//
			return courses
		}

		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return courses
}

func (m ModelCourse) Get(id int) Course {
	course := Course{}

	row := m.db.QueryRow(`SELECT id, title, description, ava, owner_id FROM courses WHERE id = ?`, id)
	err := row.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
	if err != nil {
		return Course{}
	}

	return course
}

func (m ModelCourse) Create(in CourseCreateInput) int {
	stmt, err := m.db.Prepare(`INSERT INTO courses(title, description, ava, owner_id) VALUE(?, ?, ?, ?)`)
	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.Title, in.Description, []byte{}, in.OwnerID)
	if err != nil {
		return 0
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	return int(lastID)
}

func (m ModelCourse) Delete(id int) {
	_, err := m.db.Exec(`DELETE FROM courses WHERE id = ?`, id)
	if err != nil {
		return
	}
}

func (m ModelCourse) Update(in Course) {
	stmt, err := m.db.Prepare(`
		UPDATE courses SET
			title = ?,
			description  = ?,
		    ava = ?,
		    owner_id = ?
		WHERE id = ?`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(in.Title, in.Description, in.Avatar, in.OwnerID, in.ID)
}

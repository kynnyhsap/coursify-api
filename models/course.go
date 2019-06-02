package models

import (
	"database/sql"
	"log"
)

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
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
	GetList(limit, offset int, search string) []Course
	GetListForUser(limit, offset int, userID int64, admin bool) []Course
	CountForUser(userID int64, admin bool) int
	Count() int
}

type ICourseGetter interface {
	Get(id int64) Course
}

type ICourseCreator interface {
	Create(in CourseCreateInput) int64
	ICourseGetter
}

type ICourseDeleter interface {
	Delete(id int64)
}

type ICourseUpdater interface {
	Update(in Course)
	ICourseGetter
}

const (
	TypeAllCourses   = "all"
	TypeMyCourses    = "my"
	TypeAdminCourses = "admin"
)

func NewCourseModel(db *sql.DB) ModelCourse {
	return ModelCourse{model{db}}
}

func (m ModelCourse) GetList(limit, offset int, search string) []Course {
	courses := make([]Course, 0)

	rows, err := m.db.Query(`
		SELECT
		       *
		FROM courses
		WHERE title LIKE ?
		LIMIT ? OFFSET ?
	`, "%"+search+"%", limit, offset)
	if err != nil {
		return courses
	}
	defer rows.Close()

	for rows.Next() {
		var course Course

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.OwnerID, &course.Avatar)
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

func (m ModelCourse) GetListForUser(limit, offset int, userID int64, admin bool) []Course {
	courses := make([]Course, 0)

	var rows *sql.Rows
	var err error

	if admin {
		rows, err = m.db.Query(`
			SELECT
				id, title, description, owner_id, avatar
			FROM courses c
			WHERE c.owner_id = ?
			LIMIT ? OFFSET ?
		`, userID, limit, offset)
	} else {
		rows, err = m.db.Query(`
			SELECT
				id, title, description, owner_id, avatar
			FROM courses c
				JOIN students s ON c.id = s.course_id
			WHERE s.user_id = ?
			LIMIT ? OFFSET ?
		`, userID, limit, offset)
	}

	if err != nil {
		return courses
	}

	defer rows.Close()

	for rows.Next() {
		var course Course

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.OwnerID, &course.Avatar)
		if err != nil {
			return courses
		}

		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return courses
}

func (m ModelCourse) Count() int {
	rows, err := m.db.Query(`SELECT COUNT(*) FROM courses`)

	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Println(err)
			return 0
		}
	}

	return count
}

func (m ModelCourse) CountForUser(userID int64, admin bool) int {
	var rows *sql.Rows
	var err error

	if !admin {
		rows, err = m.db.Query(`
			SELECT
			   COUNT(*)
			FROM courses c JOIN students s ON c.id = s.course_id
			WHERE s.user_id = ?
		`, userID)
	} else {
		rows, err = m.db.Query(`
			SELECT 
				COUNT(*)
			FROM courses
			WHERE owner_id = ?
		`, userID)
	}

	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Println(err)
			return 0
		}
	}

	return count
}

func (m ModelCourse) Get(id int64) Course {
	course := Course{}

	row := m.db.QueryRow(`
		SELECT
		   id, title, description, avatar, owner_id
		FROM courses
		WHERE id = ?
	`, id)
	err := row.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
	if err != nil {
		return Course{}
	}

	return course
}

func (m ModelCourse) Create(in CourseCreateInput) int64 {
	stmt, err := m.db.Prepare(`
		INSERT INTO courses(
			title, description, avatar, owner_id
		) VALUE(?, ?, ?, ?)
	`)
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

	return lastID
}

func (m ModelCourse) Delete(id int64) {
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
		    avatar = ?,
		    owner_id = ?
		WHERE id = ?`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(in.Title, in.Description, in.Avatar, in.OwnerID, in.ID)
}

package models

import (
	"database/sql"
	"log"
)

type Mentor struct {
	User
	Role string `json:"role"`
}

type CourseDetail struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Avatar        string   `json:"avatar"`
	StudentsCount int      `json:"students_count"`
	OwnerID       int      `json:"owner_id"`
	Mentors       []Mentor `json:"mentors"`
	Entered       bool     `json:"entered"`
}

type CourseCreateInput struct {
	Avatar      string   `json:"avatar"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Mentors     []Mentor `json:"mentors"`
}

type CourseUpdateInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Mentors     []Mentor `json:"mentors"`
}

type ModelCourse struct {
	model
}

type ICourseLister interface {
	GetList(limit, offset int, userID int64, search string) []CourseDetail
	GetListForUser(limit, offset int, userID int64, admin bool) []CourseDetail
	CountForUser(userID int64, admin bool) int
	Count() int
}

type ICourseGetter interface {
	Get(id int64) CourseDetail
	Leave(courseID int64, userID int64)
	Enter(courseID int64, userID int64)
}

type ICourseCreator interface {
	Create(in CourseCreateInput, ownerID int64) int64
	ICourseGetter
}

type ICourseDeleter interface {
	Delete(id int64)
}

type ICourseUpdater interface {
	Update(in CourseDetail)
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

func (m ModelCourse) Entered(courseID int64, userID int64) bool {
	row := m.db.QueryRow(`
		SELECT
       		progress
		FROM students
		WHERE course_id = ? AND user_id = ?
	`, courseID, userID)

	var progress float64 = -1.0
	err := row.Scan(&progress)
	if err != nil {
		log.Println(err)
		return false
	}

	if progress < 0 {
		return false
	}
	return true
}

func (m ModelCourse) Enter(courseID int64, userID int64) {
	stmt, err := m.db.Prepare(`
		INSERT INTO students(
			course_id, user_id
		) VALUE(?, ?)
	`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(courseID, userID)
	if err != nil {
		return
	}
}

func (m ModelCourse) Leave(courseID int64, userID int64) {
	stmt, err := m.db.Prepare(`
		DELETE FROM students WHERE course_id = ? AND user_id = ?
	`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(courseID, userID)
	if err != nil {
		return
	}
}

func (m ModelCourse) GetList(limit, offset int, userID int64, search string) []CourseDetail {
	courses := make([]CourseDetail, 0)

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
		var course CourseDetail

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.OwnerID, &course.Avatar)
		if err != nil {
			//
			return courses
		}

		course.Entered = m.Entered(course.ID, userID)
		course.Mentors = m.GetMentorsList(course.ID)
		course.StudentsCount = m.CountStudents(course.ID)
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return courses
}

func (m ModelCourse) GetListForUser(limit, offset int, userID int64, admin bool) []CourseDetail {
	courses := make([]CourseDetail, 0)

	var rows *sql.Rows
	var err error

	if admin {
		rows, err = m.db.Query(`
			SELECT DISTINCT 
				c.id, c.title, c.description, c.owner_id, c.avatar
			FROM courses c
			LEFT JOIN mentors m ON c.id = m.course_id
			WHERE c.owner_id = ? OR m.user_id = ?
			LIMIT ? OFFSET ?
		`, userID, userID, limit, offset)
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
		var course CourseDetail

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.OwnerID, &course.Avatar)
		if err != nil {
			return courses
		}

		course.Entered = m.Entered(course.ID, userID)
		course.Mentors = m.GetMentorsList(course.ID)
		course.StudentsCount = m.CountStudents(course.ID)
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

func (m ModelCourse) CountStudents(courseID int64) int {
	rows, err := m.db.Query(`SELECT COUNT(*) FROM students WHERE course_id = ?`, courseID)

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

func (m ModelCourse) GetMentorsList(courseID int64) []Mentor {
	mentors := make([]Mentor, 0)

	rows, err := m.db.Query(`
		SELECT
			u.id, u.user_name, u.full_name, u.avatar, u.about, m.role
		FROM users u
		JOIN mentors m on u.id = m.user_id
		WHERE m.course_id = ?
	`, courseID)
	if err != nil {
		return mentors
	}
	defer rows.Close()

	for rows.Next() {
		var mentor Mentor

		err = rows.Scan(&mentor.ID, &mentor.Name, &mentor.FullName, &mentor.Avatar, &mentor.About, &mentor.Role)
		if err != nil {
			//
			return mentors
		}

		mentors = append(mentors, mentor)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return mentors
}

func (m ModelCourse) Get(id int64) CourseDetail {
	course := CourseDetail{}

	row := m.db.QueryRow(`
		SELECT
		   id, title, description, avatar, owner_id
		FROM courses
		WHERE id = ?
	`, id)
	err := row.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
	if err != nil {
		return CourseDetail{}
	}

	course.Mentors = m.GetMentorsList(id)
	course.StudentsCount = m.CountStudents(id)

	return course
}

func (m ModelCourse) Create(in CourseCreateInput, ownerID int64) int64 {
	stmt, err := m.db.Prepare(`
		INSERT INTO courses(
			title, description, avatar, owner_id
		) VALUE(?, ?, ?, ?)
	`)
	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.Title, in.Description, in.Avatar, ownerID)
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

func (m ModelCourse) Update(in CourseDetail) {
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

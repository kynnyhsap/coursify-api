package models

import "database/sql"

type Lesson struct {
	ID          int    `json:"id"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Theme       string `json:"theme"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
	CourseID    int    `json:"course_id"`
}

type LessonCreateInput struct {
	Title       string `json:"title"`
	Theme       string `json:"theme"`
	Description string `json:"description"`
	Image       string `json:"image"` // string because of gin json binding
	CourseID    int    `json:"course_id"`
}

type ModelLesson struct {
	model
}

type ILessonLister interface {
	GetList(courseID int64, limit, offset int) []Lesson
}

type ILessonGetter interface {
	Get(id int64) Lesson
}

type ILessonCreator interface {
	Create(in LessonCreateInput) int64
	ILessonGetter
}

type ILessonDeleter interface {
	Delete(id int64)
}

type ILessonUpdater interface {
	Update(in Lesson)
	ILessonGetter
}

func NewLessonModel(db *sql.DB) ModelLesson {
	return ModelLesson{model{db}}
}

func (m ModelLesson) GetList(courseID int64, limit, offset int) []Lesson {
	lessons := make([]Lesson, 0)

	rows, err := m.db.Query(`SELECT * FROM lessons LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return lessons
	}
	defer rows.Close()

	for rows.Next() {
		var lesson Lesson

		err = rows.Scan(&lesson.ID, &lesson.Title, &lesson.Theme, &lesson.Number, &lesson.Image, &lesson.CourseID)
		if err != nil {
			//
			return lessons
		}

		// TODO: select lesson components?

		lessons = append(lessons, lesson)
	}

	if err = rows.Err(); err != nil {
		// log error
	}

	return lessons
}

func (m ModelLesson) Get(id int64) Lesson {
	lesson := Lesson{}

	row := m.db.QueryRow(`SELECT id, title, theme, description, number, header_ava, course_id FROM lessons WHERE id = ?`, id)
	err := row.Scan(&lesson.ID, &lesson.Title, &lesson.Theme, &lesson.Description, &lesson.Number, &lesson.Image, &lesson.CourseID)
	if err != nil {
		return Lesson{}
	}

	// TODO: select lessons components?

	return lesson
}

func (m ModelLesson) Create(in LessonCreateInput) int64 {
	// TODO: define lesson number
	stmt, err := m.db.Prepare(`INSERT INTO lessons(number, title, theme, description, header_ava, course_id) VALUE(1, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0
	}

	res, err := stmt.Exec(in.Title, in.Theme, in.Description, []byte(in.Image), in.CourseID)
	if err != nil {
		return 0
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0
	}

	return lastID
}

func (m ModelLesson) Delete(id int64) {
	_, err := m.db.Exec(`DELETE FROM lessons WHERE id = ?`, id)
	if err != nil {
		return
	}
}

func (m ModelLesson) Update(in Lesson) {
	stmt, err := m.db.Prepare(`
		UPDATE lessons SET
			title = ?,
			theme  = ?,
			description  = ?,
		    header_ava = ?,
		    course_id = ?,
		    number = ?
		WHERE id = ?`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(in.Title, in.Theme, in.Description, in.Image, in.CourseID, in.Number, in.ID)
}

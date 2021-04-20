package storage

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("storage: resource not found")

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

type StudentSubject struct {
	ID            int
	CorrelativeID int
	Status        string
	Name          string
	Type          string
	Description   *string
}

const getStudentSubjects = `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = :email AND c.id = :careerID`

func (s *Storage) GetStudentSubjects(studentEmail, careerID string) ([]StudentSubject, error) {
	stmt, err := s.db.PrepareNamed(getStudentSubjects)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	params := map[string]interface{}{"email": studentEmail, "careerID": careerID}

	var studentSubjects []struct {
		ID            int64   `db:"subject_id"`
		CorrelativeID *int64  `db:"correlative_id"`
		Description   *string `db:"description"`
		Status        string  `db:"status"`
		Name          string  `db:"name"`
		Type          string  `db:"type"`
	}

	if err := stmt.Select(&studentSubjects, params); err != nil {
		return nil, err
	}

	if studentSubjects == nil {
		return nil, ErrNotFound
	}

	response := make([]StudentSubject, 0, len(studentSubjects))
	for _, studentSubject := range studentSubjects {
		var correlativeID int
		if studentSubject.CorrelativeID != nil {
			correlativeID = int(*studentSubject.CorrelativeID)
		}

		var description *string
		if studentSubject.Description != nil && *studentSubject.Description != "" {
			description = studentSubject.Description
		}

		response = append(response, StudentSubject{
			ID:            int(studentSubject.ID),
			CorrelativeID: correlativeID,
			Description:   description,
			Status:        studentSubject.Status,
			Name:          studentSubject.Name,
			Type:          studentSubject.Type,
		})
	}

	return response, nil
}

type SubjectDetails struct {
	ID     int
	Hours  int
	Points int
	Name   string
	Type   string
	URI    *string
	Meet   *string
}

const getSubjectDetails = `SELECT s.id, s.name, s.uri, s.meet, cs.type, cs.hours, cs.points
FROM career_subject cs
         INNER JOIN subject s ON cs.subject_id = s.id
WHERE s.id = :subjectID
  AND career_id = :careerID
LIMIT 1;`

func (s *Storage) GetSubjectDetails(subjectID, careerID string) (SubjectDetails, error) {
	stmt, err := s.db.PrepareNamed(getSubjectDetails)
	if err != nil {
		return SubjectDetails{}, err
	}

	defer stmt.Close()

	params := map[string]interface{}{"subjectID": subjectID, "careerID": careerID}

	var subjectDetails struct {
		ID     int64   `db:"id"`
		Hours  int64   `db:"hours"`
		Points int64   `db:"points"`
		Name   string  `db:"name"`
		Type   string  `db:"type"`
		URI    *string `db:"uri"`
		Meet   *string `db:"meet"`
	}

	if err := stmt.Get(&subjectDetails, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SubjectDetails{}, ErrNotFound
		}

		return SubjectDetails{}, err
	}

	return SubjectDetails{
		ID:     int(subjectDetails.ID),
		Hours:  int(subjectDetails.Hours),
		Points: int(subjectDetails.Points),
		Name:   subjectDetails.Name,
		Type:   subjectDetails.Type,
		URI:    subjectDetails.URI,
		Meet:   subjectDetails.Meet,
	}, nil
}

type Professorship struct {
	Day   int
	Name  string
	Start string
	End   string
}

const getProfessorships = `SELECT p.name, s.day, s.start, s.end
FROM professorship p
         INNER JOIN schedule s on p.id = s.professorship_id
         INNER JOIN career_subject cs on p.career_subject_id = cs.id
WHERE cs.subject_id = :subjectID AND cs.career_id = :careerID
ORDER BY day;`

func (s *Storage) GetProfessorships(subjectID, careerID string) ([]Professorship, error) {
	stmt, err := s.db.PrepareNamed(getProfessorships)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	params := map[string]interface{}{"subjectID": subjectID, "careerID": careerID}

	var professorships []struct {
		Day   int    `db:"day"`
		Name  string `db:"name"`
		Start string `db:"start"`
		End   string `db:"end"`
	}

	if err := stmt.Select(&professorships, params); err != nil {
		return nil, err
	}

	if professorships == nil {
		return nil, ErrNotFound
	}

	response := make([]Professorship, 0, len(professorships))
	for _, professorship := range professorships {
		response = append(response, Professorship{
			Day:   professorship.Day,
			Name:  professorship.Name,
			Start: professorship.Start,
			End:   professorship.End,
		})
	}

	return response, nil
}

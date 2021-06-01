package storage

import (
	"database/sql"
	"errors"
	"fmt"

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

const getStudentCareerIDs = `SELECT career_id FROM student_career sc
    INNER JOIN student s ON sc.student_id = s.id
WHERE s.email = :email;`

func (s *Storage) GetStudentCareerIDs(studentEmail string) ([]int, error) {
	stmt, err := s.db.PrepareNamed(getStudentCareerIDs)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	params := map[string]interface{}{"email": studentEmail}

	var ids []int64
	if err := stmt.Select(&ids, params); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	var careersIDs []int
	for _, id := range ids {
		careersIDs = append(careersIDs, int(id))
	}

	return careersIDs, nil
}

const (
	getStudentID            = `SELECT id FROM student WHERE email = ?;`
	findCareerWithID        = `SELECT COUNT(1) FROM career WHERE id = ?;`
	createStudentWithCareer = `INSERT INTO student_career (student_id, career_id) VALUES (?, ?);`
)

func (s *Storage) AssignStudentToCareer(studentEmail, careerID string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin tx: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var studentID int
	if err := tx.Get(&studentID, getStudentID, studentEmail); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("could not find student: %w", ErrNotFound)
		}

		return err
	}

	if studentID == 0 {
		return fmt.Errorf("could not find student: %w", ErrNotFound)
	}

	var careerCount int
	if err := tx.Get(&careerCount, findCareerWithID, careerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("could not find career: %w", ErrNotFound)
		}

		return err
	}

	if careerCount == 0 {
		return fmt.Errorf("could not find career: %w", ErrNotFound)
	}

	_, err = tx.Exec(createStudentWithCareer, studentID, careerID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit tx: %v", err)
	}

	return nil
}

func (s *Storage) UpdateStudentSubject(req UpdateStudentSubjectRequest) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin tx: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	studentID, err := s.getStudentByEmail(tx, req.StudentEmail)
	if err != nil {
		return err
	}

	if err := s.checkStudentAssignedToCareer(tx, studentID, req.CareerID); err != nil {
		return err
	}

	careerSubjectID, err := s.getCareerSubjectByIDs(tx, req.CareerID, req.SubjectID)
	if err != nil {
		return err
	}

	if err := s.updateStudentSubject(tx, studentID, careerSubjectID, req.Status, req.Description); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit tx: %v", err)
	}

	return nil
}

const getStudentByEmail = `SELECT id FROM student WHERE email = ?;`

func (s *Storage) getStudentByEmail(tx *sqlx.Tx, studentEmail string) (int, error) {
	var id int
	if err := tx.Get(&id, getStudentByEmail, studentEmail); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("could not find student: %w", ErrNotFound)
		}

		return 0, err
	}

	return id, nil
}

const checkStudentAssignedToCareer = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`

func (s *Storage) checkStudentAssignedToCareer(tx *sqlx.Tx, studentID int, careerID string) error {
	var results int
	if err := tx.Get(&results, checkStudentAssignedToCareer, studentID, careerID); err != nil {
		return err
	}

	if results == 0 {
		return fmt.Errorf("could not find student assigned to career: %w", ErrNotFound)
	}

	return nil
}

const getCareerSubjectByIDs = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`

func (s *Storage) getCareerSubjectByIDs(tx *sqlx.Tx, careerID, subjectID string) (int, error) {
	var id int
	if err := tx.Get(&id, getCareerSubjectByIDs, careerID, subjectID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("could not find career and subject: %w", ErrNotFound)
		}

		return 0, err
	}

	return id, nil
}

const updateStudentSubject = `INSERT INTO student_career_subject
    (student_id, career_subject_id, status, description)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE status      = ?,
                        description = ?;`

func (s *Storage) updateStudentSubject(tx *sqlx.Tx, studentID, careerSubjectID int, status string, description *string) error {
	if _, err := tx.Exec(updateStudentSubject, studentID, careerSubjectID, status, description, status, description); err != nil {
		return err
	}

	return nil
}

type UpdateStudentSubjectRequest struct {
	StudentEmail string
	CareerID     string
	SubjectID    string
	Status       string
	Description  *string
}

const getStudentSubjects = `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = :careerID
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = :email`

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
	Name   string
	Type   string
	URI    *string
	Meet   *string
	Hours  *int
	Points *int
}

const getSubjectDetails = `SELECT s.id, s.name, s.uri, s.meet, cs.type, cs.hours, cs.points
FROM career_subject cs
         INNER JOIN subject s ON cs.subject_id = s.id
WHERE s.id = :subjectID AND career_id = :careerID
LIMIT 1;`

func (s *Storage) GetSubjectDetails(subjectID, careerID string) (SubjectDetails, error) {
	stmt, err := s.db.PrepareNamed(getSubjectDetails)
	if err != nil {
		return SubjectDetails{}, err
	}

	defer stmt.Close()

	params := map[string]interface{}{"subjectID": subjectID, "careerID": careerID}

	var subjectDetails struct {
		ID     int     `db:"id"`
		Name   string  `db:"name"`
		Type   string  `db:"type"`
		URI    *string `db:"uri"`
		Meet   *string `db:"meet"`
		Hours  *int    `db:"hours"`
		Points *int    `db:"points"`
	}

	if err := stmt.Get(&subjectDetails, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SubjectDetails{}, ErrNotFound
		}

		return SubjectDetails{}, err
	}

	return SubjectDetails{
		ID:     subjectDetails.ID,
		Name:   subjectDetails.Name,
		Type:   subjectDetails.Type,
		URI:    subjectDetails.URI,
		Meet:   subjectDetails.Meet,
		Hours:  subjectDetails.Hours,
		Points: subjectDetails.Points,
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

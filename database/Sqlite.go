package database

import (
	"database/sql"
	"log"

	"github.com/IvoGoman/portalnotifier/util"
	_ "github.com/mattn/go-sqlite3"
)

const gradeTableDef = `create table if not exists grades(
                                    name text not null primary key,
                                    examid integer,
                                    semester text,
                                    trycountexam integer,
                                    date text,
                                    prof text,
                                    form text,
                                    grade float,
                                    bonus float,
                                    status text,
                                    trycountstudent string);`

func CreateDB() {
	db, err := sql.Open("sqlite3", "./portal.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(gradeTableDef)
	if err != nil {
		log.Fatal(err)
	}
}

// Store the grades retrieved from the portal
func StoreGrades(grades map[string]util.Module) {
	db, err := sql.Open("sqlite3", "./portal.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	insertGrade := `insert into grades (name, examid, semester, trycountexam,
                                        date, prof, form, grade, bonus, status,
                                         trycountstudent)
                                         values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(insertGrade)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, v := range grades {
		_, err := stmt.Exec(v.Name, v.ExamID, v.Semester, v.TryCountExam,
			v.Date, v.Prof, v.Form, v.Grade, v.Bonus, v.Status,
			v.TryCountStudent)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func SelectGrades() map[string]util.Module {
	grades := make(map[string]util.Module)
	var module util.Module
	db, err := sql.Open("sqlite3", "./portal.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM grades;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var examID int64
		var semester string
		var tryCountExam int64
		var date string
		var name string
		var prof string
		var form string
		var grade float64
		var bonus float64
		var status string
		var tryCountStudent string
		err = rows.Scan(&name, &examID, &semester, &tryCountExam, &date,
			&prof, &form, &grade, &bonus, &status,
			&tryCountStudent)
		if err != nil {
			log.Fatal(err)
		}
		module.ExamID = examID
		module.Semester = semester
		module.TryCountExam = tryCountExam
		module.Date = date
		module.Name = name
		module.Prof = prof
		module.Form = form
		module.Grade = grade
		module.Bonus = bonus
		module.Status = status
		module.TryCountStudent = tryCountStudent
		grades[module.Name] = module
	}
	return grades
}

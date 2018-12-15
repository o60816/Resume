package models

import (
	"fmt"
	"log"
)

type Work struct {
	Id          int       `json:"id"`
	Period      string    `json:"period"`
	Logo        string    `json:"logo"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	Content     string    `json:"content"`
	ProjectList []Project `json:"ProjectList"`
}

func GetAllWork(tblWork string) ([]Work, error) {
	rows, err := db.Query("SELECT id, period, logo, company, position, content FROM " + tblWork)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	workList := make([]Work, 0)

	for rows.Next() {
		var work Work
		if err := rows.Scan(&work.Id, &work.Period, &work.Logo, &work.Company, &work.Position, &work.Content); err != nil {
			return nil, err
		}
		workList = append(workList, work)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workList, nil
}

func GetWorkByID(tblWork string, workId string) (Work, error) {
	var work Work
	rows, err := db.Query(fmt.Sprintf("SELECT id, period, logo, company, position, content FROM %s WHERE id=%s", tblWork, workId))
	if err != nil {
		return work, err
	}

	for rows.Next() {
		if err := rows.Scan(&work.Id, &work.Period, &work.Logo, &work.Company, &work.Position, &work.Content); err != nil {
			return work, err
		}
	}

	if err = rows.Err(); err != nil {
		return work, err
	}

	return work, err
}

func EditWork(tblWork string, tblProject string, method string, work Work) (bool, error) {
	queryList := make([]string, 0)
	if method == "POST" {
		queryList = append(queryList, fmt.Sprintf("INSERT INTO %s(period, logo, company, position, content) VALUES('%s','%s','%s','%s','%s')", tblWork, work.Period, work.Logo, work.Company, work.Position, work.Content))
	} else if method == "PATCH" {
		queryList = append(queryList, fmt.Sprintf("UPDATE %s SET period='%s',logo='%s',company='%s',position='%s',content='%s' WHERE id='%d'", tblWork, work.Period, work.Logo, work.Company, work.Position, work.Content, work.Id))
	} else if method == "DELETE" {
		queryList = append(queryList, fmt.Sprintf("DELETE FROM %s WHERE id=%d", tblWork, work.Id))
		queryList = append(queryList, fmt.Sprintf("DELETE FROM %s WHERE workId=%d", tblProject, work.Id))
	}

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
	}()
	for _, query := range queryList {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return false, err
		}
	}
	err = tx.Commit()

	return true, nil
}

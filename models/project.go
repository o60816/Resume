package models

import (
	"fmt"
)

type Project struct {
	Id          int    `json:"id"`
	WorkId      int    `json:"workid"`
	ProjectName string `json:"ProjectName"`
	Tech        string `json:"tech"`
}

func GetProjectByWorkId(tblProject string, workId int) ([]Project, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, workId, projectName, tech FROM %s WHERE workId=%d", tblProject, workId))
	if err != nil {
		return nil, err
	}

	projectList := make([]Project, 0)
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.Id, &project.WorkId, &project.ProjectName, &project.Tech); err != nil {
			return nil, err
		}
		projectList = append(projectList, project)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projectList, nil
}

func GetProjectById(tblProject string, projectId string) (Project, error) {
	var project Project
	rows, err := db.Query(fmt.Sprintf("SELECT id, workId, projectName, tech FROM %s WHERE id=%s", tblProject, projectId))
	if err != nil {
		return project, err
	}
	for rows.Next() {
		if err := rows.Scan(&project.Id, &project.WorkId, &project.ProjectName, &project.Tech); err != nil {
			return project, err
		}
	}
	if err = rows.Err(); err != nil {
		return project, err
	}
	return project, nil
}

func EditProject(tblProject string, method string, project Project) (bool, error) {
	queryList := make([]string, 0)
	if method == "POST" {
		queryList = append(queryList, fmt.Sprintf("INSERT INTO %s(workId, projectName, tech) VALUES('%d','%s','%s')", tblProject, project.WorkId, project.ProjectName, project.Tech))
	} else if method == "PATCH" {
		queryList = append(queryList, fmt.Sprintf("UPDATE %s SET workId='%d',projectName='%s',tech='%s' WHERE id='%d'", tblProject, project.WorkId, project.ProjectName, project.Tech, project.Id))
	} else if method == "DELETE" {
		queryList = append(queryList, fmt.Sprintf("DELETE FROM %s WHERE id=%d", tblProject, project.Id))
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

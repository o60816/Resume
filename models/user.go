package models

import (
	"fmt"
	"log"
)

type User struct {
	Id           int
	Name         string
	Title        string
	Information  string
	Introduction string
	Motto        string
	Photo        string
	Account      string
	Password     string
	Github       string
	Linkedin     string
	Facebook     string
	Email        string
}

func Login(inputAccount string, inputPassword string) (User, bool) {
	var user User
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE account='%s'", inputAccount))
	if err != nil {
		log.Panic(err)
		return user, false
	}

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.Title, &user.Information, &user.Introduction, &user.Motto, &user.Photo, &user.Account, &user.Password, &user.Github, &user.Linkedin, &user.Facebook, &user.Email); err != nil {
			log.Panic(err)
			return user, false
		}
	}

	//New user register for him/her
	if user.Password == "" {
		result, err := db.Exec(fmt.Sprintf("INSERT INTO user(name, title, information, introduction, motto, photo, account, password, github, linkedin, facebook, email) VALUE('', '', '', '', '', '', '%s','%s','','','','')", inputAccount, inputPassword))
		if err != nil {
			return user, false
		}
		userId, err := result.LastInsertId()
		if err != nil {
			return user, false
		}
		user.Id = int(userId)
	} else {
		if inputPassword != user.Password {
			return user, false
		}
	}

	return user, true
}

func GetUserByName(userName string) (User, error) {
	var user User
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE name='%s'", userName))
	if err != nil {
		return user, err
	}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Title, &user.Information, &user.Introduction, &user.Motto, &user.Photo, &user.Account, &user.Password); err != nil {
			return user, err
		}
	}
	return user, err
}

func GetUserById(userId int) (User, error) {
	var user User
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE id='%d'", userId))
	if err != nil {
		return user, err
	}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Title, &user.Information, &user.Introduction, &user.Motto, &user.Photo, &user.Account, &user.Password, &user.Github, &user.Linkedin, &user.Facebook, &user.Email); err != nil {
			return user, err
		}
	}
	return user, err
}

func UpdateUser(user User) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE user SET name='%s', title='%s', information='%s', introduction='%s', motto='%s', photo='%s', github='%s', linkedin='%s', facebook='%s', email='%s' WHERE id='%d'",
		user.Name, user.Title, user.Information, user.Introduction, user.Motto, user.Photo, user.Github, user.Linkedin, user.Facebook, user.Email, user.Id))
	return err
}

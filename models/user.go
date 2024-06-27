package models

import (
	"errors"
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/utils"
)

type User struct {
	ID       string
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users(id, email,password) VALUES (?,?,?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.ID, u.Email, hashedPassword)
	return err
}

func (u *User) ValidateCredentials() (string, error) {
	query := "SELECT id, password FROM users where email = ?"

	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return "", errors.New("Bad Credentials!!")
	}

	passwordIsValid := utils.CheckPassword(u.Password, retrievedPassword)

	if !passwordIsValid {
		return "", errors.New("Bad Credentials!!")
	}
	return string(u.ID), nil
}

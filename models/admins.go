package models

import (
	"errors"
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/utils"
)

type Admin struct {
	ID       string
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (a Admin) Save() error {
	query := `INSERT INTO admins(id, email,password) VALUES (?,?,?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(a.Password)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(a.ID, a.Email, hashedPassword)
	return err
}

func (a *Admin) ValidateCredentials() (string, error) {
	query := "SELECT id, password FROM admins where email = ?"

	row := db.DB.QueryRow(query, a.Email)

	var retrievedPassword string
	err := row.Scan(&a.ID, &retrievedPassword)

	if err != nil {
		return "", errors.New("Bad Credentials!!")
	}

	passwordIsValid := utils.CheckPassword(a.Password, retrievedPassword)

	if !passwordIsValid {
		return "", errors.New("Bad Credentials!!")
	}
	return string(a.ID), nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, email FROM users`
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users = []User{}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (a *Admin) AdminCheck(id string) (bool, error) {
	query := `SELECT COUNT(*) FROM admins WHERE ID = ?`

	var count int
	err := db.DB.QueryRow(query, id).Scan(&count)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func DeleteAllUsers() error {
	query := `DELETE FROM users`

	result, err := db.DB.Exec(query)

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

func (admin Admin) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

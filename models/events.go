package models

import (
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/utils"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      *string
}

func (e *Event) Save() error {
	return utils.RetryOnBusy(func() error {
		query := `INSERT INTO events(name, location, description, dateTime, user_id) 
		VALUES (?, ?, ? ,? ,?)`

		if e.UserID == nil {
			newUUID := uuid.New().String()
			e.UserID = &newUUID
		}

		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		e.ID = id
		return err
	})
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Location, &event.Description, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Location, &event.Description, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	return utils.RetryOnBusy(func() error {
		query := `UPDATE events SET name = ?, description = ?, location = ?, dateTime = ? WHERE id = ?`

		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
		return err
	})
}

func (event Event) Delete() error {
	return utils.RetryOnBusy(func() error {
		query := "DELETE FROM events WHERE id = ?"

		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(event.ID)
		return err
	})
}

func DeleteAllEvents() error {
	query := `DELETE FROM events`

	result, err := db.DB.Exec(query)

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

func (event Event) Register(userID string) error {
	return utils.RetryOnBusy(func() error {
		query := "INSERT INTO registrations (event_id, user_id) VALUES (?,?)"
		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(event.ID, userID)
		return err
	})
}

func (event Event) CancelRegistration(userID string) error {
	return utils.RetryOnBusy(func() error {
		query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"

		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(event.ID, userID)
		return err
	})
}

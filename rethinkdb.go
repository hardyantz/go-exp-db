package main

import (
	"errors"
	"exp-traildb/stringrand"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type rSession struct {
	session   *rethinkdb.Session
	dbName    string
	tableName string
}

type Logs struct {
	ID     string `json:"id",rethinkdb:"id,omitempty"`
	Module string `json:"module",rethinkdb:"module"`
	Action string `json:"action",rethinkdb:"action"`
	Logs   []Log  `json:"logs",rethinkdb:"logs"`
	User   User   `json:"user",rethinkdb:"user"`
	Target string `json:"target",rethinkdb:"target"`
}

type User struct {
	Role string `json:"role",rethinkdb:"role"`
}

type Log struct {
	Field    string `json:"field",rethinkdb:"field"`
	OldValue string `json:"oldValue",rethinkdb:"oldValue"`
	NewValue string `json:"newValue",rethinkdb:"newValue"`
}

type RdbResponses struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {

	var rethinkSession rSession

	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address: "localhost",
	})

	rethinkSession.session = session
	rethinkSession.dbName = "test"
	rethinkSession.tableName = "testgin"

	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.POST("/create", func(c echo.Context) error {
		ok := rethinkSession.insertTable()
		if !ok {
			return c.JSON(http.StatusBadRequest, errors.New("failed to insert"))
		}

		return c.JSON(http.StatusCreated, RdbResponses{
			Status:  http.StatusCreated,
			Message: "Data successfully insert :)",
		})
	})

	e.GET("/get", func(c echo.Context) error {
		res, err := rethinkSession.fetchAllRecord()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, res)
	})

	e.GET("/get/:id", func(c echo.Context) error {
		var err error

		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, RdbResponses{
				Status:  http.StatusBadRequest,
				Message: "ID required",
			})
		}

		res, err := rethinkSession.getOne(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, RdbResponses{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":1111"))
}

func (r *rSession) insertTable() bool {
	data := Logs{
		Module: "MEMBERSHIP",
		Action: "UPDATE",
		Logs: []Log{
			{
				Field:    "firstName",
				OldValue: stringrand.String(10),
				NewValue: stringrand.String(5),
			},
		},
		User:   User{Role: "Admin"},
		Target: "USR" + strconv.Itoa(rand.Int()),
	}

	_, err := rethinkdb.DB(r.dbName).Table(r.tableName).Insert(data).RunWrite(r.session)
	if err != nil {
		return false
	}

	return true
}

func (r *rSession) getOne(id string) (Logs, error) {
	var logs Logs
	cursor, err := rethinkdb.DB(r.dbName).Table(r.tableName).Get(id).Run(r.session)
	if err != nil {
		return logs, err
	}

	cursor.One(&logs)
	cursor.Close()

	if logs.ID == "" {
		return logs, errors.New("data not found")
	}

	return logs, nil
}

func (r *rSession) delete(id string) error {
	err := rethinkdb.DB(r.dbName).Table(r.tableName).Get(id).Delete().Exec(r.session)
	if err != nil {
		return err
	}

	return nil
}

func (r *rSession) fetchAllRecord() ([]Logs, error) {
	var result []Logs
	rows, err := rethinkdb.Table(r.tableName).Run(r.session)
	if err != nil {
		return result, err
	}

	// Read records into persons slice
	var logs []Logs
	err2 := rows.All(&logs)
	if err2 != nil {
		fmt.Println(err2)
		return result, err
	}

	for _, p := range logs {
		result = append(result, p)
	}

	return result, nil
}

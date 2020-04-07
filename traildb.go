package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
	tdb "github.com/traildb/traildb-go"
)

type PostData struct {
	Timestamp int    `tdb:"timestamp"`
	IpAddress string `json:"ip",tdb:"ip"`
	Title     string `json:"title",tdb:"title"`
	User      string `json:"user",tdb:"user"`
}

type TrailResponses struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const SessionLimit = uint64(30 * 60)
const tdbPath = "data/pydata-tutorial.tdb"
const tdbNewPath = "data/mydata"

func main() {
	e := echo.New()

	e.GET("/get", func(c echo.Context) error {
		q := c.QueryString()
		fmt.Println(q)

		res, err := getAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, TrailResponses{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, res)
	})

	e.GET("/get-wiki", func(c echo.Context) error {
		q := c.QueryString()
		fmt.Println(q)

		res, err := getWiki()
		if err != nil {
			return c.JSON(http.StatusBadRequest, TrailResponses{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, res)
	})

	e.POST("/create", func(c echo.Context) error {
		u := new(PostData)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, TrailResponses{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		err := create(u)
		if err != nil {
			return c.JSON(http.StatusBadRequest, TrailResponses{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, TrailResponses{
			Status:  http.StatusCreated,
			Message: "Data successfully insert :)",
		})
	})

	e.Logger.Fatal(e.Start(":1212"))
}

func create(p *PostData) error {
	db, err := tdb.Open(tdbNewPath)
	if err != nil {
		return err
	}
	defer db.Close()
	cookie := random.String(32, "1234567890")
	cons, err := tdb.NewTrailDBConstructor(tdbNewPath, "ip", "title", "user")
	if err != nil {
		return err
	}

	defer cons.Close()

	err = cons.Add(cookie, time.Now().Unix(), []string{p.IpAddress, p.Title, p.User})
	if err != nil {
		return err
	}

	err = cons.Append(db)
	if err != nil {
		return err
	}

	err = cons.Finalize()
	if err != nil {
		return err
	}

	return nil
}

func getAll() ([]map[string]string, error) {

	var responses []map[string]string

	db, err := tdb.Open(tdbNewPath)
	if err != nil {
		return responses, err
	}
	defer db.Close()

	for i := uint64(0); i < db.NumTrails; i++ {
		trail, err := tdb.NewTrail(db, i)
		if err != nil {
			panic(err.Error())
		}
		for {
			evt := trail.NextEvent()
			if evt == nil {
				trail.Close()
				break
			}

			s := evt.ToMap()
			s["timestamp"] = strconv.FormatUint(evt.Timestamp, 10)

			responses = append(responses, s)
		}
	}

	return responses, nil
}

func getWiki() ([]map[string]string, error) {

	var responses []map[string]string

	db, err := tdb.Open(tdbPath)
	if err != nil {
		return responses, err
	}
	defer db.Close()

	for i := uint64(0); i < db.NumTrails; i++ {
		trail, err := tdb.NewTrail(db, i)
		if err != nil {
			panic(err.Error())
		}
		for {
			evt := trail.NextEvent()
			if evt == nil {
				trail.Close()
				break
			}

			s := evt.ToMap()
			s["timestamp"] = strconv.FormatUint(evt.Timestamp, 10)

			responses = append(responses, s)
		}
	}

	return responses, nil
}

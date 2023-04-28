package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

// migrate runs a Sqlite migration script
func migrate(db *sql.DB) error {
	if db == nil {
		return errors.New("nil DB connection passed, please fix")
	}
	query := `
	CREATE TABLE IF NOT EXISTS events(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    date TEXT NOT NULL,
	    name TEXT NOT NULL,
	    description TEXT
	);
	CREATE TABLE IF NOT EXISTS event_lang(
	    event_id INTEGER,
	    language TEXT NOT NULL,
	    FOREIGN KEY(event_id) REFERENCES events(id)
	);
	CREATE TABLE IF NOT EXISTS event_video_quality(
	    event_id INTEGER,
	    quality TEXT NOT NULL,
		FOREIGN KEY(event_id) REFERENCES events(id)
	);
	CREATE TABLE IF NOT EXISTS event_audio_quality(
	    event_id INTEGER,
	    quality TEXT NOT NULL,
	    FOREIGN KEY(event_id) REFERENCES events(id)
	);
	CREATE TABLE IF NOT EXISTS event_invitees (
	    event_id INTEGER,
	    Email TEXT NOT NULL,
	    FOREIGN KEY(event_id) REFERENCES events(id)
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error on db schema migration %w", err)
	}
	return nil
}

// initDB uses the db driver to open a connection to a sqlite db file
func initDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	return db, err
}

// event Data Model
// The event data model is used for our internal representation of an event, for purposes of creating an event.
type event struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name" validate:"required"`
	Date         string   `json:"date" validate:"required"`
	Languages    []string `json:"languages" validate:"required"`
	VideoQuality []string `json:"VideoQuality" validate:"required"`
	AudioQuality []string `json:"AudioQuality" validate:"required"`
	Invitees     []string `json:"invitees" validate:"required"`
	Description  string   `json:"description,omitempty"`
}

// For the response types the invitees field is omitted as it seems inappropriate to return a list of user emails to any
// user who requests the events.

// eventResponse is the response type used for retrieving an individual event.
type eventResponse struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Languages    []string `json:"language"`
	VideoQuality string   `json:"videoQuality"`
	AudioQuality string   `json:"audioQuality"`
	Description  string   `json:"description,omitempty"`
}

// eventsResponse is the response type used for retrieving all available events.
type eventsResponse struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Languages    []string `json:"language"`
	VideoQuality []string `json:"videoQuality"`
	AudioQuality []string `json:"audioQuality"`
	Description  string   `json:"description,omitempty"`
}

// eventQueryStringFields is a struct used to store and access the query string data internally.
type eventQueryStringFields struct {
	VideoQuality string `json:"VideoQuality" query:"videoQuality"`
	AudioQuality string `json:"AudioQuality" query:"audioQuality"`
}

// Event Repository

// initRepo is a convenience function used at the server startup to allow for access to the db.
func initRepo() error {
	db, err := initDB("sqlite.db")
	if err != nil {
		return fmt.Errorf("error opening database: %w\n", err)
	}
	if err := migrate(db); err != nil {
		return err
	}
	newRepo(db)
	return nil
}

type Repo struct {
	db *sql.DB
}

var EventRepo *Repo

// newRepo initializes the EventRepo used by requests to create and access events in the db.
func newRepo(db *sql.DB) *Repo {
	EventRepo = &Repo{
		db: db,
	}
	return EventRepo
}

// create an event
func (r *Repo) create(ev event) (*event, error) {
	id, err := r.createEvent(ev.Name, ev.Date, ev.Description)
	if err != nil {
		return nil, err
	}
	ev.Id = id
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(id int64, languages []string) {
		defer wg.Done()
		err := r.createEventLang(id, languages)
		if err != nil {
			panic(err)
		}
	}(id, ev.Languages)
	wg.Add(1)
	go func(id int64, quality []string) {
		defer wg.Done()
		err := r.createEventVideoQuality(id, quality)
		if err != nil {
			panic(err)
		}
	}(id, ev.VideoQuality)
	wg.Add(1)
	go func(id int64, quality []string) {
		defer wg.Done()
		err := r.createEventAudioQuality(id, quality)
		if err != nil {
			panic(err)
		}
	}(id, ev.AudioQuality)
	wg.Add(1)
	go func(id int64, invitees []string) {
		defer wg.Done()
		err := r.createEventInvitees(id, invitees)
		if err != nil {
			panic(err)
		}
	}(id, ev.Invitees)
	wg.Wait()
	return &ev, nil
}

// Helper functions which create entries into the corresponding database tables.

func (r *Repo) createEvent(name, date, description string) (int64, error) {
	res, err := r.db.Exec("INSERT INTO events(name, date, description) VALUES (?, ? ,?)", name, date, description)
	if err != nil {
		err := fmt.Errorf("error inserting into events table: %w\n", err)
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repo) createEventLang(id int64, languages []string) error {
	stmt, err := r.db.Prepare("INSERT INTO event_lang(event_id, language) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, l := range languages {
		_, err := stmt.Exec(id, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) createEventVideoQuality(id int64, quality []string) error {
	stmt, err := r.db.Prepare("INSERT INTO event_video_quality(event_id, quality) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, q := range quality {
		_, err := stmt.Exec(id, q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) createEventAudioQuality(id int64, quality []string) error {
	stmt, err := r.db.Prepare("INSERT INTO event_audio_quality(event_id, quality) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, q := range quality {
		_, err := stmt.Exec(id, q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) createEventInvitees(id int64, invitees []string) error {
	stmt, err := r.db.Prepare("INSERT INTO event_invitees(event_id, Email) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, i := range invitees {
		_, err := stmt.Exec(id, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get an event
func (r *Repo) getOneEvent(id int64, qs eventQueryStringFields) (*eventResponse, error) {
	row := r.db.QueryRow("SELECT id, name, date, description FROM events WHERE id = ?", id)

	var res eventResponse
	// Get event id, name, date, description
	if err := row.Scan(&res.Id, &res.Name, &res.Date, &res.Description); err != nil {
		return nil, err
	}

	// Get event languages
	langRows, err := r.db.Query("SELECT language FROM event_lang WHERE event_id = ?", id)
	if err != nil {
		return nil, err
	}
	for langRows.Next() {
		var language string
		if err := langRows.Scan(&language); err != nil {
			return nil, err
		}
		res.Languages = append(res.Languages, language)
	}

	// Get event video and audio quality
	videoRows, err := r.db.Query("SELECT quality FROM event_video_quality WHERE event_id = ? ", id)
	if err != nil {
		return nil, err
	}
	for videoRows.Next() {
		var videoq string
		if err := videoRows.Scan(&videoq); err != nil {
			return nil, err
		}
		if videoq == qs.VideoQuality {
			res.VideoQuality = videoq
			break
		} else {
			res.VideoQuality = ""
		}

	}
	if res.VideoQuality == "" {
		res.VideoQuality = defaultVideoQuality
	}
	audioRows, err := r.db.Query("SELECT quality FROM event_audio_quality WHERE event_id = ?", id)
	if err != nil {
		return nil, err
	}
	for audioRows.Next() {
		var audioq string
		if err := audioRows.Scan(&audioq); err != nil {
			return nil, err
		}
		if audioq == qs.AudioQuality {
			res.AudioQuality = audioq
			break
		} else {
			res.AudioQuality = ""
		}
	}
	if res.AudioQuality == "" {
		res.AudioQuality = defaultAudioQuality
	}

	return &res, nil
}

// Get all available events
func (r *Repo) getAllEvents() ([]eventsResponse, error) {
	rows, err := r.db.Query("SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []eventsResponse
	for rows.Next() {
		var event eventsResponse
		// Get event id, name, date, description
		if err := rows.Scan(&event.Id, &event.Name, &event.Date, &event.Description); err != nil {
			return nil, err
		}
		// Get event languages
		langRows, err := r.db.Query("SELECT language FROM event_lang WHERE event_id = ?", event.Id)
		if err != nil {
			return nil, err
		}
		event.Languages = make([]string, 0)
		for langRows.Next() {
			var language string
			if err := langRows.Scan(&language); err != nil {
				return nil, err
			}
			event.Languages = append(event.Languages, language)
		}
		// Get event audio and video quality
		audioRows, err := r.db.Query("SELECT quality FROM event_audio_quality WHERE event_id = ?", event.Id)
		if err != nil {
			return nil, err
		}
		event.AudioQuality = make([]string, 0)
		for audioRows.Next() {
			var audioq string
			if err := audioRows.Scan(&audioq); err != nil {
				return nil, err
			}
			event.AudioQuality = append(event.AudioQuality, audioq)
		}
		videoRows, err := r.db.Query("SELECT quality FROM event_video_quality WHERE event_id = ? ", event.Id)
		if err != nil {
			return nil, err
		}
		event.VideoQuality = make([]string, 0)
		for videoRows.Next() {
			var videoq string
			if err := videoRows.Scan(&videoq); err != nil {
				return nil, err
			}
			event.VideoQuality = append(event.VideoQuality, videoq)
		}
		all = append(all, event)
	}
	return all, nil
}

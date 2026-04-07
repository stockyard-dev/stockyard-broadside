package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type PressRelease struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Contact     string `json:"contact"`
	Status      string `json:"status"`
	PublishDate string `json:"publish_date"`
	Outlet      string `json:"outlet"`
	Tags        string `json:"tags"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "broadside.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS press_releases(id TEXT PRIMARY KEY,title TEXT NOT NULL,body TEXT DEFAULT '',contact TEXT DEFAULT '',status TEXT DEFAULT 'draft',publish_date TEXT DEFAULT '',outlet TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
	resource TEXT NOT NULL,
	record_id TEXT NOT NULL,
	data TEXT NOT NULL DEFAULT '{}',
	PRIMARY KEY(resource, record_id)
)`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(e *PressRelease) error {
	e.ID = genID()
	e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO press_releases(id,title,body,contact,status,publish_date,outlet,tags,created_at)VALUES(?,?,?,?,?,?,?,?,?)`, e.ID, e.Title, e.Body, e.Contact, e.Status, e.PublishDate, e.Outlet, e.Tags, e.CreatedAt)
	return err
}
func (d *DB) Get(id string) *PressRelease {
	var e PressRelease
	if d.db.QueryRow(`SELECT id,title,body,contact,status,publish_date,outlet,tags,created_at FROM press_releases WHERE id=?`, id).Scan(&e.ID, &e.Title, &e.Body, &e.Contact, &e.Status, &e.PublishDate, &e.Outlet, &e.Tags, &e.CreatedAt) != nil {
		return nil
	}
	return &e
}
func (d *DB) List() []PressRelease {
	rows, _ := d.db.Query(`SELECT id,title,body,contact,status,publish_date,outlet,tags,created_at FROM press_releases ORDER BY created_at DESC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []PressRelease
	for rows.Next() {
		var e PressRelease
		rows.Scan(&e.ID, &e.Title, &e.Body, &e.Contact, &e.Status, &e.PublishDate, &e.Outlet, &e.Tags, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}
func (d *DB) Update(e *PressRelease) error {
	_, err := d.db.Exec(`UPDATE press_releases SET title=?,body=?,contact=?,status=?,publish_date=?,outlet=?,tags=? WHERE id=?`, e.Title, e.Body, e.Contact, e.Status, e.PublishDate, e.Outlet, e.Tags, e.ID)
	return err
}
func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM press_releases WHERE id=?`, id)
	return err
}
func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM press_releases`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []PressRelease {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ? OR body LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(`SELECT id,title,body,contact,status,publish_date,outlet,tags,created_at FROM press_releases WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []PressRelease
	for rows.Next() {
		var e PressRelease
		rows.Scan(&e.ID, &e.Title, &e.Body, &e.Contact, &e.Status, &e.PublishDate, &e.Outlet, &e.Tags, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Stats() map[string]any {
	m := map[string]any{"total": d.Count()}
	rows, _ := d.db.Query(`SELECT status,COUNT(*) FROM press_releases GROUP BY status`)
	if rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}
	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}

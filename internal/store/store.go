package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Feedback is a single feature request, idea, or bug report. Status flow
// is typically: open → planned → in_progress → completed → closed.
type Feedback struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Author    string `json:"author"`
	Category  string `json:"category"`
	Votes     int    `json:"votes"`
	Status    string `json:"status"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "podium.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS feedback(
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT DEFAULT '',
		author TEXT DEFAULT '',
		category TEXT DEFAULT '',
		votes INTEGER DEFAULT 0,
		status TEXT DEFAULT 'open',
		tags TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feedback_status ON feedback(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feedback_category ON feedback(category)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feedback_votes ON feedback(votes)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Feedback) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "open"
	}
	_, err := d.db.Exec(
		`INSERT INTO feedback(id, title, body, author, category, votes, status, tags, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Title, e.Body, e.Author, e.Category, e.Votes, e.Status, e.Tags, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Feedback {
	var e Feedback
	err := d.db.QueryRow(
		`SELECT id, title, body, author, category, votes, status, tags, created_at
		 FROM feedback WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Title, &e.Body, &e.Author, &e.Category, &e.Votes, &e.Status, &e.Tags, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

// List returns feedback ordered by votes descending (most popular first),
// matching the typical feedback board UX.
func (d *DB) List() []Feedback {
	rows, _ := d.db.Query(
		`SELECT id, title, body, author, category, votes, status, tags, created_at
		 FROM feedback ORDER BY votes DESC, created_at DESC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Feedback
	for rows.Next() {
		var e Feedback
		rows.Scan(&e.ID, &e.Title, &e.Body, &e.Author, &e.Category, &e.Votes, &e.Status, &e.Tags, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Feedback) error {
	_, err := d.db.Exec(
		`UPDATE feedback SET title=?, body=?, author=?, category=?, votes=?, status=?, tags=?
		 WHERE id=?`,
		e.Title, e.Body, e.Author, e.Category, e.Votes, e.Status, e.Tags, e.ID,
	)
	return err
}

// Upvote atomically increments the votes counter. Avoids the race that
// would occur if the dashboard read-modified-wrote via the regular
// update endpoint.
func (d *DB) Upvote(id string) (int, error) {
	_, err := d.db.Exec(`UPDATE feedback SET votes = votes + 1 WHERE id=?`, id)
	if err != nil {
		return 0, err
	}
	var v int
	d.db.QueryRow(`SELECT votes FROM feedback WHERE id=?`, id).Scan(&v)
	return v, nil
}

// Downvote decrements votes but never below zero.
func (d *DB) Downvote(id string) (int, error) {
	_, err := d.db.Exec(`UPDATE feedback SET votes = MAX(votes - 1, 0) WHERE id=?`, id)
	if err != nil {
		return 0, err
	}
	var v int
	d.db.QueryRow(`SELECT votes FROM feedback WHERE id=?`, id).Scan(&v)
	return v, nil
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM feedback WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM feedback`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Feedback {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ? OR body LIKE ? OR tags LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, title, body, author, category, votes, status, tags, created_at
		 FROM feedback WHERE `+where+`
		 ORDER BY votes DESC, created_at DESC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Feedback
	for rows.Next() {
		var e Feedback
		rows.Scan(&e.ID, &e.Title, &e.Body, &e.Author, &e.Category, &e.Votes, &e.Status, &e.Tags, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns aggregate metrics: total feedback items, total votes
// across all items, by_status and by_category breakdowns, plus the
// top item by vote count for a quick "most wanted" indicator.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":       d.Count(),
		"total_votes": 0,
		"by_status":   map[string]int{},
		"by_category": map[string]int{},
		"top_votes":   0,
	}

	var totalVotes int
	d.db.QueryRow(`SELECT COALESCE(SUM(votes), 0) FROM feedback`).Scan(&totalVotes)
	m["total_votes"] = totalVotes

	var top int
	d.db.QueryRow(`SELECT COALESCE(MAX(votes), 0) FROM feedback`).Scan(&top)
	m["top_votes"] = top

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM feedback GROUP BY status`); rows != nil {
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

	if rows, _ := d.db.Query(`SELECT category, COUNT(*) FROM feedback WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_category"] = by
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

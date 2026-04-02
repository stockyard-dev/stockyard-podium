package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Feedback struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	Author string `json:"author"`
	Category string `json:"category"`
	Votes int `json:"votes"`
	Status string `json:"status"`
	Tags string `json:"tags"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"podium.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS feedback(id TEXT PRIMARY KEY,title TEXT NOT NULL,body TEXT DEFAULT '',author TEXT DEFAULT '',category TEXT DEFAULT '',votes INTEGER DEFAULT 0,status TEXT DEFAULT 'open',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Feedback)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO feedback(id,title,body,author,category,votes,status,tags,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Body,e.Author,e.Category,e.Votes,e.Status,e.Tags,e.CreatedAt);return err}
func(d *DB)Get(id string)*Feedback{var e Feedback;if d.db.QueryRow(`SELECT id,title,body,author,category,votes,status,tags,created_at FROM feedback WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Body,&e.Author,&e.Category,&e.Votes,&e.Status,&e.Tags,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Feedback{rows,_:=d.db.Query(`SELECT id,title,body,author,category,votes,status,tags,created_at FROM feedback ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Feedback;for rows.Next(){var e Feedback;rows.Scan(&e.ID,&e.Title,&e.Body,&e.Author,&e.Category,&e.Votes,&e.Status,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM feedback WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM feedback`).Scan(&n);return n}

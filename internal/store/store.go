package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Feedback struct{ID int64 `json:"id"`;Page string `json:"page"`;Category string `json:"category"`;Body string `json:"body"`;Rating int `json:"rating"`;Status string `json:"status"`;CreatedAt time.Time `json:"created_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"podium.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS feedback(id INTEGER PRIMARY KEY AUTOINCREMENT,page TEXT DEFAULT '/',category TEXT DEFAULT 'general',body TEXT NOT NULL,rating INTEGER DEFAULT 0,status TEXT DEFAULT 'new',created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)Create(f *Feedback)error{res,err:=db.Exec(`INSERT INTO feedback(page,category,body,rating)VALUES(?,?,?,?)`,f.Page,f.Category,f.Body,f.Rating);if err!=nil{return err};f.ID,_=res.LastInsertId();return nil}
func(db *DB)List(status string)([]Feedback,error){q:=`SELECT id,page,category,body,rating,status,created_at FROM feedback WHERE 1=1`;args:=[]interface{}{};if status!=""{q+=` AND status=?`;args=append(args,status)};q+=` ORDER BY created_at DESC`;rows,err:=db.Query(q,args...);if err!=nil{return nil,err};defer rows.Close();var out[]Feedback;for rows.Next(){var f Feedback;rows.Scan(&f.ID,&f.Page,&f.Category,&f.Body,&f.Rating,&f.Status,&f.CreatedAt);out=append(out,f)};return out,nil}
func(db *DB)UpdateStatus(id int64,status string){db.Exec(`UPDATE feedback SET status=? WHERE id=?`,status,id)}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM feedback WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var total,newCount int;var avgRating float64;db.QueryRow(`SELECT COUNT(*) FROM feedback`).Scan(&total);db.QueryRow(`SELECT COUNT(*) FROM feedback WHERE status='new'`).Scan(&newCount);db.QueryRow(`SELECT COALESCE(AVG(rating),0) FROM feedback WHERE rating>0`).Scan(&avgRating);return map[string]interface{}{"total":total,"new":newCount,"avg_rating":avgRating},nil}

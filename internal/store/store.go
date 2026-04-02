package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type PressRelease struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	Contact string `json:"contact"`
	Status string `json:"status"`
	PublishDate string `json:"publish_date"`
	Outlet string `json:"outlet"`
	Tags string `json:"tags"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"broadside.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS press_releases(id TEXT PRIMARY KEY,title TEXT NOT NULL,body TEXT DEFAULT '',contact TEXT DEFAULT '',status TEXT DEFAULT 'draft',publish_date TEXT DEFAULT '',outlet TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *PressRelease)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO press_releases(id,title,body,contact,status,publish_date,outlet,tags,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Body,e.Contact,e.Status,e.PublishDate,e.Outlet,e.Tags,e.CreatedAt);return err}
func(d *DB)Get(id string)*PressRelease{var e PressRelease;if d.db.QueryRow(`SELECT id,title,body,contact,status,publish_date,outlet,tags,created_at FROM press_releases WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Body,&e.Contact,&e.Status,&e.PublishDate,&e.Outlet,&e.Tags,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]PressRelease{rows,_:=d.db.Query(`SELECT id,title,body,contact,status,publish_date,outlet,tags,created_at FROM press_releases ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []PressRelease;for rows.Next(){var e PressRelease;rows.Scan(&e.ID,&e.Title,&e.Body,&e.Contact,&e.Status,&e.PublishDate,&e.Outlet,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM press_releases WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM press_releases`).Scan(&n);return n}

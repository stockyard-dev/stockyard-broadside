package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Template struct{ID int64 `json:"id"`;Name string `json:"name"`;HTML string `json:"html"`;Width int `json:"width"`;Height int `json:"height"`;CreatedAt time.Time `json:"created_at"`}
func Open(dataDir string)(*DB,error){if err:=os.MkdirAll(dataDir,0755);err!=nil{return nil,fmt.Errorf("mkdir: %w",err)};dsn:=filepath.Join(dataDir,"broadside.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);if err:=migrate(db);err!=nil{return nil,fmt.Errorf("migrate: %w",err)};return &DB{db},nil}
func migrate(db *sql.DB)error{_,err:=db.Exec(`CREATE TABLE IF NOT EXISTS templates(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL,html TEXT NOT NULL DEFAULT '',width INTEGER DEFAULT 1200,height INTEGER DEFAULT 630,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);CREATE TABLE IF NOT EXISTS render_log(id INTEGER PRIMARY KEY AUTOINCREMENT,template_id INTEGER,rendered_at DATETIME DEFAULT CURRENT_TIMESTAMP);`);return err}
func(db *DB)ListTemplates()([]Template,error){rows,err:=db.Query(`SELECT id,name,html,width,height,created_at FROM templates ORDER BY created_at DESC`);if err!=nil{return nil,err};defer rows.Close();var out[]Template;for rows.Next(){var t Template;rows.Scan(&t.ID,&t.Name,&t.HTML,&t.Width,&t.Height,&t.CreatedAt);out=append(out,t)};return out,nil}
func(db *DB)CreateTemplate(t *Template)error{res,err:=db.Exec(`INSERT INTO templates(name,html,width,height)VALUES(?,?,?,?)`,t.Name,t.HTML,t.Width,t.Height);if err!=nil{return err};t.ID,_=res.LastInsertId();return nil}
func(db *DB)GetTemplate(id int64)(*Template,error){t:=&Template{};err:=db.QueryRow(`SELECT id,name,html,width,height,created_at FROM templates WHERE id=?`,id).Scan(&t.ID,&t.Name,&t.HTML,&t.Width,&t.Height,&t.CreatedAt);if err==sql.ErrNoRows{return nil,nil};return t,err}
func(db *DB)UpdateTemplate(t *Template)error{_,err:=db.Exec(`UPDATE templates SET name=?,html=?,width=?,height=? WHERE id=?`,t.Name,t.HTML,t.Width,t.Height,t.ID);return err}
func(db *DB)DeleteTemplate(id int64)error{_,err:=db.Exec(`DELETE FROM templates WHERE id=?`,id);return err}
func(db *DB)LogRender(tmplID int64){db.Exec(`INSERT INTO render_log(template_id)VALUES(?)`,tmplID)}
func(db *DB)CountTemplates()(int,error){var n int;db.QueryRow(`SELECT COUNT(*) FROM templates`).Scan(&n);return n,nil}
func(db *DB)CountRenders()(int,error){var n int;db.QueryRow(`SELECT COUNT(*) FROM render_log`).Scan(&n);return n,nil}

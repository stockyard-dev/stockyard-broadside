package server
import("encoding/json";"fmt";"net/http";"strconv";"strings";"github.com/stockyard-dev/stockyard-broadside/internal/store")
func(s *Server)handleListTemplates(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListTemplates();if list==nil{list=[]store.Template{}};writeJSON(w,200,list)}
func(s *Server)handleGetTemplate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);t,_:=s.db.GetTemplate(id);if t==nil{writeError(w,404,"not found");return};writeJSON(w,200,t)}
func(s *Server)handleCreateTemplate(w http.ResponseWriter,r *http.Request){
    if !s.limits.IsPro(){n,_:=s.db.CountTemplates();if n>=3{writeError(w,403,"free tier: 3 templates max");return}}
    var t store.Template;json.NewDecoder(r.Body).Decode(&t)
    if t.Name==""{writeError(w,400,"name required");return}
    if t.HTML==""{t.HTML=defaultHTML(t.Name)};if t.Width==0{t.Width=1200};if t.Height==0{t.Height=630}
    if err:=s.db.CreateTemplate(&t);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,t)}
func(s *Server)handleUpdateTemplate(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    existing,_:=s.db.GetTemplate(id);if existing==nil{writeError(w,404,"not found");return}
    json.NewDecoder(r.Body).Decode(existing)
    existing.ID=id;s.db.UpdateTemplate(existing);writeJSON(w,200,existing)}
func(s *Server)handleDeleteTemplate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteTemplate(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleRender(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    t,_:=s.db.GetTemplate(id);if t==nil{writeError(w,404,"template not found");return}
    html:=t.HTML
    for k,vs:=range r.URL.Query(){if len(vs)>0{html=strings.ReplaceAll(html,"{{"+k+"}}",vs[0])}}
    s.db.LogRender(id)
    svg:=fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d"><foreignObject width="100%%" height="100%%"><div xmlns="http://www.w3.org/1999/xhtml" style="width:%dpx;height:%dpx;overflow:hidden;box-sizing:border-box">%s</div></foreignObject></svg>`,t.Width,t.Height,t.Width,t.Height,html)
    w.Header().Set("Content-Type","image/svg+xml");w.Header().Set("Cache-Control","public,max-age=3600");w.WriteHeader(200);w.Write([]byte(svg))}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){t,_:=s.db.CountTemplates();rn,_:=s.db.CountRenders();writeJSON(w,200,map[string]interface{}{"templates":t,"renders":rn})}
func defaultHTML(name string)string{return fmt.Sprintf(`<div style="background:#1a1410;color:#f5e6c8;width:100%%;height:100%%;display:flex;flex-direction:column;justify-content:center;align-items:center;font-family:sans-serif;padding:60px;text-align:center"><h1 style="font-size:64px;margin:0;color:#c4622d">{{title}}</h1><p style="font-size:28px;color:#8b5e3c;margin-top:24px">{{subtitle}}</p><p style="font-size:18px;color:#7a6550;margin-top:40px">%s</p></div>`,name)}

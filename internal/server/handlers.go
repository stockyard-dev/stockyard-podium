package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-podium/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){status:=r.URL.Query().Get("status");list,_:=s.db.List(status);if list==nil{list=[]store.Feedback{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var f store.Feedback;json.NewDecoder(r.Body).Decode(&f);if f.Body==""{writeError(w,400,"body required");return};s.db.Create(&f);writeJSON(w,201,f)}
func(s *Server)handleUpdate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Status string `json:"status"`};json.NewDecoder(r.Body).Decode(&req);s.db.UpdateStatus(id,req.Status);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}

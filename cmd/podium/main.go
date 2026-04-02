package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-podium/internal/server";"github.com/stockyard-dev/stockyard-podium/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./podium-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("podium: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Podium — Self-hosted product feedback tool\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("podium: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}

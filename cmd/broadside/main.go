package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-broadside/internal/server";"github.com/stockyard-dev/stockyard-broadside/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./broadside-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("broadside: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Broadside — Self-hosted press release tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n",port,port)
log.Printf("broadside: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}

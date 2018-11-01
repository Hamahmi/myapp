package main

import (
 "database/sql"
 "fmt"
 "log"
 "net/http"
 "os"
 "time"
 "math/rand"
 _ "github.com/lib/pq"
  "github.com/go-redis/redis"
)

var (
 dbHost     = envOrDefault("MYAPP_DATABASE_HOST", "localhost")
 dbPort     = envOrDefault("MYAPP_DATABASE_PORT", "5432")
 dbUser     = envOrDefault("MYAPP_DATABASE_USER", "root")
 dbPassword = envOrDefault("MYAPP_DATABASE_PASSWORD", "secret")
 dbName     = envOrDefault("MYAPP_DATABASE_NAME", "myapp")

 webHost = envOrDefault("MYAPP_WEB_HOST", "")
 webPort = envOrDefault("MYAPP_WEB_PORT", "8080")

 db *sql.DB
 redisPool *redis.Client
)

func envOrDefault(key, defaultValue string) string {
 if value := os.Getenv(key); value != "" {
  return value
 }
 return defaultValue
}

func myHandler(w http.ResponseWriter, r *http.Request) {
 rows, err := db.Query("SELECT * FROM users")
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 defer rows.Close()

 fmt.Fprintln(w, "ID | Name")
 fmt.Fprintln(w, "---+--------")
 for rows.Next() {
  var (
    id   int
    name string
  )

  rows.Scan(&id, &name)

  fmt.Fprintf(w, "%2d | %s\n", id, name)
 }
}



func redisHandler(w http.ResponseWriter, r *http.Request) {

 n, err := redisPool.Get("n").Result()


 if err != nil {
  http.Error(w, "ERROR YBN EL SARMA", http.StatusInternalServerError)
  return
 }

 if(err == redis.Nil){
  randon := rand.Intn(100)
  redisPool.Set("n" , randon, 5*time.Second)
 }




 fmt.Fprintf(w, "n : %d", n)
}



func main() {
 redisHost := os.Getenv("REDISHOST")
 redisPort := os.Getenv("REDISPORT")

 redisPool = redis.NewClient(&redis.Options{
  Addr:         redisHost + ":" +redisPort,
 })


 http.HandleFunc("/cache", redisHandler)
 log.Fatal(http.ListenAndServe(":8080", nil))




 dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
  dbHost, dbPort, dbUser, dbPassword, dbName)
 var err error
 db, err = sql.Open("postgres", dbInfo)
 if err != nil {
  log.Fatal(err)
 }
 if err = db.Ping(); err != nil {
  log.Fatal(err)
 }

 http.HandleFunc("/", myHandler)
 log.Print("Listening on " + webHost + ":" + webPort + "...")
 http.ListenAndServe(webHost+":"+webPort, nil)
}

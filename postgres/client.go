package main

import (
    "fmt"
    "bufio"
    "database/sql"
    "log"
    "os"
    "os/signal"
    "time"

    _ "github.com/lib/pq"
)

const (
    DB_DRIVER   = "postgres"
    DB_USER     = "zaid"
    DB_NAME     = "zaid"
    DB_QUERY    = "INSERT INTO logs(message) VALUES($1);"
)
var (
    id int
    when time.Time
    message string
)

func Insert(db *sql.DB, stop chan os.Signal) (err error){
    fmt.Println("Inserting values")
    reader := bufio.NewReader(os.Stdin)
    for {
        select {
        case <- stop:
            return nil
        default:
            fmt.Println("Enter a Message:")
            text, _ := reader.ReadString('\n')
            rows, err := db.Query(DB_QUERY, text)
            if err != nil {
                return err
            }
            defer rows.Close()
        }
    }
}

func main() {
    dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", DB_USER, DB_NAME)
    db, err := sql.Open(DB_DRIVER, dbinfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()


    // set interrupt channel
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    defer close(stop)

    err = Insert(db, stop)
    if err != nil {
        log.Fatal(err)
        return
    }

    rows, err := db.Query("SELECT * FROM logs;")
    if err != nil {
        log.Fatal()
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id, &when, &message)
        if err != nil {
            log.Fatal(err)
        }
        log.Println(id, when, message)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}

package db

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// InitDB attempts to connect to the MySQL database with retry logic.
func InitDB() (*sql.DB, error) {
    var db *sql.DB
    var err error
    dsn := fmt.Sprintf("root:%s@tcp(mysql:%d)/%s", "password", 3306, "todo_db")

    // Retry loop for connecting to MySQL
    for i := 0; i < 5; i++ {
        db, err = sql.Open("mysql", dsn)
        if err != nil {
            log.Printf("Error connecting to MySQL (attempt %d): %v", i+1, err)
            time.Sleep(2 * time.Second) // Wait 2 seconds before retrying
            continue
        }

        // Ping the database to ensure it's reachable
        err = db.Ping()
        if err == nil {
            log.Println("Successfully connected to MySQL")
            return db, nil
        }

        log.Printf("Failed to ping MySQL (attempt %d): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }

    return nil, fmt.Errorf("could not connect to MySQL: %v", err)
}


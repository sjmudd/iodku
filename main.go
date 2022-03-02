package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

const WriteProbeQuery = `
	INSERT INTO test.heartbeat
		(id, master_ts, master_csec, update_by, master_id)
	VALUES
		(1, UTC_TIMESTAMP(6), ROUND(100 * @@timestamp), 'mysql_availability_collector', @@global.server_id)
	ON DUPLICATE KEY UPDATE
		master_ts=UTC_TIMESTAMP(6),
		master_csec=ROUND(100 * @@timestamp),
		update_by=VALUES(update_by),
		master_id=@@global.server_id
	;`

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHostname := os.Getenv("MYSQL_HOSTNAME")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHostname, mysqlPort, mysqlDatabase)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(fmt.Sprintf("Failed to open DB connections: %+v", err))

	}
	defer db.Close()
	for true {
		func() {
			started := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err = db.ExecContext(ctx, WriteProbeQuery)
			elapsed := time.Since(started)
			if err == nil {
				log.Printf("OK: took: %v", elapsed)
			} else {
				log.Printf("Error after %v: %v", elapsed, err)
			}
		}()

		time.Sleep(1 * time.Second)
	}
}

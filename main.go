package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

const (
	CreateTableIfMissing = `
	CREATE TABLE IF NOT EXISTS heartbeat (
		id bigint not null auto_increment primary key,
		master_ts timestamp(6) NOT NULL,
		master_csec bigint NOT NULL DEFAULT 0,
		update_by varchar(255) NOT NULL DEFAULT "",
		master_id bigint unsigned NOT NULL DEFAULT 0
	) ENGINE=InnoDB
	`

	WriteProbeQuery = `
	INSERT INTO heartbeat
		(id, master_ts, master_csec, update_by, master_id)
	VALUES
		(1, UTC_TIMESTAMP(6), ROUND(100 * @@timestamp), 'mysql_availability_collector', @@global.server_id)
	ON DUPLICATE KEY UPDATE
		master_ts=UTC_TIMESTAMP(6),
		master_csec=ROUND(100 * @@timestamp),
		update_by=VALUES(update_by),
		master_id=@@global.server_id
	;`
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	insertInterval := flag.String("insert-interval", "1s", "insert interval as understood by go")
	flag.Parse()

	parsedInterval, err := time.ParseDuration(*insertInterval)
	if err != nil {
		log.Printf("Unable to parse insert-interval")
		os.Exit(1)
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Printf("MYSQL_DSN needed to connect to server not provided or empty.  See: https://github.com/go-sql-driver/mysql#dsn-data-source-name")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to open DB connections: %+v", err))
		os.Exit(1)

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

		time.Sleep(parsedInterval)
	}
}

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

// Result contains the result of a single insert attempt
type Result struct {
	started  time.Time
	err      error
	duration time.Duration
}

// durationMetrics returns the number of successful attempts and the minimum, maximum and average duration of the successful insert attempts
func durationMetrics(results []Result) (int, time.Duration, time.Duration, time.Duration) {
	var (
		successful int
		total      int64         // time.Duration
		min        time.Duration = 1<<63 - 1
		max        time.Duration
	)

	for _, v := range results {
		if v.err == nil {
			successful++
		}
		total += int64(v.duration)
		if v.duration < min {
			min = v.duration
		}
		if v.duration > max {
			max = v.duration
		}
	}

	return successful, min, max, time.Duration(total / int64(len(results)))
}

func printSummary(interval time.Duration, results []Result) {
	log.Printf("Summary:")
	log.Printf("- interval:   %v", interval)
	log.Printf("- attempts:   %d", len(results))
	if len(results) > 0 {
		success, min, max, average := durationMetrics(results)
		log.Printf("- successful: %d, %6.2f%%", success, 100.0*float64(success)/float64(len(results)))
		log.Printf("- min:        %v", min)
		log.Printf("- average:    %v", average)
		log.Printf("- max:        %v", max)
	}
}

// common exit code handling
func exitIfCondition(condition bool, message string, exitCode int) {
	if condition {
		log.Printf(message)
		os.Exit(exitCode)
	}
}

func main() {
	var results []Result

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	insertInterval := flag.String("insert-interval", "1s", "insert interval as understood by Go")
	count := flag.Int("count", -1, "number of inserts to perform, default -1 = infinite")
	summary := flag.Bool("summary", false, "provide a summary when the insert count has been reached")
	maxWait := flag.String("max-wait", "1s", "maximum time to wait for the insert to complete as understood by Go")
	flag.Parse()

	parsedInterval, err := time.ParseDuration(*insertInterval)
	exitIfCondition(err != nil, "Unable to parse insert-interval", 1)

	maxWaitInterval, err := time.ParseDuration(*maxWait)
	exitIfCondition(err != nil, "Unable to parse max-wait", 1)

	dsn := os.Getenv("MYSQL_DSN")
	exitIfCondition(dsn == "", "MYSQL_DSN needed to connect to server not provided or empty.  See: https://github.com/go-sql-driver/mysql#dsn-data-source-name", 1)

	db, err := sql.Open("mysql", dsn)
	exitIfCondition(err != nil, fmt.Sprintf("Failed to open DB connections: %+v", err), 1)
	defer db.Close()

	startTime := time.Now()
	for *count == -1 || *count > 0 {
		{
			started := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), maxWaitInterval)
			defer cancel()
			_, err = db.ExecContext(ctx, WriteProbeQuery)
			elapsed := time.Since(started)
			if err == nil {
				log.Printf("OK: took: %v", elapsed)
			} else {
				log.Printf("Error after %v: %v", elapsed, err)
			}
			results = append(
				results,
				Result{
					started:  started,
					err:      err,
					duration: elapsed,
				})
		}

		time.Sleep(parsedInterval)
		if *count != -1 {
			*count--
		}
	}

	if *summary {
		printSummary(time.Since(startTime), results)
	}
}

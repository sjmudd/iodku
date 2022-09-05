# iodku

Testing inserts to a table and the time it takes.

Related to: https://bugs.mysql.com/bug.php?id=106526

Required table definition

```
CREATE TABLE `heartbeat` (
  `id` int NOT NULL,
  `master_ts` timestamp(6) NULL DEFAULT NULL,
  `master_csec` bigint DEFAULT NULL,
  `update_by` varchar(100) DEFAULT NULL,
  `master_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```

## Environment variables:

You can use `MYSQL_DSN=<user>:<pass>@tcp(<host>:<port>)/<db>` as the `MYSQL_DSN`
to configure how `iodku` will connect to the server.

## Command line options

* `--insert-interval=XXX` is used to determine the insert interval. Default is `1s`. See: https://github.com/go-sql-driver/mysql#dsn-data-source-name for the exact format.
* `--count=xx` number of inserts to perform, default is -1 (insert forever)

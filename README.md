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

## Sample output

```
[user@myhost ~]$ MYSQL_DSN='user_test:user_pass@tcp(myhost:3306)/iodku' ./iodku --insert-interval=50ms --count=10
2022/09/05 23:58:34.766807 OK: took: 21.639759ms
2022/09/05 23:58:34.836831 OK: took: 19.713192ms
2022/09/05 23:58:34.906791 OK: took: 19.738075ms
2022/09/05 23:58:34.976900 OK: took: 19.855682ms
2022/09/05 23:58:35.046594 OK: took: 19.529551ms
2022/09/05 23:58:35.116654 OK: took: 19.832476ms
2022/09/05 23:58:35.186455 OK: took: 19.600244ms
2022/09/05 23:58:35.256271 OK: took: 19.614538ms
2022/09/05 23:58:35.326060 OK: took: 19.569987ms
2022/09/05 23:58:35.395703 OK: took: 19.462197ms
[user@myhost ~]$
```

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
* `--summary` provide summary statistics of the inserts which were performed.

## Sample output

```
[user@myhost ~]$ MYSQL_DSN='user_test:user_pass@tcp(myhost:3306)/user_iodku' ./iodku --insert-interval=1s --count=10 --summary
2022/09/06 01:19:18.999876 OK: took: 22.752079ms
2022/09/06 01:19:20.020215 OK: took: 19.926164ms
2022/09/06 01:19:21.040269 OK: took: 19.734472ms
2022/09/06 01:19:22.060740 OK: took: 20.249455ms
2022/09/06 01:19:23.080779 OK: took: 19.707931ms
2022/09/06 01:19:24.100869 OK: took: 19.885405ms
2022/09/06 01:19:25.109734 OK: took: 8.583185ms
2022/09/06 01:19:26.118434 OK: took: 8.578377ms
2022/09/06 01:19:27.127168 OK: took: 8.646657ms
2022/09/06 01:19:28.141451 OK: took: 14.213081ms
2022/09/06 01:19:29.141781 Summary:
2022/09/06 01:19:29.141821 - interval:   10.164657709s
2022/09/06 01:19:29.141826 - attempts:   10
2022/09/06 01:19:29.141833 - successful: 10, 100.00%
2022/09/06 01:19:29.141837 - min:        8.578377ms
2022/09/06 01:19:29.141841 - average:    16.22768ms
2022/09/06 01:19:29.141844 - max:        22.752079ms
[user@myhost ~]
```

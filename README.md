# iodku

Testing inserts to a table and the time it takes.

Related to: https://bugs.mysql.com/bug.php?id=106526

## Table definition

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

Normal sample output might be:

```
$ MYSQL_DSN='user_test:user_pass@tcp(myhost:3306)/iodku' ./iodku --insert-interval=1s --count=10 --summary
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
$
```

A sample output with errors might look like this:
```
$ MYSQL_DSN='user_test:user_pass@tcp(myhost:3306)/iodku' ./iodku --insert-interval=1s --count=10 --summary --max-wait=19ms
2022/09/06 01:51:11.398199 OK: took: 9.171335ms
2022/09/06 01:51:12.407491 OK: took: 8.341306ms
2022/09/06 01:51:13.421910 OK: took: 14.060208ms
2022/09/06 01:51:14.436470 OK: took: 14.200061ms
2022/09/06 01:51:15.456453 Error after 19.323708ms: context deadline exceeded
2022/09/06 01:51:16.476552 Error after 19.960234ms: context deadline exceeded
2022/09/06 01:51:17.497449 Error after 19.847565ms: context deadline exceeded
2022/09/06 01:51:18.517580 Error after 20.011119ms: context deadline exceeded
2022/09/06 01:51:19.537625 Error after 19.952122ms: context deadline exceeded
2022/09/06 01:51:20.558140 Error after 20.01018ms: context deadline exceeded
2022/09/06 01:51:21.559110 Summary:
2022/09/06 01:51:21.559172 - interval:   10.170091554s
2022/09/06 01:51:21.559188 - attempts:   10
2022/09/06 01:51:21.559197 - successful: 4,  40.00%
2022/09/06 01:51:21.559201 - min:        8.341306ms
2022/09/06 01:51:21.559206 - average:    16.487783ms
2022/09/06 01:51:21.559209 - max:        20.011119ms
$
```

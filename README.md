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

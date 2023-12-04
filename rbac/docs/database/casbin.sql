CREATE DATABASE IF NOT EXISTS `casbin`;
use casbin;
CREATE TABLE IF NOT EXISTS `casbin_rule`
(
    `id`    bigint unsigned NOT NULL AUTO_INCREMENT,
    `ptype` varchar(100) DEFAULT NULL,
    `v0`    varchar(40)  DEFAULT NULL,
    `v1`    varchar(40)  DEFAULT NULL,
    `v2`    varchar(40)  DEFAULT NULL,
    `v3`    varchar(40)  DEFAULT NULL,
    `v4`    varchar(40)  DEFAULT NULL,
    `v5`    varchar(40)  DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`),
    UNIQUE KEY `unique_index` (`v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 20014
  DEFAULT CHARSET = utf8mb4;


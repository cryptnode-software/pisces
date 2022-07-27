-- +migrate Up
create schema if not exists pisces CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
-- create user if not exists pisces@'localhost' identified by piscespw;
create user if not exists pisces@'%' identified by 'piscespw';
-- grant all on pisces.* to pisces@'localhost';
grant all on pisces.* to pisces@'%';

-- +migrate Down
revoke all on pisces.*  from pisces@'%';
-- revoke all on pisces.*  from pisces@'localhost';
drop user if exists pisces@'%';
-- drop user if exists pisces@'localhost';
drop schema if exists pisces;
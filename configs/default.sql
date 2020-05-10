create table `sshkeys` (
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `addr`      TEXT,
    `login`     TEXT,
    `cert`      TEXT
);

create table `dockers` (
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `name`      TEXT,
    `path`      TEXT
);

create table `tasks` (
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `name`      TEXT,
    `path`      TEXT
);

create table `artifacts` (
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `task`      INTEGER,
    `path`      TEXT,
    `log`       TEXT,
    `datetime`  NUMERIC
);
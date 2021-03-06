-- Database: kindle_quotes

-- DROP DATABASE kindle_quotes;

-- CREATE DATABASE kindle_quotes
--  WITH OWNER = postgres
--       ENCODING = 'UTF8'
--       TABLESPACE = pg_default
--       LC_COLLATE = 'en_US.utf8'
--       LC_CTYPE = 'en_US.utf8'
--       CONNECTION LIMIT = -1;

--DROP TABLE tbl_authors;
--DROP TABLE tbl_sources;
--DROP TABLE tbl_quotes;

--use kindle_quotes;

create table tbl_users
(
	id bigserial,
	username text not null,
	password text not null
);

create unique index tbl_users_id_uindex
	on tbl_users (id);

create unique index tbl_users_username_uindex
	on tbl_users (username);

alter table tbl_users
	add constraint tbl_users_pk
		primary key (id);


CREATE TABLE  tbl_authors (
    user_id BIGINT NOT NULL,
    author_id serial PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES tbl_users (id),
    author_name VARCHAR(256) UNIQUE NOT NULL
);

CREATE TABLE  tbl_sources (
    source_id serial PRIMARY KEY,
    source_title VARCHAR(256) UNIQUE NOT NULL,
    author_id INT not null,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES tbl_users (id),
    FOREIGN KEY (author_id) REFERENCES tbl_authors (author_id)
);

CREATE TABLE tbl_quotes (
    quote_id serial PRIMARY KEY,
    source_id INT NOT NULL,
    quote VARCHAR(4096) NOT NULL,
    date_taken timestamp,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES tbl_users (id),
    FOREIGN KEY (source_id) REFERENCES tbl_sources (source_id)
);


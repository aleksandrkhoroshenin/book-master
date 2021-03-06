CREATE TABLE books (
	id uuid PRIMARY KEY,
	name varchar(50),
    page_count      int,
    release_year    int,
	check (page_count < 10000),
	check (release_year > 1900)
);

create table session (
  id uuid primary key,
  user_id uuid not null unique
);

create table users (
  id uuid primary key,
  login varchar(50) not null unique,
  password varchar(50) not null
);
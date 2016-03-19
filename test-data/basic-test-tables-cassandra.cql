--
-- create a keyspace
--
CREATE KEYSPACE testgdbc with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
USE testgdbc;

--
-- create tables
--
CREATE TABLE first (
  id int,
  first_name text,
  middle_name text,
  last_name text,
  age int,
  birthdate text ,
  description text,
  more_info text,
  addr text,
  city text,
  PRIMARY KEY(id)
);

CREATE TABLE second (
  id int,
  first_name text,
  middle_name text,
  last_name text,
  age int,
  birthdate text ,
  description text,
  more_info text,
  addr text,
  city text,
  PRIMARY KEY(id)
);

--
-- test the schema
--

DESC first;

DESC second;

--
-- insert test data
--

INSERT INTO first (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (1, 'Bob', 'A', 'Count', 20, '1994-07-17', 'a little info', 'more info', '123 main st', 'cityville');
INSERT INTO first (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (2, 'Carl', 'Ben', 'Gone', 21, '1995-06-20', 'car salesman', NULL, NULL, 'pleasantville');
INSERT INTO first (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (3, 'Sid', 'Timothy', 'Thomas', 30, '1990-03-13', NULL, NULL, NULL, '');

INSERT INTO second (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (1, 'Bob', 'Allen', 'Count', 20, '1994-07-17', 'a little info', 'more info', '123 main st', 'cityville');
INSERT INTO second (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (2, 'Carl', 'Ben', 'Gone', 23, '1995-06-20', 'cashier', 'dentist', 'wacky circle', 'pleasantville');
INSERT INTO second (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (3, 'Sid', 'Timothy', 'Thomas', 30, '1990-03-13', NULL, NULL, NULL, '');
INSERT INTO second (id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city) VALUES (4, 'Andrew', 'Dishional', 'Guy', 18, '2000-01-03', 'helicopter pilot', 'pink hair', '555 one street', 'onesville');

--
-- test data
--

SELECT * FROM first;

SELECT * FROM second;

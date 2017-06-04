/*
* Database definition script for nfl_app
*
* Author: Kyle Ames
* Last Updated: December 24, 2015
*/

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL UNIQUE,
    admin boolean NOT NULL DEFAULT FALSE,
    last_login timestamp,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS pvs (
    id integer PRIMARY KEY,
    type varchar(1) NOT NULL UNIQUE,
    seven integer NOT NULL,
    five integer NOT NULL,
    three integer NOT NULL,
    one integer NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id integer PRIMARY KEY,
    city varchar(64) NOT NULL,
    nickname varchar(64) NOT NULL,
    stadium varchar(64) NOT NULL,
    abbreviation varchar(4) NOT NULL
);

CREATE TABLE IF NOT EXISTS years (
    id integer PRIMARY KEY,
    year integer NOT NULL UNIQUE,
    year_start integer NOT NULL
);

CREATE TABLE IF NOT EXISTS weeks (
    id integer PRIMARY KEY,
    year_id integer REFERENCES years(id) ON DELETE CASCADE,
    pvs_id integer REFERENCES pvs(id),
    week integer NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id integer PRIMARY KEY,
    week_id integer REFERENCES weeks(id) ON DELETE CASCADE,
    date integer NOT NULL,
    home_id integer REFERENCES teams(id),
    away_id integer REFERENCES teams(id),
    home_score integer DEFAULT -1,
    away_score integer DEFAULT -1
);

CREATE TABLE IF NOT EXISTS picks (
    id integer PRIMARY KEY,
    user_id integer REFERENCES users(id),
    game_id integer REFERENCES games(id),
    selection integer REFERENCES teams(id) DEFAULT NULL,
    points integer DEFAULT 0, 
);

CREATE TABLE IF NOT EXISTS statistics (
    id integer PRIMARY KEY,
    user_id integer REFERENCES users(id),
    week_id integer REFERENCES weeks(id),
    zero integer,
    one integer,
    three integer,
    five integer,
    seven integer,
    winner boolean,
    lowest boolean
);

INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Buffalo', 'Bills', 'Ralph Wilson Stadium', 'BUF');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Miami', 'Dolphins', 'Sun Life Stadium', 'MIA');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('New England', 'Patriots', 'Gilette Stadium', 'NE');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('New York', 'Jets', 'MetLife Stadium', 'NYJ');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Baltimore', 'Ravens', 'M&T Bank Stadium', 'BAL');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Cincinatti', 'Bengals', 'Paul Brown Stadium', 'CIN');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Cleveland', 'Browns', 'First Energy Stadium', 'CLE');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Pittsburgh', 'Steelers', 'Heinz Field', 'PIT');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Houston', 'Texans', 'Reliant Stadium', 'HOU');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Indianapolis', 'Colts', 'Lucas Oil Stadium', 'IND');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Jacksonville', 'Jaguars', 'EverBank Field', 'JAC');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Tennessee', 'Titans', 'LP Field', 'TEN');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Denver', 'Broncos', 'Mile High Stadium', 'DEN');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Kansas City', 'Chiefs', 'Arrowhead Stadium', 'KC');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Oakland', 'Raiders', 'O.co Coliseum', 'OAK');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('San Diego', 'Chargers', 'Qualcomm Stadium', 'SD');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Dallas', 'Cowboys', 'AT&T Stadium', 'DAL');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('New York', 'Giants', 'MetLife Stadium', 'NYG');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Philadelphia', 'Eagles', 'Lincoln Financial Field', 'PHI');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Washington', 'Redskins', 'FedEx Field', 'WAS');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Chicago', 'Bears', 'Soldier Field', 'CHI');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Detroit', 'Lions', 'Ford Field', 'DET');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Green Bay', 'Packers', 'Lambeau Field', 'GB');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Minnesota', 'Vikings', 'TCF Bank Stadium', 'MIN');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Atlanta', 'Falcons', 'Georiga Dome', 'ATL');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Carolina', 'Panthers', 'Bank of America Stadium', 'CAR');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('New Orleans', 'Saints', 'Mercedes-Benz Superdome', 'NO');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Tampa Bay', 'Buccaneers', 'Raymond James Stadium', 'TB');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Arizona', 'Cardinals', 'University of Phoenix Stadium', 'ARI');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('St. Louis', 'Rams', 'Edward Jones Dome', 'STL');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('San Francisco', '49ers', 'Candlestick Park', 'SF');
INSERT INTO teams(city, nickname, stadium, abbreviation) VALUES('Seattle', 'Seahawks', 'CenturyLink Field', 'SEA');

INSERT INTO pvs(type, seven, five, three, one) VALUES('A', 1, 2, 5, 8);
INSERT INTO pvs(type, seven, five, three, one) VALUES('B', 1, 2, 5, 7);
INSERT INTO pvs(type, seven, five, three, one) VALUES('C', 1, 2, 5, 6);
INSERT INTO pvs(type, seven, five, three, one) VALUES ('D', 1, 2, 5, 5);

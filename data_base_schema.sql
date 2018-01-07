
SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;


CREATE TABLE IF NOT EXISTS journeys (
	origin varchar(20),
    destiny varchar(20),
    departure_date timestamp,
    arrival_date timestamp,
    check_date date default now(),
    class varchar(20),
    price real,
    CONSTRAINT journeys_pkey PRIMARY KEY (origin, destiny, departure_date, class, check_date)
);

GRANT SELECT,UPDATE,INSERT,DELETE ON TABLE public.journeys TO renfe;


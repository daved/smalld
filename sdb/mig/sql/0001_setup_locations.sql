-- +migrate Up
-- +migrate StatementBegin

CREATE SEQUENCE locations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE locations (
    id integer NOT NULL DEFAULT nextval('locations_id_seq'),
    label text,
    acc numeric,
    geom geometry(Point,4326),
    received timestamp without time zone DEFAULT now()
);

ALTER SEQUENCE locations_id_seq OWNED BY locations.id;

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE locations;

-- +migrate StatementEnd

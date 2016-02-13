-- +migrate Up
-- +migrate StatementBegin

CREATE SEQUENCE adminareas_gid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE adminareas (
    gid integer NOT NULL DEFAULT nextval('adminareas_gid_seq'),
    osm_id character varying(20),
    lastchange character varying(19),
    code smallint,
    fclass character varying(40),
    postalcode character varying(10),
    name character varying(100),
    geom geometry(MultiPolygon,4326)
);

ALTER SEQUENCE adminareas_gid_seq OWNED BY adminareas.gid;

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE adminareas;

-- +migrate StatementEnd

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

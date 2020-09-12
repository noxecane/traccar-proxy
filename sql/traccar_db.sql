CREATE TABLE public.tc_devices (
    id integer NOT NULL,
    name character varying(128) NOT NULL,
    uniqueid character varying(128) NOT NULL,
    lastupdate timestamp without time zone,
    positionid integer,
    groupid integer,
    attributes character varying(4000),
    phone character varying(128),
    model character varying(128),
    contact character varying(512),
    category character varying(128),
    disabled boolean DEFAULT false
);

CREATE TABLE public.tc_events (
    id integer NOT NULL,
    type character varying(128) NOT NULL,
    servertime timestamp without time zone NOT NULL,
    deviceid integer,
    positionid integer,
    geofenceid integer,
    attributes character varying(4000),
    maintenanceid integer
);

CREATE TABLE public.tc_positions (
    id integer NOT NULL,
    protocol character varying(128),
    deviceid integer NOT NULL,
    servertime timestamp without time zone DEFAULT now() NOT NULL,
    devicetime timestamp without time zone NOT NULL,
    fixtime timestamp without time zone NOT NULL,
    valid boolean NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    altitude double precision NOT NULL,
    speed double precision NOT NULL,
    course double precision NOT NULL,
    address character varying(512),
    attributes character varying(4000),
    accuracy double precision DEFAULT 0 NOT NULL,
    network character varying(4000)
);
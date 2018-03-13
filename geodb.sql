BEGIN;

CREATE TABLE geo (
    zip VARCHAR(5)
,    city VARCHAR(64)
,    state CHAR(2)
,    country VARCHAR(64)
,    timezone VARCHAR(64)
,    latitude DOUBLE
,    longitude DOUBLE
);

CREATE INDEX idx_zip ON geo (zip);
CREATE INDEX idx_lat_long ON geo (latitude, longitude);

COMMIT;

/*
zip,city,state,country,timezone,latitude,longitude
00501,Holtsville,NY,"United States",America/New_York,40.922326,-72.637078
*/

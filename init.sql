CREATE SCHEMA segmenting
    AUTHORIZATION postgres;

CREATE TABLE segmenting.user (
    id BIGINT PRIMARY KEY,
    creationdate TIMESTAMPTZ NOT NULL
);

CREATE TABLE segmenting.segment (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    creationdate TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX name_idx ON segmenting.segment (name);

CREATE TABLE segmenting.assignment (
    userid BIGSERIAL NOT NULL,
    segmentid BIGINT NOT NULL,
    untildate TIMESTAMPTZ,
    PRIMARY KEY (userid, segmentid)
);
CREATE INDEX assignment_segment_idx ON segmenting.assignment (segmentid);


CREATE TABLE segmenting.operation (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR,
    creationdate TIMESTAMPTZ NOT NULL,
    updatedate TIMESTAMPTZ NOT NULL,
    isactive BOOLEAN NOT NULL
);

INSERT INTO segmenting.operation (name, description, creationdate, updatedate, isactive)
    VALUES ('ADDED', 'Добавление', now(), now(), TRUE),
           ('REMOVED', 'Удаление', now(), now(), TRUE),
           ('EXPIRED', 'Удаление по истечении времени', now(), now(), TRUE);

CREATE TABLE segmenting.history (
    id BIGSERIAL PRIMARY KEY,
    userid BIGINT NOT NULL,
    segmentid BIGINT NOT NULL,
    operationid BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);
CREATE INDEX history_user_idx ON segmenting.history (userid);
CREATE INDEX history_segment_idx ON segmenting.history (segmentid);

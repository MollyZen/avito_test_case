CREATE SCHEMA segmenting
    AUTHORIZATION postgres;

CREATE TABLE segmenting.user (
    id BIGINT PRIMARY KEY,
    creationdate TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE segmenting.segment (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR NOT NULL CHECK (slug ~ '^[a-zA-z0-9\-_]{4,}$'),
    isactive BOOL NOT NULL DEFAULT TRUE,
    creationdate TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX name_idx ON segmenting.segment (slug);

CREATE TABLE segmenting.assignment (
    userid BIGSERIAL NOT NULL,
    segmentid BIGINT NOT NULL,
    untildate TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (userid, segmentid),
    CONSTRAINT fk_assig_userid
        FOREIGN KEY(userid)
            REFERENCES segmenting.user(id),
    CONSTRAINT fk_assig_segmentid
        FOREIGN KEY(segmentid)
            REFERENCES segmenting.segment(id)
);
CREATE INDEX assignment_segment_idx ON segmenting.assignment (segmentid);


CREATE TABLE segmenting.operation (
    id BIGINT PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR,
    creationdate TIMESTAMPTZ NOT NULL DEFAULT now(),
    updatedate TIMESTAMPTZ NOT NULL DEFAULT now(),
    isactive BOOLEAN NOT NULL
);

INSERT INTO segmenting.operation (id, name, description, creationdate, updatedate, isactive)
    VALUES (0, 'ADDED', 'Добавление', now(), now(), TRUE),
           (1, 'REMOVED', 'Удаление', now(), now(), TRUE),
           (2, 'EXPIRED', 'Удаление по истечении времени', now(), now(), TRUE),
           (3, 'UPDATED', 'Обновление значения', now(), now(), TRUE),
           (4, 'SEGMENT_DELETED', 'Удален вместе с сегментом', now(), now(), TRUE);

CREATE TABLE segmenting.history (
    id BIGSERIAL PRIMARY KEY,
    userid BIGINT NOT NULL,
    segmentid BIGINT NOT NULL,
    operationid BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX history_user_idx ON segmenting.history (userid);
CREATE INDEX history_segment_idx ON segmenting.history (segmentid);

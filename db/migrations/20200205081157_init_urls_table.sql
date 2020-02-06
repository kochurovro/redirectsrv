-- migrate:up
CREATE TABLE urls (
     id MEDIUMINT NOT NULL AUTO_INCREMENT,
     name CHAR(30) NOT NULL,
     url VARCHAR(2083) DEFAULT '',
     PRIMARY KEY (id)
);

-- migrate:down

DROP TABLE urls
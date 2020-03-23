CREATE TABLE `svc-article`.articles
(
    id    INT UNSIGNED,
    title VARCHAR(255) NOT NULL,
    date  DATE         NOT NULL,
    body  VARCHAR(1024),
    PRIMARY KEY (id),
    UNIQUE KEY `ID_UNIQUE` (id)
) ENGINE = InnoDB;

CREATE TABLE `svc-article`.tags
(
    id        INT UNSIGNED AUTO_INCREMENT,
    tag_title VARCHAR(30) UNIQUE,
    PRIMARY KEY (id),
    UNIQUE KEY `ID_UNIQUE` (id)
) ENGINE = InnoDB;

CREATE TABLE `svc-article`.article_tags
(
    article_id INT UNSIGNED,
    tag_id     INT UNSIGNED,
    CONSTRAINT `fk_article_id` FOREIGN KEY
        (article_id) REFERENCES `svc-article`.articles (id),
    CONSTRAINT `fk_tag_id` FOREIGN KEY
        (tag_id) REFERENCES `svc-article`.tags (id)
) ENGINE = InnoDB;

CREATE TABLE `svc-article`.articles_test
(
    id    INT UNSIGNED,
    title VARCHAR(255) NOT NULL,
    date  DATE         NOT NULL,
    body  VARCHAR(1024),
    PRIMARY KEY (id),
    UNIQUE KEY `ID_UNIQUE` (id)
) ENGINE = InnoDB;

CREATE TABLE `svc-article`.tags_test
(
    id        INT UNSIGNED AUTO_INCREMENT,
    tag_title VARCHAR(30) UNIQUE,
    PRIMARY KEY (id),
    UNIQUE KEY `ID_UNIQUE` (id)
) ENGINE = InnoDB;

CREATE TABLE `svc-article`.article_tags_test
(
    article_id INT UNSIGNED,
    tag_id     INT UNSIGNED,
    CONSTRAINT `fk_article_id_test` FOREIGN KEY
        (article_id) REFERENCES `svc-article`.articles_test (id),
    CONSTRAINT `fk_tag_id_test` FOREIGN KEY
        (tag_id) REFERENCES `svc-article`.tags_test (id)
) ENGINE = InnoDB;
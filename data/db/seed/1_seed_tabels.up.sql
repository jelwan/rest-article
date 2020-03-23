INSERT INTO `svc-article`.tags(tag_title)
VALUES ('sports'),
       ('tech'),
       ('cooking'),
       ('science'),
       ('fitness'),
       ('health');

SELECT `id`, `tag_title`
FROM `svc-article`.tags
WHERE `tag_title` IN ('sports', 'tech', 'cooking', 'science', 'fitness', 'health');

Select `id`, `title`, `date`, `body`
FROM `svc-article`.articles
where id = 1;

SELECT `tag_id`
FROM `svc-article`.article_tags
WHERE `article_id` = 1;

INSERT INTO `svc-article`.article_tags(article_id, tag_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (1, 4);

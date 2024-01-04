BEGIN;

INSERT INTO "users"(
    "username",
    "email",
    "password",
    "bio"
)
VALUES
('jake','jake@j.com','$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6','I work at statefarm'),
('joedoe','joedoe@j.com','$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6','I work at home');

INSERT INTO "articles"(
    "author_id",
    "slug",
    "title",
    "description",
    "body"
)
VALUES
(1,'how-to-train-your-dragon','How to train your dragon','Ever wonder how?','It takes a Jacobian'),
(2,'Today is friday','wish everyday is holiday','hope','what is your hobby');

INSERT INTO "comments"(
    "body",
    "article_id",
    "author_id"
)
VALUES
('It takes a Jacobian',1,2),
('Wow WOw wow',2,1);

INSERT INTO "tags"(
    "name"
)
VALUES
('sun'),
('set');

INSERT INTO "article_tags"(
    "article_id",
    "tag_id"
)
VALUES
(1,1),
(1,2),
(2,2);

INSERT INTO "article_favorites"(
    "user_id",
    "article_id"
)
VALUES
(1,2),
(1,1),
(2,2);

INSERT INTO "user_follows"(
    "follower_id",
    "following_id"
)
VALUES
(1,2),
(2,1);


COMMIT;
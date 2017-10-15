CREATE TABLE post(
    id serial NOT NULL PRIMARY KEY,
    blog_id INT NOT NULL REFERENCES blog,
    title TEXT,
    body TEXT
);

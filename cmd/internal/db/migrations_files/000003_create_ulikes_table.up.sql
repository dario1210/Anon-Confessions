CREATE TABLE posts_likes (
    post_id INTEGER NOT NULL,              
    user_id INTEGER NOT NULL,               
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) 
);

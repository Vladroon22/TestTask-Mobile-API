CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(20),
    hash VARCHAR(70) NOT NULL,
    email VARCHAR(30) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    author VARCHAR(30) NOT NULL,
    title VARCHAR(15) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    post_id INT NOT NULL, 
    author VARCHAR(30) NOT NULL,
    content TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    post_id INT,
    comment_id INT,
    liker VARCHAR(20) NOT NULL,
    type_of_like VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE (post_id, comment_id, liker) -- Уникальный индекс для предотвращения дублирования лайков
);

CREATE INDEX idx_nickname ON users (nickname);

CREATE INDEX idx_email ON users (email);

-- Вставка данных в таблицу users
INSERT INTO users (nickname, name, hash, email) VALUES
('user1', 'John Doe', '$2a$10$3qvjn1aja61SmS8YfKYkieSSkR2bG/Hw4DRZ2IV8LZuDbRi4Y5PzS', 'user1@example.com'),
('user2', 'Jane Smith', '$2a$10$WmKZ/ivq.FmoZ9pCDiR4WepKTK0zZTNLU1O9/JgkQJkoHw6EHcMQq', 'user2@example.com'),
('user3', 'Alice Johnson', '$2a$10$x.HRt0cFlJ7BMLjlguNA2ONlrMeD.UIasANO1a9oB8BK/loiirxE2', 'user3@example.com'),
('user4', 'Bob Brown', '$2a$10$EZ5DSPY6lk1WWni.yrP5w.AgJ87d.Se0X7bFEINkAFNBKnt9j8VhK', 'user4@example.com'),
('user5', 'Charlie Davis', '$2a$10$BDXiS./vNvrbDxnw5ElVmOzabB0SWKMNpK/aXDu3wERU/qGNs1w5W', 'user5@example.com');

-- Вставка данных в таблицу posts
INSERT INTO posts (user_id, author, title, content) VALUES
(1, 'user1', 'First Post', 'This is the content of the first post.'),
(2, 'user2', 'Second Post', 'This is the content of the second post.'),
(3, 'user3', 'Third Post', 'This is the content of the third post.'),
(4, 'user4', 'Fourth Post', 'This is the content of the fourth post.'),
(5, 'user5', 'Fifth Post', 'This is the content of the fifth post.');

-- Вставка данных в таблицу comments
INSERT INTO comments (user_id, post_id, author, content) VALUES
(1, 1, 'user1', 'This is a comment on the first post.'),
(2, 1, 'user2', 'Another comment on the first post.'),
(3, 2, 'user3', 'This is a comment on the second post.'),
(4, 3, 'user4', 'This is a comment on the third post.'),
(5, 4, 'user5', 'This is a comment on the fourth post.');

-- Вставка данных в таблицу likes
INSERT INTO likes (post_id, comment_id, liker, type_of_like) VALUES
(1, 1, 'user1', 'like'),
(1, 2, 'user2', 'dislike'),
(2, 3, 'user3', 'like'),
(3, 4, 'user4', 'like'),
(4, 5, 'user5', 'dislike');
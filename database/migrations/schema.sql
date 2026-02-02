
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    bio VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tweets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content VARCHAR(280) NOT NULL, -- Required by the PDF
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index necessary: To quickly find "all tweets of Pedro"
CREATE INDEX idx_tweets_user_id ON tweets(user_id);

CREATE TABLE followers (
    id SERIAL PRIMARY KEY,
    follower_id INT NOT NULL, -- The one who follows (me)
    followed_id INT NOT NULL, -- The one who is followed (the famous one)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_followed FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,

    -- Business rule 1: You cannot follow the same person twice
    CONSTRAINT unique_following UNIQUE (follower_id, followed_id),

    -- Business rule 2: You cannot follow yourself (Optional but recommended)
    CONSTRAINT check_no_self_follow CHECK (follower_id <> followed_id)
);

-- Index 1: To know "Who I follow" (Build my timeline)
CREATE INDEX idx_followers_follower ON followers(follower_id);

-- Index 2: To know "Who follows me" (Fan-out or notifications)
CREATE INDEX idx_followers_followed ON followers(followed_id);

CREATE DATABASE IF NOT EXISTS investordb;

USE investordb;

CREATE TABLE IF NOT EXISTS users (
		id BIGINT NOT NULL UNIQUE AUTO_INCREMENT,
		username VARCHAR(40) NOT NULL UNIQUE,
		password VARCHAR(40) NOT NULL,
		rating BIGINT NOT NULL,
		created_at INT NOT NULL,
		PRIMARY KEY (id)
	);

CREATE TABLE IF NOT EXISTS threads (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

CREATE TABLE IF NOT EXISTS category (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
CREATE TABLE IF NOT EXISTS stocks (
		symbol VARCHAR(10) NOT NULL UNIQUE,
		name TEXT,
		info TEXT,
		image_url TEXT,
		PRIMARY KEY (symbol)
	);
	
CREATE TABLE IF NOT EXISTS comments (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		thread_id BIGINT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (thread_id) REFERENCES threads(id)
	);

INSERT stocks (symbol, name, info, image_url)
VALUES ("AAPL", "Apple Inc.", "Apple Bla bla bla bla bla", "https://example.com/");

INSERT users (username, password, rating, created_at)
VALUES ("admin", "password", 0, 0);
		
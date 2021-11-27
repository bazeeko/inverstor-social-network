package main

import (
	"database/sql"
	"fmt"
	"log"

	_userDeliveryHttp "github.com/bazeeko/investor-social-network/user/delivery/http"
	_userRepoMysql "github.com/bazeeko/investor-social-network/user/repository/mysql"
	_userUsecase "github.com/bazeeko/investor-social-network/user/usecase"

	_stockDeliveryHttp "github.com/bazeeko/investor-social-network/stock/delivery/http"
	_stockRepoMysql "github.com/bazeeko/investor-social-network/stock/repository/mysql"
	_stockUsecase "github.com/bazeeko/investor-social-network/stock/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func connectDB(config string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", config)
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	// _, err = conn.Exec(`CREATE DATABASE IF NOT EXISTS investordb`)
	// if err != nil {
	// 	return nil, fmt.Errorf("connectDB: %w", err)
	// }

	_, err = conn.Exec(`USE heroku_449b7b52dda3dc9`)
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT NOT NULL UNIQUE AUTO_INCREMENT,
		username VARCHAR(40) NOT NULL UNIQUE,
		password VARCHAR(40) NOT NULL,
		rating BIGINT NOT NULL,
		created_at INT NOT NULL,
		PRIMARY KEY (id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS threads (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS category (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS stocks (
		symbol VARCHAR(10) NOT NULL UNIQUE,
		name TEXT,
		info TEXT,
		image_url TEXT,
		PRIMARY KEY (symbol)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		thread_id BIGINT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (thread_id) REFERENCES threads(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`INSERT stocks (symbol, name, info, image_url) VALUES (?, ?, ?, ?)`,
		"AAPL", "Apple Inc.", "Apple Bla bla bla bla bla", "https://example.com/")

	if err != nil {
		log.Println(err)
	}

	_, err = conn.Exec(`INSERT users (username, password, rating, created_at) VALUES (?, ?, ?, ?)`,
		"DWAdmin", "12345", 0, 0)

	if err != nil {
		log.Println(err)
	}

	return conn, nil
}

func main() {
	// tcp(127.0.0.1:3306)
	// dbURL := os.Getenv("CLEARDB_DATABASE_URL")

	config1 := "bd57b99c55080f:d3327b71@tcp(eu-cdbr-west-01.cleardb.com)/heroku_449b7b52dda3dc9"

	config := "root:password@tcp(mysqldb)/"
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("main: %s\n", err)
	}

	e := echo.New()

	userRepoMysql := _userRepoMysql.NewMysqlUserRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepoMysql)
	_userDeliveryHttp.NewUserHandler(e, userUsecase)

	stockRepoMysql := _stockRepoMysql.NewMysqlStockRepository(db)
	stockUsecase := _stockUsecase.NewStockUsecase(stockRepoMysql)
	_stockDeliveryHttp.NewStockHandler(e, stockUsecase)

	log.Fatalln(e.Start(":8080"))
	// conn, err := sql.Open("mysql", config)
}

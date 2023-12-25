package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"spam-telegram-bot/configs"
	"spam-telegram-bot/internal/bot"
	"spam-telegram-bot/internal/repository/migrations"
	"spam-telegram-bot/internal/repository/models"
)

func main() {
	flag.Parse()

	cfg, err := configs.NewConfig("config.toml")

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := migrations.SqliteMigration(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	fmt.Print(db)

	spambot := &bot.SpamBot{
		Bot:  bot.InitBot(cfg.BotToken),
		User: models.UserModel{DB: db},
	}

	router := http.NewServeMux()
	bot.SetRoutes(router, spambot)
	srv := &http.Server{
		Addr:         ":5000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	spambot.Bot.Handle("/start", spambot.StartHandler)
	log.Println("Server is running on port 5000")
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	spambot.Bot.Start()
}

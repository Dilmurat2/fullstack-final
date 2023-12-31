package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gopkg.in/telebot.v3"
	"spam-telegram-bot/internal/repository/models"
	"spam-telegram-bot/internal/web"
)

type SpamBot struct {
	Bot  *telebot.Bot
	User models.UserModel
}

type Message struct {
	Message string `json:"message"`
}

func (b *SpamBot) SendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(web.StatusMessage(http.StatusMethodNotAllowed))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(web.StatusMessage(http.StatusBadRequest))
		return
	}

	var msg Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(web.StatusMessage(http.StatusBadRequest))
		return
	}

	users, err := b.User.SelectAll()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(web.StatusMessage(http.StatusInternalServerError))
		return
	}

	for _, user := range users {
		u := &telebot.User{
			ID:        user.TelegramId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
		}
		_, err := b.Bot.Send(u, msg.Message)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(web.StatusMessage(http.StatusInternalServerError))
			return
		}
	}
	w.Write(web.StatusMessage(200))
}

func (b *SpamBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		TelegramId: ctx.Chat().ID,
		Username:   ctx.Sender().Username,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
	}

	existUser, err := b.User.FindOne(newUser)
	if err != nil {
		_ = fmt.Errorf("ошибка поиска юзера %v", err)
	}

	if existUser == nil {
		err := b.User.AddUser(newUser)
		if err != nil {
			_ = fmt.Errorf("ошибка создания юзера %v", err)
		}
	}

	return ctx.Send("Привет " + ctx.Sender().FirstName)
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}

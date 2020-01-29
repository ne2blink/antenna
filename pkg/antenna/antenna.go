package antenna

import (
	"fmt"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ne2blink/antenna/pkg/storage"
)

// Antenna is a telegram bot implemented in Golang, broadcasting message to subscribers.
type Antenna struct {
	bot   *tgbotapi.BotAPI
	store storage.Store
	log   *zap.SugaredLogger
}

// New creates a new Antenna instance.
func New(bot *tgbotapi.BotAPI, store storage.Store, log *zap.SugaredLogger) *Antenna {
	return &Antenna{
		bot:   bot,
		store: store,
		log:   log,
	}
}

// Listen listens telegram
func (a *Antenna) Listen() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go a.handleMessage(update.Message)
	}
	return nil
}

func (a *Antenna) handleMessage(msg *tgbotapi.Message) {
	log := a.log.With(
		"username", msg.From.UserName,
		"chat_id", msg.Chat.ID,
	)
	defer log.Sync()
	defer func() {
		if r := recover(); r != nil {
			log.Errorw("panic", "error", fmt.Sprint(r))
		}
	}()

	h := &handler{
		base: a,
		msg:  msg,
		log:  log,
	}
	if err := h.handle(); err != nil {
		log.Errorw(msg.Text, "error", err.Error())
	} else {
		log.Info(msg.Text)
	}
}
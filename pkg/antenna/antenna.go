package antenna

import (
	"fmt"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ne2blink/antenna/pkg/storage"
)

// Antenna is a telegram bot implemented in Golang, broadcasting message to subscribers.
type Antenna struct {
	bot    *tgbotapi.BotAPI
	store  storage.Store
	log    *zap.SugaredLogger
	admins map[string]struct{}
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
		switch {
		case update.CallbackQuery != nil:
			go a.handleCallback(update.CallbackQuery)
		case update.Message != nil:
			go a.handleMessage(update.Message)
		}
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
			log.Errorw("panic", "err", fmt.Sprint(r))
		}
	}()

	h := &handler{
		base: a,
		msg:  msg,
		log:  log,
	}
	logEvent(log, msg.Text, h.handle())
}

func (a *Antenna) handleCallback(cb *tgbotapi.CallbackQuery) {
	log := a.log.With(
		"username", cb.From.UserName,
		"chat_id", cb.Message.Chat.ID,
	)
	defer log.Sync()
	defer func() {
		if r := recover(); r != nil {
			log.Errorw("panic", "err", fmt.Sprint(r))
		}
	}()

	h := &callback{
		base: a,
		cb:   cb,
		log:  log,
	}
	logEvent(log, cb.Data, h.handle())
}

// AddAdmin enables admin feature and adds add admins.
func (a *Antenna) AddAdmin(usernames ...string) {
	if a.admins == nil {
		a.admins = make(map[string]struct{})
	}
	for _, username := range usernames {
		a.admins[username] = struct{}{}
	}
}

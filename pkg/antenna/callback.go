package antenna

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type callback struct {
	base *Antenna
	cb   *tgbotapi.CallbackQuery
	log  *zap.SugaredLogger
}

func (c *callback) handle() error {
	args := strings.SplitN(c.cb.Data, "|", 3)
	if len(args) != 3 || args[0] != "antenna" {
		return nil
	}
	cmd, data := args[1], args[2]
	switch cmd {
	case "broadcast":
		return c.handleBroadcast(data)
	}
	return nil
}

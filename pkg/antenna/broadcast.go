package antenna

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Broadcast sends text to all subscribers
func (a *Antenna) Broadcast(appID, text, parseMode string) error {
	app, err := a.store.GetApp(appID)
	if err != nil {
		return err
	}
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(app.Name, "antenna|broadcast|"+app.ID+"|"+app.Name),
		),
	)

	ids, err := a.store.ListSubscribers(appID)
	if err != nil {
		return err
	}

	texts := make([]string, 0, len(text)/4096+1)
	for offset := 0; offset < len(text); offset += 4096 {
		if offset+4096 < len(text) {
			texts = append(texts, text[offset:offset+4096])
		} else {
			texts = append(texts, text[offset:])
		}
	}

	for _, id := range ids {
		for i, text := range texts {
			msg := tgbotapi.NewMessage(id, text)
			msg.ParseMode = parseMode
			msg.ReplyMarkup = replyMarkup
			_, err := a.bot.Send(msg)
			if err != nil {
				a.log.Warnw(
					"failed to send message",
					"text", text,
					"text_id", i,
					"app_id", appID,
					"chat_id", id,
					"parse_mode", parseMode,
				)
				if err, ok := err.(tgbotapi.Error); ok && err.Message == "Forbidden: bot was blocked by the user" {
					if err := a.store.UnsubscribeAll(id); err != nil {
						a.log.Errorw(
							"failed to unsubscribe for blocked user",
							"chat_id", id,
							"err", err.Error(),
						)
					}
					break
				}
			}
		}
	}

	return nil
}

func (c *callback) handleBroadcast(data string) error {
	parts := strings.SplitN(data, "|", 2)
	if len(parts) != 2 {
		return errors.New("invalid broadcast callback data")
	}
	appID, appName := parts[0], parts[1]
	text := "Message was sent from app `" + appName + "` with ID `" + appID + "`."
	msg := tgbotapi.NewMessage(c.cb.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = c.cb.Message.MessageID
	_, err := c.base.bot.Send(msg)
	return err
}

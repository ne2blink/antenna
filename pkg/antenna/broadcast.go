package antenna

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Broadcast sends text to all subscribers
func (a *Antenna) Broadcast(appID, text, parseMode string) error {
	ids, err := a.store.ListSubscribers(appID)
	if err != nil {
		return err
	}

	for _, id := range ids {
		msg := tgbotapi.NewMessage(id, text)
		msg.ParseMode = parseMode
		_, err := a.bot.Send(msg)
		if err != nil {
			a.log.Warnw(
				"failed to send message",
				"text", text,
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
			}
		}
	}

	return nil
}

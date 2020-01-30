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

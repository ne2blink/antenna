package antenna

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type handler struct {
	base *Antenna
	msg  *tgbotapi.Message
	log  *zap.SugaredLogger
}

func (h *handler) handle() error {
	cmd := h.msg.Command()
	switch cmd {
	case "start":
		return h.replyMarkdown("Hi, I'm `Antenna`. Try /all to subscribe applications.")
	case "stop":
		return h.handleStop()
	case "all":
		return h.handleAll()
	case "list":
		return h.handleList()
	case "add":
		return h.replyText("Try /all command. Then try /add_<app_id> commands.")
	case "del":
		return h.replyText("Try /del_<app_id> to unsubscribe an application.")
	}
	switch {
	case strings.HasPrefix(cmd, "add_"):
		return h.handleAdd(strings.TrimPrefix(cmd, "add_"))
	case strings.HasPrefix(cmd, "del_"):
		return h.handleDelete(strings.TrimPrefix(cmd, "del_"))
	}
	return nil
}

func (h *handler) handleStop() error {
	if err := h.base.store.UnsubscribeAll(h.msg.Chat.ID); err != nil {
		return err
	}
	return h.replyText("Unsubscribed all applications.")
}

func (h *handler) handleAll() error {
	apps, err := h.base.store.ListApps()
	if err != nil {
		return err
	}

	lines := []string{
		"Click one application to subscribe:",
	}
	for _, app := range apps {
		lines = append(lines, fmt.Sprintf("/add_%s %s", app.ID, app.Name))
	}

	return h.replyText(strings.Join(lines, "\n"))
}

func (h *handler) handleList() error {
	apps, err := h.base.store.ListSubscribedApps(h.msg.Chat.ID)
	if err != nil {
		return err
	}

	lines := []string{
		"Click one application to unsubscribe:",
	}
	for _, app := range apps {
		lines = append(lines, fmt.Sprintf("/del_%s %s", app.ID, app.Name))
	}

	return h.replyText(strings.Join(lines, "\n"))
}

func (h *handler) handleAdd(ID string) error {
	app, err := h.base.store.GetApp(ID)
	if err != nil {
		return err
	}
	if err := h.base.store.Subscribe(h.msg.Chat.ID, ID); err != nil {
		return err
	}
	return h.replyTextWithQuote(fmt.Sprintf(
		"Successfully subscribed %s.",
		app.Name,
	))
}

func (h *handler) handleDelete(ID string) error {
	app, err := h.base.store.GetApp(ID)
	if err != nil {
		return err
	}
	if err := h.base.store.Unsubscribe(h.msg.Chat.ID, ID); err != nil {
		return err
	}
	return h.replyTextWithQuote(fmt.Sprintf(
		"Successfully unsubscribed %s.",
		app.Name,
	))
}

func (h *handler) replyText(text string) error {
	return h.replyMessage(text, "", 0)
}

func (h *handler) replyTextWithQuote(text string) error {
	return h.replyMessage(text, "", h.msg.MessageID)
}

func (h *handler) replyMarkdown(text string) error {
	return h.replyMessage(text, "Markdown", 0)
}

func (h *handler) replyMessage(text, parseMode string, msgID int) error {
	msg := tgbotapi.NewMessage(h.msg.Chat.ID, text)
	msg.ParseMode = parseMode
	msg.ReplyToMessageID = msgID
	_, err := h.base.bot.Send(msg)
	return err
}

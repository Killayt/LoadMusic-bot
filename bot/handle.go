package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TelegramBot) handleUpdates(update tgbotapi.Update) {
	if m := update.Message; m != nil {
		if m.IsCommand() && m.Command() == "start" {
			t.send(m.Chat.ID, "Wats up? I can download youtube video and convert it to audio mp3 format")
			return
		}

		if t.downloadService.IsValidURL(m.Text) {
			t.queue.Enqueue(m)
			return
		}

		t.send(m.Chat.ID, "Invalid message text. I'm waiting for youtube link.")
	}
}

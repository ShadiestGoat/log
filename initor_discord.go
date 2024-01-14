package log

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Send warnings/errors through a discord webhook
func NewLoggerDiscordWebhook(cPrefix, webhook string) LogCB {
	return func() (logger DoLog, closer Closer) {
		logger = func(lvl LogLevel, prefix, msg string) {
			if lvl <= LL_SUCCESS {
				return
			}

			content := prefix + "\n```\n" + msg + "```"

			if cPrefix != "" {
				content = cPrefix + " " + content
			}

			buf, _ := json.Marshal(discordWebhook{
				Content: content,
			})

			http.Post(webhook, "application/json", bytes.NewReader(buf))
		}

		return
	}
}

type discordWebhook struct {
	Content string `json:"content"`
}

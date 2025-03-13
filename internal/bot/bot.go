package bot

import (
    "fmt"
    "log"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
    bot       *tgbotapi.BotAPI
    token     string
    channelID int64
}

func New(token string, channelID int64) Bot {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Fatal("Error creating bot:", err)
    }
    log.Printf("Authorized on bot %s", bot.Self.UserName)

    return Bot{
        bot:       bot,
        token:     token,
        channelID: channelID,
    }
}

func (b *Bot) SendMessage(message string) {
    timer := time.NewTimer(10 * time.Minute)
    defer timer.Stop()

    msg := tgbotapi.NewMessage(b.channelID, message)
    msg.ParseMode = tgbotapi.ModeMarkdown
    msg.DisableWebPagePreview = true

    _, err := b.bot.Send(msg)
    if err != nil {
        fmt.Println("Error sending message:", err)
        for err != nil {
            time.Sleep(15 * time.Second)

            select {
            case <-timer.C:
                return
            default:
                _, err = b.bot.Send(msg)
            }
        }
    }
}

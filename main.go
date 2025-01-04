package main

import (
    "context"
    "time"
    "os"

    "github.com/fxckcode/BussinessBot/ai"
    "github.com/fxckcode/BussinessBot/cmd"
    "github.com/fxckcode/BussinessBot/env"
    "github.com/sirupsen/logrus"
    tele "gopkg.in/telebot.v4"
)

var (
    ctx            = context.Background()
    modelAI string = "gemini-2.0-flash-exp"
    log            = logrus.New()
)

func init() {
    // Configura logrus
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
    })
    log.SetOutput(os.Stdout)
    log.SetLevel(logrus.InfoLevel)
}

func main() {
    BotToken := env.ViperEnvVariable("TELEGRAM_BOT_TOKEN")
    pref := tele.Settings{
        Token:  BotToken,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    }

    b, err := tele.NewBot(pref)
    if err != nil {
        log.Fatalf("Error creating bot: %v", err)
        return
    }

    b.Handle(tele.OnText, func(c tele.Context) error {
        log.Infof("Received message: %s", c.Text())
        start := time.Now()

        res := ai.SearchGemini(ctx, c.Text(), modelAI)

        duration := time.Since(start)
        log.Infof("Processed message in %s", duration)


        return c.Reply(res, tele.ModeMarkdown)
    })

    cmd.ClearConsole()

    log.Info("Bot is running")
    b.Start()
}
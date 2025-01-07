package main

import (
	"context"
	"os"
	"time"

	"github.com/fxckcode/BussinessBot/ai"
	"github.com/fxckcode/BussinessBot/api/models"
	"github.com/fxckcode/BussinessBot/api/routes"
	"github.com/fxckcode/BussinessBot/cmd"
	"github.com/fxckcode/BussinessBot/db"
	"github.com/fxckcode/BussinessBot/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v4"
)

var (
    ctx            = context.Background()
    modelAI string = "gemini-2.0-flash-exp"
    log            = logrus.New()
    PORT = env.ViperEnvVariable("PORT")
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

func startBot() {
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

    log.Info("Bot is running")
    b.Start()
}

func startAPI() {
    app := fiber.New()
    app.Use(logger.New())
    routes.TasksRoutes(app)

    db.DBConnection()
    db.DB.AutoMigrate(models.Task{})

    log.Infof("API is running on port %s", PORT)
    app.Listen(PORT)
}

func main() {
    cmd.ClearConsole()

    go startBot()
    startAPI()
}
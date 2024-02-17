package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/has1985/news-go-bot/internal/bot"
	"github.com/has1985/news-go-bot/internal/bot/middleware"
	"github.com/has1985/news-go-bot/internal/botkit"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/has1985/news-go-bot/internal/fetcher"
	"github.com/has1985/news-go-bot/internal/notifier"
	"github.com/has1985/news-go-bot/internal/storage"
	"github.com/has1985/news-go-bot/internal/summary"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := NewLogger()

	botAPI, err := tgbotapi.NewBotAPI(telegramBotToken())
	if err != nil {
		logger.Printf("[ERROR] failed to create botAPI: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", databaseDSN())
	if err != nil {
		logger.Printf("[ERROR] failed to connect to db: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			fetchInterval(),
			filterKeywords(),
		)
		summarizer = summary.NewOpenAISummarizer(
			openAIKey(),
			openAIModel(),
			openAIPrompt(),
		)
		notifier = notifier.New(
			articleStorage,
			summarizer,
			botAPI,
			notificationInterval(),
			2000*fetchInterval(),
			telegramChannelID(),
		)
	)

	newsBot := botkit.NewBot(botAPI)
	//newsBot.RegisterCmdView("start",middleware.AdminsOnly(
	//	config.Get().TelegramChannelID,
	//	bot.ViewCmdStart(),
	//),)
	newsBot.RegisterCmdView(
		"addsource",
		middleware.AdminsOnly(
			telegramChannelID(),
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"setpriority",
		middleware.AdminsOnly(
			telegramChannelID(),
			bot.ViewCmdSetPriority(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"getsource",
		middleware.AdminsOnly(
			telegramChannelID(),
			bot.ViewCmdGetSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"listsources",
		middleware.AdminsOnly(
			telegramChannelID(),
			bot.ViewCmdListSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"deletesource",
		middleware.AdminsOnly(
			telegramChannelID(),
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctxl := ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))

	///////////////////////

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Printf("[ERROR] failed to run fetcher: %v", err)
				return
			}

			logger.Printf("[INFO] fetcher stopped")
		}
	}(ctxl)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to run notifier: %v", err)
				return
			}

			logger.Printf("[INFO] notifier stopped")
		}
	}(ctxl)

	go func(ctx context.Context) {
		if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Printf("[ERROR] failed to run http server: %v", err)
				return
			}

			logger.Printf("[INFO] http server stopped")
		}
	}(ctxl)

	if err := newsBot.Run(ctxl); err != nil {
		logger.Printf("[ERROR] failed to run botkit: %v", err)
	}
}

func NewLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Set the log level on the default logger based on command line flag
	logLevels := map[string]logrus.Level{
		"debug":   logrus.DebugLevel,
		"info":    logrus.InfoLevel,
		"warning": logrus.WarnLevel,
		"error":   logrus.ErrorLevel,
		"fatal":   logrus.FatalLevel,
		"panic":   logrus.PanicLevel,
	}
	if level, ok := logLevels[loggingLevel()]; !ok {
		logger.Errorf("Invalid %q provided for log level", loggingLevel())
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return logger
}

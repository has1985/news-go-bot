package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

const (
	defaultConfigDirectory = "."
	defaultConfigFile      = ""

	// Logging
	defaultLoggingLevel = "debug"

	defaultTelegramBotToken     = "6845544854:AAFIIdeLVjFG7BEgtKyUKrsFLjGnBkSXdBA"
	defaultTelegramChannelID    = -1002090424799
	defaultDatabaseDSN          = "postgres://postgres:postgres@localhost:5435/news-go-bot?sslmode=disable"
	defaultFetchInterval        = time.Minute * 10
	defaultNotificationInterval = time.Minute * 1
	defaultOpenAIKey            = ""
	defaultOpenAIPrompt         = ""
	defaultOpenAIModel          = "gpt-3.5-turbo"
)

var (
	defaultFilterKeywords = []string{}
)

var (
	flagConfigDirectory = pflag.String("config.source", defaultConfigDirectory, "directory of the configuration file")
	flagConfigFile      = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")

	flagLoggingLevel = pflag.String("logging.level", defaultLoggingLevel, "log level of application")

	flagTelegramBotToken     = pflag.String("telegram.bot.token", defaultTelegramBotToken, "token of the telegram bot")
	flagTelegramChannelID    = pflag.Int("telegram.channel.id", defaultTelegramChannelID, "id of the telegram channel")
	flagDatabaseDSN          = pflag.String("database.dsn", defaultDatabaseDSN, "databaseDSN")
	flagFetchInterval        = pflag.Duration("fetch.interval", defaultFetchInterval, "fetch interval")
	flagNotificationInterval = pflag.Duration("notification.interval", defaultNotificationInterval, "notification interval")
	flagFilterKeywords       = pflag.StringSlice("filter.keywords", defaultFilterKeywords, "filter keywords")
	flagOpenAIKey            = pflag.String("openai.key", defaultOpenAIKey, "openAI key")
	flagOpenAIPrompt         = pflag.String("openai.prompt", defaultOpenAIPrompt, "openAI prompt")
	flagOpenAIModel          = pflag.String("openai.model", defaultOpenAIModel, "openAI model")
)

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(configSource())
	if configFile() != "" {
		log.Printf("Serving from configuration file: %s", configFile())
		viper.SetConfigName(configFile())
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("cannot load configuration: %v", err)
		}
	} else {
		log.Printf("Serving from default values, environment variables, and/or flags")
	}
}

func configSource() string {
	return viper.GetString("config.source")
}

func configFile() string {
	return viper.GetString("config.file")
}

func loggingLevel() string {
	return viper.GetString("logging.level")
}

func telegramBotToken() string {
	return viper.GetString("telegram.bot.token")
}

func telegramChannelID() int64 {
	return viper.GetInt64("telegram.channel.id")
}

func databaseDSN() string {
	return viper.GetString("database.dsn")
}

func fetchInterval() time.Duration {
	return viper.GetDuration("fetch.interval")
}

func notificationInterval() time.Duration {
	return viper.GetDuration("notification.interval")
}

func filterKeywords() []string {
	return viper.GetStringSlice("filter.keywords")
}

func openAIKey() string {
	return viper.GetString("openai.key")
}

func openAIPrompt() string {
	return viper.GetString("openai.prompt")
}

func openAIModel() string {
	return viper.GetString("openai.model")
}

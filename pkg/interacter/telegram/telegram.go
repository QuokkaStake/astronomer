package telegram

import (
	"gopkg.in/telebot.v3/middleware"
	datafetcher "main/pkg/data_fetcher"
	databasePkg "main/pkg/database"
	"main/pkg/templates"
	"main/pkg/types"
	"time"

	"main/pkg/utils"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type Interacter struct {
	Token  string
	Admins []int64

	Version string

	TelegramBot     *tele.Bot
	Logger          zerolog.Logger
	DataFetcher     *datafetcher.DataFetcher
	Database        *databasePkg.Database
	Chains          types.Chains
	TemplateManager templates.Manager
}

const (
	MaxMessageSize = 4096
)

func NewInteracter(
	config types.TelegramConfig,
	version string,
	logger *zerolog.Logger,
	dataFetcher *datafetcher.DataFetcher,
	database *databasePkg.Database,
	chains types.Chains,
) *Interacter {
	return &Interacter{
		Token:           config.Token,
		Admins:          config.Admins,
		Logger:          logger.With().Str("component", "telegram_interacter").Logger(),
		Version:         version,
		DataFetcher:     dataFetcher,
		Database:        database,
		Chains:          chains,
		TemplateManager: templates.NewTelegramTemplatesManager(logger),
	}
}

func (interacter *Interacter) Init() {
	if interacter.Token == "" {
		interacter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram interacter")
		return
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  interacter.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		interacter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return
	}

	// bot.Handle("/start", interacter.HandleHelp)
	// bot.Handle("/help", interacter.HandleHelp)
	// bot.Handle("/subscribe", interacter.HandleSubscribe)
	// bot.Handle("/unsubscribe", interacter.HandleUnsubscribe)
	// bot.Handle("/status", interacter.HandleStatus)
	//bot.Handle("/validators", interacter.HandleListValidators)
	//bot.Handle("/missing", interacter.HandleMissingValidators)
	//bot.Handle("/notifiers", interacter.HandleNotifiers)
	//bot.Handle("/params", interacter.HandleParams)
	bot.Handle("/validator", interacter.HandleValidator)

	if len(interacter.Admins) > 0 {
		interacter.Logger.Debug().Msg("Using admins whitelist")
		bot.Use(middleware.Whitelist(interacter.Admins...))
	}

	bot.Handle("/chain_bind", interacter.HandleChainBind)

	interacter.TelegramBot = bot
}

func (interacter *Interacter) Start() {
	interacter.TelegramBot.Start()
}

func (interacter *Interacter) Enabled() bool {
	return interacter.Token != ""
}

func (interacter *Interacter) Name() string {
	return "telegram"
}

func (interacter *Interacter) BotReply(c tele.Context, msg string) error {
	messages := utils.SplitStringIntoChunks(msg, MaxMessageSize)

	for _, message := range messages {
		if err := c.Reply(message, tele.ModeHTML, tele.NoPreview); err != nil {
			interacter.Logger.Error().Err(err).Msg("Could not send Telegram message")
			return err
		}
	}
	return nil
}

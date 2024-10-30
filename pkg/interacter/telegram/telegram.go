package telegram

import (
	datafetcher "main/pkg/data_fetcher"
	databasePkg "main/pkg/database"
	"main/pkg/metrics"
	"main/pkg/templates"
	"main/pkg/types"
	"strconv"
	"time"

	"gopkg.in/telebot.v3/middleware"

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
	MetricsManager  *metrics.Manager
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
	metricsManager *metrics.Manager,
) *Interacter {
	return &Interacter{
		Token:           config.Token,
		Admins:          config.Admins,
		Logger:          logger.With().Str("component", "telegram_interacter").Logger(),
		Version:         version,
		DataFetcher:     dataFetcher,
		Database:        database,
		TemplateManager: templates.NewTelegramTemplatesManager(logger),
		MetricsManager:  metricsManager,
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
		interacter.Logger.Panic().Err(err).Msg("Could not create Telegram bot")
	}

	interacter.AddCommand("/start", bot, interacter.GetHelpCommand())
	interacter.AddCommand("/help", bot, interacter.GetHelpCommand())
	interacter.AddCommand("/validator", bot, interacter.GetValidatorCommand())
	interacter.AddCommand("/validators", bot, interacter.GetValidatorsCommand())
	interacter.AddCommand("/params", bot, interacter.GetParamsCommand())
	interacter.AddCommand("/proposal", bot, interacter.GetSingleProposalCommand())
	interacter.AddCommand("/proposals", bot, interacter.GetActiveProposalsCommand())
	interacter.AddCommand("/wallet_link", bot, interacter.GetWalletLinkCommand())
	interacter.AddCommand("/wallet_unlink", bot, interacter.GetWalletUnlinkCommand())
	interacter.AddCommand("/validator_link", bot, interacter.GetValidatorLinkCommand())
	interacter.AddCommand("/validator_unlink", bot, interacter.GetValidatorUnlinkCommand())
	interacter.AddCommand("/wallets", bot, interacter.GetWalletsCommand())
	interacter.AddCommand("/chains", bot, interacter.GetChainsListCommand())
	interacter.AddCommand("/balance", bot, interacter.GetBalanceCommand())
	interacter.AddCommand("/supply", bot, interacter.GetSupplyCommand())

	if len(interacter.Admins) > 0 {
		interacter.Logger.Debug().Msg("Using admins whitelist")
		bot.Use(middleware.Whitelist(interacter.Admins...))
	}

	interacter.AddCommand("/chain_bind", bot, interacter.GetChainBindCommand())
	interacter.AddCommand("/chain_unbind", bot, interacter.GetChainUnbindCommand())
	interacter.AddCommand("/chain_add", bot, interacter.GetChainAddCommand())
	interacter.AddCommand("/chain_update", bot, interacter.GetChainUpdateCommand())
	interacter.AddCommand("/chain_delete", bot, interacter.GetChainDeleteCommand())
	interacter.AddCommand("/explorer_add", bot, interacter.GetExplorerAddCommand())
	interacter.AddCommand("/explorer_delete", bot, interacter.GetExplorerDeleteCommand())
	interacter.AddCommand("/denom_add", bot, interacter.GetDenomAddCommand())
	interacter.AddCommand("/denom_delete", bot, interacter.GetDenomDeleteCommand())
	interacter.AddCommand("/lcd_add", bot, interacter.GetLCDAddCommand())
	interacter.AddCommand("/lcd_delete", bot, interacter.GetLCDDeleteCommand())

	interacter.TelegramBot = bot
}

func (interacter *Interacter) AddCommand(query string, bot *tele.Bot, command Command) {
	bot.Handle(query, func(c tele.Context) error {
		interacter.Logger.Info().
			Str("sender", c.Sender().Username).
			Str("text", c.Text()).
			Str("command", command.Name).
			Msg("Got query")

		interacter.MetricsManager.LogReporterQuery(interacter.Name(), command.Name)

		userID := strconv.FormatInt(c.Sender().ID, 10)
		chatID := strconv.FormatInt(c.Chat().ID, 10)

		queryToInsert := &types.Query{
			Reporter: interacter.Name(),
			UserID:   userID,
			Username: c.Sender().Username,
			ChatID:   chatID,
			Command:  command.Name,
			Query:    c.Text(),
		}

		if err := interacter.Database.InsertQuery(queryToInsert); err != nil {
			interacter.Logger.Error().Err(err).Msg("Error inserting query info")
			return interacter.BotReply(c, "Internal error!")
		}

		chainBinds, err := interacter.Database.GetAllChainBinds(chatID)
		if err != nil {
			interacter.Logger.Error().Err(err).Msg("Error getting chain binds")
			return interacter.BotReply(c, "Internal error!")
		}

		result, err := command.Execute(c, chainBinds)
		if err != nil {
			interacter.Logger.Error().
				Err(err).
				Str("command", command.Name).
				Msg("Error processing command")
			if result != "" {
				return interacter.BotReply(c, result)
			} else {
				return interacter.BotReply(c, "Internal error!")
			}
		}

		return interacter.BotReply(c, result)
	})
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

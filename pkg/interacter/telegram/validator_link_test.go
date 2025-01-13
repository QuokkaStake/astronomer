package telegram

import (
	"errors"
	"main/assets"
	converterPkg "main/pkg/converter"
	datafetcher "main/pkg/data_fetcher"
	databasePkg "main/pkg/database"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/tendermint"
	"main/pkg/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkInvalidInvocation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Usage: /validator_link &lt;chain&gt; &lt;address&gt;"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	dataFetcher := datafetcher.NewDataFetcher(logger, database, nil, metricsManager, nil)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkErrorFetchingChain(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Internal error!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	dataFetcher := datafetcher.NewDataFetcher(logger, database, nil, metricsManager, nil)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnError(errors.New("custom error"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkChainNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/chain-not-found.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	dataFetcher := datafetcher.NewDataFetcher(logger, database, nil, metricsManager, nil)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE").
		WillReturnRows(sqlmock.NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkErrorFetchingValidator(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators/valoper",
		httpmock.NewErrorResponder(errors.New("custom error")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error linking validator: could not get data after 3 attempts"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	converter := converterPkg.NewConverter()
	nodesManager := tendermint.NewNodeManager(logger, database, converter, metricsManager)
	dataFetcher := datafetcher.NewDataFetcher(logger, database, converter, metricsManager, nodesManager)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkErrorInserting(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators/valoper",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validator.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Internal error!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	converter := converterPkg.NewConverter()
	nodesManager := tendermint.NewNodeManager(logger, database, converter, metricsManager)
	dataFetcher := datafetcher.NewDataFetcher(logger, database, converter, metricsManager, nodesManager)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))

	mock.ExpectExec("INSERT INTO validator_links").
		WillReturnError(errors.New("custom error"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkDuplicate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators/valoper",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validator.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("You have already linked this validator!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	converter := converterPkg.NewConverter()
	nodesManager := tendermint.NewNodeManager(logger, database, converter, metricsManager)
	dataFetcher := datafetcher.NewDataFetcher(logger, database, converter, metricsManager, nodesManager)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))

	mock.ExpectExec("INSERT INTO validator_links").
		WillReturnError(errors.New("duplicate key value violates unique constraint \"validator_links_validator_address_chain_key\""))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramValidatorLinkOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators/valoper",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validator.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/validator-link.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, types.MetricsConfig{})
	database := databasePkg.NewDatabase(logger, types.DatabaseConfig{})
	converter := converterPkg.NewConverter()
	nodesManager := tendermint.NewNodeManager(logger, database, converter, metricsManager)
	dataFetcher := datafetcher.NewDataFetcher(logger, database, converter, metricsManager, nodesManager)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains WHERE").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))

	mock.ExpectExec("INSERT INTO validator_links").
		WillReturnResult(sqlmock.NewResult(1, 1))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator_link chain valoper",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator_link", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

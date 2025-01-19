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
	timePkg "main/pkg/time"
	"main/pkg/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestValidatorSearchInvalidInvocation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Usage: /validator [chain] &lt;query&gt;"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}))
	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestValidatorSearchErrorFetchingChain(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error searching for validator: custom error"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}))

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
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator chain quokka",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestValidatorSearchErrorFetchingExplorers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error searching for validator: custom error"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	logger := loggerPkg.GetNopLogger()
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers").
		WillReturnError(errors.New("custom error"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator chain quokka",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestValidatorSearchNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("<strong>Chain</strong>\nNo validator found."),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators?pagination.count_total=true&pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validators.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/slashing/v1beta1/signing_infos?pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("signing-infos.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/slashing/v1beta1/params",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("slashing-params.json")))

	logger := loggerPkg.GetNopLogger()
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain",
				"name",
				"proposal_link_pattern",
				"wallet_link_pattern",
				"validator_link_pattern",
				"main_link",
			}).AddRow("chain", "Ping", "", "", "https://example.com/validators/%s", ""),
		)

	for range 3 {
		mock.ExpectQuery("SELECT host FROM lcd").
			WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))
	}

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator chain asdascaxcasda",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestValidatorSearchOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/validator.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators?pagination.count_total=true&pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validators.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/slashing/v1beta1/signing_infos?pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("signing-infos.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/slashing/v1beta1/params",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("slashing-params.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=cosmos&vs_currencies=usd",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("coingecko.json")))

	logger := loggerPkg.GetNopLogger()
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain",
				"name",
				"proposal_link_pattern",
				"wallet_link_pattern",
				"validator_link_pattern",
				"main_link",
			}).AddRow("chain", "Ping", "", "", "https://example.com/validators/%s", ""),
		)

	for range 3 {
		mock.ExpectQuery("SELECT host FROM lcd").
			WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))
	}

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/validator chain quokka",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/validator", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

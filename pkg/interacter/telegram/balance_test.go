package telegram

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
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
	"time"
)

//nolint:paralleltest // disabled
func TestTelegramBalanceErrorFetchingWallets(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error getting wallets balances: custom error"),
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
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
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/balance", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramBalanceErrorFetchingChains(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error getting wallets balances: custom error"),
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
		WillReturnRows(sqlmock.NewRows([]string{"chain", "reporter", "user_id", "address", "alias"}))

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
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/balance", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramBalanceErrorFetchingExplorers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error getting wallets balances: custom error"),
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain", "reporter", "user_id", "address", "alias"}).
			AddRow("chain", "reporter", "1", "address", "alias"),
		)

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "ustake", "chainvaloper"))

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
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/balance", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramBalanceChainNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("❌ Error getting wallets balances: custom error"),
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain", "reporter", "user_id", "address", "alias"}).
			AddRow("otherchain", "reporter", "1", "address", "alias"),
		)

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "ustake", "chainvaloper"))

	mock.ExpectQuery("SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers").
		WillReturnRows(sqlmock.NewRows([]string{"chain", "name", "proposal_link_pattern", "wallet_link_pattern", "validator_link_pattern", "main_link"}))

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
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	interacter.TelegramBot.Trigger("/balance", ctx)
}

//nolint:paralleltest // disabled
func TestTelegramBalanceNoWallets(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("You are not subscribed to any wallets."),
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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
		WillReturnRows(sqlmock.NewRows([]string{"chain", "reporter", "user_id", "address", "alias"}))

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "ustake", "chainvaloper"))

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
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/balance", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramBalanceOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/balance.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=cosmos&vs_currencies=usd",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("coingecko.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/bank/v1beta1/balances/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("balance.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/delegations/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2?pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("delegation.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/delegators/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2/redelegations?pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("redelegation.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/delegators/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2/unbonding_delegations?pagination.limit=1000",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("unbond.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/distribution/v1beta1/validators/cosmosvaloper1xqz9pemz5e5zycaa89kys5aw6m8rhgsvw4328e/commission",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("commission.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/distribution/v1beta1/delegators/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2/rewards",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("rewards.json")))

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/validators/cosmosvaloper1xqz9pemz5e5zycaa89kys5aw6m8rhgsvw4328e",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("validator.json")))

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
		WillReturnRows(sqlmock.NewRows([]string{"chain"}).AddRow("chain1").AddRow("chain2"))

	mock.ExpectQuery("SELECT chain, reporter, user_id, address, alias FROM wallet_links").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain", "reporter", "user_id", "address", "alias"}).
			AddRow("chain", "reporter", "1", "cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2", "Wallet without error").
			AddRow("chain", "reporter", "1", "cosmos1rxvkwfw3467nxgs6r7yav6cnygkjzkkc0edu0f", "Wallet with error").
			AddRow("chain", "reporter", "1", "notok", "Wrong Bech2 prefix wallet"),
		)

	mock.ExpectQuery("SELECT name, pretty_name, base_denom, bech32_validator_prefix FROM chains").
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"))

	mock.ExpectQuery("SELECT chain, name, proposal_link_pattern, wallet_link_pattern, validator_link_pattern, main_link FROM explorers").
		WillReturnRows(sqlmock.
			NewRows([]string{
				"chain",
				"name",
				"proposal_link_pattern",
				"wallet_link_pattern",
				"validator_link_pattern",
				"main_link",
			}).
			AddRow("chain", "Ping", "", "https://example.com/wallet/%s", "https://example.com/validator/%s", ""))

	// 3x6 per each wallet - 1 when bech32 conversion failed
	for _ = range 17 {
		mock.ExpectQuery("SELECT host FROM lcd").
			WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))
	}

	mock.ExpectQuery("SELECT chain, denom, display_denom, denom_exponent, coingecko_currency, ignored FROM denoms").
		WillReturnRows(sqlmock.
			NewRows([]string{"chain", "denom", "display_denom", "denom_exponent", "coingecko_currency", "ignored"}).
			AddRow("chain", "uatom", "ATOM", 6, "cosmos", false),
		)

	for _ = range 4 {
		mock.ExpectQuery("SELECT host FROM lcd").
			WillReturnRows(sqlmock.NewRows([]string{"host"}).AddRow("https://example.com"))
	}

	database.SetClient(db)

	renderTime, err := time.Parse(time.RFC3339, "2025-01-17T23:49:00Z")
	require.NoError(t, err)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		dataFetcher,
		database,
		metricsManager,
		&timePkg.StubTime{NowTime: renderTime},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/balance",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/balance", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

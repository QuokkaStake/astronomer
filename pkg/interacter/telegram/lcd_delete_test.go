package telegram

import (
	"errors"
	"main/assets"
	datafetcher "main/pkg/data_fetcher"
	databasePkg "main/pkg/database"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	timePkg "main/pkg/time"
	"main/pkg/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteNotEnoughArgs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Usage: /lcd_delete &lt;chain name&gt; &lt;host&gt;"),
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
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/lcd_delete",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteErrorFetchingChain(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error finding chain!"),
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
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteErrorFetchingLCDHosts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error finding LCD hosts!"),
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
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
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
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteDeletingSingleLCDHost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Cannot remove the only chain LCD!"),
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
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser", ID: 1},
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteErrorDeleting(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error deleting LCD host!"),
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
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.
			NewRows([]string{"host"}).
			AddRow("https://example.com").
			AddRow("https://example2.com"),
		)

	mock.ExpectExec("DELETE FROM lcd").
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
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteLCDNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Chain LCD host was not found!"),
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
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.
			NewRows([]string{"host"}).
			AddRow("https://example.com").
			AddRow("https://example2.com"),
		)

	mock.ExpectExec("DELETE FROM lcd").
		WillReturnResult(sqlmock.NewResult(1, 0))

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
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramLCDDeleteOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/lcd-delete.html")),
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
		WillReturnRows(sqlmock.
			NewRows([]string{"name", "pretty_name", "base_denom", "bech32_validator_prefix"}).
			AddRow("chain", "Chain", "uatom", "cosmosvaloper"),
		)

	mock.ExpectQuery("SELECT host FROM lcd").
		WillReturnRows(sqlmock.
			NewRows([]string{"host"}).
			AddRow("https://example.com").
			AddRow("https://example2.com"),
		)

	mock.ExpectExec("DELETE FROM lcd").
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			Text:   "/lcd_delete chain https://example.com",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/lcd_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

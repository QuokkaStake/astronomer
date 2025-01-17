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
func TestTelegramChainDeleteNotEnoughArgs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Usage: /chain_delete &lt;chain name&gt;"),
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
			Text:   "/chain_delete",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/chain_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramChainDeleteErrorDeleting(t *testing.T) {
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

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM lcd").WillReturnError(errors.New("custom error"))

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
			Text:   "/chain_delete chain",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/chain_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramChainDeleteChainNotExisting(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Chain was not found!"),
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

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM lcd").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM chains").WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()

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
			Text:   "/chain_delete chain",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/chain_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramChainDeleteOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Successfully deleted chain!"),
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

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM lcd").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM chains").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

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
			Text:   "/chain_delete chain",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/chain_delete", ctx)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

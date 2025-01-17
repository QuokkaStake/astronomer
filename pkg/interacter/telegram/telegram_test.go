package telegram

import (
	"errors"
	"main/assets"
	databasePkg "main/pkg/database"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	timePkg "main/pkg/time"
	"main/pkg/types"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"

	"github.com/jarcoal/httpmock"
)

func TestTelegramInitNoTokenProvided(t *testing.T) {
	t.Parallel()

	interacter := NewInteracter(
		types.TelegramConfig{},
		"v1.2.3",
		loggerPkg.GetNopLogger(),
		nil,
		nil,
		nil,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	require.False(t, interacter.Enabled())
	require.Equal(t, "telegram", interacter.Name())
}

//nolint:paralleltest // disabled
func TestTelegramInitCannotFetchBot(t *testing.T) {
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
		httpmock.NewErrorResponder(errors.New("custom error")))

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy"},
		"v1.2.3",
		loggerPkg.GetNopLogger(),
		nil,
		nil,
		nil,
		&timePkg.SystemTime{},
	)
	interacter.Init()
}

//nolint:paralleltest // disabled
func TestTelegramStartOkay(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		loggerPkg.GetNopLogger(),
		nil,
		nil,
		nil,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	go interacter.Start()
	interacter.Stop()
}

//nolint:paralleltest // disabled
func TestTelegramSendMultilineFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewErrorResponder(errors.New("custom error")))

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		loggerPkg.GetNopLogger(),
		nil,
		nil,
		nil,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err := interacter.BotReply(ctx, strings.Repeat("a", 5000))
	require.Error(t, err)

	err = interacter.BotReply(ctx, strings.Repeat("a", 10))
	require.Error(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramSendMultilineOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")))

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		loggerPkg.GetNopLogger(),
		nil,
		nil,
		nil,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err := interacter.BotReply(ctx, strings.Repeat("a", 5000))
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramAddCommandFailedToInsertQuery(t *testing.T) {
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

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").WillReturnError(errors.New("custom error"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		nil,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/help", ctx)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramAddCommandFailedToFetchChains(t *testing.T) {
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

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("INSERT INTO queries").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT chain FROM chain_binds").
		WillReturnError(errors.New("custom error"))

	database.SetClient(db)

	interacter := NewInteracter(
		types.TelegramConfig{Token: "xxx:yyy", Admins: []int64{1, 2}},
		"v1.2.3",
		logger,
		nil,
		database,
		metricsManager,
		&timePkg.SystemTime{},
	)
	interacter.Init()

	ctx := interacter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = interacter.TelegramBot.Trigger("/help", ctx)
	require.NoError(t, err)
}

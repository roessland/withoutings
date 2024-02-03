package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/roessland/withoutings/pkg/logging"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"net/url"
)

type FetchNotificationData struct {
	Notification *subscription.Notification
}

type FetchNotificationDataHandler interface {
	Handle(ctx context.Context, cmd FetchNotificationData) error
}

func (h fetchNotificationDataHandler) Handle(ctx context.Context, cmd FetchNotificationData) error {
	log := logging.MustGetLoggerFromContext(ctx)
	log = log.WithField("notification_uuid", cmd.Notification.UUID())
	log.WithField("event", "info.command.FetchNotificationData.started").Info()

	if cmd.Notification.DataStatus() != subscription.NotificationDataStatusAwaitingFetch {
		log.WithField("event", "error.command.FetchNotificationData.invalidDataStatus").
			WithField("data_status", cmd.Notification.DataStatus()).
			Error()
		return nil
	}

	// Get account corresponding to the raw notification
	parsedParams, err := subscription.ParseNotificationParams(cmd.Notification.Params())
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.ParseNotificationParams.failed").
			Error()
		return nil
	}
	acc, err := h.accountRepo.GetAccountByUUID(ctx, cmd.Notification.AccountUUID())
	if errors.Is(err, account.ErrAccountNotFound) {
		log.WithField("event", "warn.command.FetchNotificationData.account_not_found").
			WithField("account_uuid", cmd.Notification.AccountUUID()).
			Warn()
		return fmt.Errorf("cannot fetch data for account that does not exist: %w", err)
	} else if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.GetAccountByUUID(.failed").
			Error()
		return fmt.Errorf("failed to get account by uuid: %w", err)
	}

	data, err := h.getAvailableData(ctx, acc, parsedParams)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.getAvailableData.failed").
			Error()
		return fmt.Errorf("failed to get available data: %w", err)
	}

	// Update notification with fetched data
	err = h.subscriptionRepo.UpdateNotification(
		ctx,
		cmd.Notification.UUID(),
		func(ctx context.Context, notification *subscription.Notification) (*subscription.Notification, error) {
			notification.FetchedData(data)
			return notification, nil
		},
	)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.UpdateNotification.failed").
			Error()
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// getAvailableData fetches data from Withings API corresponding to the notification category.
// See https://developer.withings.com/developer-guide/v3/data-api/keep-user-data-up-to-date/
func (h fetchNotificationDataHandler) getAvailableData(ctx context.Context, acc *account.Account, parsedParams subscription.ParsedNotificationParams) ([]byte, error) {
	switch parsedParams.Appli {
	case 1:
		return h.getAvailableData1(ctx, acc, parsedParams)
	case 50:
		return h.getAvailableData50(ctx, acc, parsedParams)
	case 51:
		return h.getAvailableData51(ctx, acc, parsedParams)
	default:
		return nil, fmt.Errorf("not yet able to handle appli: %d", parsedParams.Appli)
	}
}

// getAvailableData1 fetches data from Withings API for appli 1:
// New weight-related data
func (h fetchNotificationDataHandler) getAvailableData1(ctx context.Context, acc *account.Account, parsedParams subscription.ParsedNotificationParams) ([]byte, error) {
	if parsedParams.Appli != 1 {
		panic("getAvailableData1 called with wrong application ID")
	}

	log := logging.MustGetLoggerFromContext(ctx)

	// Fetch data from Withings API
	params := url.Values{
		"action":    {"getmeas"},
		"startdate": {parsedParams.StartDateStr},
		"enddate":   {parsedParams.EndDateStr},
	}
	getmeasResponse, err := h.withingsSvc.MeasureGetmeas(ctx, acc,
		withings.MeasureGetmeasParams(params.Encode()),
	)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.measuregetmeas.failed").
			WithField("getmeasresponse", getmeasResponse).
			Error()
		return nil, fmt.Errorf("failed to call Withings API measuregetmeas: %w", err)
	}

	return getmeasResponse.Raw, nil
}

// getAvailableData50 fetches data from Withings API for appli 50:
// New bed in event (user lies on bed)
func (h fetchNotificationDataHandler) getAvailableData50(_ context.Context, _ *account.Account, parsedParams subscription.ParsedNotificationParams) ([]byte, error) {
	if parsedParams.Appli != 50 {
		panic("getAvailableData50 called with wrong application ID")
	}
	// No service to call
	return []byte("{}"), nil
}

// getAvailableData51 fetches data from Withings API for appli 51:
// New bed out event (user gets out of bed)
func (h fetchNotificationDataHandler) getAvailableData51(_ context.Context, _ *account.Account, parsedParams subscription.ParsedNotificationParams) ([]byte, error) {
	if parsedParams.Appli != 51 {
		panic("getAvailableData51 called with wrong application ID")
	}
	// No service to call
	return []byte("{}"), nil
}

func NewFetchNotificationDataHandler(
	subscriptionsRepo subscription.Repo,
	withingsSvc withingsSvc.Service,
	accountRepo account.Repo,
	publisher message.Publisher,
) FetchNotificationDataHandler {
	return fetchNotificationDataHandler{
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
		accountRepo:      accountRepo,
		publisher:        publisher,
	}
}

type fetchNotificationDataHandler struct {
	subscriptionRepo subscription.Repo
	withingsSvc      withingsSvc.Service
	accountRepo      account.Repo
	publisher        message.Publisher
}

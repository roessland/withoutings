package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/logging"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"net/url"
	"time"
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

	datas, err := h.getAvailableData(ctx, acc, parsedParams)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.getAvailableData.failed").
			Error()
		return fmt.Errorf("failed to get available data: %w", err)
	}

	for _, data := range datas {
		notificationData, err := subscription.NewNotificationData(subscription.NewNotificationDataParams{
			NotificationDataUUID: uuid.New(),
			NotificationUUID:     cmd.Notification.UUID(),
			AccountUUID:          cmd.Notification.AccountUUID(),
			Service:              subscription.NotificationDataService(data.service),
			Data:                 data.data,
			FetchedAt:            data.fetchedAt,
		})
		if err != nil {
			log.WithError(err).
				WithField("event", "error.command.FetchNotificationData.NewNotificationData.failed").
				Error()
			return fmt.Errorf("failed to create notification data: %w", err)
		}

		err = h.subscriptionRepo.StoreNotificationData(ctx, notificationData)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.command.FetchNotificationData.StoreNotificationData.failed").
				Error()
			return fmt.Errorf("failed to store notification data: %w", err)
		}
	}

	// Update notification, marking data status as fetched
	err = h.subscriptionRepo.UpdateNotification(
		ctx,
		cmd.Notification.UUID(),
		func(ctx context.Context, notification *subscription.Notification) (*subscription.Notification, error) {
			notification.FetchedData()
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

type availableDatas []availableData

type availableData struct {
	data      []byte
	fetchedAt time.Time
	service   string
}

// getAvailableData fetches data from Withings API corresponding to the notification category.
// See https://developer.withings.com/developer-guide/v3/data-api/keep-user-data-up-to-date/
func (h fetchNotificationDataHandler) getAvailableData(ctx context.Context, acc *account.Account, parsedParams subscription.ParsedNotificationParams) (availableDatas, error) {
	switch parsedParams.Appli {
	case 1:
		return h.getAvailableData1(ctx, acc, parsedParams)
	case 4:
		return h.getAvailableData4(ctx, acc, parsedParams)
	case 44:
		return h.getAvailableData44(ctx, acc, parsedParams)
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
func (h fetchNotificationDataHandler) getAvailableData1(
	ctx context.Context,
	acc *account.Account,
	parsedParams subscription.ParsedNotificationParams,
) (availableDatas, error) {
	if parsedParams.Appli != 1 {
		panic("getAvailableData1 called with wrong application ID")
	}

	log := logging.MustGetLoggerFromContext(ctx)
	log = log.WithField("appli", parsedParams.Appli)

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

	ad := availableDatas{
		{
			data:      getmeasResponse.Raw,
			fetchedAt: time.Now(),
			service:   "Measure - Getmeas",
		},
	}

	return ad, nil
}

// getAvailableData4 fetches data from Withings API for appli 4:
// New pressure related data
func (h fetchNotificationDataHandler) getAvailableData4(
	ctx context.Context,
	acc *account.Account,
	parsedParams subscription.ParsedNotificationParams,
) (availableDatas, error) {
	if parsedParams.Appli != 4 {
		panic("getAvailableData4 called with wrong application ID")
	}

	log := logging.MustGetLoggerFromContext(ctx)
	log = log.WithField("appli", parsedParams.Appli)

	// Fetch data from Withings API
	params := url.Values{
		"action":    {"getmeas"},
		"startdate": {parsedParams.StartDateStr},
		"enddate":   {parsedParams.EndDateStr},
		"meastypes": {"9,10,11,54"},
		// 9	Diastolic Blood Pressure (mmHg)
		// 10	Systolic Blood Pressure (mmHg)
		// 11	Heart Pulse (bpm) - only for BPM and scale devices
		// 54	SP02 (%)
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

	ad := availableDatas{
		{
			data:      getmeasResponse.Raw,
			fetchedAt: time.Now(),
			service:   "Measure - Getmeas",
		},
	}

	return ad, nil
}

// getAvailableData44 fetches data from Withings API for appli 44:
// New sleep-related data
func (h fetchNotificationDataHandler) getAvailableData44(
	ctx context.Context,
	acc *account.Account,
	parsedParams subscription.ParsedNotificationParams,
) (availableDatas, error) {
	if parsedParams.Appli != 44 {
		panic("getAvailableData44 called with wrong application ID")
	}

	log := logging.MustGetLoggerFromContext(ctx)
	log = log.WithField("appli", parsedParams.Appli)
	ad := availableDatas{}

	getSummaryParams := withings.NewSleepGetsummaryParams()
	getSummaryParams.Startdateymd = parsedParams.StartDate.Format("2006-01-02")
	getSummaryParams.Enddateymd = parsedParams.EndDate.Format("2006-01-02")
	sleepGetsummaryResponse, err := h.withingsSvc.SleepGetsummary(ctx, acc, getSummaryParams)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.SleepGetsummary.failed").
			WithField("SleepGetsummaryResponse", sleepGetsummaryResponse).
			Error()
		return nil, fmt.Errorf("failed to call Withings API SleepGetsummary: %w", err)
	}
	ad = append(ad, availableData{
		data:      sleepGetsummaryResponse.Raw,
		fetchedAt: time.Now(),
		service:   "Sleep v2 - Getsummary", // todo use const
	})

	getParams := withings.NewSleepGetParams()
	getParams.Startdate = parsedParams.StartDate.Unix()
	getParams.Enddate = parsedParams.EndDate.Unix()
	sleepGetResponse, err := h.withingsSvc.SleepGet(ctx, acc, getParams)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.FetchNotificationData.SleepGet.failed").
			WithField("SleepGetResponse", sleepGetResponse).
			Error()
		return nil, fmt.Errorf("failed to call Withings API SleepGet: %w", err)
	}

	ad = append(ad, availableData{
		data:      sleepGetResponse.Raw,
		fetchedAt: time.Now(),
		service:   "Sleep v2 - Get",
	})

	return ad, nil
}

// getAvailableData50 fetches data from Withings API for appli 50:
// New bed in event (user lies on bed)
func (h fetchNotificationDataHandler) getAvailableData50(
	_ context.Context,
	_ *account.Account,
	parsedParams subscription.ParsedNotificationParams,
) (availableDatas, error) {
	if parsedParams.Appli != 50 {
		panic("getAvailableData50 called with wrong application ID")
	}
	// No service to call
	return availableDatas{}, nil
}

// getAvailableData51 fetches data from Withings API for appli 51:
// New bed out event (user gets out of bed)
func (h fetchNotificationDataHandler) getAvailableData51(
	_ context.Context,
	_ *account.Account,
	parsedParams subscription.ParsedNotificationParams,
) (availableDatas, error) {
	if parsedParams.Appli != 51 {
		panic("getAvailableData51 called with wrong application ID")
	}
	// No service to call
	return availableDatas{}, nil
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

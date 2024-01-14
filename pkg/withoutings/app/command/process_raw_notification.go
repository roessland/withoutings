package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/logging"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"time"
)

type ProcessRawNotification struct {
	RawNotification *subscription.RawNotification
}

type ProcessRawNotificationHandler interface {
	Handle(ctx context.Context, cmd ProcessRawNotification) error
}

func (h processRawNotificationHandler) Handle(ctx context.Context, cmd ProcessRawNotification) error {
	log := logging.MustGetLoggerFromContext(ctx)

	if cmd.RawNotification.Status() != subscription.RawNotificationStatusPending {
		log.WithField("event", "error.command.ProcessRawNotification.invalidStatus").
			WithField("status", cmd.RawNotification.Status()).
			Error()
		return nil
	}

	// Get account corresponding to the raw notification
	parsedParams, err := cmd.RawNotification.ParsedData()
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.ProcessRawNotification.parseData.failed").
			Error()
		return nil
	}
	acc, err := h.accountRepo.GetAccountByWithingsUserID(ctx, parsedParams.WithingsUserID)
	if errors.Is(err, account.ErrAccountNotFound) {
		log.WithField("event", "warn.command.ProcessRawNotification.account_not_found").
			WithField("withings_user_id", parsedParams.WithingsUserID).
			Warn()
		return nil
	} else if err != nil {
		log.WithError(err).
			WithField("event", "error.command.ProcessRawNotification.GetAccountByWithingsUserID.failed").
			Error()
		return fmt.Errorf("failed to get account by withings user id: %w", err)
	}

	// Fetch data from Withings API
	getmeasResponse, err := h.withingsSvc.MeasureGetmeas(ctx, acc,
		withings.MeasureGetmeasParams(cmd.RawNotification.Data()),
	)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.ProcessRawNotification.measuregetmeas.failed").
			WithField("getmeasresponse", getmeasResponse).
			Error()
		return fmt.Errorf("failed to call Withings API measuregetmeas: %w", err)
	}

	// type NewNotificationParams struct {
	//	NotificationUUID    uuid.UUID
	//	AccountUUID         uuid.UUID
	//	ReceivedAt          time.Time
	//	Params              NotificationParams
	//	Data                []byte
	//	FetchedAt           time.Time
	//	RawNotificationUUID uuid.UUID
	//	Source              string
	//}

	// Make notification
	notification, err := subscription.NewNotification(
		subscription.NewNotificationParams{
			NotificationUUID:    uuid.New(),
			AccountUUID:         acc.UUID(),
			ReceivedAt:          cmd.RawNotification.ReceivedAt(),
			Params:              cmd.RawNotification.Data(),
			Data:                getmeasResponse.Raw,
			FetchedAt:           time.Now(),
			RawNotificationUUID: cmd.RawNotification.UUID(),
			Source:              cmd.RawNotification.Source(),
		},
	)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.ProcessRawNotification.NewNotification.failed").
			Error()
		return fmt.Errorf("failed to make notification: %w", err)
	}

	// Persist notification
	err = h.subscriptionRepo.CreateNotification(ctx, notification)
	if err != nil {
		log.WithError(err).
			WithField("event", "error.command.ProcessRawNotification.CreateNotification.failed").
			Error()
		return fmt.Errorf("failed to persist notification: %w", err)
	}

	// TODO: Emit NotificationReceived event to withings_notification_received topic

	return nil
}

func NewProcessRawNotificationHandler(
	subscriptionsRepo subscription.Repo,
	withingsSvc withingsSvc.Service,
) ProcessRawNotificationHandler {
	return processRawNotificationHandler{
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
	}
}

type processRawNotificationHandler struct {
	subscriptionRepo subscription.Repo
	withingsSvc      withingsSvc.Service
	accountRepo      account.Repo
}

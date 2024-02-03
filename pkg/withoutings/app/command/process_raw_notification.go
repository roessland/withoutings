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
)

type ProcessRawNotification struct {
	RawNotification *subscription.RawNotification
}

type ProcessRawNotificationHandler interface {
	Handle(ctx context.Context, cmd ProcessRawNotification) error
}

func (h processRawNotificationHandler) Handle(ctx context.Context, cmd ProcessRawNotification) error {
	log := logging.MustGetLoggerFromContext(ctx)
	log = log.WithField("raw_notification_uuid", cmd.RawNotification.UUID())
	log.WithField("event", "info.command.ProcessRawNotification.started").Info()

	if cmd.RawNotification.Status() != subscription.RawNotificationStatusPending {
		log.WithField("event", "error.command.ProcessRawNotification.invalidStatus").
			WithField("status", cmd.RawNotification.Status()).
			Error()
		return nil
	}

	// Get account corresponding to the raw notification
	parsedParams, err := subscription.ParseNotificationParams(cmd.RawNotification.Data())
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

	// Make notification
	notification, err := subscription.NewNotification(
		subscription.NewNotificationParams{
			NotificationUUID:    uuid.New(),
			AccountUUID:         acc.UUID(),
			ReceivedAt:          cmd.RawNotification.ReceivedAt(),
			Params:              cmd.RawNotification.Data(),
			Data:                nil,
			DataStatus:          subscription.NotificationDataStatusAwaitingFetch,
			FetchedAt:           nil,
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

	return nil
}

func NewProcessRawNotificationHandler(
	subscriptionsRepo subscription.Repo,
	withingsSvc withingsSvc.Service,
	accountRepo account.Repo,
) ProcessRawNotificationHandler {
	return processRawNotificationHandler{
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
		accountRepo:      accountRepo,
	}
}

type processRawNotificationHandler struct {
	subscriptionRepo subscription.Repo
	withingsSvc      withingsSvc.Service
	accountRepo      account.Repo
}

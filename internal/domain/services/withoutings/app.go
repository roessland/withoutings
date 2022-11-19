package withoutings

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roessland/withoutings/internal/config"
	"github.com/roessland/withoutings/internal/domain/services/sleep"
	"github.com/roessland/withoutings/web/sessions"
	"github.com/roessland/withoutings/web/templates"
	"github.com/roessland/withoutings/withingsapi"
	"github.com/sirupsen/logrus"
	"time"
)

// Service holds all application resources.
type Service struct {
	Log       logrus.FieldLogger
	Withings  *withingsapi.Client
	Sessions  *sessions.Manager
	Templates templates.Templates
	Sleep     *sleep.Sleep
	DB        *pgxpool.Pool
	Config    *config.Config
	//Async     *asynq.Client
}

// const redisAddr = "127.0.0.1:6379"

// NewService creates a new Withoutings service.
func NewService(ctx context.Context) (*Service, error) {
	svc := &Service{}

	var err error

	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger := logrus.New()
	svc.Log = logger

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return svc, fmt.Errorf("load config: %w", err)
	}
	svc.Config = cfg

	svc.Sessions = sessions.NewManager(cfg.SessionSecret)

	svc.DB, err = pgxpool.New(initCtx, cfg.DatabaseURL)
	if err != nil {
		return svc, fmt.Errorf("create connection pool: %w", err)
	}

	svc.Withings = withingsapi.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	svc.Templates = templates.LoadTemplates()

	svc.Sleep = sleep.NewSleep(svc.Withings)

	//
	//svc.Async = asynq.NewClient(asynq.RedisClientOpt{
	//	Addr: redisAddr,
	//})

	return svc, nil
}

func (app *Service) Close() {
	//err := app.Async.Close()
	//if err != nil {
	//	app.Log.Print(err)
	//}

	app.DB.Close()
}

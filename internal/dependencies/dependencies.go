package dependencies

import (
	"os"
	"os/signal"
	"syscall"

	authAdpt "github.com/warehouse/user-service/internal/adapter/auth"
	mailAdpt "github.com/warehouse/user-service/internal/adapter/mail"
	randomAdpt "github.com/warehouse/user-service/internal/adapter/random"
	timeAdpt "github.com/warehouse/user-service/internal/adapter/time"
	"github.com/warehouse/user-service/internal/broker"
	"github.com/warehouse/user-service/internal/config"
	"github.com/warehouse/user-service/internal/db"
	"github.com/warehouse/user-service/internal/handler/grpc"
	"github.com/warehouse/user-service/internal/handler/http"
	"github.com/warehouse/user-service/internal/handler/middlewares"
	"github.com/warehouse/user-service/internal/pkg/logger"
	transactionsRepo "github.com/warehouse/user-service/internal/repository/operations/transactions"
	userRepo "github.com/warehouse/user-service/internal/repository/operations/user"
	"github.com/warehouse/user-service/internal/server"
	userSvc "github.com/warehouse/user-service/internal/service/user"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Dependencies interface {
		Close()
		Cfg() *config.Config
		Internal() dependencies
		WaitForInterrupr()

		GrpcServer() server.Server
	}

	dependencies struct {
		cfg                     *config.Config
		log                     logger.Logger
		warehouseRequestHandler http.WarehouseRequestHandler
		handlerMiddleware       middlewares.Middleware

		userHandler *grpc.UserHandler

		psqlClient   *db.PostgresClient
		rabbitClient *broker.RabbitClient

		userService userSvc.Service

		pgxTransactionRepo transactionsRepo.Repository
		userRepo           userRepo.Repository

		timeAdapter   timeAdpt.Adapter
		randomAdapter randomAdpt.Adapter
		authAdapter   authAdpt.Adapter
		mailAdapter   mailAdpt.Adapter

		grpcServer server.Server

		shutdownChannel chan os.Signal
		closeCallbacks  []func()
	}
)

func NewDependencies(cfgPath string) (Dependencies, error) {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil && err.Error() == "Config File \"config\" Not Found in \"[]\"" {
		cfg, err = config.NewConfig("./configs/local")
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "1"
	encoderCfg.TimeKey = "t"

	z := zap.New(
		&logger.WarehouseZapCore{
			Core: zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.Lock(os.Stdout),
				zap.NewAtomicLevel(),
			),
		},
		zap.AddCaller(),
	)

	return &dependencies{
		cfg:             cfg,
		log:             logger.NewLogger(z),
		shutdownChannel: make(chan os.Signal),
	}, nil
}

func (d *dependencies) Close() {
	for i := len(d.closeCallbacks) - 1; i >= 0; i-- {
		d.closeCallbacks[i]()
	}
	d.log.Zap().Sync()
}

func (d *dependencies) Internal() dependencies {
	return *d
}

func (d *dependencies) Cfg() *config.Config {
	return d.cfg
}

func (d *dependencies) WarehouseJsonRequestHandler() http.WarehouseRequestHandler {
	if d.warehouseRequestHandler == nil {
		d.warehouseRequestHandler = http.NewWarehouseJsonRequestHandler(d.log, d.cfg.Timeouts.AccCookie)
	}

	return d.warehouseRequestHandler
}

func (d *dependencies) GrpcServer() server.Server {
	if d.grpcServer == nil {
		var err error
		msg := "initialize grpc server"
		if d.grpcServer, err = server.NewGrpcServer(
			d.log,
			*d.cfg,
			d.UserHandler(),
		); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			msg := "shutting down grpc server"
			if err := d.grpcServer.Stop(); err != nil {
				d.log.Zap().Warn(msg, zap.Error(err))
			}
			d.log.Zap().Info(msg)
		})
	}
	return d.grpcServer
}

func (d *dependencies) WaitForInterrupr() {
	signal.Notify(d.shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	d.log.Zap().Info("Wait for receive interrupt signal")
	<-d.shutdownChannel // ждем когда сигнал запишется в канал и сразу убираем его, значит, что сигнал получен
	d.log.Zap().Info("Receive interrupt signal")
}

package main

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	accountInfrastructure "UnpakSiamida/modules/account/infrastructure"

	accountPresentation "UnpakSiamida/modules/account/presentation"

	banksoalInfrastructure "UnpakSiamida/modules/banksoal/infrastructure"

	banksoalPresentation "UnpakSiamida/modules/banksoal/presentation"

	kategoriInfrastructure "UnpakSiamida/modules/kategori/infrastructure"

	kategoriPresentation "UnpakSiamida/modules/kategori/presentation"

	templatepertanyaanInfrastructure "UnpakSiamida/modules/templatepertanyaan/infrastructure"

	templatepertanyaanPresentation "UnpakSiamida/modules/templatepertanyaan/presentation"

	// userInfrastructure "UnpakSiamida/modules/user/infrastructure"

	// userPresentation "UnpakSiamida/modules/user/presentation"

	/////////

	commoninfra "UnpakSiamida/common/infrastructure"

	commonpresentation "UnpakSiamida/common/presentation"

	//////////

	// eventUser "UnpakSiamida/modules/user/event"

	_ "UnpakSiamida/docs"

	"github.com/gofiber/swagger"
	_ "github.com/swaggo/files"
)

var startupErrors []fiber.Map

func mustStart(name string, fn func() error) {
	if err := fn(); err != nil {
		startupErrors = append(startupErrors, fiber.Map{
			"module": name,
			"error":  err.Error(),
		})
	}
}

var (
	dbMain   *gorm.DB
	dbSimak  *gorm.DB
	dbSimpeg *gorm.DB

	onceMain   sync.Once
	onceSimak  sync.Once
	onceSimpeg sync.Once
)

func NewMySQL() (*gorm.DB, error) {
	var err error
	onceMain.Do(func() {
		dsn := "root:@tcp(127.0.0.1:3306)/unpak_simonev?charset=utf8mb4&parseTime=true&loc=Local"

		dbMain, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, _ := dbMain.DB()
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	})
	return dbMain, err
}

func NewMySQLSimak() (*gorm.DB, error) {
	var err error
	onceSimak.Do(func() {
		dsn := "root:@tcp(127.0.0.1:3306)/unpak_simak?charset=utf8mb4&parseTime=true&loc=Local"

		dbSimak, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, _ := dbSimak.DB()
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	})
	return dbSimak, err
}

func NewMySQLSimpeg() (*gorm.DB, error) {
	var err error
	onceSimpeg.Do(func() {
		dsn := "root:@tcp(127.0.0.1:3306)/unpak_simpeg?charset=utf8mb4&parseTime=true&loc=Local"

		dbSimpeg, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, _ := dbSimpeg.DB()
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	})
	return dbSimpeg, err
}

// @title UnpakSiamidaV2 API
// @version 1.0
// @description All Module Siamida
// @host localhost:3000
// @BasePath /
func main() {
	cfg := commonpresentation.DefaultHeaderSecurityConfig()
	cfg.ResolveAndCheck = false

	app := fiber.New(fiber.Config{
		// DisableStartupMessage: true,
		ReadBufferSize: 16 * 1024,
		// Prefork:        true, // gunakan semua CPU cores
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})
	// app.Use(recover())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",     // X-Content-Type-Options
		XFrameOptions:             "DENY",        // X-Frame-Options
		ReferrerPolicy:            "no-referrer", // Referrer-Policy
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'none'",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}))
	app.Use(commonpresentation.LoggerMiddleware)
	app.Use(commonpresentation.HeaderSecurityMiddleware(cfg))
	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Del("X-Powered-By")
		return c.Next()
	})

	mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

	var db *gorm.DB
	var dbSimak *gorm.DB
	var dbSimpeg *gorm.DB
	// var redis commondomain.IRedisStore
	mustStart("Database", func() error {
		var err error
		db, err = NewMySQL()
		return err
	})
	mustStart("Database Simak", func() error {
		var err error
		dbSimak, err = NewMySQLSimak()
		return err
	})
	mustStart("Database Simpeg", func() error {
		var err error
		dbSimpeg, err = NewMySQLSimpeg()
		return err
	})
	// mustStart("Redis", func() error {
	// 	var err error
	// 	redis = NewRedisStore()
	// 	return err
	// })

	// var tg commoninfra.TelegramSender
	// modeTelegram := os.Getenv("TELEGRAM_MODE")

	// mustStart("Telegram Service", func() error {
	// 	factory := &commoninfra.DefaultTelegramFactory{
	// 		UseFake: modeTelegram != "dev",
	// 	}

	// 	client, err := factory.Create()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	tg = client
	// 	return nil
	// })

	//berlaku untuk startup bukan hot reload
	// mustStart("User Module", func() error {
	// 	return userInfrastructure.RegisterModuleUser(db, tg)
	// })

	mustStart("Account Module", func() error {
		return accountInfrastructure.RegisterModuleAccount(db, dbSimak, dbSimpeg)
	})

	mustStart("BankSoal Module", func() error {
		return banksoalInfrastructure.RegisterModuleBankSoal(db)
	})

	mustStart("Kategori Module", func() error {
		return kategoriInfrastructure.RegisterModuleKategori(db)
	})

	mustStart("TaemplatePertanyaan Module", func() error {
		return templatepertanyaanInfrastructure.RegisterModuleTemplatePertanyaan(db)
	})

	if len(startupErrors) > 0 {
		app.Use(func(c *fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{
				"Code":    "INTERNAL_SERVER_ERROR",
				"Message": "Startup module failed",
				"Trace":   startupErrors,
			})
		})
	}

	// dispatcher := commoninfra.NewEventDispatcher()
	// commoninfra.RegisterEvent[eventUser.UserCreatedEvent](dispatcher)
	// commoninfra.RegisterEvent[eventUser.UserUpdatedEvent](dispatcher)

	// userPresentation.ModuleUser(app)
	accountPresentation.ModuleAccount(app)
	banksoalPresentation.ModuleBankSoal(app)
	kategoriPresentation.ModuleKategori(app)
	templatepertanyaanPresentation.ModuleTemplatePertanyaan(app)

	// ctx, stop := signal.NotifyContext(
	// 	context.Background(),
	// 	os.Interrupt,
	// 	syscall.SIGTERM,
	// )
	// defer stop()

	// outboxProcessor := &commoninfra.OutboxProcessor{
	// 	DB:         db,
	// 	Dispatcher: dispatcher,
	// }

	app.Get("/swagger/*", swagger.HandlerDefault)
	// go commoninfra.StartOutboxWorker(ctx, outboxProcessor)
	app.Listen(":3000")
}

type ValidationBehavior struct{}

func NewValidationBehavior() *ValidationBehavior {
	return &ValidationBehavior{}
}

func (b *ValidationBehavior) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	if err := commoninfra.Validate(request); err != nil {
		return nil, err
	}

	return next(ctx)
}

// func NewRedisStore() *commoninfra.RedisStore {
// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 		DB:   0,
// 	})

// 	return commoninfra.NewRedisStore(client)
// }

package dawn

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/dawn/daemon"
	"github.com/go-dawn/dawn/fiberx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Version of current dawn package
const Version = "0.3.8"

// Config is a struct holding the sloop settings.
type Config struct {
	// App indicates to fiber app instance
	App *fiber.App

	// Daemon indicates if dawn app run in daemon mode
	Daemon bool
	// Tries is the max count to start dawn app in daemon
	// mode when error occurs. Default value is 10.
	Tries int
	// StdoutLogFile is the path to write stdout logs
	// in daemon mode. If not set, all stdout logs is
	// directed to os.Stdout
	StdoutLogFile string
	// StderrLogFile is the path to write stderr logs.
	// in daemon mode. If not set, all stdout logs is
	// directed to os.Stderr
	StderrLogFile string
}

// Sloop denotes Dawn application
type Sloop struct {
	// Config is the embedded config
	Config

	app      *fiber.App
	mods     []Moduler
	cleanups []Cleanup
	sigCh    chan os.Signal
}

// New returns a new Sloop with options.
func New(config ...Config) *Sloop {
	s := &Sloop{
		sigCh: make(chan os.Signal, 1),
	}

	if len(config) > 0 {
		s.Config = config[0]
	}

	s.app = s.Config.App

	return s
}

// Default returns an Sloop instance with the
// - RequestID
// - Logger
// - Recovery
// - Pprof
// middleware already attached in default fiber app.
func Default(cfg ...fiber.Config) *Sloop {
	c := fiber.Config{}
	if len(cfg) > 0 {
		c = cfg[0]
	}
	if c.ErrorHandler == nil {
		c.ErrorHandler = fiberx.ErrHandler
	}
	app := fiber.New(c)
	app.Use(
		requestid.New(),
		fiberx.Logger(),
		recover.New(),
	)

	if config.GetBool("debug") {
		app.Use(pprof.New())
	}

	return &Sloop{
		app:   app,
		sigCh: make(chan os.Signal, 1),
	}
}

// AddModulers appends more Modulers
func (s *Sloop) AddModulers(m ...Moduler) *Sloop {
	s.mods = append(s.mods, m...)

	return s
}

// Run runs a web server
func (s *Sloop) Run(addr string) error {
	defer s.Cleanup()
	if s.app == nil {
		return errors.New("dawn: app is nil")
	}

	s.Setup().registerRoutes()

	if config.GetBool("daemon.enable") {
		daemon.Run()
	}

	return s.app.Listen(addr)
}

// RunTls runs a tls web server
func (s *Sloop) RunTls(addr, certFile, keyFile string) error {
	defer s.Cleanup()

	if s.app == nil {
		return errors.New("dawn: app is nil")
	}

	// Create tls certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	s.Setup().registerRoutes()

	if config.GetBool("daemon.enable") {
		daemon.Run()
	}

	return s.app.Listener(ln)
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Sloop) Shutdown() error {
	if s.app == nil {
		return fmt.Errorf("shutdown: fiber app is not found")
	}
	return s.app.Shutdown()
}

// Router returns the server router
func (s *Sloop) Router() fiber.Router {
	return s.app
}

// Setup initializes all modules and then boots them
func (s *Sloop) Setup() *Sloop {
	return s.init().boot()
}

func (s *Sloop) init() *Sloop {
	for _, mod := range s.mods {
		if cleanup := mod.Init(); cleanup != nil {
			s.cleanups = append(s.cleanups, cleanup)
		}
	}
	return s
}

func (s *Sloop) boot() *Sloop {
	for _, mod := range s.mods {
		mod.Boot()
	}

	return s
}

func (s *Sloop) registerRoutes() *Sloop {
	for _, mod := range s.mods {
		mod.RegisterRoutes(s.app)
	}
	return s
}

// Cleanup releases resources
func (s *Sloop) Cleanup() {
	for _, fn := range s.cleanups {
		fn()
	}
}

// Watch listens to signal to exit
func (s *Sloop) Watch() {
	if config.GetBool("daemon.enable") {
		daemon.Run()
	}

	signal.Notify(s.sigCh,
		syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGQUIT)

	<-s.sigCh
}

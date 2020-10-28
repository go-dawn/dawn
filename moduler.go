package dawn

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Cleanup is a function does cleanup works
type Cleanup func()

// Moduler is the interface that wraps the module's method.
type Moduler interface {
	// Stringer indicates module's name
	fmt.Stringer

	// Init does initialization works and should return
	// a cleanup function.
	Init() Cleanup

	// Boot boots the module.
	Boot()

	// RegisterRoutes add routes to fiber router
	RegisterRoutes(fiber.Router)
}

// Module is an empty struct implements Moduler interface
// and can be embedded into custom struct as a Moduler
type Module struct{}

// String indicates module's name
func (Module) String() string { return "anonymous" }

// Init does initialization works and should return a cleanup function.
func (Module) Init() Cleanup { return func() {} }

// Boot boots the module.
func (Module) Boot() {}

// RegisterRoutes add routes to fiber router
func (Module) RegisterRoutes(fiber.Router) {}

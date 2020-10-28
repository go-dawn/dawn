package dawn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockModule struct {
	Module
}

// go test -run Test_Moduler_Embed_Empty_Module -race
func Test_Moduler_Embed_Empty_Module(t *testing.T) {
	t.Parallel()

	module := mockModule{Module{}}

	assert.Implements(t, (*Moduler)(nil), module)

	assert.Equal(t, "anonymous", module.String())

	assert.NotNil(t, module.Init())

	module.Boot()

	module.RegisterRoutes(nil)
}

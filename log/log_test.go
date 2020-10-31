package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	SetOutput(new(bytes.Buffer))
	Errorln("errorln")
	Errorf("%s", "errorf")
	Warningln("warningln")
	Warningf("%s", "warningf")
	Infoln(0, "infoln level 0")
	Infof(0, "%s", "infof level 0")
	Infof(1, "%s", "infof level 1")
}

func Test_Moduler(t *testing.T) {
	m := New(nil)

	assert.Equal(t, "dawn:log", m.String())

	m.Init()()
}

package log

import (
	"bytes"
	"testing"
)

func Test_All(t *testing.T) {
	InitFlags(nil)
	SetOutput(new(bytes.Buffer))
	Errorln("errorln")
	Errorf("%s", "errorf")
	Warningln("warningln")
	Warningf("%s", "warningf")
	Infoln(0, "infoln level 0")
	Infof(0, "%s", "infof level 0")
	Infof(1, "%s", "infof level 1")
	Flush()
}

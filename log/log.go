package log

import (
	"flag"
	"io"

	"github.com/kiyonlin/klog"
)

// InitFlags is for explicitly initializing the flags.
// Default to use logs as log dir.
func InitFlags(flagset *flag.FlagSet) {
	klog.InitFlags(flagset)
}

// SetOutput sets the output destination for all severities
func SetOutput(w io.Writer) {
	klog.SetOutput(w)
}

// Flush flushes all pending log I/O.
func Flush() {
	klog.Flush()
}

// Warningln logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Warningln(args ...interface{}) {
	klog.WarningDepth(1, args...)
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	klog.WarningDepthf(1, format, args...)
}

// Errorln logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Errorln(args ...interface{}) {
	klog.ErrorDepth(1, args...)
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	klog.ErrorDepthf(1, format, args...)
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
func Infoln(level int, args ...interface{}) {
	l := klog.Level(level)
	if klog.V(l).Enabled() {
		klog.V(l).InfoDepth(1, args...)
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
func Infof(level int, format string, args ...interface{}) {
	l := klog.Level(level)
	if klog.V(l).Enabled() {
		klog.V(l).InfoDepthf(1, format, args...)
	}
}

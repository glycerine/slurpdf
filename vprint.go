package slurpdf

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	//"runtime/pprof"
	"sync"
	"time"
	//"github.com/glycerine/cryrand"
	//"4d63.com/tz"
)

var gtz *time.Location
var utc *time.Location
var chicago *time.Location

func init() {

	// do this is ~/.bashrc so we get the default.
	os.Setenv("TZ", "America/Chicago")

	var err error
	//Chicago, err = tz.LoadLocation("America/Chicago")
	utc, err = time.LoadLocation("UTC")
	panicOn(err)
	chicago, err = time.LoadLocation("America/Chicago")
	panicOn(err)
	//gtz = utc
	gtz = chicago
}

const RFC3339MsecTz0 = "2006-01-02T15:04:05.000Z07:00"

// for tons of debug output
var verboseVerbose bool = false

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func pp(format string, a ...interface{}) {
	if verboseVerbose {
		tsPrintf(format, a...)
	}
}

func vv(format string, a ...interface{}) {
	tsPrintf(format, a...)
}

func alwaysPrintf(format string, a ...interface{}) {
	tsPrintf(format, a...)
}

var tsPrintfMut sync.Mutex

// time-stamped printf
func tsPrintf(format string, a ...interface{}) {
	tsPrintfMut.Lock()
	printf("\n%s %s ", fileLine(3), ts())
	printf(format+"\n", a...)
	tsPrintfMut.Unlock()
}

// get timestamp for logging purposes
func ts() string {
	return time.Now().Format(RFC3339MsecTz0)
}

// so we can multi write easily, use our own printf
var ourStdout io.Writer = os.Stdout

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(ourStdout, format, a...)
}

func fileLine(depth int) string {
	_, fileName, fileLine, ok := runtime.Caller(depth)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", path.Base(fileName), fileLine)
	} else {
		s = ""
	}
	return s
}

func caller(upStack int) string {
	// elide ourself and runtime.Callers
	target := upStack + 2

	pc := make([]uintptr, target+2)
	n := runtime.Callers(0, pc)

	f := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(pc[:n])
		for i := 0; i <= target; i++ {
			contender, more := frames.Next()
			if i == target {
				f = contender
			}
			if !more {
				break
			}
		}
	}
	return f.Function
}

/*
func startProfilingCPU(path string) {
	// add randomness so two tests run at once don't overwrite each other.
	rnd8 := cryrand.RandomStringWithUp(8)
	fn := path + ".cpuprof." + rnd8
	f, err := os.Create(fn)
	panicOn(err)
	AlwaysPrintf("will write cpu profile to '%v'", fn)

	pprof.StartCPUProfile(f)
}
*/

func stack() string {
	return string(debug.Stack())
}

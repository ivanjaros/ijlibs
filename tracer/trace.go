package misc

import (
	"runtime"
	"strings"
)

func BackTrace(wrap int) (file string, fc string, line int) {
	function, file, line, _ := runtime.Caller(1 + wrap)

	i := strings.LastIndex(file, "/")
	if i != -1 {
		file = file[i+1:]
	}

	return file, runtime.FuncForPC(function).Name(), line
}

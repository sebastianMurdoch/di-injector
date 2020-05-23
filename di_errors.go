package di_injector

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrAddDependency = "cannot add the dependency %v because %v"
	ErrInjectInner = "coudnt inject inner dependency -- %v"
	ErrFoundImplementation = "None implementation was found for dependency: %v"
	ErrAtInjection       = "fatal error at injection"
	ErrInjectInterface	= "cannot inject into interface{}"
)

func DiError(errFormat string, values ...interface{}) error {
	var message string
	if values == nil {
		message =errFormat
	}else {
		message = fmt.Sprintf(errFormat,values)
	}
	errLog.Print(message)
	return errors.New(message)
}

func DiFatal(errFormat string, values ...interface{}) {
	errLog.Fatal(fmt.Sprintf(errFormat,values))
}

var errLog = Logger(log.New(os.Stderr, "[di-injector] ", log.Ldate|log.Ltime|log.Lshortfile))

// Logger is used to log critical error messages.
type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}

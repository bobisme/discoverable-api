package errors

import (
	"fmt"
)

type ErrClass int

const (
	Internal ErrClass = iota
	Authentication
	Permission
	Client
	Validation
)

type Error interface {
	Error() string
	Class() ErrClass
	Cause() error
}

type Err struct {
	msg   string
	cause error
	class ErrClass
}

func New(msg string) *Err { return &Err{msg, nil, Internal} }
func NewAuthentication(msg string) *Err {
	return &Err{msg, nil, Authentication}
}
func NewPermission(msg string) *Err { return &Err{msg, nil, Permission} }
func NewClient(msg string) *Err     { return &Err{msg, nil, Client} }

func Wrap(message string, e error) *Err {
	msg := fmt.Sprintf("%s: %s", message, e.Error())
	class := Internal
	if err, ok := e.(Error); ok {
		class = err.Class()
	}
	return &Err{msg, e, class}
}

func (e *Err) Error() string {
	return e.msg
}

func (e *Err) Class() ErrClass {
	return e.class
}

type ValidationErr struct {
	Err
	path string
}

func NewValidation(msg, path string) *ValidationErr {
	return &ValidationErr{Err{msg, nil, Validation}, path}
}

// Path returns the JSON path to the field in question.
func (e *ValidationErr) Path() string {
	return e.path
}

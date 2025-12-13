package application

import "fmt"

type ErrUser struct {
	Message string
}

func (e ErrUser) Error() string {
	return fmt.Sprintf("incorrect use of the program: %s", e.Message)
}

func NewErrUser(message string) error {
	return ErrUser{Message: message}
}

type ErrGenerating struct {
	Message string
}

func (e ErrGenerating) Error() string {
	return fmt.Sprintf("error while generating fractal: %s", e.Message)
}

func NewErrGenerating(message string) error {
	return ErrGenerating{Message: message}
}

type ErrFile struct {
	Message string
}

func (e ErrFile) Error() string {
	return fmt.Sprintf("error while saving file: %s", e.Message)
}

func NewErrFile(message string) error {
	return ErrFile{Message: message}
}

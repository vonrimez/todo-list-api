package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

type Validate struct {
	validate *validator.Validate
}

func NewValidate() *Validate {
	return &Validate{validate: validator.New()}
}

func (v *Validate) validateStruct(obj any) error {
	return v.validate.Struct(obj)
}

// return filled map if error consists of fields parsed in ShouldBindJSON
// or nil otherwise
func (v *Validate) ValError(err error) map[string]any {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		fieldErrors := make(map[string]any)
		for _, fe := range valErr {
			fieldErrors[fe.Field()] = fmt.Sprintf("failed on '%s' rule", fe.Tag())
		}
		return fieldErrors
	}
	return nil
}

// ShouldBindJSON return internal json decoder error or structed error
// which can be parsed into a map[`string`]any using the ValError
func (v *Validate) ShouldBindJSON(reader io.Reader, obj any) error {
	err := json.NewDecoder(reader).Decode(obj)
	if err != nil {
		return err
	}
	return v.validateStruct(obj)
}

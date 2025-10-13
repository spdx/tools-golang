package ld

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// JoinErrors returns errors.Join'd errors, taking into account nested joined errors, flattening these to a single joined set
func JoinErrors(errs ...error) error {
	var out []error
	for _, err := range errs {
		out = append(out, flattenErrors(err)...)
	}
	switch len(out) {
	case 0:
		return nil
	case 1:
		return out[0]
	default:
		return errors.Join(out...)
	}
}

func flattenErrors(err error) []error {
	var out []error
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range joined.Unwrap() {
			out = append(out, flattenErrors(e)...)
		}
	} else {
		if err != nil {
			return []error{err}
		}
	}
	return out
}

var validatorInterface = reflect.TypeOf((*Validator)(nil)).Elem()

type validationError struct {
	Path []any
	Err  error
}

func (v *validationError) String() string {
	path := ""
	for i := 0; i < len(v.Path); i++ {
		part := v.Path[i]
		switch p := part.(type) {
		case int:
			path += "[" + strconv.Itoa(p) + "]"
		case reflect.StructField:
			if !p.Anonymous {
				path += "." + p.Name
			}
		case reflect.Type:
			path += "<" + p.Name() + ">"
		default:
			path += "/" + fmt.Sprint(p)
		}
	}
	return path + ": " + v.Err.Error()
}

func (v *validationError) Error() string {
	return v.String()
}

func newValidationError(err error, path ...any) *validationError {
	// if the error is a validation error, prepend the path
	if vErr, ok := err.(*validationError); ok {
		vErr.Path = append(path, vErr.Path...)
		return vErr
	}
	return &validationError{
		Path: path,
		Err:  err,
	}
}

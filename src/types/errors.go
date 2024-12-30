// Package types provides the type definitions for the project
package types

var _ error = (Error)(nil)

type Error interface {
	error
	Name() string
	Code() string
	Message() string
}

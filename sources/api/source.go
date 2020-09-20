package api

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

/*
Source implements the source.Source interface for the "api" source.
*/
type Source struct {
	options *source.Options
}

/*
New returns a valid Blacksmith source.

There is no need for a schedule on this source since it only handles HTTP triggers.
*/
func New() source.Source {
	return &Source{
		options: &source.Options{},
	}
}

/*
String returns the string representation of the source.
*/
func (s *Source) String() string {
	return "api"
}

/*
Options returns common source options. They will be shared across every triggers
of this source, except when overridden.
*/
func (s *Source) Options() *source.Options {
	return s.options
}

/*
Triggers return a list of triggers the source is able to handle.
*/
func (s *Source) Triggers() map[string]source.Trigger {
	return map[string]source.Trigger{
		"register": TriggerRegister{},
	}
}

package postgres

import (
	"github.com/nunchistudio/blacksmith/flow/destination"
)

/*
Destination implements the destination.Destination interface for the "postgres"
destination.
*/
type Destination struct {
	options *destination.Options
}

/*
New returns a valid Blacksmith destination.

For the purpose of the demo, we will use PostgreSQL as a data warehouse with a
little interval. We do not want events to be loaded in realtime. In case of
failure, we specify to retry every 2 minutes with a limit of 20 retries. After
that, if the jobs still fail they will be marked as "discarded".
*/
func New() destination.Destination {
	return &Destination{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   false,
				Interval:   "@every 2h",
				MaxRetries: 20,
			},
		},
	}
}

/*
String returns the string representation of the destination.
*/
func (postgres *Destination) String() string {
	return "postgres"
}

/*
Options returns common destination options. They will be shared across every actions
of this destination, except when overridden.
*/
func (postgres *Destination) Options() *destination.Options {
	return postgres.options
}

/*
Actions return a list of actions the destination is able to handle.
*/
func (postgres *Destination) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"register": ActionRegister{},
	}
}

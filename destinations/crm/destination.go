package crm

import (
	"github.com/nunchistudio/blacksmith/flow/destination"
)

/*
Destination implements the destination.Destination interface for the "crm"
destination.
*/
type Destination struct {
	options *destination.Options
}

/*
New returns a valid Blacksmith destination.

For the purpose of the demo, we want events to be loaded into the CRM in realtime.
In case of failure, we specify to retry every 20 seconds with a limit of 50 retries.
After that, if the job still failed it will be marked as "discarded".
*/
func New() destination.Destination {
	return &Destination{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   true,
				Interval:   "@every 20s",
				MaxRetries: 50,
			},
		},
	}
}

/*
String returns the string representation of the destination.
*/
func (crm *Destination) String() string {
	return "crm"
}

/*
Options returns common destination options. They will be shared across every actions
of this destination, except when overridden.
*/
func (crm *Destination) Options() *destination.Options {
	return crm.options
}

/*
Actions return a list of actions the destination is able to handle.
*/
func (crm *Destination) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"register": ActionRegister{},
		"notify":   ActionNotify{},
	}
}

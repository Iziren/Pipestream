package postgres

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

/*
Source implements the source.Source interface for the "postgres" source.
*/
type Source struct {
	options *source.Options
}

/*
New returns a valid Blacksmith source.

For the purpose of the demo, we set the interval to every 30 minutes. Which means,
if not overridden, every trigger in CRON mode of this source will run every 30
minutes.
*/
func New() source.Source {
	return &Source{
		options: &source.Options{
			DefaultSchedule: &source.Schedule{
				Interval: "@every 30m",
			},
		},
	}
}

/*
String returns the string representation of the source.
*/
func (postgres *Source) String() string {
	return "postgres"
}

/*
Options returns common source options. They will be shared across every triggers
of this source, except when overridden.
*/
func (postgres *Source) Options() *source.Options {
	return postgres.options
}

/*
Triggers return a list of triggers the source is able to handle.
*/
func (postgres *Source) Triggers() map[string]source.Trigger {
	return map[string]source.Trigger{
		"dummy-cron":    TriggerDummyCRON{},
		"dummy-forever": TriggerDummyCDC{},
	}
}

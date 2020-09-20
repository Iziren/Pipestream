package postgres

import (
	"time"

	"github.com/nunchistudio/blacksmith/flow/source"

	"github.com/nunchistudio/smithy/sources"
)

/*
TriggerDummyCRON is the payload structure sent by an event and that will be received
by the gateway Blacksmith needs "Context", "Data", and "SentAt" keys to ensure
consistency across triggers.
*/
type TriggerDummyCRON struct {

	// Context is a shared context across events. It is used to save common properties
	// about events such as timezone, location, language, IP address, etc.
	Context *sources.Context `json:"context"`

	// Data is the data specific to this trigger.
	Data *DummyCRON `json:"data"`

	// SentAt is the registered timestamp the event was sent at.
	SentAt *time.Time `json:"sent_at"`
}

/*
DummyCRON is the data payload specific to this trigger.
*/
type DummyCRON struct{}

/*
String returns the string representation of the trigger.
*/
func (t TriggerDummyCRON) String() string {
	return "dummy-cron"
}

/*
Mode allows to register the trigger as a CRON task. Since we set a schedule,
this will override the destination's schedule. Which means this CRON task will
run every 10 seconds and not every 30 minutes.

Everytime the schedule is met, the Extract function will be triggered.
*/
func (t TriggerDummyCRON) Mode() *source.Mode {
	return &source.Mode{
		Mode: source.ModeCRON,
		UsingCRON: &source.Schedule{
			Interval: "@every 1m",
		},
	}
}

/*
Extract is the event listening function. This is this function that is triggered
everytime the schedule is met.

Since no payload, no errors, and no flows are returned, the event will be created
in the store with no context, no data, and no associated jobs.
*/
func (t TriggerDummyCRON) Extract(tk *source.Toolkit) (*source.Payload, error) {
	tk.Logger.Info("postgres/dummy-cron: Running a CRON task...")

	return nil, nil
}

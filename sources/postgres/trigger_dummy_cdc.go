package postgres

import (
	"time"

	"github.com/nunchistudio/blacksmith/flow/source"

	"github.com/nunchistudio/smithy/sources"
)

/*
TriggerDummyCDC is a dummy source's trigger to demonstrate how Blacksmith works
with an infinite loop. These ongoing loops are useful to listen for CDC notifications
such as database ones.
*/
type TriggerDummyCDC struct {

	// Context is a shared context across events. It is used to save common properties
	// about events such as timezone, location, language, IP address, etc.
	Context *sources.Context `json:"context"`

	// Data is the data specific to this trigger.
	Data *DummyCDC `json:"data"`

	// SentAt is the registered timestamp the event was sent at.
	SentAt *time.Time `json:"sent_at"`
}

/*
DummyCDC is the data payload specific to this trigger.
*/
type DummyCDC struct{}

/*
String returns the string representation of the trigger.
*/
func (t TriggerDummyCDC) String() string {
	return "dummy-cdc"
}

/*
Mode allows to register the trigger as an ongoing notification listener. No
additional details are needed for this mode.
*/
func (t TriggerDummyCDC) Mode() *source.Mode {
	return &source.Mode{
		Mode: source.ModeCDC,
	}
}

/*
Extract is function being run by the gateway. It is up to the function body to
include the forever loop. When set to this mode, the toolkit gives you access
to channels to either return the payload or an error whenever needed.

Also, since this mode is asynchronous, there is no way for the gateway to know
when the trigger is done. To gracefully shutdown, the function receives a message
on "IsShuttingDown" and must write to "IsDone" whenever the function is ready to
exit. Otherwise, the gateway will block until "true" is received on "IsDone".
*/
func (t TriggerDummyCDC) Extract(tk *source.Toolkit, notifier *source.Notifier) {

	for {
		select {
		// case <-notification:
		// 	notifier.Payload <- &source.Payload{}
		// 	notifier.Error <- &errors.Error{}
		case <-notifier.IsShuttingDown:
			notifier.Done <- true
		}
	}
}

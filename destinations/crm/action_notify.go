package crm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/smithy/sources"
)

/*
ActionNotify is the payload structure received by this action and that will be
sent to the destination by the scheduler. Blacksmith needs "Context", "Data",
and "SentAt" keys to ensure consistency across actions.
*/
type ActionNotify struct {

	// Context is a shared context across events. It is used to save common properties
	// about events such as timezone, location, language, IP address, etc.
	Context *sources.Context `json:"context"`

	// Data is the data specific to this action.
	Data *Notification `json:"data"`

	// SentAt is the registered timestamp the event was sent at.
	SentAt *time.Time `json:"sent_at"`
}

/*
Notification is the data payload specific to this action.
*/
type Notification struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

/*
String returns the string representation of the action.
*/
func (a ActionNotify) String() string {
	return "notify"
}

/*
Schedule allows the action to override the schedule options of its destination.

For the purpose of the demo we override the destination's schedule because we
do not want to be notified in realtime when a new user is registered in the CRM.
We will only receive notification every 20 seconds.
*/
func (a ActionNotify) Schedule() *destination.Schedule {
	return &destination.Schedule{
		Realtime: false,
	}
}

/*
Marshal is the function being run when the action receive data in the ActionNotify
receiver. Like for a source's trigger, it is also in charge of the "T" in the ETL
process: it can Transform (if needed) the payload to the given data structure.
*/
func (a ActionNotify) Marshal(tk *destination.Toolkit) (*destination.Payload, error) {

	// Try to marshal the data passed directly to the receiver.
	data, err := json.Marshal(&a.Data)
	if err != nil {
		return nil, err
	}

	// Create a payload with the data. Since the "Context" key is not set, the one
	// from the event will automatically be applied by the scheduler.
	p := &destination.Payload{
		Data:   data,
		SentAt: a.SentAt,
	}

	// Return the payload with the marshaled data.
	return p, nil
}

/*
Load is the function being run when the action is run by the scheduler. It is in
charge of the "L" in the ETL process: it Loads the data to the destination.

It received a queue of events containing jobs related to this action only.
In this case, since there is no error returned, the load will be a success.
*/
func (a ActionNotify) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their related
	// jobs. The jobs present in the events are specific to this action only.
	var notifications = []Notification{}
	for _, event := range queue.Events {
		for _, job := range event.Jobs {
			var n Notification
			json.Unmarshal(job.Data, &n)

			notifications = append(notifications, n)
		}
	}

	// Simply print a message.
	fmt.Println("Batch of notifications:", notifications)

	// Whenever we are ready, we inform the scheduler about the jobs status.
	// Since we do not do anything but to print a message, we inform the scheduler
	// only once, with no error and no job IDs. When no job IDs are provided, the
	// scheduler will mark every jobs from the queue as "succeeded", "failed", or
	// "discarded".
	//
	// In this case since Error is nil every jobs will be marked as "succeeded".
	// Otherwise, the scheduler will mark each job as "failed" or "discarded"
	// given the current attempt number of the job and the max retries allowed by
	// the action.
	then <- destination.Then{
		Error: nil,
	}
}

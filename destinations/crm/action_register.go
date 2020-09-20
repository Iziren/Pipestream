package crm

import (
	"encoding/json"
	"time"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"github.com/nunchistudio/smithy/sources"
)

/*
ActionRegister is the payload structure received by this action and that will be
sent to the destination by the scheduler. Blacksmith needs "Context", "Data",
and "SentAt" keys to ensure consistency across actions.
*/
type ActionRegister struct {

	// Context is a shared context across events. It is used to save common properties
	// about events such as timezone, location, language, IP address, etc.
	Context *sources.Context `json:"context"`

	// Data is the data specific to this action.
	Data *User `json:"data"`

	// SentAt is the registered timestamp the event was sent at.
	SentAt *time.Time `json:"sent_at"`
}

/*
User is the data payload specific to this action.
*/
type User struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

/*
String returns the string representation of the action.
*/
func (a ActionRegister) String() string {
	return "register"
}

/*
Schedule allows the action to override the schedule options of its destination.

Here we do not override the destination's schedule. The action will run in
realtime (if Pub / Sub is enabled), with retries every 20 seconds and maximum 50
retries.
*/
func (a ActionRegister) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receive data in the ActionRegister
receiver. Like for a source's trigger, it is also in charge of the "T" in the ETL
process: it can Transform (if needed) the payload to the given data structure.
*/
func (a ActionRegister) Marshal(tk *destination.Toolkit) (*destination.Payload, error) {

	// Try to marshal the data passed directly to the receiver.
	data, err := json.Marshal(&a.Data)
	if err != nil {
		return nil, err
	}

	// Create a payload with the data. Since the 'Context' key is not set, the one
	// from the event will automatically be applied.
	p := &destination.Payload{
		Data:   data,
		SentAt: a.SentAt,
	}

	// Return the payload with the marshaled data.
	return p, nil
}

/*
Load is the function being run by the scheduler to load the data into the destination.
It is in charge of the "L" in the ETL process.

It received a queue of events containing jobs related to this action only.
In this case, since there is no error returned, the load will be a success.
*/
func (a ActionRegister) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their related
	// jobs. The jobs present in the events are specific to this action only.
	for _, event := range queue.Events {
		for _, job := range event.Jobs {
			var u User
			json.Unmarshal(job.Data, &u)

			// Whenever we are ready, we inform the scheduler about the job status.
			// Here we inform the scheduler for each job individually.
			//
			// In this case Error is not nil and we force the job to be discarded,
			// because it is pointless to retry the job another time. In other scenarios,
			// the scheduler will mark the job as "succeeded" if Error is nil, or as
			// "failed" or "discarded" given the current attempt number of the job and
			// the max retries allowed by the action.
			then <- destination.Then{
				Jobs: []string{job.ID},
				Error: &errors.Error{
					StatusCode: 401,
					Message:    "Not authorized",
					Validations: []errors.Validation{
						{
							Message: "Email address not authorized",
							Path:    []string{"request", "payload", "data", "email"},
						},
					},
				},
				ForceDiscard: true,
				OnFailed:     []destination.Action{},
				OnDiscarded: []destination.Action{
					ActionNotify{
						Data: &Notification{
							Email:   u.Email,
							Message: "Failed to register",
						},
					},
				},
				OnSucceeded: []destination.Action{
					ActionNotify{
						Data: &Notification{
							Email:   u.Email,
							Message: "Successfully registered",
						},
					},
				},
			}
		}
	}
}

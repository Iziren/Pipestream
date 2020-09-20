package postgres

import (
	"encoding/json"
	"time"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"

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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

/*
String returns the string representation of the action.
*/
func (a ActionRegister) String() string {
	return "register"
}

/*
Schedule allows the action to override the schedule options of its destination.

For the purpose of the demo, we override the original schedule of the destination
for a specific one.
*/
func (a ActionRegister) Schedule() *destination.Schedule {
	return &destination.Schedule{
		Realtime:   true,
		Interval:   "@every 1m",
		MaxRetries: 3,
	}
}

/*
Marshal is the function being run when the action receive data in the ActionRegister
receiver. Like for a source's trigger, it is also in charge of the "T" in the ETL
process: it can Transform (if needed) the payload to the given data structure.
*/
func (a ActionRegister) Marshal(tk *destination.Toolkit) (*destination.Payload, error) {

	// Try to marshal the action data passed directly to the struct.
	buff, err := json.Marshal(&a.Data)
	if err != nil {
		return nil, err
	}

	// Create a payload with the data. Since the "Context" key is not set, the one
	// from the event will automatically be applied.
	p := &destination.Payload{
		Data:   buff,
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

	// Do something...

	// Whenever we are ready, we inform the scheduler about the jobs status.
	// Since we do not do anything, we inform the scheduler only once, with no
	// error and no job IDs. When no job IDs are provided, the scheduler will mark
	// every jobs from the queue as "succeeded", "failed", or "discarded".
	//
	// In this case since Error is nil every jobs will be marked as "succeeded".
	// Otherwise, the scheduler will mark each job as "failed" or "discarded"
	// given the current attempt number of the job and the max retries allowed by
	// the action.
	then <- destination.Then{
		Error: nil,
	}
}

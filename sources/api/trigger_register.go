package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/source"

	"github.com/nunchistudio/smithy/flows"
	"github.com/nunchistudio/smithy/sources"
)

/*
TriggerRegister is the payload structure sent by an event and that will be received
by the gateway Blacksmith needs "Context", "Data", and "SentAt" keys to ensure
consistency across triggers.
*/
type TriggerRegister struct {

	// Context is a shared context across events. It is used to save common properties
	// about events such as timezone, location, language, IP address, etc.
	Context *sources.Context `json:"context"`

	// Data is the data specific to this trigger.
	Data *User `json:"data"`

	// SentAt is the registered timestamp the event was sent at.
	SentAt *time.Time `json:"sent_at"`
}

/*
User is the data payload specific to this trigger.
*/
type User struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

/*
String returns the string representation of the trigger.
*/
func (t TriggerRegister) String() string {
	return "register"
}

/*
Mode allows to register the trigger as a HTTP route. This means, every time a
"POST" request is executed against the route "/register" the Extract function
will run.
*/
func (t TriggerRegister) Mode() *source.Mode {
	return &source.Mode{
		Mode: source.ModeHTTP,
		UsingHTTP: &source.Route{
			Methods:  []string{"POST"},
			Path:     "/register",
			ShowMeta: true,
			ShowData: false,
		},
	}
}

/*
Extract is the function being run when the HTTP route is triggered. It is in
charge of the "E" in the ETL process: Extract the data from the source given
the event.

The function allows to return data to flows. It is the "T" in the ETL process:
it transforms the payload from the source's trigger to given destinations' actions.
*/
func (t TriggerRegister) Extract(tk *source.Toolkit, req *http.Request) (*source.Payload, error) {

	// Create an empty payload, catch unwanted fields, and unmarshal it.
	// Return an error if any occured.
	var payload TriggerRegister
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	// Try to marshal the context from the request payload.
	ctx, err := json.Marshal(&payload.Context)
	if err != nil {
		return nil, err
	}

	// Try to marshal the data from the request payload.
	data, err := json.Marshal(&payload.Data)
	if err != nil {
		return nil, err
	}

	// Return the context, data, and a collection of flows to run.
	return &source.Payload{
		Context: ctx,
		Data:    data,
		SentAt:  payload.SentAt,
		Flows: []flow.Flow{
			&flows.OnRegister{
				Username:  payload.Data.Username,
				FullName:  payload.Data.FirstName + " " + strings.ToUpper(payload.Data.LastName),
				FirstName: payload.Data.FirstName,
				LastName:  strings.ToUpper(payload.Data.LastName),
				Email:     strings.ToLower(payload.Data.Email),
			},
		},
	}, nil
}

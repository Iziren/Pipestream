package flows

import (
	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/smithy/destinations/crm"
	"github.com/nunchistudio/smithy/destinations/postgres"
)

/*
OnRegister implements the flow.Flow interface.
*/
type OnRegister struct {
	options *flow.Options

	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

/*
Options returns the fow options. This flow is enabled but can be disabled
whenever you want.
*/
func (f *OnRegister) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by the scheduler when receiving the flow from the
actions. It is up to the flow to receive the data from sources and match it
against the desired actions.
*/
func (f *OnRegister) Transform(tk *flow.Toolkit) destination.Actions {
	return map[string][]destination.Action{
		"crm": []destination.Action{
			&crm.ActionRegister{
				Data: &crm.User{
					FullName: f.FullName,
					Email:    f.Email,
				},
			},
		},
		"postgres": []destination.Action{
			&postgres.ActionRegister{
				Data: &postgres.User{
					FirstName: f.FirstName,
					LastName:  f.LastName,
					Email:     f.Email,
				},
			},
		},
	}
}

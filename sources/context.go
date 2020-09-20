package sources

import (
	"net"
)

/*
Context is the context shared across the whole application. Every events will share
this context for consistency across triggers and actions.
*/
type Context struct {
	IP       net.IP   `json:"ip,omitempty"`
	Locale   string   `json:"locale,omitempty"`
	Timezone string   `json:"timezone,omitempty"`
	Library  *Library `json:"library,omitempty"`
}

/*
Library contains data about the library making the request.
*/
type Library struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

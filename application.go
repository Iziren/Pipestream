package main

import (
	"github.com/nunchistudio/blacksmith"
	"github.com/nunchistudio/blacksmith/adapter/pubsub"
	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/flow/source"
	"github.com/nunchistudio/blacksmith/service"

	"github.com/nunchistudio/smithy/sources/api"
	spg "github.com/nunchistudio/smithy/sources/postgres"

	"github.com/nunchistudio/smithy/destinations/crm"
	dpg "github.com/nunchistudio/smithy/destinations/postgres"
)

/*
Init is used by the Blacksmith CLI to build the application as a Go plugin. This
is the entrypoint for Blacksmith to validate the application and version migrations
from the events and adapters.
*/
func Init() *blacksmith.Options {

	var options = &blacksmith.Options{

		Gateway: &service.Options{
			// KeyFile:  "server.key",
			// CertFile: "server.crt",
		},
		Scheduler: &service.Options{
			// KeyFile:  "server.key",
			// CertFile: "server.crt",
		},

		Store: &store.Options{
			From: "postgres",
		},
		PubSub: &pubsub.Options{
			From: "nats",
		},

		// Supervisor: &supervisor.Options{
		// 	From:    "consul",
		// 	Join: &supervisor.Node{
		// 		Address: "localhost:8500",
		// 	},
		// },

		// Wanderer: &wanderer.Options{
		// 	From:    "postgres",
		// },

		Sources: []*source.Options{
			{
				Load: api.New(),
			},
			{
				Load: spg.New(),
			},
		},

		Destinations: []*destination.Options{
			{
				Load: crm.New(),
			},
			{
				Load: dpg.New(),
			},
		},
	}

	return options
}

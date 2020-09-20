# Smithy

This is the sample application we build from to demonstrate the features and benefits
of Blacksmith, the data engineering platform we use at [Nunchi](https://nunchi.studio/).

## Overview

### Sources and triggers

| Sources    | Triggers     | Mode | Details                           | Flows to execute when triggered |
|------------|--------------|------|-----------------------------------|---------------------------------|
| `api`      | `register`   | HTTP | Method: `POST`, Path: `/register` | `OnRegister`                    |
| `postgres` | `dummy_cron` | CRON | Interval: `@every 1m`             |                                 |
| `postgres` | `dummy_cdc`  | CDC  |                                   |                                 |

### Flows

| Flows        | Actions to run                      |
|--------------|-------------------------------------|
| `OnRegister` | `crm.register`, `postgres.register` |

### Destinations and actions

| Destinations | Actions    | Realtime | On success       | On failure | On discard       |
|--------------|------------|----------|------------------|------------|------------------|
| `crm`        | `register` | Yes      | New job `notify` |            | New job `notify` |
| `crm`        | `notify`   | No       |                  |            |                  |
| `postgres`   | `register` | Yes      |                  |            |                  |

## Usage

### Start the server

First, the Blacksmith application must be up and running:
```go
$ docker-compose up
```

### Your first event

If you want to test the project without any modification, you can make a `POST`
request at `http://localhost:8080/register` with the following JSON payload:
```js
{
	"context": {
		"locale": "fr-FR",
		"ip": "127.0.0.1"
	},
	"data": {
		"first_name": "john",
		"last_name": "doe"
	}
}
```

The HTTP response will contains meta data about the event, the marshaled payload
sent, and the jobs that will be executed.

Also, you can take a look at the data persisted in your PostgreSQL store. You
will see the events, their jobs, and related transitions.

## Links

- [Blacksmith repository on GitHub](https://github.com/nunchistudio/blacksmith)
- [Blacksmith guides on Nunchi website](https://nunchi.studio/blacksmith)
- [API reference on Go developer portal](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc)

## Professional services

Along consulting and training, we provide different product offerings as well as
different levels of support.

- [Discover our services](https://nunchi.studio/support)

## License

Repository licensed under the [MIT License](./LICENSE).

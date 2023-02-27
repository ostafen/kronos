# Kronos

Kronos allows you to periodically invoke your webhooks using cron expressions.

## Features:

- :zap: Easy to use REST API to schedule your webhooks;
- :alarm_clock: Complex scheduling using cron expressions;
- :mailbox_with_mail: **Prometheus** integration for getting failures notifications.

## Build (Go 1.9+)

Run the following command

```bash
foo@bar$ make build
```
to build an executable which will be output to the `bin` subfolder.

## Yaml file configuration

```yaml
logging:
  level: INFO
  format: JSON

port: 9175

store:
  path: "/path/to/db/file" # default is kronos.bolt
```

## Docker compose configuration

```yaml
services:
  kronos:
    image: ghcr.io/ostafen/kronos
    ports:
      - '9175:9175'
    environment:
      - PORT=9175 # configuration properties can be overridden through environment variables
      - STORE_PATH=/data/kronos.bolt
    volumes:
      - ./data:/data
```

## Registering a periodic schedule
To start getting some webhook notifications, let's add a new schedule which will be notified every minute:
```bash
curl -X POST localhost:9175/schedules -H 'Content-Type: application/json' -d \
'{
    "title": "sample-schedule",
    "description": "a sample schedule description",
    "cronExpr": "0/1 * * * *",
    "url": "your-webhook-address",
    "isRecurring": true,
    "startAt": "2023-02-19T11:34:00Z",
    "endAt": "2023-02-19T11:38:00Z"
}'
```
On success, the response of the server will be similar to the following:
```json
{
    "id": "1e6d146b-e3b7-4e5c-b7ce-b7b2860f461b",
    "title": "sample-schedule",
    "status": "not_started",
    "description": "a sample schedule description",
    "cronExpr": "0/1 * * * *",
    "url": "your-webhook-address",
    "metadata": null,
    "isRecurring": true,
    "createdAt": "2023-02-19T12:32:30.788562107+01:00",
    "runAt": "0001-01-01T00:00:00Z",
    "startAt": "2023-02-19T11:34:00Z",
    "endAt": "2023-02-19T11:38:00Z",
    "nextScheduleAt": "2023-02-19T11:34:00Z",
}
```

The above table contains the full list of supported fields:

| Parameter   |      Required      | Description |
|-------------|:------------------:|:------------|
| title |  true | the name of your schedule. It must be unique. |
| description |  false   | an optional description of your schedule. |
| isRecurring | false | whether the schedule is recurring or not. |
| cronExpr | if isRecurring = true | cron expression for recurring schedules. |
| url | true | webhook notification endpoint. |
| email | false | email address for notifying repeated failures. |
| runAt | if isRecurring = false | for non-recurring schedules, it indicates the instant the schedule will be triggered at. |
| startAt | false | UTC start date of the schedule. Must be equal to runAt if isRecurring = false. |
| endAt | false | UTC end date of the schedule. Must be equal to runAt if isRecurring = false. |
| metadata | false | optional metadata which will be sent when triggering a webhook. |


## REST API

- **POST** `/schedules` - Register a new schedule
- **GET** `/schedules/{id}` - Get details about an already existing schedule
- **DELETE** `/schedules/{id}` - Delete a schedule
- **POST** `/schedules/{id}/pause` - Pause an active schedule
- **POST** `/schedules/{id}/resume` - Resume a paused schedule
- **POST** `/schedules/{id}/trigger` - Immediately trigger a notification for a given schedule

## Contact
Stefano Scafiti @ostafen

## License
Kronos source code is available under the **MIT** License.
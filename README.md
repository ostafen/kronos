# Kronos

Kronos allows you to periodically invoke your webhooks using cron expressions.

## Features:

- :zap: Easy to use REST API to schedule your webhooks;
- :alarm_clock: Complex scheduling using cron expressions;
- :mailbox_with_mail: Get email notification on repeated failures.

## Build (Go 1.9+)

Run the following command

```bash
foo@bar$ make build
```
to build an executable which will be output to the `bin` subfolder.

# Sample configuration

```yaml
logging:
  level: INFO
  format: JSON
alert:
  email:
    server: smtp-server-address:port
    address: yuor-email-address
    password: your-email-password # if you use gmail with 2FA enabled, you can use app password

port: 9175

store:
  driver: sqlite3 # currently, the only driver supported
```

# Registering a periodic schedule
To start getting some webhook notifications, let's add a new schedule which will be notified every minute:
```bash
curl -X POST localhost:9175/schedules -H 'Content-Type: application/json' -d \
'{
    "title": "sample-schedule",
    "description": "a sample schedule description",
    "cronExpr": "0/1 * * * *",
    "email": "your-notification-email",
    "url": "your-webhook-address"
}'
```
On success, the response of the server will be similar to the following:
```json
{
    "id": "1e6d146b-e3b7-4e5c-b7ce-b7b2860f461b",
    "title": "sample-schedule",
    "status": "active",
    "description": "a sample schedule description",
    "cronExpr": "0/1 * * * *",
    "email": "your-notification-email",
    "url": "your-webhook-address",
    "metadata": null,
    "createdAt": "2023-02-18T09:38:08.72077066Z",
    "nextScheduleAt": "0001-01-01T00:00:00Z"
}
```

## REST API

- **POST** `/schedules` - Register a new schedule
- **GET** `/schedules/{id}` - Get details about an already existing schedule
- **DELETE** `/schedules/{id}` - Delete a schedule
- **POST** `/schedules/{id}/pause` - Pause an active schedule
- **POST** `/schedules/{id}/resume` - Resume a paused schedule
- **POST** `/schedules/{id}/trigger` - Immediately trigger a notification for a given schedule

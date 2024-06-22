# User News API

With this tool, it is possible send email notifications. Also, it is integrated with a request rate limiter that protects the users.

Currently, the rate limiter allows to sent messages with the following rules:

- Status: not more than 2 per minute for each user
- News: not more than 1 per day for each user
- Marketing: not more than 3 per hour for each user

But they can be changed in /ratelimiter/config.go file.

## Software requirements

It is necessary to have installed:

- Docker
- Docker-compose

## How it is work?

This can be consumed throughout the curl:

`
curl --location 'http://localhost:8080/notifications' \
--header 'Content-Type: application/json' \
--data-raw '{
"user_email": "example@gmail.com",
"message_type": "Status"
}'
`

Then, the user email will receive a new message (check the spam).

The message_type must be configured previously. By default, only "Status", "News" and "Marketing" types are allowed.   

Also, at the project root, you can find an importable postman collection named postman_collection.jon

## How does it launch the application?

You only need to go to the root of the project and do:

`
docker-compose up
`

This will create two Docker containers (and Docker images), one for the Golang application that will be available throughout port 8080 and another for the Redis database that will be available throughout port 6379. If some of these ports are not available, you need to change them in the docker-compose.yml file.

### Before launching the application

It is necessary to set up the environment variables for the application container, you can find them in docker-compose.yml:

- NOTIFIER_SENDER: It is the email address that will send the report.
- NOTIFIER_PASSWORD: It is the password associated with NOTIFIER_SENDER. For Gmail, it has to be an app password ([how do I create one?](https://support.google.com/mail/answer/185833?hl=en)), but for others, you must find out.
- NOTIFIER_HOST: Host of the email address. By default, the Gmail host is established.
- NOTIFIER_PORT: Port of the email address. By default, the Gmail port is established.
- REDIS_ADDRESS: Address asked by Redis, for docker-compose example is already set.
- REDIS_PASSWORD: Password asked by Redis, for docker-compose example is already set.

dp-identity-api
================
An API used to manage the authorisation of users accessing data publishing services.

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :25600    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 20s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.


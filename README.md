consul-ate
=======
A tool for tracking a processes life in consul.

[![Build Status](https://drone.io/github.com/JeremyOT/consul-ate/status.png)](https://drone.io/github.com/JeremyOT/consul-ate/latest)


Usage
-----

```bash
Usage:
  Registers a service with the local consul agent and tracks its health. If the command
  exits without an error, the corresponding service will be deregistered. Otherwise, it
  will fail its status checks.
  To monitor an external command add the command argurments after all consulate options.
  E.g. consulate -name my-service -- my_script.sh arg1 arg2
  stdin, stdout and stderr will be piped through from the monitored command.

Options:
  -check="": JSON containing a Check. If specified, will replace the default TTL check.
  -consul="localhost:8500": The consul agent to connect to.
  -id="": The ID to register with.
  -interval=10ns: The health check update interval.
  -name="": The service name to register with.
  -port=0: The port the service operates on.
  -tags="": Comma delimited tags to add when registering the service.
  -ttl=30ns: The health check ttl.
```

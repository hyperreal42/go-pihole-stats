# go-pihole-stats

WIP

Get stats from your Pi-hole instance at the command-line.

To use, one must have a running Pi-hole instance on one's network. Export the following environment variables:

Note: the auth token is stored in the `WEBPASSWORD` variable in `/etc/pihole/setupVars.conf` on the Pi-hole machine.

```bash
export PIHOLE_STATS_URL="http://pi.hole/admin"
export PIHOLE_STATS_AUTH="auth token"
```

Build and/or run the program:
```bash
go build main.go
./main

# or
go run main.go
```
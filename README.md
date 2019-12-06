# god

Lightweight supervisord-like thing for use as pid-eins in containers

CAUTION: this code have not seen much battletesting

# Goals

* Simple limited featureset
* If any of god's processes exits, try to shutdown other processes or timeout
* ~~Reap zombies~~ (Use `--init` for this, available for [docker-run](https://docs.docker.com/engine/reference/run/), [docker-service-create](https://docs.docker.com/engine/reference/commandline/service_create/) and the [docker-compose format](https://docs.docker.com/compose/compose-file/#init))
* Basic dependency-based start of processes

# Known issues

* If shutdown signal is received while god waits for dependencies to fulfill - these processes may start even when supposed not to.

# Signal logic

If SIGTERM, SIGQUIT or SIGINT is received, god will try to shutdown by sending SIGQUIT (can be overwritten with `stop_signal`) to
all processes.

# config.yml

```yaml
shutdown_timeout: 9s # defaults to 1 minute
processes:
  - name: nginx
    cmd: nginx -g 'daemon off;' -c $NGINXCONF
    # Setting bash: true allows one to use environment variables
    # and other bash hackery by passing `cmd` to `bash -c` for convenience
    bash: true

  - name: cloudflared
    cmd: cloudflared tunnel --no-autoupdate --unix-socket /var/run/website.sock
    stop_signal: SIGQUIT # if for whatever reason the default SIGTERM wont do 
    # this process will not start until a `/var/run/website.sock` 
    requires:
      - exists: /var/run/website.sock
        timeout: 10s

  - name: php
    cmd: php-fpm

  # this "process" exits for demo purposes making the daemon send stop signals
  # to other processes and eventually exit
  - name: stop-everything-after-a-minute
    cmd: sleep 60
```

# Docker example

Build
```
docker build -t cego/god .
```

Then use the image in a another docker multistage build

```
.... 

COPY --from=cego/god /god /god
COPY config.yml /config.yml

....

CMD ["/god"]
```


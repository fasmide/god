# god

Lightweight supervisord-like thing for use as pid-eins in containers

CAUTION: this code have not seen much battletesting

# Goals

* Simple limited featureset
* If any of god's processes exits, pull the entire daemon down
* Reap zombies
* Basic dependency-based start of processes

# Known issues

* If shutdown signal is received while god waits for dependencies to fulfill - these processes will start anyways.
* Reaper and go's os/exec "races" for exit-codes. 

# Signal logic

We are only acting on SIGTERM (the default docker swarm, docker stop signal), which in turn will 
send SIGTERM (by default) to all processes and wait forever for them to exit. 
Docker will take action if this is too long

# config.yml

```yaml
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

  # this "process" eventually pulls the entire daemon down
  - name: stop-everything-after-a-minute
    cmd: sleep 60
```

# example for dockerfiles


```
.... 

COPY --from=registry.cego.dk/cego/god:v1 /god /god
COPY config.yml /config.yml

....

CMD ["/god"]
```


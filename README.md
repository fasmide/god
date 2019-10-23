# god

Lightweight supervisord-like thing for use as pid-eins in containers

CAUTION: this code have not seen much battletesting

# Goals

* Simple limited featureset
* If any of god's processes exits, pull the entire daemon down
* Basic dependency-based start of processes

# Needed stuff

* use the slice of requirements instead of hardcoding requirements[0]
* take care of signals and propagate them to processes (or ?)

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


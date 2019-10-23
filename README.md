# god

Lightweight supervisord-like thing for use as pid-eins in containers

# Goals

* Simple limited featureset
* If any of god's processes exits, pull the entire container down
* Basic dependency-based start of processes

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
      exists: /var/run/website.sock
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


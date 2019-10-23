package god

import "fmt"

// ExampleYml outputs an example configuration file
func ExampleYml() {
	fmt.Println(`processes:

  - name: nginx
    cmd: nginx -g 'daemon off;' -c $NGINXCONF
    # Setting bash: true allows one to use environment variables
    # and other bash hackery by passing cmd to bash -c for convenience
    bash: true

  - name: cloudflared
    cmd: cloudflared tunnel --no-autoupdate --unix-socket /var/run/website.sock
    # this process will not start until a /var/run/website.sock 
    requires:
      exists: /var/run/website.sock
      timeout: 10s

  - name: php
    cmd: php-fpm

  # this "process" eventually pulls this entire runtime down
  - name: stop-everything-after-a-minute
    cmd: sleep 60`)
}

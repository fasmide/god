package god

import "fmt"

// ExampleYml outputs an example configuration file
func ExampleYml() {
	fmt.Println(`processes:
    - 
      name: Sleeper
      cmd: sleep infinity
    
    # This service will make god fail after 20 seconds
    - 
      name: no-so-sleeper
      cmd: sleep 20
    - 
      name: just_dmesg
      cmd: dmesg -w

    # Set bash: true to pass command to bash
    -
      name: random
      cmd: pv -q -L 100 /dev/urandom | xxd
      bash: true

    # All environment variables are passed to everyone
    - 
      name: hello-service
      cmd: while true; do echo "hej hej $NAME"; sleep 1; done
      bash: true`)
}

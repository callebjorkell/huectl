# huectl
My own take on a command line control for Philips Hue lights. Small CLI tool using cobra, and huego

```shell
huectl controls a Philips Hue installation

Usage:
  huectl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  connect     Connects to a bridge, storing the new user in a config. The link button must be pressed before running this command.
  help        Help about any command
  lights      Manage individual lights
  version     Prints the version of huectl

Flags:
      --debug   Turn on debug logging.
  -h, --help    help for huectl

Use "huectl [command] --help" for more information about a command.
```
## Credits
As the saying goes:
> Give credit where credit is due

* https://cobra.dev/
* https://github.com/amimof/huego is an excellent library for controlling the Hue lights, and this is what is used in huectl 
* https://github.com/skwair/huectl Has a similar tool written in go (and unfortunately was the first to use the name :laughing:). I thought 
  about extending it, but I also wanted to add a different structure to the commands and port it to huego, and hence started pretty much from 
  scratch. Some elements of this tool has definitely been borrowed from skwair though.

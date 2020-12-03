<div class="info" align="center">
  <h1 class="name">📬<br>postaco</h1>
  Convert postman collections to local documentation server
  <br>
  <br>
</div>


## Installation

Download latest binary release from release page. You can use this install script to download the latest version:

```sh
# install latest release to /usr/local/bin/
curl https://i.jpillora.com/avrebarra/postaco! | *remove_this* bash

# install specific version
curl https://i.jpillora.com/avrebarra/postaco@{version} | *remove_this* bash

# start postaco
postaco -help

```

## Usage
### CLI

```bash
$ postaco -help

postaco v1 - service

Available commands:

   build   build documentation folder 

Flags:

  -help
        Get help on the 'postaco' command.
  -port int
        port to bind (default: 8877) (default 8877)
  -quiet
        perform quiet operation
```

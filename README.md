<div class="info" align="center">
  <h1 class="name">📬<br>postaco</h1>
  postman collections as local documentation server
  <br>
  <br>

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]

</div>


## Installation

Download latest binary release from release page.

You can use this install script to download the latest version:

```sh
# install latest release to /usr/local/bin/
curl https://i.jpillora.com/avrebarra/postaco! | *remove_this* bash
```

```sh
# install specific version
curl https://i.jpillora.com/avrebarra/postaco@{version} | *remove_this* bash
```

## Usage
### CLI

```bash
$ postaco -help

postaco v0 - postman collections as local documentation server

Available commands:

   build   build documentation folder 

Flags:

  -help
        Get help on the 'postaco' command.
  -quiet
        perform quiet operation
```

[godoc-image]: https://godoc.org/github.com/avrebarra/postaco?status.svg
[godoc-url]: https://godoc.org/github.com/avrebarra/postaco
[report-image]: https://goreportcard.com/badge/github.com/avrebarra/postaco
[report-url]: https://goreportcard.com/report/github.com/avrebarra/postaco
[tests-image]: https://cloud.drone.io/api/badges/avrebarra/postaco/status.svg
[tests-url]: https://cloud.drone.io/avrebarra/postaco
[coverage-image]: https://codecov.io/gh/avrebarra/postaco/graph/badge.svg
[coverage-url]: https://codecov.io/gh/avrebarra/postaco
[sponsor-image]: https://img.shields.io/badge/github-donate-green.svg
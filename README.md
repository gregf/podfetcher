[![Build Status](https://travis-ci.org/gregf/podfetcher.svg?branch=master)](https://travis-ci.org/gregf/podfetcher)
# podfetcher

## Description

Podfetcher is a simple Podcast Fetcher written in Go. It will download all your favorite podcasts including youtube subscriptions for later viewing.

## Requirements

* [youtube-dl](https://rg3.github.io/youtube-dl/) (Optional if you are not using youtube feeds)

## Install

If you are on `linux amd64` there are binary releases. Check the latest [release](https://github.com/gregf/podfetcher/releases), extract the tarball cd and run make install (sudo not needed).

Alterntively you can run `go get github.com/gregf/podfetcher`. See the [Go Install Guide](https://golang.org/doc/install#install) for more information. After which you'll want to copy the example configuration files to your podfetcher config directory. 

```
mkdir ~/.config/podfetcher
cp $GOPATH/src/github.com/gregf/podfetcher/examples/* ~/.config/podfetcher/
```

# Usage

```
Usage:
  podfetcher [command]

Available Commands:
  version     Print the version number of Hugo
  update      Updates the database with the latest episodes to be fetched.
  fetch       Fetches podcast episodes that have not been downloaded.
  catchup     Marks all podcast episodes as downloaded
  lsnew       Display new episodes to be downloaded.
  import      Import feeds from a opml file.
  help        Help about any command

Flags:
  -h, --help=false: help for podfetcher


Use "podfetcher help [command]" for more information about a command.
```

You can add feeds to the feeds file `~/.config/podfetcher/feeds` run podfetcher update and podfetcher fetch to start downloading. It will only grab the last 10 episodes by default. You can change the `episodes` setting in `config.yml` if so desired.

You can also add filters to `config.yml` that will allow you to skip podcast episodes you do not want to see. To add a filter enter the podcast title under the filters section. Followed by a list of filters. You can get the podcast title from lsnew usually.

```
config.yml
...
filters:
  Some Show:
    - "Filter me"
    - "Another Filter"

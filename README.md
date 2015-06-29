# podfetcher

## Description

Podfetcher is a simple Podcast Fetcher written in Go. It will download all your favorite podcasts including youtube subscriptions for later viewing.

You can also add filters to `config.yml` that will allow you to skip podcast episodes you do not want to see. To add a filter enter the podcast title under the filters section. Followed by a list of filters. You can get the podcast title from lsnew usually.

```
config.yml
...
filters:
  Some Show:
    - "Filter me"
    - "Another Filter"
```

## Requirements

* [youtube-dl](https://rg3.github.io/youtube-dl/) (Optional if you are not using youtube feeds)

## Usage

Fetch the latest [release](https://github.com/gregf/podfetcher/releases), extract the tarball cd and run make install.

You can add feeds to the feeds file `~/.podfetcher/feeds` run podfetcher update and podfetcher fetch to start downloading. It will only grab the last 10 episodes by default. You can change the `episodes` setting in `config.yml` if so desired.


# TODO

* Use libnotify to notify when a download completes
* Transaction support so we hopefully don't have to use the lockfile in the future
* Inclusive filter
* Pause/UnPause podcast
  - type in name of podcast
  - get the name we think you want back ask if this is it
  - if yes change paused to true
* lspodcasts (would this be useful?)
* lsepisodes (takes a podcast name)
  - shows filtered
  - shows downloaded
  - date downloaded?
* add (appends feed to feeds file)
* import should optionally take a url to a opml file
* browse (opens podcast dir with main.browse-program which will default to ranger)
  - learn how to detach the process
* Have a minimal update time
  - respect headers
  - update podcasts on update command
* Tests! Start with something simple like helpers
* Improve comments
* Setup Command to self generate config on startup if none present.
* Add some color? Make sure users can turn this off.
* Fix bug where filenames include /'s like the /dev/hell podcast

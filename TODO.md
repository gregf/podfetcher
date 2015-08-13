# TODO

* Move database calls into commands, This should simplify things a bit
* Sanitize filepaths/names with https://github.com/kennygrant/sanitize
* Get rid of termtables and display tables manually
* Display podcast id and episode ids where it makes sense
* Fetch command should take a castid
* Transaction support so we hopefully don't have to use the lockfile in the future
* Inclusive filter
* lsepisodes (takes a podcast name)
  - shows filtered
  - shows downloaded
  - date downloaded?
* import should optionally take a url to a opml file
* browse (opens podcast dir with main.browse-program which will default to ranger)
  - learn how to detach the process
* Have a minimal update time
  - respect headers
  - update podcasts on update command
* Tests! Start with something simple like helpers
* Improve comments
* Add some color? Make sure users can turn this off.
* Compare .config/podfetcher/feeds with list in database & remove feeds no longer in feeds file

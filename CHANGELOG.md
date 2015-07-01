# Unreleased

* Add a optional notification after each download
* Log more useful errors with run commands
* More useful output for lsnew
* Auto generate a config file if one is not present

# v0.3 / 2015-06-29

* Added the ability to filter out unwanted podcast episodes.
* Add version command
* Switched from toml to yml for configuration
* Add lockfile to protect database
* Call ToLower() when we are comparing urls
* Custom downloader function so we can control the output format

# v0.2 / 2015-06-26

* use XDG_CONFIG_HOME and XDG_CACHE_HOME
* expandPath on ~ for config.toml
* Added Guid to the item table for another unique_index
* Add PubDate to episodes table
* New command lsnew displays episodes to be downloaded
* New func readFeeds that allows the feeds file to contain blank lines and comments
* Went over all the code with golint
* Add Import command, to import feeds from a opml file.

# v0.1 / 2015-06-23

* Initial Release

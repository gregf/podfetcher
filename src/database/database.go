package database

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/jinzhu/gorm"
	// required by gorm
	_ "github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/mattn/go-sqlite3"
)

// Podcast struct
type Podcast struct {
	ID       int `sql:"index"`
	Title    string
	RssURL   string `sql:"unique_index"`
	Episodes []Episode
}

// Episode struct
type Episode struct {
	ID           int `sql:"index"`
	PodcastID    int
	Title        string
	EnclosureURL string `sql:"unique_index"`
	Downloaded   bool
	GUID         string `sql:"unique_index"`
	PubDate      string
}

func databasePath() (path string) {
	if len(os.Getenv("XDG_CACHE_HOME")) > 0 {
		path = filepath.Join(os.Getenv("XDG_CACHE_HOME"), "podfetcher")
		os.MkdirAll(path, 0755)
		return filepath.Join(path, "cache.db")
	}
	path = filepath.Join(os.Getenv("HOME"), ".cache", "podfatcher")
	os.MkdirAll(path, 0755)
	return filepath.Join(path, "cache.db")
}

func init() {
	db := DBSession()
	db.LogMode(false)

	db.CreateTable(&Podcast{})
	db.CreateTable(&Episode{})
	db.AutoMigrate(&Podcast{}, &Episode{})
}

// DBSession sets up a database session
func DBSession() (db gorm.DB) {
	sqliteSession, err := gorm.Open("sqlite3", databasePath())
	if err != nil {
		log.Fatal(err)
	}

	return sqliteSession
}

// SetDownloadedByURL updates all downloaded columns to be true
func SetDownloadedByURL(url string) {
	db := DBSession()
	db.LogMode(false)

	db.Table("episodes").
		Where("enclosure_url = ?", url).
		UpdateColumn("downloaded", true)
}

// FindEpisodeTitleByURL finds episode titles by url
func FindEpisodeTitleByURL(url string) (title string) {
	db := DBSession()
	db.LogMode(false)

	row := db.Table("episodes").
		Where("enclosure_url = ?", url).
		Select("title").
		Row()
	row.Scan(&title)

	return title
}

// FindPodcastTitleByURL finds podcast titles by URL
func FindPodcastTitleByURL(url string) (title string) {
	db := DBSession()
	db.LogMode(false)

	var podcastID int
	row := db.Table("episodes").Where("enclosure_url = ?", url).Select("podcast_id").Row()
	row.Scan(&podcastID)

	prow := db.Table("podcasts").
		Where("id = ?", podcastID).
		Select("title").
		Row()
	prow.Scan(&title)

	return title
}

// FindNewEpisodes finds episodes where downloaded = false
func FindNewEpisodes() (urls []string, err error) {
	db := DBSession()
	db.LogMode(false)

	rows, err := db.Table("episodes").
		Where("downloaded = ?", false).
		Select("enclosure_url").
		Rows()
	defer rows.Close()

	for rows.Next() {
		var enclosureURL string
		rows.Scan(&enclosureURL)
		urls = append(urls, enclosureURL)
	}
	return urls, err
}

// findPodcastID locates podcast ID by rssURL
func findPodcastID(rssurl string) (podcastID int) {
	db := DBSession()
	db.LogMode(false)

	row := db.Table("podcasts").
		Where("rss_url = ?", rssurl).
		Select("id").
		Row()
	row.Scan(&podcastID)
	return podcastID
}

// AddPodcast Inserts a new podcast into the database
func AddPodcast(title, rssurl string) {
	db := DBSession()
	db.LogMode(false)

	podcast := Podcast{
		Title:  title,
		RssURL: strings.ToLower(rssurl),
	}
	if db.NewRecord(&podcast) {
		db.Create(&podcast)
	}
}

// AddItem takes a map[string]string of episode items to be inserted into the
// database.
//
// item[rssURL] item[title], item[enclosureURL], item[guid], items[pubdate]
func AddItem(items map[string]string) {
	db := DBSession()
	db.LogMode(false)

	podcastID := findPodcastID(items["rssURL"])

	episode := Episode{
		Title:        items["title"],
		EnclosureURL: strings.ToLower(items["enclosureURL"]),
		Downloaded:   false,
		PodcastID:    podcastID,
		GUID:         items["guid"],
		PubDate:      items["pubdate"],
	}

	if db.NewRecord(&episode) {
		db.Create(&episode)
	}
}

// CatchUp Marks all downloaded = false to be downloaded = true
func CatchUp() {
	db := DBSession()
	db.LogMode(false)

	db.Table("episodes").Where("downloaded = ?", false).
		UpdateColumn("downloaded", true)
}

// FindEpisodesWithPodcastTitle Finds episodes with their podcast title and
// returns a map[string]string
func FindEpisodesWithPodcastTitle() (m map[string][]string) {
	db := DBSession()
	db.LogMode(false)

	rows, err := db.Table("Episodes").
		Where("downloaded = ?", false).
		Select("title, podcast_id").
		Rows()
	if err != nil {
		log.Fatal(err)
	}

	m = make(map[string][]string)

	for rows.Next() {
		var eptitle string
		var podcastID int
		var title string
		rows.Scan(&eptitle, &podcastID)
		row := db.Table("podcasts").
			Where("id =?", podcastID).
			Select("title").
			Row()
		row.Scan(&title)
		m[title] = append(m[title], eptitle)
	}

	return m
}

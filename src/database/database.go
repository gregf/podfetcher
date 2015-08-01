package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/gohome"
	"github.com/jinzhu/gorm"
	// required by gorm
	_ "github.com/mattn/go-sqlite3"
)

var db gorm.DB

const appName = "podfetcher"

// Podcast struct
type Podcast struct {
	ID       int `sql:"index"`
	Title    string
	RssURL   string `sql:"unique_index"`
	Paused   bool
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
	path = gohome.Cache(appName)
	os.MkdirAll(path, 0755)
	return filepath.Join(path, "cache.db")
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", databasePath())
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)
	db.CreateTable(&Podcast{Paused: false})
	db.CreateTable(&Episode{})
	db.AutoMigrate(&Podcast{}, &Episode{})
}

// SetDownloadedByURL updates all downloaded columns to be true
func SetDownloadedByURL(url string) {
	db.Table("episodes").
		Where("enclosure_url = ?", url).
		UpdateColumn("downloaded", true)
}

// FindEpisodeTitleByURL finds episode titles by url
func FindEpisodeTitleByURL(url string) (title string) {
	row := db.Table("episodes").
		Where("enclosure_url = ?", url).
		Select("title").
		Row()
	row.Scan(&title)

	return title
}

// FindPodcastTitleByURL finds podcast titles by URL
func FindPodcastTitleByURL(url string) (title string) {
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
	rows, err := db.Table("episodes").
		Where("downloaded = ?", false).
		Select("podcast_id, enclosure_url").
		Rows()
	defer rows.Close()

	for rows.Next() {
		var enclosureURL string
		var podcastID int
		rows.Scan(&podcastID, &enclosureURL)
		paused := FindPodcastPausedState(podcastID)
		if paused {
			return
		}
		urls = append(urls, enclosureURL)
	}
	return urls, err
}

// FindAllPodcasts Find all podcasts and their IDs
func FindAllPodcasts() (ids []int, titles []string) {
	rows, err := db.Table("podcasts").Select("id, title").Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ID int
		var title string
		rows.Scan(&ID, &title)
		ids = append(ids, ID)
		titles = append(titles, title)
	}

	return ids, titles
}

// findPodcastID locates podcast ID by rssURL
func findPodcastID(rssurl string) (podcastID int) {
	row := db.Table("podcasts").
		Where("rss_url = ?", rssurl).
		Select("id").
		Row()
	row.Scan(&podcastID)
	return podcastID
}

// AddPodcast Inserts a new podcast into the database
func AddPodcast(title, rssurl string) {
	podcast := Podcast{
		Title:  title,
		RssURL: rssurl,
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
	podcastID := findPodcastID(items["rssURL"])

	episode := Episode{
		Title:        items["title"],
		EnclosureURL: items["enclosureURL"],
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
func CatchUp(id int) {
	if id == 0 {
		db.Table("episodes").Where("downloaded = ?", false).
			UpdateColumn("downloaded", true)
	} else {
		db.Table("episodes").Where("podcast_id = ?", id).
			Where("downloaded = ?", false).
			UpdateColumn("downloaded", true)
	}
}

// FindEpisodesWithPodcastTitle Finds episodes with their podcast title and
// returns a map[string]string
func FindEpisodesWithPodcastTitle() (m map[string][]string) {
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

		paused := FindPodcastPausedState(podcastID)
		if paused {
			return
		}

		m[title] = append(m[title], eptitle)
	}

	return m
}

//FindPodcastPausedState finds out wether or not a podcast is paused
func FindPodcastPausedState(id int) (paused bool) {
	row := db.Table("podcasts").Where("id = ?", id).Select("paused").Row()
	row.Scan(&paused)

	return paused
}

//TogglePause toggles between paused states true and false
func TogglePause(id int) (paused bool) {
	state := FindPodcastPausedState(id)
	paused = false
	if state == false {
		paused = true
	}

	db.Table("podcasts").Where("id = ?", id).
		UpdateColumn("paused", paused)

	return paused
}

// FindPodcastTitle looks up a podcast title by its id
func FindPodcastTitle(id int) (title string) {
	row := db.Table("podcasts").Where("id = ?", id).Select("title").Row()
	row.Scan(&title)

	return title
}

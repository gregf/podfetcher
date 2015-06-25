package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/mattn/go-sqlite3"
)

type Podcast struct {
	Id       int `sql:"index"`
	Title    string
	RssUrl   string `sql:"unique_index"`
	Episodes []Episode
}

type Episode struct {
	Id           int `sql:"index"`
	PodcastID    int
	Title        string
	EnclosureUrl string `sql:"unique_index"`
	Downloaded   bool
	Guid         string `sql:"unique_index"`
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
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)

	if !dbexists(databasePath()) {
		db.CreateTable(&Podcast{})
		db.CreateTable(&Episode{})
	}
	db.AutoMigrate(&Podcast{}, &Episode{})
}

func dbexists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func DBSession() (db gorm.DB, err error) {
	sqliteSession, err := gorm.Open("sqlite3", databasePath())
	if err != nil {
		log.Fatal(err)
	}

	return sqliteSession, err
}

func SetDownloadedByUrl(url string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	db.Table("episodes").Where("enclosure_url = ?", url).UpdateColumn("downloaded", true)
}

func FindEpisodeTitleByUrl(url string) (title string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	row := db.Table("episodes").Where("enclosure_url = ?", url).Select("title").Row()
	row.Scan(&title)

	return title
}

func FindPodcastTitleByUrl(url string) (title string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	var podcastId int
	row := db.Table("episodes").Where("enclosure_url = ?", url).Select("podcast_id").Row()
	row.Scan(&podcastId)

	prow := db.Table("podcasts").Where("id = ?", podcastId).Select("title").Row()
	prow.Scan(&title)

	return title
}

func FindNewEpisodes() (urls []string, err error) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	var total_count int
	rows, err := db.Table("episodes").Where("downloaded = ?",
		false).Select("enclosure_url").Count(&total_count).Rows()
	defer rows.Close()

	urls = make([]string, 0, total_count)
	for rows.Next() {
		var enclosure_url string
		rows.Scan(&enclosure_url)
		urls = append(urls, enclosure_url)
	}
	return urls, err
}

func findPodcastID(rssurl string) (podcastId int) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	row := db.Table("podcasts").Where("rss_url = ?", rssurl).Select("id").Row()
	row.Scan(&podcastId)
	return podcastId
}

func AddPodcast(title string, rssurl string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(false)

	podcast := Podcast{
		Title:  title,
		RssUrl: rssurl,
	}
	if db.NewRecord(&podcast) {
		db.Create(&podcast)
	}
}

// item[rssUrl] item[title], item[enclosureUrl], item[guid], items[pubdate]
func AddItem(items map[string]string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(false)

	podcastId := findPodcastID(items["rssUrl"])

	episode := Episode{
		Title:        items["title"],
		EnclosureUrl: items["enclosureUrl"],
		Downloaded:   false,
		PodcastID:    podcastId,
		Guid:         items["guid"],
		PubDate:      items["pubdate"],
	}

	if db.NewRecord(&episode) {
		db.Create(&episode)
	}
}

func CatchUp() {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(false)

	db.Table("episodes").Where("downloaded = ?", false).
		UpdateColumn("downloaded", true)
}

func FindEpisodesWithPodcastTitle() (m map[string][]string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(false)

	rows, err := db.Table("Episodes").Where("downloaded = ?",
		false).Select("title, podcast_id").Rows()

	m = make(map[string][]string)

	for rows.Next() {
		var eptitle string
		var podcastId int
		var title string
		rows.Scan(&eptitle, &podcastId)
		row := db.Table("podcasts").Where("id =?",
			podcastId).Select("title").Row()
		row.Scan(&title)
		m[title] = append(m[title], eptitle)
	}

	return m
}

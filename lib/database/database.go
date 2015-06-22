package database

import (
	"log"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
}

var database = path.Join(os.Getenv("HOME"), "/.podfetcher/cache.db")

func init() {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)

	if !dbexists(database) {
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
	sqliteSession, err := gorm.Open("sqlite3", database)
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

func FindTitleByUrl(url string) (title string) {
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

func AddItem(title string, rssurl string, enclosureurl string) {
	db, err := DBSession()
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(false)

	podcastId := findPodcastID(rssurl)

	episode := Episode{
		Title:        title,
		EnclosureUrl: enclosureurl,
		Downloaded:   false,
		PodcastID:    podcastId,
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

	db.Table("episodes").Where("downloaded = ?", false).UpdateColumn("downloaded", true)
}

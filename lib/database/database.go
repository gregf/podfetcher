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
	db.CreateTable(&Podcast{})
	db.CreateTable(&Episode{})
	db.AutoMigrate(&Podcast{}, &Episode{})
}

func DBSession() (db gorm.DB, err error) {
	sqliteSession, err := gorm.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	return sqliteSession, err
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

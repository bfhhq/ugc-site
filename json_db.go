package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type json_db struct {
	fileName string
	videos   []*Video
}

func NewJsonDb(fileName string) (*json_db, error) {
	return &json_db{
		fileName: fileName,
	}, nil
}

func (db *json_db) LoadVideos() error {

	if fd, err := os.Open(db.fileName); err == nil {
		jsonStr, _ := ioutil.ReadAll(fd)
		json.Unmarshal(jsonStr, &db.videos)
	}

	log.Println(db.videos)

	return nil
}

func (db *json_db) SaveVideos() error {

	fd, err := os.Create(db.fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	jsonStr, err := json.Marshal(db.videos)

	fd.WriteString(string(jsonStr))

	return err
}

func (db *json_db) InsertVideo(v *Video) error {

	db.videos = append([]*Video{v}, db.videos...)

	return nil
}

func (db *json_db) DeleteVideo(v *Video) error {

	for i, vv := range db.videos {
		if vv == v {
			db.videos = append(db.videos[:i], db.videos[i+1:]...)
			return nil
		}
	}

	return nil
}

func (db *json_db) FindVideos(sql string) (*[]*Video, error) {

	return &db.videos, nil
}

func (db *json_db) GetVideo(name string) (*Video, bool) {
	for _, v := range db.videos {
		if v.Name == name {
			return v, true
		}
	}

	return nil, false
}

func (db *json_db) Close() error {
	return nil
}

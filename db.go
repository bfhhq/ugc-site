package main

type DB interface {
	LoadVideos() error
	SaveVideos() error
	InsertVideo(v *Video) error
	DeleteVideo(v *Video) error
	FindVideos(sql string) (*[]*Video, error)
	GetVideo(name string) (*Video, bool)
	Close() error
}

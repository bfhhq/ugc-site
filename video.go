package main

type Video struct {
	Name        string
	Title       string
	Description string `json:"desc"`
	Url         string
}

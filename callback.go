package main

import (
	"encoding/json"
	"fmt"
	"github.com/baofengcloud/go-sdk/src/baofengcloud"
	"github.com/bitly/go-nsq"
	"log"
)

type CallbackData struct {
	Status   int
	Type     string
	FileName string
	FileType baofengcloud.FileType `json:"ifpublic"`
	Url      string
}

/*
{"status":0,"type":"upload","fileid":"C547C83
285A0ADEF72A459B29C9B04D0","servicetype":1,"filename":"vod_hls_mp3.ts","showname
":"vod_hls_mp3.ts","filekey":"","filesize":974404,"duration":9641,"uploadtime":"
2015-03-12 10:49:23","publishtime":"2015-03-12 10:49:57","ifpublic":1,"url":"ser
vicetype=1&uid=5119278&fid=C547C83285A0ADEF72A459B29C9B04D0"}
*/

func checkCallback() {

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("bfcloud", "site", config)
	q.AddHandler(nsq.HandlerFunc(recvCallback))
	err := q.ConnectToNSQD(confFile.NsqdAddress)
	if err != nil {
		log.Println("Could not connect NSQ server")
	}
}

func recvCallback(message *nsq.Message) error {
	log.Printf("Got a message: %s", string(message.Body))

	data := CallbackData{}
	json.Unmarshal(message.Body, &data)

	if data.Type != "upload" {
		return nil
	}

	v, ok := db.GetVideo(data.FileName)
	if ok == false {
		return nil
	}

	if data.Status == 0 {
		v.Url = data.Url
		v.SwfUrl, _ = buildSwfPlayerUrl(data.FileType, data.Url, false)
	} else {
		v.Title += "(status)"
		v.Title = fmt.Sprintf("%s (status=%d)", v.Title, data.Status)
	}

	db.SaveVideos()

	return nil
}

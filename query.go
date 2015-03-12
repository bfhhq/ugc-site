package main

import (
	"fmt"
	"github.com/baofengcloud/go-sdk/src/baofengcloud"
	"time"
)

func getSwfPlayerUrl(fileName string, autoPlay bool) (string, error) {

	var fileInfo *baofengcloud.FileInfo
	var err error
	if fileInfo, err = baofengcloud.QueryFile(&conf, baofengcloud.Saas, fileName, ""); err != nil {
		return "", err
	}

	return buildSwfPlayerUrl(fileInfo.FileType, fileInfo.Url, autoPlay)
}

func buildSwfPlayerUrl(fileType baofengcloud.FileType, url string, autoPlay bool) (string, error) {
	var playurl string
	var err error
	if playurl, err = baofengcloud.BuildSwfPlayUrl(&conf, fileType, url, "", 1*time.Hour); err != nil {
		return "", err
	}

	autoFlag := 0
	if autoPlay {
		autoFlag = 1
	}

	return fmt.Sprintf("%s&auto=%d", playurl, autoFlag), nil

}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/baofengcloud/go-sdk/src/baofengcloud"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ConfigFile struct {
	AccessKey   string
	SecretKey   string
	CallbackUrl string
	NsqdAddress string
	DataPath    string
}

var confFile ConfigFile
var confFilePath = "conf.json"

var db DB

var conf baofengcloud.Configure

func main() {

	if err := readConf(); err != nil {
		log.Fatal(err)
	}

	if err := openDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LoadVideos()

	go checkCallback()

	http.HandleFunc("/", root)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/delete", delete)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/api/token/upload", createUploadToken)
	http.HandleFunc("/api/video/all", videosJson)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readConf() error {

	if fd, err := os.Open(confFilePath); err == nil {
		jsonStr, _ := ioutil.ReadAll(fd)
		json.Unmarshal(jsonStr, &confFile)
	} else {
		return err
	}

	conf.AccessKey = confFile.AccessKey
	conf.SecretKey = confFile.SecretKey

	fmt.Printf("AK:%s \nSK:%s \nCallbackUrl:%s\n", conf.AccessKey, conf.SecretKey, confFile.CallbackUrl)

	return nil
}

func openDB() error {

	var err error
	db, err = NewJsonDb("db.json")
	return err
}

func root(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	videos, _ := db.FindVideos("")

	data := map[string]interface{}{}
	data["Videos"] = videos

	t.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		t, err := template.ParseFiles("upload.html")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		t.Execute(w, nil)
	} else if r.Method == "POST" {

		fileData, fileHeader, err := r.FormFile("file")

		if err != nil {
			fmt.Fprint(w, "upload file fail. ", err)
			return
		}

		if _, ok := db.GetVideo(fileHeader.Filename); ok == true {
			fmt.Fprintf(w, "%v already exist!", fileHeader.Filename)
			return
		}

		v := &Video{
			Name:        fileHeader.Filename,
			Title:       r.FormValue("title"),
			Description: r.FormValue("desc"),
		}

		db.InsertVideo(v)

		success := false
		defer func() {
			if success == false {
				deleteFileFromCloud(v.Name)
				db.DeleteVideo(v)
			}

			db.SaveVideos()
		}()

		localFilePath := filepath.Join(confFile.DataPath, v.Name)
		if fd, err := os.Create(localFilePath); err == nil {

			io.Copy(fd, fileData)

			fd.Close()

		} else {
			fmt.Fprint(w, "save file fail! ", err)
			return
		}

		err = baofengcloud.UploadFile2(&conf, baofengcloud.Saas,
			baofengcloud.Public, localFilePath, v.Name, "", confFile.CallbackUrl)

		if err != nil {
			fmt.Fprint(w, "upload fail, ", err)
		} else {
			success = true
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	v, ok := db.GetVideo(r.FormValue("name"))
	if ok == false {
		fmt.Fprint(w, "video not found!")
		return
	}

	err := deleteFileFromCloud(v.Name)

	db.DeleteVideo(v)
	db.SaveVideos()

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteFileFromCloud(name string) error {
	_, err := baofengcloud.DeleteFile(&conf, baofengcloud.Saas, name, "", confFile.CallbackUrl)

	return err
}

func createUploadToken(w http.ResponseWriter, r *http.Request) {

	fileName := r.FormValue("name")
	fileSize, _ := strconv.ParseInt(r.FormValue("size"), 10, 64)

	if len(fileName) == 0 || fileSize == 0 {
		http.Error(w, "name or size is empty!", http.StatusInternalServerError)
		return
	}

	if _, ok := db.GetVideo(fileName); ok == true {
		http.Error(w, "file already exist.", http.StatusInternalServerError)
		return
	}

	//insert to db, so we can trace upload progress and delete file in case upload fail.
	v := &Video{
		Name:        fileName,
		Title:       r.FormValue("title"),
		Description: r.FormValue("desc"),
	}

	db.InsertVideo(v)

	log.Printf("request upload token: %s  %s", v.Title, v.Name)

	token := baofengcloud.CreateUploadToken(conf.AccessKey, conf.SecretKey, baofengcloud.Saas,
		baofengcloud.Public, baofengcloud.Partial, fileName, "", fileSize, 1*time.Hour, confFile.CallbackUrl)

	//	result := map[string]interface{}{}
	//	result["token"] = token

	//	b, _ := json.Marshal(result)

	//	w.Header().Add("Content-Type", "application/json")

	//	fmt.Fprint(w, string(b))

	fmt.Fprint(w, token)
}

func videosJson(w http.ResponseWriter, r *http.Request) {
	videos, _ := db.FindVideos("")
	w.Header().Add("Content-Type", "application/json")

	b, _ := json.Marshal(videos)
	fmt.Fprint(w, string(b))
}

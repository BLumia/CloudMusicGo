package main

import "fmt"
import "log"
import "flag"
import "strconv"
import "net/http"
import "io/ioutil"
import "path/filepath"
import "encoding/json"

type MusicFileStruct struct {
	FileName     string `json:"fileName"`
	FileSize     int64  `json:"fileSize"`
	ModifiedTime string `json:"modifiedTime"`
}

type DataStruct struct {
	MusicList     []MusicFileStruct `json:"musicList"`
	SubFolderList []string          `json:"subFolderList"`
}

type ResultStruct struct {
	Type string     `json:"type"`
	Data DataStruct `json:"data"`
}

type ResponseStruct struct {
	Status  int32         `json:"status"`
	Message string        `json:"message"`
	Result  *ResultStruct `json:"result,omitempty"`
}

func Fire(res http.ResponseWriter, status int32, message string, result *ResultStruct) {
	// res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(int(status))

	resStruct := &ResponseStruct{
		Status:  status,
		Message: message,
		Result:  result,
	}
	resBtye, _ := json.Marshal(resStruct)
	fmt.Fprintf(res, string(resBtye))
}

func GetPlaylist(response http.ResponseWriter, request *http.Request) {
	var folderSlice []string
	var musicSlice []MusicFileStruct

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if f.IsDir() {
			folderSlice = append(folderSlice, f.Name())
		} else {
			var extension = filepath.Ext(f.Name())
			if extension == ".mp3" || extension == ".wav" {
				musicSlice = append(musicSlice, MusicFileStruct{
					FileName:     f.Name(),
					FileSize:     f.Size(),
					ModifiedTime: "3",
				})
			}
		}
	}

	dataStruct := &DataStruct{
		MusicList:     musicSlice,
		SubFolderList: folderSlice,
	}
	resStruct := &ResultStruct{
		Type: "fileList",
		Data: *dataStruct,
	}
	Fire(response, 200, "OK", resStruct)
}

func Controller(response http.ResponseWriter, request *http.Request) {
	log.Println(request.Method)
	if request.Method == "POST" {
		request.ParseForm()
		switch do := request.Form["do"][0]; do {
		case "getplaylist":
			GetPlaylist(response, request)
		default:
			Fire(response, 400, "Illegal request!", nil)
		}
	} else {
		Fire(response, 400, "Illegal request!", nil)
	}
}

func main() {
	var PortNumber int
	const (
		defaultPortNumber = 4000
		portUsage         = "The port number that the `Private Cloud Music - Go` should listen."
	)
	//var SongRoot string

	flag.IntVar(&PortNumber, "port", defaultPortNumber, portUsage)
	flag.IntVar(&PortNumber, "p", defaultPortNumber, portUsage+" (shorthand)")

	flag.Parse()
	fmt.Println("Now listening port: " + strconv.Itoa(PortNumber))

	http.HandleFunc("/", Controller)
	err := http.ListenAndServe(":"+strconv.Itoa(PortNumber), nil)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import "os"
import "fmt"
import "log"
import "flag"
import "strings"
import "strconv"
import "net/http"
import "io/ioutil"
import "path/filepath"
import "encoding/json"

var isFileServerEnabled bool
var SongRoot string

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

	SongFolder := request.PostFormValue("folder")

	log.Println(request.Method + " GetPlaylist, Folder:" + SongFolder)

	files, _ := ioutil.ReadDir(SongRoot + "/" + SongFolder)
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

func FileServer(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	log.Println(request.Method + " FileSrv, Req Path:" + path)

	if path == "/" {
		path = "/index.html"
	}

	requestType := path[strings.LastIndex(path, "."):]
	// Todo: media content-type
	switch requestType {
	case ".css":
		response.Header().Set("content-type", "text/css")
	case ".js":
		response.Header().Set("content-type", "text/javascript")
	case ".htm":
		fallthrough
	case ".html":
		response.Header().Set("content-type", "text/html")
	default:
	}

	fin, err := os.Open(SongRoot + path)
	if err != nil {
		log.Println("Request file [" + path + "] not exist.")
		http.NotFound(response, request)
	} else {
		fd, _ := ioutil.ReadAll(fin)
		response.Write(fd)
	}
}

func Controller(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		//request.ParseForm()
		switch do := request.PostFormValue("do"); do {
		case "getplaylist":
			GetPlaylist(response, request)
		default:
			log.Println(request.Method + " Missing or wrong `do` param.")
			Fire(response, 400, "Illegal request!", nil)
		}
	} else {
		if isFileServerEnabled {
			FileServer(response, request)
		} else {
			log.Println(request.Method + " FileServer not enabled, request illegal.")
			Fire(response, 400, "Illegal request!", nil)
		}
	}
}

func main() {
	var PortNumber int
	const (
		defaultPortNumber      = 4000
		portUsage              = "The port number that the `Private Cloud Music - Go` should listen."
		defaultFileServerState = false
		fileServerUsage        = "Enable a built-in file server for the audio files at root of song folder."
		defaultSongRoot        = "."
		SongRootUsage          = "Root of song folder path."
	)

	flag.IntVar(&PortNumber, "port", defaultPortNumber, portUsage)
	flag.IntVar(&PortNumber, "p", defaultPortNumber, portUsage+" (shorthand)")
	flag.BoolVar(&isFileServerEnabled, "fileserver", defaultFileServerState, fileServerUsage)
	flag.BoolVar(&isFileServerEnabled, "f", defaultFileServerState, fileServerUsage+" (shorthand)")
	flag.StringVar(&SongRoot, "root", defaultSongRoot, SongRootUsage)
	flag.StringVar(&SongRoot, "r", defaultSongRoot, SongRootUsage+" (shorthand)")

	flag.Parse()
	fmt.Println("Now listening port: " + strconv.Itoa(PortNumber))
	if isFileServerEnabled {
		fmt.Println("Built-in file server enabled, root: [" + SongRoot + "] .")
	}
	http.HandleFunc("/", Controller)
	err := http.ListenAndServe(":"+strconv.Itoa(PortNumber), nil)
	if err != nil {
		log.Fatal(err)
	}
}

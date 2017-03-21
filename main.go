package main

import "fmt"
import "log"
import "flag"
import "strconv"
import "net/http"
import "encoding/json"

type MusicFileStruct struct {
	FileName     string `json:"fileName"`
	FileSize     string `json:"fileSize"`
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
	Status  int32        `json:"status"`
	Message string       `json:"message"`
	Result  ResultStruct `json:"result"`
}

func Fire(res http.ResponseWriter, status int32, message string) {
	res.WriteHeader(int(status))
	// res.Header().Set("Access-Control-Allow-Origin", "*")
	resStruct := &ResponseStruct{
		Status:  status,
		Message: message} // 智障吧这个，`}` 必须写这儿？
	resBtye, _ := json.Marshal(resStruct)
	fmt.Fprintf(res, string(resBtye))
}

func GetPlaylist(response http.ResponseWriter, request *http.Request) {
	Fire(response, 501, "Not Implemented!")
}

func Controller(response http.ResponseWriter, request *http.Request) {
	log.Println(request.Method)
	if request.Method == "POST" {
		request.ParseForm()
		switch do := request.Form["do"][0]; do {
		case "getplaylist":
			GetPlaylist(response, request)
		default:
			Fire(response, 400, "Illegal request!")
		}
	} else {
		Fire(response, 400, "Illegal request!")
	}
}

func main() {
	var PortNumber int
	const (
		defaultPortNumber = 4000
		usage             = "The port number that the `Private Cloud Music - Go` should listen."
	)
	//var SongRoot string

	flag.IntVar(&PortNumber, "port", defaultPortNumber, usage)
	flag.IntVar(&PortNumber, "p", defaultPortNumber, usage+" (shorthand)")

	flag.Parse()
	fmt.Println("Now listening port: " + strconv.Itoa(PortNumber))

	http.HandleFunc("/", Controller)
	err := http.ListenAndServe(":"+strconv.Itoa(PortNumber), nil)
	if err != nil {
		log.Fatal(err)
	}
}

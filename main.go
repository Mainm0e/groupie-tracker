package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"text/template"
)

var temp *template.Template

type errorData struct {
    Num  int
    Text string
}

type ArtistsData struct {
    Id           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

var allData []ArtistsData

func gatherDataUp(link string) []ArtistsData {
    data1 := getData(link)
    Artists := []ArtistsData{}
    err := json.Unmarshal(data1, &Artists)
    if err != nil {
        log.Fatal(err)
        return nil
    }
    return Artists
}

func main() {
	allData = gatherDataUp("https://groupietrackers.herokuapp.com/api/artists")
    if allData == nil {
		fmt.Println("Failed to gather Data from API")
        os.Exit(1)
    }
	
	PORT := ":8080"
	http.HandleFunc("/", HomeHandler)
	fileServer := http.FileServer(http.Dir("./docs"))
	http.Handle("/docs/",http.StripPrefix("/docs/", fileServer))
	http.ListenAndServe(PORT,nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
	g := allData
	temp = template.Must(template.ParseGlob("docs/static/*.html"))
	temp.ExecuteTemplate(w,"index.html",g)
}


func getData(link string) []byte {
    data1, err1 := http.Get(link)
    if err1 != nil {
        log.Fatal(err1)
    }
    data2, err2 := ioutil.ReadAll(data1.Body)
    if err2 != nil {
        log.Fatal(err2)
        return nil
    }
    return data2
}
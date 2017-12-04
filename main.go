package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grafov/m3u8"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type MasterPlaylistWrapper struct {
	*m3u8.MasterPlaylist
	FullUrl string
	BaseUrl string
}

type MediaPlaylistWrapper struct {
	*m3u8.MediaPlaylist
	BaseUrl     string
	VariantInfo string
}

// FuncMap for template
var fmap = template.FuncMap{
	"byteToMb": byteToMb,
}

// Handler for getting m3u8
func handlerRawData(w http.ResponseWriter, r *http.Request) {
	m3u8Url := r.FormValue("m3u8Url")
	// ignore the certificate verification for https request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(m3u8Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", keepLines(string(body)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	m3u8Url := r.FormValue("m3u8Url")
	m3u8BaseUrl := m3u8Url[:strings.LastIndex(m3u8Url, "/")+1]

	// ignore the certificate verification for https request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(m3u8Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		panic(err)
	}

	switch listType {
	case m3u8.MEDIA:
		var mediaPlWrapper MediaPlaylistWrapper
		mediaPlWrapper.MediaPlaylist = p.(*m3u8.MediaPlaylist)
		mediaPlWrapper.BaseUrl = m3u8BaseUrl
		mediaPlWrapper.MediaPlaylist.Segments = mediaPlWrapper.MediaPlaylist.Segments[0 : mediaPlWrapper.MediaPlaylist.Count()-1]
		mediaPlWrapper.VariantInfo = r.FormValue("variantInfo")
		responseMediaPlaylist(w, &mediaPlWrapper)
	case m3u8.MASTER:
		var masterPlWrapper MasterPlaylistWrapper
		masterPlWrapper.MasterPlaylist = p.(*m3u8.MasterPlaylist)
		masterPlWrapper.FullUrl = m3u8Url
		masterPlWrapper.BaseUrl = m3u8BaseUrl
		responseMasterPlaylist(w, &masterPlWrapper)
	}
}

func responseMasterPlaylist(w http.ResponseWriter, masterPlaylistWrapper *MasterPlaylistWrapper) {
	fmap := template.FuncMap{
		"byteToMb": byteToMb,
	}

	t := template.Must(template.New("master-playlist.html").Funcs(fmap).ParseFiles("static/template/master-playlist.html"))
	err := t.Execute(w, masterPlaylistWrapper)
	if err != nil {
		panic(err)
	}
}

func responseMediaPlaylist(w http.ResponseWriter, mediaPlaylistWrapper *MediaPlaylistWrapper) {
	t := template.Must(template.New("media-playlist.html").Funcs(fmap).ParseFiles("static/template/media-playlist.html"))
	err := t.Execute(w, mediaPlaylistWrapper)
	if err != nil {
		panic(err)
	}
}

func byteToMb(bandwidth uint32) string {
	// bps단위는 1024가 아닌 1000 단위로 한다
	kbps := (bandwidth / 1000)
	if kbps <= 1000 {
		return fmt.Sprintf("%d Kbps(%d bps)", kbps, bandwidth)
	} else {
		return fmt.Sprintf("%0.1f Mbps(%d bps)", (float64(bandwidth) / (1000.0 * 1000.0)), bandwidth)
	}
}

func main() { //"https://devimages.apple.com.edgekey.net/streaming/examples/bipbop_4x3/bipbop_4x3_variant.m3u8"
	//"https://devstreaming-cdn.apple.com/videos/streaming/examples/img_bipbop_adv_example_ts/master.m3u8"
	//"https://devstreaming-cdn.apple.com/videos/streaming/examples/bipbop_16x9/bipbop_16x9_variant.m3u8"
	corsObj := handlers.AllowedOrigins([]string{"*"})
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	r.HandleFunc("/rawdata", handlerRawData).Methods("POST")
	r.Handle("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("./static/js")))
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir("./static/css")))
	http.ListenAndServe(":8080", handlers.CORS(corsObj)(r))

}

func keepLines(s string) string {
	return strings.Replace(s, "\n", "<br>", -1)
}

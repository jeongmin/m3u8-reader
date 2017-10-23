package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grafov/m3u8"
	"html/template"
	"net/http"
	"strings"
)

var m3u8BaseUrl string
var fmap = template.FuncMap{
	"byteToMb": byteToMb,
}

func handler(w http.ResponseWriter, r *http.Request) {
	m3u8Url := r.FormValue("m3u8Url")
	m3u8BaseUrl = m3u8Url[:strings.LastIndex(m3u8Url, "/")+1]

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
		mediapl := p.(*m3u8.MediaPlaylist)
		mediapl.BaseUrl = m3u8BaseUrl
		mediapl.Segments2 = mediapl.Segments[0 : mediapl.Count()-1]
		mediapl.VariantInfo = r.FormValue("variantInfo")
		responseMediaPlaylist(w, mediapl)
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		masterpl.BaseUrl = m3u8BaseUrl
		responseMasterPlaylist(w, masterpl)
	}
}

func responseMasterPlaylist(w http.ResponseWriter, masterPlaylist *m3u8.MasterPlaylist) {
	fmap := template.FuncMap{
		"byteToMb": byteToMb,
	}

	t := template.Must(template.New("master-playlist.html").Funcs(fmap).ParseFiles("master-playlist.html"))
	err := t.Execute(w, masterPlaylist)
	if err != nil {
		panic(err)
	}
}

func responseMediaPlaylist(w http.ResponseWriter, mediaPlaylist *m3u8.MediaPlaylist) {
	t := template.Must(template.New("media-playlist.html").Funcs(fmap).ParseFiles("media-playlist.html"))
	err := t.Execute(w, mediaPlaylist)
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

func main() {
	//"http://dev.p.naverrmc.edgesuite.net/global/read/wav_2017_03_14_1/657b3mqX75JTqHFiOEowbaejFA_rmcvideo_360P_640_1228_128_adoptive.m3u8"
	//"https://devimages.apple.com.edgekey.net/streaming/examples/bipbop_4x3/bipbop_4x3_variant.m3u8"
	//"https://devstreaming-cdn.apple.com/videos/streaming/examples/img_bipbop_adv_example_ts/master.m3u8"
	//"https://devstreaming-cdn.apple.com/videos/streaming/examples/bipbop_16x9/bipbop_16x9_variant.m3u8"
	corsObj := handlers.AllowedOrigins([]string{"*"})
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	http.ListenAndServe(":9000", handlers.CORS(corsObj)(r))
}

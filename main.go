package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-mjpeg"
)

var streamList []*mjpeg.Stream
var (
	cameraList = flag.String("cameraList", "localhost:8080/mjpeg,localhost:8080/mjpeg", "if more than one camera, use commas to delimit") //variable number of MJPEG cameras
	addr       = flag.String("addr", ":8888", "Server address")                                                                           //must start with sudo if it runs on low port number like 80
	interval   = flag.Duration("interval", 200*time.Millisecond, "interval")
	directory  = flag.String("d", "images", "relative path of static files to save images to")
	header     = flag.String("header", "Yolo-Coordinates", "optional header value read from the mjpeg stream to add as part of the filename")
)

func getMjpegStream(cameraUrl string) (*mjpeg.Decoder, error) {
	req, err := http.NewRequest("GET", cameraUrl, nil)
	if err != nil {
		log.Printf("error fetching cameraUrl %v error=%v\n", cameraUrl, err.Error())
		time.Sleep(time.Second)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	println(":")
	if err != nil {
		log.Printf("error fetching url:%v\n", err.Error())
		time.Sleep(time.Second)
		return nil, err
	}
	println("here")
	dec, err := mjpeg.NewDecoderFromResponse(res)
	if err != nil {
		log.Printf("Could not get new mjpeg Decoder from req err=%v\n", err)
	}
	return dec, err
}
func proxy(wg *sync.WaitGroup, stream *mjpeg.Stream, cameraLink string) {
	defer wg.Done()
	fmt.Printf("Camera-------------------------------------------------=%v\n", cameraLink)
	cameraUrl := "http://" + cameraLink
	dec, _ := getMjpegStream(cameraUrl)
	for {
		if dec == nil {
			time.Sleep(time.Second)
			log.Printf("nil caught for %v\n", cameraUrl)
			dec, _ = getMjpegStream(cameraUrl) //restart HTTP request
			continue
		}
		p, err := dec.Part()
		tag := p.Header.Get(*header)

		if err != nil {
			log.Printf("error decoding Part=%v err=%v\n", cameraUrl, err.Error())
			time.Sleep(time.Second)
			dec, _ = getMjpegStream(cameraUrl) //restart HTTP request
			continue
		}

		img, err := jpeg.Decode(p)
		if err != nil {
			log.Printf("error decoding image=%v err=%v\n", cameraUrl, err.Error())
			time.Sleep(time.Second)
			dec, _ = getMjpegStream(cameraUrl) //restart HTTP request
			continue
		}
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			log.Printf("error writing jpeg %v\n", err.Error())
		}
		cameraName := cameraUrl[7:] //skip http:// chars
		now := time.Now()
		weekday := now.Weekday()
		hour := now.Hour()
		hr, min, sec := now.Clock()
		path := *directory + "/all/" + weekday.String() + "/" + strconv.Itoa(hour) //hourly rotation for images
		filename := fmt.Sprintf("%v-%d%02d%02d%v.jpg", strings.Replace(cameraName, "/", "-", -1), hr, min, sec, tag)
		writeImage(buf.Bytes(), path, filename) //write all images
		stream.Update(buf.Bytes())
	}
}

func writeImage(b []byte, path string, filename string) {
	err := os.MkdirAll(path, 0777) //"Sunday"=0
	if err != nil {
		println("mkdir error=", err.Error())
		return
	}
	err = os.WriteFile(path+"/"+filename, b, 0644)
	if err != nil {
		println("file write error:", err.Error())
	}
}

func main() {
	flag.Parse()
	if *cameraList == "" {
		flag.Usage()
		os.Exit(1)
	}
	var wg sync.WaitGroup
	cameras := strings.Split(*cameraList, ",")
	println("cameras=", len(cameras))
	streamList = make([]*mjpeg.Stream, len(cameras))

	for i, camera := range cameras {
		streamList[i] = mjpeg.NewStreamWithInterval(*interval)
		wg.Add(1)
		go proxy(&wg, streamList[i], camera)
	}

	for i, stream := range streamList {
		http.HandleFunc("/"+strconv.Itoa(i), stream.ServeHTTP)

	}

	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(*directory))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>`))
		w.Write([]byte(`<body>`))
		for i := range streamList {
			line := "<img src=\"http://" + r.Host + "/" + strconv.Itoa(i) + "\" />"
			w.Write([]byte(line))
		}
		w.Write([]byte(`</body>`))
		w.Write([]byte(`</html>`))
	})

	server := &http.Server{Addr: *addr}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		server.Shutdown(context.Background())
	}()

	err := server.ListenAndServe()
	if err != nil {
		println("Cannot start server:", err.Error())
		os.Exit(1)
	}
	for _, stream := range streamList {
		stream.Close()
	}
	wg.Wait()
}

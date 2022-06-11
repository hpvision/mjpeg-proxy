# mjpeg-proxy

This program can read an mjpeg stream and write the images to disk in a heirarchy of folders by date.
Also proxies the images (images can be modified) into a new mjpeg stream onto a new port. 

Usage of ./mjpeg-proxy:
  -addr string
    	Server address (default ":8888")
  -cameraList string
    	if more than one camera, use commas to delimit (default "localhost:8080/mjpeg,localhost:8080/mjpeg")
  -d string
    	relative path of static files to save images to (default "images")
  -interval duration
    	interval (default 200ms)


To test from a macbook (run these commands from another window):

Install gocv:
https://gocv.io/getting-started/macos/

git clone https://github.com/dougwatson/go-mjpeg.git
cd go-mjpeg
go build

Then start the webcam:


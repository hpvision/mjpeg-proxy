# mjpeg-proxy

This program can read an mjpeg stream and write the images to disk in a heirarchy of folders by date.

Also proxies the images (images can be modified) into a new mjpeg stream onto a new port. 


# Usage of ./mjpeg-proxy:


  **-addr** string
    	Server address (default ":8888")
        
  **-cameraList** string
    	if more than one camera, use commas to delimit (default "localhost:8080/mjpeg,localhost:8080/mjpeg")

  **-d** string
    	relative path of static files to save images to (default "images")

  **-interval** duration
    	interval (default 200ms)


# To test from a macbook (run these commands from another window):


Install gocv:

https://gocv.io/getting-started/macos/


git clone https://github.com/dougwatson/go-mjpeg.git

cd go-mjpeg

go build


# In one terminal, start the webcam to write an MJPEG stream on port 8080:

cd go-mjpeg/_example/camera/

go build

./camera

# Then in another teminal, start the proxy:

(if you have not already downloaded and built this code)

git clone https://github.com/hpvision/mjpeg-proxy.git

cd mjpeg-proxy

go build


./mjpeg-proxy

It will read the mjpeg stream from port 8080 and proxy it to port 8888- also writing the images to disk in the data directory.

# To test:

open a web browser to:

http://localhost:8080 (to see the original webcam mjpeg feed)

and

http://localhost:8888 (to see the proxied feed on port 8888)

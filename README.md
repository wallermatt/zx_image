# zx_image
A command line tool written in Golang that takes a ZX Spectrum snapshot file (.sna) and creates a .png screenshot of what the computer would be displaying on the screen at the time.

## Usage
Once Go is installed:
```
go build zx_image.go
./zx_image <optional snapshot filename + path> <optional image name>
```
or
```
go run zx_image.go <optional snapshot filename + path> <optional image name>
```

e.g.
```
go run zx_image.go elite.sna elite_image.png
./zx_image /home/user/downloads/be.sna screenshot.png
```

Note: 
- If the second argument is omitted the image file name will be that of the snapshot file, e.g. be.sna -> be.png.
- If the first argument is omitted the snapshot in the testData folder will be used.


## Documentation

### ZX Spectrum Display
A very useful guide in three parts was written by David Black:
1. [Part I](http://www.overtakenbyevents.com/lets-talk-about-the-zx-specrum-screen-layout/)
2. [Part II](http://www.overtakenbyevents.com/lets-talk-about-the-zx-specrum-screen-layout-part-two/)
3. [Part III](http://www.overtakenbyevents.com/lets-talk-about-the-zx-specrum-screen-layout-part-three/)

### Snapshot File Structure
Essentially a dump of the ZX Spectrum's memory, prefixed by 27 bytes that contain the contents of the Z80's registers:
- [SNA Format](https://sinclair.wiki.zxnet.co.uk/wiki/SNA_format)

# Converting heic files to jpeg

This all started out by trying to rapidly browse on a Mac the image files in the `Photos` application and turning them into `jpeg` files for later use in a web page.

## Accessing the picture files

The image files are stored in various subdirectories of `"~/Pictures/Photos Library.photoslibrary/originals"`. It contains subdirectories named `1`, `2`, ... which contain `jpeg`` and `heic` files. According to [MacWorld: What is HEIC?](https://www.macworld.com/article/672609/what-is-heic.html)

> Apple has replaced the JPEG image format with the new HEIC alternative in iOS.
> ...
> HEIC is the file format name Apple has chosen for the new HEIF standard. HEIF stands for High Efficiency Image Format, and, as the name suggests, is a more streamlined way of storing image files.

Unfortunately, the standard [image](https://pkg.go.dev/image) libraries of Go cannot handle this format.


## Installation

Furtunately, there is a C library for HEIF images called [libheif].

> libheif is an ISO/IEC 23008-12:2017 HEIF and AVIF (AV1 Image File Format) file format decoder and encoder. There is partial support for ISO/IEC 23008-12:2022 (2nd Edition) capabilities.

It also has a go API, but requires that the C library is installed, see [libheif installation](https://github.com/strukturag/libheif#macos):

    brew install cmake make pkg-config x265 libde265 libjpeg libtool
    brew install libheif

I'm not sure whether the first line actually is required. 

Then continue to install the go packages as usual:

    go get "github.com/strukturag/libheif/go/heif"


## Package initialization
I was wandering why `image.Decode` worked only after I had installed `libheif/heif`. It is because the initialization registers the format `heif` with image (see [source code](https://github.com/strukturag/libheif/blob/master/go/heif/heif.go)) in the package `init()` function.

~~~go
func init() {
	image.RegisterFormat("heif", "????ftypheic", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheim", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheis", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheix", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevc", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevm", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevs", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypmif1", decodeImage, decodeConfig)
	image.RegisterFormat("avif", "????ftypavif", decodeImage, decodeConfig)
	image.RegisterFormat("avif", "????ftypavis", decodeImage, decodeConfig)
}
~~~

To read the heif file, use:

~~~go
import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/strukturag/libheif/go/heif"
)

ver := heif.GetVersion()
println(ver)

f, err := os.Open("sample.heic")
defer f.Close()
img, fmt, err := image.Decode(f)
println(fmt)
~~~

To convert it to jpeg:

~~~go
outf, _ := os.Create("sample.jpeg")
defer outf.Close()

var out bytes.Buffer
if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: 80}); err != nil {
    panic(err)
}
_, out := outf.Write(out.Bytes())
outf.Close()
~~~

## Exif

Another question is to get as much information about the image as possible. This information is stored in the image file in the **Exchangeable image file format** Exif. According to [wikipedia][exif], 

> **Exchangeable image file format** (officially **Exif**, according to JEIDA/JEITA/CIPA specifications)[5] is a standard that specifies formats for images, sound, and ancillary tags used by digital cameras (including smartphones), scanners and other systems handling image and sound files recorded by digital cameras. 

I found the Go library [imagemeta] to read this data.

~~~go
import (
    "github.com/evanoberholster/imagemeta"
)
exif, err := imagemeta.Decode(f)
println(exif.String())

/*
Exif
ImageType: 	image/heif
Make: 		Apple
Model: 		iPhone 13
LensMake: 	Apple
LensModel: 	iPhone 13 back dual wide camera 5.1mm f/1.6
CameraSerial: 	
LensSerial: 	
Image Size: 	4032x3024
Orientation: 	Horizontal
ShutterSpeed: 	1/577
Aperture: 	1.60
ISO: 		50
Flash: 		Off, Did not fire
Focal Length: 	5.10mm
Fl 35mm Eqv: 	26.00mm
Exposure Prgm: 	Program AE
Metering Mode: 	Multi-segment
Exposure Mode: 	Auto
Date Modified: 	2023-10-18 15:39:18 -0500 -05:00
Date Created: 	2023-10-18 15:39:18.681 -0500 -05:00
Date Original: 	2023-10-18 15:39:18.681 -0500 -05:00
Date GPS: 	2023-10-18 00:00:00 +0000 UTC
Artist: 	
Copyright: 	
Software: 	16.6.1
Image Desc: 	
GPS Altitude: 	185.04
GPS Latitude: 	41.947917
GPS Longitude: 	-87.656342
*/
~~~

## References

- **libheif**: C library with Go API for heif images [libheif]
- **go-libheif**: Go Wrapper for [libheif] at [go-libheif](https://github.com/MaestroError/go-libheif/tree/maestro)
- **imagemetag**: Go package to handle image EXIF metadata at [imagemeta]

[libheif]: https://github.com/strukturag/libheif
[imagemeta]: https://github.com/evanoberholster/imagemeta
[exif]: https://en.wikipedia.org/wiki/Exif


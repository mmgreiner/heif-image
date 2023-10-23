/*
small program to test libheif
see https://github.com/strukturag/libheif/blob/master/go/heif/heif.go

See also https://github.com/MaestroError/go-libheif/blob/maestro/libheif.go#L102
*/
package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/strukturag/libheif/go/heif"

	"github.com/evanoberholster/imagemeta"
)

const (
	SampleFile  = "/Users/mmgreiner/Pictures/Photos Library.photoslibrary/originals/F/F6480269-9529-4327-9A08-6A01FC4C518F.heic"
	SampleFile1 = "sample.heic"
)

func main() {

	ver := heif.GetVersion()
	println(ver)

	f, err := os.Open(SampleFile1)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// look at the metadata
	exif, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	println(exif.String())

	// reset the reader
	f.Seek(0, 0)

	// decode the image
	img, fmt, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	println(fmt)

	// now save as jpeg
	outfn := "sample.jpeg"
	outf, _ := os.Create(outfn)
	defer outf.Close()

	var out bytes.Buffer
	if err := jpeg.Encode(&out, img, nil); err != nil { // &jpeg.Options{Quality: 80}
		panic(err)
	}
	println("writing to", outfn)
	if _, err := outf.Write(out.Bytes()); err != nil {
		panic(err)
	}

}

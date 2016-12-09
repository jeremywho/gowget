package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var bytesToMegaBytes = 1048576.0

// PassThru code originally from
// http://stackoverflow.com/a/22422650/613575
type PassThru struct {
	io.Reader
	curr  int64
	total float64
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.curr += int64(n)

	// last read will have EOF err
	if err == nil || (err == io.EOF && n > 0) {
		printProgress(float64(pt.curr), pt.total)
	}

	return n, err
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s url outfile\n", os.Args[0])
		os.Exit(1)
	}

	//"http://ipv4.download.thinkbroadband.com/5MB.zip"
	resp, _ := http.Get(os.Args[1])
	defer resp.Body.Close()

	out, _ := os.Create(os.Args[2])
	defer out.Close()

	src := &PassThru{Reader: resp.Body, total: float64(resp.ContentLength)}

	size, err := io.Copy(out, src)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nFile Transferred. (%.1f MB)\n", float64(size)/bytesToMegaBytes)
}

func printProgress(curr, total float64) {
	width := 40.0
	output := ""
	threshold := (curr / total) * float64(width)
	for i := 0.0; i < width; i++ {
		if i < threshold {
			output += "="
		} else {
			output += " "
		}
	}

	fmt.Printf("\r[%s] %.1f of %.1fMB", output, curr/bytesToMegaBytes, total/bytesToMegaBytes)
}

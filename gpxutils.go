//
//  gpxutils.go
//  gpxutis
//
//  Created by Mario Martelli on 11.10.2019.
//  Copyright © 2019 Mario Martelli. All rights reserved.
//
//  This file is part of gpxutils.
//
//  Everorg is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  EverOrg is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with EverOrg.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var (
	readFile  = ""
	ls, split bool
	outDir    = ""
)

func main() {

	// command line parsing
	inPtr := flag.String("in", "none", "relative path to GPX file")
	outPtr := flag.String("out", "out", "relative path to output directory")
	flag.BoolVar(&split, "split", false, "write single files for tracks")
	flag.BoolVar(&ls, "ls", false, "list tracks ")
	flag.Parse()
	if inPtr == nil || *inPtr == "" {
		panic("input file is missing")
	}

	if outPtr == nil || *outPtr == "" {
		outDir = "out"
	} else {
		outDir = *outPtr
	}

	// Open the file given at commandline
	readFile = *inPtr
	xmlFile, err := os.Open(readFile)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer func() { _ = xmlFile.Close() }()

	b, _ := ioutil.ReadAll(xmlFile)

	// Progress spinner
	s := spinner.New(spinner.CharSets[38], 200*time.Millisecond) // Build our new spinner
	s.Prefix = "Processing input-file: "
	s.Writer = os.Stderr

	// Unmarshall gpx file
	s.Start() // Start the spinner
	var q Query
	_ = xml.Unmarshal(b, &q)
	s.Stop()

	// GPX Header
	header := `<?xml version="1.0" encoding="UTF-8" ?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1"
	    version="1.1"
	    creator="rubiTrack - https://www.rubitrack.com"
	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	    xmlns:gpxdata="http://www.cluetrust.com/XML/GPXDATA/1/0"
	    xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.cluetrust.com/XML/GPXDATA/1/0 http://www.cluetrust.com/Schemas/gpxdata10.xsd">`

	//	Parse the contained xml
	if split {

		if _, err = os.Stat(outDir); os.IsNotExist(err) {
			_ = os.Mkdir(outDir, 0711)
		}
		s.Prefix = "Creating single files: "
		s.Start() // Start the spinner

		for _, track := range q.Tracks {
			outfile := getTimestamp(track)
			b, _ := xml.MarshalIndent(track, "  ", "  ")
			f, err := os.Create(outDir + "/" + outfile + ".gpx")
			if err != nil {
				fmt.Println(err)
				continue
			}

			f.WriteString(header)
			f.WriteString(string(b))
			f.WriteString("\n</gpx>\n")
			_ = f.Close()
		}
		s.Stop()
		return
	}

	if ls {
		for _, track := range q.Tracks {
			fmt.Println(getTimestamp(track))
		}

		fmt.Println("Total amout of tracks in file: ", len(q.Tracks))
	}
}

// Returns first timestamp found in track
func getTimestamp(track Track) string {

	for _, trksegs := range track.Trksegs {
		for _, trkpt := range trksegs.Trkpts {
			if len(trkpt.Time) > 0 {
				return trkpt.Time
			}
		}
	}
	return "unknown"
}

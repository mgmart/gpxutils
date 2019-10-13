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
)

func main() {

	// command line parsing
	wordPtr := flag.String("input", "GPX File", "relative path to GPX file")

	flag.BoolVar(&split, "split", false, "write single files for tracks")
	flag.BoolVar(&ls, "ls", false, "list tracks ")
	flag.Parse()
	if wordPtr == nil || *wordPtr == "" {
		panic("input file is missing")
	}

	// Open the file given at commandline
	readFile = *wordPtr
	xmlFile, err := os.Open(readFile)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer func() { _ = xmlFile.Close() }()

	b, _ := ioutil.ReadAll(xmlFile)

	s := spinner.New(spinner.CharSets[38], 200*time.Millisecond) // Build our new spinner
	s.Prefix = "Processing input-file: "
	s.Writer = os.Stderr
	s.Start() // Start the spinner

	var q Query
	_ = xml.Unmarshal(b, &q)

	s.Stop()
	header := `<?xml version="1.0" encoding="UTF-8" ?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1"
	    version="1.1"
	    creator="rubiTrack - https://www.rubitrack.com"
	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	    xmlns:gpxdata="http://www.cluetrust.com/XML/GPXDATA/1/0"
	    xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.cluetrust.com/XML/GPXDATA/1/0 http://www.cluetrust.com/Schemas/gpxdata10.xsd">`
	//	Parse the contained xml

	if split {

		if _, err = os.Stat("out"); os.IsNotExist(err) {
			_ = os.Mkdir("out", 0711)
		}
		s.Prefix = "Creating single files: "
		s.Start() // Start the spinner

		for _, track := range q.Tracks {
			outfile := ""
			for _, trksegs := range track.Trksegs {
				for _, trkpt := range trksegs.Trkpts {
					if len(trkpt.Time) > 0 {
						outfile = trkpt.Time
						break
					}
				}
			}

			b, _ := xml.MarshalIndent(track, "  ", "  ")
			f, err := os.Create("out/" + outfile + ".gpx")
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

	fmt.Println("Amount of tracks: ", len(q.Tracks))

	if ls {

		for i := 0; i < len(q.Tracks); i++ {
			fmt.Println("Track description: " + q.Tracks[i].Desc)
		}
	}
}

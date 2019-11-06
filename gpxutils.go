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
	"path"

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
	s := spinner.New(spinner.CharSets[38], 200*time.Millisecond)
	s.Prefix = "Processing input-file: "
	s.Writer = os.Stderr

	// Unmarshall gpx file
	s.Start() // Start the spinner
	var q Query
	_ = xml.Unmarshal(b, &q)
	s.Stop() // Stop Spinner

	//	Write each track to a single file
	if split {

		// Create out directory
		if _, err = os.Stat(outDir); os.IsNotExist(err) {
			_ = os.Mkdir(outDir, 0711)
		}
		s.Prefix = "Creating single files: "
		s.Start() // Start the spinner
		outfile := ""
		for i, track := range q.Tracks {
			if track.getTimestamp() != "unknown" {
				outfile = track.getTimestamp()
			} else {
				basename := path.Base(*inPtr)
				extension := path.Ext(basename)
				name := basename[0:len(basename) - len(extension)]
				outfile = name + fmt.Sprintf("-%02d", i)
			}
			f, err := os.Create(outDir + "/" + outfile + ".gpx")
			if err != nil {
				fmt.Println(err)
				continue
			}
			f.WriteString(track.gpx())
			_ = f.Close()
		}
		s.Stop()
		return
	}

	// List all tracks contained in file
	if ls {
		for _, line := range q.listOfTracks() {
			fmt.Println(line)
		}
		fmt.Println("Total amout of tracks in file: ", q.numberOfTracks())
	}
}

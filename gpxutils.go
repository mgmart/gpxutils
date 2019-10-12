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

	"golang.org/x/net/html"
)

// Globals for filehandling
var (
	readFile       = ""
	attFolderExt   = "-attachments"
	attachmentPath = ""

	isMerged bool
)

// Get Attributes for html tag
func getAttr(attribute string, token html.Token) string {
	for _, attr := range token.Attr {
		if attr.Key == attribute {
			return attr.Val
		}
	}
	return ""
}

// Org mode representation of Node

func main() {
	wordPtr := flag.String("input", "gpx File", "relative path to enex file")

	flag.BoolVar(&isMerged, "split", false, "whether to write single files for tracks")
	flag.Parse()
	if wordPtr == nil || *wordPtr == "" {
		panic("input file is missing")
	}
	fmt.Println("input:", *wordPtr)

	// Open the file given at commandline
	readFile = *wordPtr
	xmlFile, err := os.Open(readFile)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer func() { _ = xmlFile.Close() }()

	fmt.Println("Reading file")
	b, _ := ioutil.ReadAll(xmlFile)
	fmt.Println("Finished reading file")

	fmt.Println("Marshalling file")
	var q Query
	_ = xml.Unmarshal(b, &q)
	fmt.Println("Finished marshalling file")

	// Parse the contained xml
	if isMerged {
		return
	}

	fmt.Println("Amount of tracks: ", len(q.Tracks))

	for i := 0; i < len(q.Tracks); i++ {
		fmt.Println("Track description: " + q.Tracks[i].Desc)

	}
}

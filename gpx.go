//
//  gpx.go
//  GPXutils
//
//  Created by Mario Martelli on 11.10.2019
//  Copyright © 2019 Mario Martelli. All rights reserved.
//
//  This file is part of GPXutils
//
//  EverOrg is free software: you can redistribute it and/or modify
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
	"fmt"
	"strings"
	"time"
	//	"time"
)

//
// Data Structures
//

// Query ...
type Query struct {
	Tracks []Track `xml:"trk"`
}

// Note ...
type Track struct {
	XMLName xml.Name `xml:"trk"`
	Name    string   `xml:"name"`
	Desc    string   `xml:"desc"`
	Trksegs []Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpts []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Lat  string          `xml:"lat,attr"`
	Lon  string          `xml:"lon,attr"`
	Ele  string          `xml:"ele"`
	Time string          `xml:"time"`
	Exts *ExtensionsType `xml:"extensions"`
}

type ExtensionsType struct {
	XML []byte `xml:",innerxml"`
}

// Which tracks are in file?
func (q Query) listOfTracks() []string {
	var ret []string
	for _, track := range q.Tracks {
		ret = append(ret, track.getTimestamp())
	}
	return ret
}

// How many tracks are in file?
func (q Query) numberOfTracks() int {
	return len(q.Tracks)
}

// Returns first timestamp found in track
func (track Track) getTimestamp() string {

	for _, trksegs := range track.Trksegs {
		for _, trkpt := range trksegs.Trkpts {
			if len(trkpt.Time) > 0 {
				return strings.Replace(trkpt.Time, ":", "-", -1)
			}
		}
	}
	return "unknown"
}

// GPX representation of track
func (t Track) gpx() string {

	// GPX Header
	header := `<?xml version="1.0" encoding="UTF-8" ?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1"
	    version="1.1"
	    creator="rubiTrack - https://www.rubitrack.com"
	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	    xmlns:gpxdata="http://www.cluetrust.com/XML/GPXDATA/1/0"
	    xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.cluetrust.com/XML/GPXDATA/1/0 http://www.cluetrust.com/Schemas/gpxdata10.xsd">`
	footer := "\n</gpx>\n"

	b, _ := xml.MarshalIndent(t, "  ", "  ")
	return header + string(b) + footer
}

// Returns something
func (track *Track) setTimeStamps() string{
	track2 := Track{}
	track2.Name = track.Name
	track2.Desc = track.Desc
	t := time.Date(2022, 10, 2, 7, 13, 10, 10, time.UTC)
	for _, trksegs := range track.Trksegs {
		trksegs2 := Trkseg{}
		for _, trkpt := range trksegs.Trkpts {
			trkpt2 := Trkpt{}
			trkpt2.Lat = trkpt.Lat
			trkpt2.Lon = trkpt.Lon
			trkpt2.Ele = trkpt.Ele
			trkpt2.Time = t.Format(time.RFC3339)
			t = t.Add(time.Second * 52)
			trksegs2.Trkpts = append(trksegs2.Trkpts, trkpt2)
		}
		track2.Trksegs = append(track2.Trksegs, trksegs2)

	}
	// GPX Header
	header := `<?xml version="1.0" encoding="UTF-8" ?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1"
	    version="1.1"
	    creator="rubiTrack - https://www.rubitrack.com"
	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	    xmlns:gpxdata="http://www.cluetrust.com/XML/GPXDATA/1/0"
	    xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.cluetrust.com/XML/GPXDATA/1/0 http://www.cluetrust.com/Schemas/gpxdata10.xsd">`
	footer := "\n</gpx>\n"

	b, _ := xml.MarshalIndent(track2, "  ", "  ")
	fmt.Println(string(b))
	return header + string(b) + footer
	//	return header + footer
}

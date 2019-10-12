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

import "encoding/xml"

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

// type GpxExt struct {
// 	Hr  int `xml:"gpxdata:hr"`
// 	Cad int `xml:"gpxdata:cadence"`
// }

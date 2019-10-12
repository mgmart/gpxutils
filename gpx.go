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

//
// Data Structures
//

// Query ...
type Query struct {
	Tracks []Track `xml:"trk"`
}

// Note ...
type Track struct {
	Name   string `xml:"name"`
	Desc   string `xml:"desc"`
	Trkseg struct {
		Trkpt struct {
			Lat  string `xml:"lat"`
			Lon  string `xml:"lon"`
			Ele  string `xml:"ele"`
			Time string `xml:"time"`
			Ext  struct {
				Hr  string `xml:"gpxdata:hr"`
				Cad string `xml:"gpxdata:cadence"`
			} `xml:"extensions"`
		} `xml:"trkpt"`
	} `xml:"trkseg"`
}

// Resource ...
// type Resource struct {
// 	Mime string `xml:"mime"`
// 	Data struct {
// 		Content  string `xml:",chardata"`
// 		Encoding string `xml:"encoding,attr"`
// 	} `xml:"data"`
// }

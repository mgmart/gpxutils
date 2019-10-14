//
//  gpx_test.go
//  GPXutils
//
//  Created by Mario Martelli on 13.10.2019
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
	"testing"
)

const xmlTrack string = `
<trk>
    <name>Bonn ---- Radfahren</name>
    <desc>tescht</desc>
    <trkseg>
      <trkpt lat="50.83326700" lon="7.13071600">
        <ele>85.7</ele>
        <time>2011-04-09T06:03:40Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>1</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83318700" lon="7.13064500">
        <ele>86.2</ele>
        <time>2011-04-09T06:03:46Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>9</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83307600" lon="7.13060100">
        <ele>80.9</ele>
        <time>2011-04-09T06:03:51Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>0</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83304600" lon="7.13059300">
        <ele>85.7</ele>
        <time>2011-04-09T06:03:52Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>0</gpxdata:cadence>
                </extensions>
      </trkpt>
    </trkseg>
  </trk>
`

func TestGetTimstamp(t *testing.T) {

	var track Track
	_ = xml.Unmarshal([]byte(xmlTrack), &track)

	if track.getTimestamp() != "2011-04-09T06-03-40Z" {
		t.Error("Expected 2011-04-09T06-03-40Z and got ", track.getTimestamp())
	}
}

func TestGpx(t *testing.T) {

	var track Track
	_ = xml.Unmarshal([]byte(xmlTrack), &track)

	if track.gpx() != `<?xml version="1.0" encoding="UTF-8" ?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1"
	    version="1.1"
	    creator="rubiTrack - https://www.rubitrack.com"
	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	    xmlns:gpxdata="http://www.cluetrust.com/XML/GPXDATA/1/0"
	    xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.cluetrust.com/XML/GPXDATA/1/0 http://www.cluetrust.com/Schemas/gpxdata10.xsd">  <trk>
    <name>Bonn ---- Radfahren</name>
    <desc>tescht</desc>
    <trkseg>
      <trkpt lat="50.83326700" lon="7.13071600">
        <ele>85.7</ele>
        <time>2011-04-09T06:03:40Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>1</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83318700" lon="7.13064500">
        <ele>86.2</ele>
        <time>2011-04-09T06:03:46Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>9</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83307600" lon="7.13060100">
        <ele>80.9</ele>
        <time>2011-04-09T06:03:51Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>0</gpxdata:cadence>
                </extensions>
      </trkpt>
      <trkpt lat="50.83304600" lon="7.13059300">
        <ele>85.7</ele>
        <time>2011-04-09T06:03:52Z</time>
        <extensions>
                    <gpxdata:hr>0</gpxdata:hr>
                    <gpxdata:cadence>0</gpxdata:cadence>
                </extensions>
      </trkpt>
    </trkseg>
  </trk>
</gpx>
` {
		t.Error("track.gpx gave faulty return")
	}
}

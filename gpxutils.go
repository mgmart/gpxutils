//
//  gpxutils.go
//  gpxutis
//
//  Created by Mario Martelli on 11.10.2019.
//  Copyright © 2019 Mario Martelli. All rights reserved.
//
//  This file is part of gpxutils.
//
//  gpxutils is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  gpxutils is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with gpxutils. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"

	"os"
	"path"
	"time"

	"github.com/briandowns/spinner"
)

const version = "1.1.0"

var (
	readFile      = ""
	ls, split, tm bool
	showVersion   bool
	outDir        = ""
	timezone      = "+02:00"
	begin, end    time.Time
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "gpxutils v%s - GPS Track Utilities\n\n", version)
	fmt.Fprintf(os.Stderr, "USAGE:\n")
	fmt.Fprintf(os.Stderr, "  gpxutils [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "DESCRIPTION:\n")
	fmt.Fprintf(os.Stderr, "  Extract tracks from multi-track GPX files, list tracks, or add timestamps.\n\n")
	fmt.Fprintf(os.Stderr, "OPTIONS:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
	fmt.Fprintf(os.Stderr, "  List all tracks in a GPX file:\n")
	fmt.Fprintf(os.Stderr, "    gpxutils -in track.gpx -ls\n\n")
	fmt.Fprintf(os.Stderr, "  Extract all tracks to individual files:\n")
	fmt.Fprintf(os.Stderr, "    gpxutils -in track.gpx -out ./tracks -split\n\n")
	fmt.Fprintf(os.Stderr, "  Add timestamps to waypoints:\n")
	fmt.Fprintf(os.Stderr, "    gpxutils -in track.gpx -time -begin 09:00 -end 17:00\n\n")
}

func validateInputFile(filepath string) error {
	if filepath == "" || filepath == "none" {
		return fmt.Errorf("input file is required (use -in flag)")
	}
	
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("input file '%s' does not exist", filepath)
	}
	
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("cannot open input file '%s': %v", filepath, err)
	}
	file.Close()
	
	return nil
}

func main() {
	// Configure flag usage function
	flag.Usage = printUsage

	// Define command line flags with proper descriptions
	inPtr := flag.String("in", "", "Path to input GPX file (required)")
	outPtr := flag.String("out", "out", "Output directory for generated files")
	
	flag.BoolVar(&split, "split", false, "Extract each track to a separate GPX file")
	flag.BoolVar(&ls, "ls", false, "List all tracks found in the GPX file")
	flag.BoolVar(&tm, "time", false, "Add timestamps to waypoints (requires -begin and -end)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.StringVar(&timezone, "tz", "+02:00", "Timezone offset for timestamp operations (e.g. +02:00, -05:00)")

	flag.Func("begin", "Start time for timestamp operation (format: HH:MM)", func(s string) error {
		var err error
		ds := time.Now().Format("2006-01-02") + "T" + s + ":00" + timezone
		begin, err = time.Parse(time.RFC3339, ds)
		if err != nil {
			return fmt.Errorf("invalid start time format '%s': expected HH:MM (e.g. 09:30)", s)
		}
		return nil
	})

	flag.Func("end", "End time for timestamp operation (format: HH:MM)", func(s string) error {
		var err error
		ds := time.Now().Format("2006-01-02") + "T" + s + ":00" + timezone
		end, err = time.Parse(time.RFC3339, ds)
		if err != nil {
			return fmt.Errorf("invalid end time format '%s': expected HH:MM (e.g. 17:30)", s)
		}
		return nil
	})

	// Parse command line arguments
	flag.Parse()

	// Handle version flag
	if showVersion {
		fmt.Printf("gpxutils version %s\n", version)
		return
	}

	// Validate input file
	readFile = *inPtr
	if err := validateInputFile(readFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Set output directory
	outDir = *outPtr

	// Validate that at least one operation is specified
	if !ls && !split && !tm {
		fmt.Fprintf(os.Stderr, "Error: No operation specified. Use -ls, -split, or -time\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Open the input file
	xmlFile, err := os.Open(readFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open input file '%s': %v\n", readFile, err)
		os.Exit(1)
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

	// Write each track to a single file
	if split {
		// Create output directory
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			if err := os.Mkdir(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to create output directory '%s': %v\n", outDir, err)
				os.Exit(1)
			}
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
				name := basename[0 : len(basename)-len(extension)]
				outfile = name + fmt.Sprintf("-%02d", i)
			}
			outPath := outDir + "/" + outfile + ".gpx"
			f, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to create file '%s': %v\n", outPath, err)
				continue
			}
			if _, err := f.WriteString(track.gpx()); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to write to file '%s': %v\n", outPath, err)
				f.Close()
				continue
			}
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to close file '%s': %v\n", outPath, err)
			}
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

	if tm {
		// Validate time parameters for timestamp operation
		if begin.IsZero() {
			fmt.Fprintf(os.Stderr, "Error: Start time is required for timestamp operation (use -begin HH:MM)\n")
			os.Exit(1)
		}

		if end.IsZero() {
			fmt.Fprintf(os.Stderr, "Error: End time is required for timestamp operation (use -end HH:MM)\n")
			os.Exit(1)
		}

		if end.Before(begin) || end.Equal(begin) {
			fmt.Fprintf(os.Stderr, "Error: End time (%s) must be after start time (%s)\n", 
				end.Format("15:04"), begin.Format("15:04"))
			os.Exit(1)
		}

		// Create output directory if it doesn't exist
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			if err := os.Mkdir(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to create output directory '%s': %v\n", outDir, err)
				os.Exit(1)
			}
		}

		outfile := ""

		for i := range q.Tracks {

			basename := path.Base(*inPtr)
			extension := path.Ext(basename)
			name := basename[0 : len(basename)-len(extension)]
			outfile = name + fmt.Sprintf("-%02d", i)

			outPath := outDir + "/" + outfile + ".gpx"
			f, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to create file '%s': %v\n", outPath, err)
				continue
			}
			if _, err := f.WriteString(q.Tracks[0].setTimeStamps(begin, end)); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to write to file '%s': %v\n", outPath, err)
				f.Close()
				continue
			}
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to close file '%s': %v\n", outPath, err)
			}
		}
		return
	}

}

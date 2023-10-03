package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/packetflinger/libq2/bsp"
)

var (
	checkMissing = flag.Bool("check_missing", false, "Find missing textures")
	sourceDir    = flag.String("source", "", "Root director of our textures")
	source       = &os.File{}
	sourceFiles  = []fs.DirEntry{}
)

// Remove any duplipcates
func Deduplicate(in []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range in {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [flags] <q2map.bsp> [q2map.bsp...]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	flag.Parse()

	if ok, err := argsOkay(); !ok {
		fmt.Println(err)
		return
	}

	allTextures := []string{}

	for _, bspname := range flag.Args() {
		bspfile, err := bsp.OpenBSPFile(bspname)
		if err != nil {
			fmt.Println(err)
			return
		}

		textures := bspfile.FetchTextures()
		for _, t := range textures {
			// all names are still padded with nulls, strip those away
			allTextures = append(allTextures, strings.Trim(t.File, "\x00"))
		}
		bspfile.Close()
	}

	dedupedtextures := Deduplicate(allTextures)
	sort.Strings(dedupedtextures)

	for _, t := range dedupedtextures {
		fmt.Println(t)
	}

	if len(sourceFiles) > 0 {
		for _, sf := range sourceFiles {
			fmt.Println(sf.Name())
		}
	}
}

// Make sure all the args we need are specified. Some are dependent on others.
// Returns true if args are okay
func argsOkay() (bool, error) {
	if *checkMissing && len(*sourceDir) == 0 {
		return false, errors.New("source flag required in check_missing mode")
	}

	// if --source flag provided, check folder actually exists
	if len(*sourceDir) > 0 {
		src, err := os.Open(*sourceDir)
		if err != nil {
			return false, err
		}

		sourceFiles, err = src.ReadDir(-1)
		if err != nil {
			return false, err
		}
		if len(sourceFiles) == 0 {
			return false, errors.New("source directory contains no files")
		}
		source = src
	}

	return true, nil
}

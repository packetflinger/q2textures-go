package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/packetflinger/libq2/bsp"
)

var (
	list         = flag.Bool("list", true, "Output textures found in maps")
	checkMissing = flag.Bool("check_missing", false, "Find missing textures")
	sourceDir    = flag.String("source", "", "Root director of our textures")
	sourceFiles  = []string{}
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

	if *list {
		for _, t := range dedupedtextures {
			fmt.Println(t)
		}
	}

	if *checkMissing {
		missing := findMissing(dedupedtextures, sourceFiles)
		for _, t := range missing {
			fmt.Println(t)
		}
	}
}

// Make sure all the args we need are specified. Some are dependent on others.
// Returns true if args are okay
func argsOkay() (bool, error) {
	if *checkMissing && len(*sourceDir) == 0 {
		return false, errors.New("source flag required in check_missing mode")
	}

	// only do one
	if *checkMissing && *list {
		*list = false
	}

	// if --source flag provided, check folder actually exists
	if len(*sourceDir) > 0 {
		src, err := os.Open(*sourceDir)
		if err != nil {
			return false, err
		}
		defer src.Close()

		err = filepath.Walk(*sourceDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				sourceFiles = append(sourceFiles, path)

				return nil
			},
		)
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		if len(sourceFiles) == 0 {
			return false, errors.New("source directory contains no files")
		}
	}

	return true, nil
}

// Find any textures that exist in the map file that are not in the source
// directory.
//
// Keep in mind, textures listed in the map don't have extensions, so
// we need to match them as a prefix.
func findMissing(bspTextures []string, diskTextures []string) []string {
	found := make(map[string]bool)
	for _, t := range bspTextures {
		found[t] = false
	}

	for _, bt := range bspTextures {
		for _, dt := range diskTextures {
			if strings.Contains(dt, bt+".") {
				found[bt] = true
			}
		}
	}

	var missing []string
	for k, v := range found {
		if !v {
			missing = append(missing, k)
		}
	}
	return missing
}

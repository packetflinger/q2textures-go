package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/packetflinger/libq2/bsp"
)

var (
	checkMissing = flag.Bool("check_missing", false, "Find missing textures")
	sourceDir    = flag.String("source", "", "Root director of our textures")
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

	for _, bspname := range os.Args[1:] {
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
}

// Make sure all the args we need are specified. Some are dependent on others.
// Returns true if args are okay
func argsOkay() (bool, error) {
	if *checkMissing && len(*sourceDir) == 0 {
		return false, errors.New("source flag required in check_missing mode")
	}
	return true, nil
}

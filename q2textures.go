package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/packetflinger/libq2/bsp"
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
		fmt.Printf("Usage: %s <q2map.bsp> [q2map.bsp...]\n", os.Args[0])
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

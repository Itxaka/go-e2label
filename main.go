package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/itxaka/go-e2label/superblock"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		fmt.Println("Usage: e2label <filename> [New label]")
		os.Exit(1)
	}
	device := args[0]
	_, err := os.Stat(device)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Device does not exist")
		os.Exit(1)
	}
	sb, err := superblock.GetSuperBlock(device)
	if err != nil {
		return
	}
	if len(args) == 1 {
		fmt.Println(string(sb.VolumeName[:]))
	}
	if len(args) == 2 {
		oldname := sb.VolumeName[:]
		cleanedUp := strings.Trim(string(oldname), "\x00")
		if cleanedUp == strings.TrimSpace(args[1]) {
			fmt.Println("Old and new label are the same, not doing anything")
			os.Exit(0)
		}
		newName := [16]byte{}
		copy(newName[:], args[1])
		sb.VolumeName = newName
		err = sb.CalculateNewChecksumAndWriteIt(device)
		if err != nil {
			fmt.Println("Error calculating new checksum")
			os.Exit(1)
		}
		fmt.Printf("Label changed from %s to %s\n", cleanedUp, newName[:])
	}
}

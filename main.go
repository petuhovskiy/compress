package main

import (
	"fmt"
	"github.com/petuhovskiy/compress/tool"
	"io/ioutil"
	"os"
	"time"
)

func CmdHelp() {
	algos := ""

	for _, algo := range tool.Supported {
		algos += "\n"
		algos += fmt.Sprintf("- %s\t => %s", algo.ID, algo.Description)
	}

	algos += "\n" + fmt.Sprintf("- best\t => Encode with every supported algorithm and choose best encoding. Very very slow.")

	help := `Compress is a tool for lossless compression and decompression.

Usage:
	./compress <command> [arguments]

File compression:

	Usage:
		./compress c <in> <out> <algo>

	Example:
		./compress c in_file.txt out_file.cmp ppm

	Supported algos:` +
		"\n" + algos + "\n\n" +
		`File decompression:
	
	Usage:
		./compress d <archive> <out>

	Example:
		./compress d out_file.cmp file.txt
`

	fmt.Print(help)
}

func CmdCompress(in, out, algoID string) {
	src, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Printf("Failed to read in file, err=%v\n", err)
		return
	}

	startedAt := time.Now()

	var compressed []byte

	if algoID == "best" {
		compressed, err = tool.CompressBest(src)
	} else {
		compressed, err = tool.Compress(src, algoID)
	}

	if err != nil {
		fmt.Printf("Failed to compress, err=%v\nRun ./compress help, to get help.", err)
		return
	}

	finishedAt := time.Now()

	fmt.Printf("Finished in %v.\n", finishedAt.Sub(startedAt))

	beforeSize := len(src)
	afterSize := len(compressed)

	if beforeSize > 0 {
		factor := float64(afterSize) / float64(beforeSize) * 100
		fmt.Printf("Deflated %.2f%%\n", 100-factor)
	}

	err = ioutil.WriteFile(out, compressed, 0644)
	if err != nil {
		fmt.Printf("Failed to write file, err=%v\n", err)
		return
	}
}

func CmdDecompress(archive, out string) {
	data, err := ioutil.ReadFile(archive)
	if err != nil {
		fmt.Printf("Failed to read archive file, err=%v\n", err)
		return
	}

	decompressed, err := tool.Decompress(data)
	if err != nil {
		fmt.Printf("Failed to compress, err=%v\n", err)
		return
	}

	err = ioutil.WriteFile(out, decompressed, 0644)
	if err != nil {
		fmt.Printf("Failed to write file, err=%v\n", err)
		return
	}
}

func main() {
	args := os.Args[1:]

	switch {
	case len(args) < 1:
		CmdHelp()

	case args[0] == "c" && len(args) == 4:
		CmdCompress(args[1], args[2], args[3])

	case args[0] == "d" && len(args) == 3:
		CmdDecompress(args[1], args[2])

	default:
		CmdHelp()
	}
}

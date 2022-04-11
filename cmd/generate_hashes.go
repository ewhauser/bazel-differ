package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var SeedFilepaths string
var DisplayElapsedTime bool

// generateHashesCmd represents the generateHashes command
var generateHashesCmd = &cobra.Command{
	Use:   "generate-hashes",
	Short: "Writes to a file the SHA256 hashes for each Bazel Target in the provided workspace.",
	Long:  `Writes to a file the SHA256 hashes for each Bazel Target in the provided workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetHasher := internal.NewTargetHashingClient(GetBazelClient(), internal.Filesystem,
			internal.NewRuleProvider())

		var seedfilePaths = make(map[string]bool)
		if SeedFilepaths != "" {
			seedfilePaths = readSeedFile()
		}

		hashes, err := targetHasher.HashAllBazelTargetsAndSourcefiles(seedfilePaths)
		if err != nil {
			panic(err)
		}

		var buffer bytes.Buffer
		err = prettyEncode(hashes, &buffer)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(args[0], buffer.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println(buffer.String())
	},
}

func readSeedFile() map[string]bool {
	readFile, err := os.Open(SeedFilepaths)

	if err != nil {
		panic(err)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	seedFilepaths := make(map[string]bool)
	for fileScanner.Scan() {
		seedFilepaths[fileScanner.Text()] = true
	}
	return seedFilepaths
}

func prettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateHashesCmd)
	generateHashesCmd.Flags().StringVarP(&SeedFilepaths, "seed-filepaths", "", "",
		"A text file containing a newline separated list of filepaths, "+
			"each of these filepaths will be read and used as a seed for all targets.")
	generateHashesCmd.Flags().BoolVarP(&DisplayElapsedTime, "displayElapsedTime", "d", false,
		"This flag controls whether to print out elapsed time for bazel query and content hashing")
}

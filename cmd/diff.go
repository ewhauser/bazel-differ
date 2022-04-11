package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Writes to a file the impacted targets between two Bazel graph JSON files",
	Long:  `Writes to a file the impacted targets between two Bazel graph JSON files`,
	Run: func(cmd *cobra.Command, args []string) {
		if StartingHashes == "" {
			fmt.Println("Starting hashes is required")
			os.Exit(1)
		}
		if _, err := os.Stat(StartingHashes); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Starting hashes is required")
			os.Exit(1)
		}
		if FinalHashes == "" {
			fmt.Println("Final hashes is required")
			os.Exit(1)
		}
		if _, err := os.Stat(FinalHashes); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Final hashes is required")
			os.Exit(1)
		}
		if Output == "" {
			fmt.Println("Output path is required")
			os.Exit(1)
		}
		targetHasher := internal.NewTargetHashingClient(GetBazelClient(), internal.Filesystem,
			internal.NewRuleProvider())

		startingHashes := readHashFile(StartingHashes)
		finalHashes := readHashFile(FinalHashes)
		targets, err := targetHasher.GetImpactedTargets(startingHashes, finalHashes)
		if err != nil {
			panic(err)
		}
		writeFile(targets, Output)
	},
}

func writeFile(targets map[string]bool, output string) {
	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()

	datawriter := bufio.NewWriter(file)
	defer datawriter.Flush()

	for k := range targets {
		_, _ = datawriter.WriteString(k + "\n")
	}
}

func readHashFile(filename string) map[string]string {
	x := map[string]string{}
	startingContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(startingContent, &x)
	if err != nil {
		panic(err)
	}
	return x
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

package cmd

import (
	"errors"
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"github.com/spf13/cobra"
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

		startingHashes, err := internal.ReadHashFile(StartingHashes)
		ExitIfError(err, "")
		finalHashes, err := internal.ReadHashFile(FinalHashes)
		ExitIfError(err, "")
		targets, err := targetHasher.GetImpactedTargets(startingHashes, finalHashes)
		ExitIfError(err, "")
		internal.WriteTargetsFile(targets, Output)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().StringVarP(&StartingHashes, "startingHashes", "s", "",
		"The path to the JSON file of target hashes for the initial revision. Run 'generate-hashes' to get this value.")
	diffCmd.PersistentFlags().StringVarP(&FinalHashes, "finalHashes", "f", "",
		"The path to the JSON file of target hashes for the final revision. Run 'generate-hashes' to get this value.")
	diffCmd.PersistentFlags().StringVarP(&Output, "output", "o", "",
		"Filepath to write the impacted Bazel targets to, "+
			"newline separated")
}

package cmd

import (
	"bufio"
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"os"

	"github.com/spf13/cobra"
)

var SeedFilepaths string
var displayElapsedTime bool

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
		ExitIfError(err, "")

		res, err := internal.WriteHashFile(args[0], hashes)
		ExitIfError(err, "")
		fmt.Println(res)
	},
}

func readSeedFile() map[string]bool {
	readFile, err := os.Open(SeedFilepaths)

	ExitIfError(err, fmt.Sprintf("Error reading seedfile: %s", SeedFilepaths))

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	seedFilepaths := make(map[string]bool)
	for fileScanner.Scan() {
		seedFilepaths[fileScanner.Text()] = true
	}
	return seedFilepaths
}

func init() {
	rootCmd.AddCommand(generateHashesCmd)
	generateHashesCmd.Flags().StringVarP(&SeedFilepaths, "seed-filepaths", "", "",
		"A text file containing a newline separated list of filepaths, "+
			"each of these filepaths will be read and used as a seed for all targets.")
	generateHashesCmd.Flags().BoolVarP(&displayElapsedTime, "displayElapsedTime", "d", false,
		"This flag controls whether to print out elapsed time for bazel query and content hashing")
	generateHashesCmd.PersistentFlags().StringVarP(&StartingHashes, "startingHashes", "s", "",
		"The path to the JSON file of target hashes for the initial revision. Run 'generate-hashes' to get this value.")
	generateHashesCmd.PersistentFlags().StringVarP(&FinalHashes, "finalHashes", "f", "",
		"The path to the JSON file of target hashes for the final revision. Run 'generate-hashes' to get this value.")
}

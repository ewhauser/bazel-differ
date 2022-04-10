package cmd

import (
	"github.com/ewhauser/bazel-differ/internal"
	"os"

	"github.com/spf13/cobra"
)

var WorkspacePath string
var BazelPath string
var StartingHashes string
var FinalHashes string
var Output string
var BazelStartupOptions string
var BazelCommandOptions string
var KeepGoing bool
var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bazel-differ",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&WorkspacePath, "workspacePath", "w", "", "Path to Bazel workspace directory.")
	rootCmd.PersistentFlags().StringVarP(&BazelPath, "bazelPath", "b", "", "Path to Bazel binary")
	rootCmd.PersistentFlags().StringVarP(&StartingHashes, "startingHashes", "s", "",
		"The path to the JSON file of target hashes for the initial revision. Run 'generate-hashes' to get this value.")
	rootCmd.PersistentFlags().StringVarP(&FinalHashes, "finalHashes", "f", "",
		"he path to the JSON file of target hashes for the final revision. Run 'generate-hashes' to get this value.")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "",
		"Filepath to write the impacted Bazel targets to, "+
			"newline separated")
	rootCmd.PersistentFlags().StringVarP(&BazelStartupOptions, "bazelStartupOptions", "y", "",
		"Additional space separated Bazel client startup options used when invoking Bazel")
	rootCmd.PersistentFlags().StringVarP(&BazelCommandOptions, "bazelCommandOptions", "z", "",
		"Additional space separated Bazel command options used when invoking Bazel")
	rootCmd.PersistentFlags().BoolVarP(&KeepGoing, "keep_ooing", "k", true,
		"This flag controls if `bazel query` will be executed with the `--keep_going` flag or not. Disabling this flag allows you to catch configuration issues in your Bazel graph, but may not work for some Bazel setups. Defaults to `true`")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enables verbose output")
}

func GetBazelClient() internal.BazelClient {
	return internal.NewBazelClient(internal.Filesystem, WorkspacePath, BazelPath, Verbose, KeepGoing,
		DisplayElapsedTime,
		BazelStartupOptions, BazelCommandOptions)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

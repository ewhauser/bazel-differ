package internal

import (
	"crypto/sha256"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//go:generate mockgen -destination=../mocks/bazel_client_mock.go -package=mocks github.com/ewhauser/bazel-differ/internal BazelClient
type BazelClient interface {
	QueryAllTargets() ([]*Target, error)
	QueryAllSourceFileTargets() (map[string]*BazelSourceFileTarget, error)
}

type bazelClient struct {
	filesystem         fs.FS
	bazel              Bazel
	workingDirectory   string
	bazelPath          string
	verbose            bool
	keepGoing          bool
	displayElapsedTime bool
	startupOptions     []string
	commandOptions     []string
}

func NewBazelClient(filesystem fs.FS, workingDirectory string, bazelPath string, verbose bool, keepGoing bool,
	displayElapsedTime bool, startupOptions string, commandOptions string) BazelClient {
	return &bazelClient{
		bazel:              NewBazel(bazelPath, workingDirectory),
		workingDirectory:   workingDirectory,
		bazelPath:          bazelPath,
		verbose:            verbose,
		keepGoing:          keepGoing,
		displayElapsedTime: displayElapsedTime,
		startupOptions:     strings.Split(startupOptions, " "),
		commandOptions:     strings.Split(commandOptions, " "),
		filesystem:         filesystem,
	}
}

func (b bazelClient) QueryAllTargets() ([]*Target, error) {
	return b.performBazelQuery("'//external:all-targets' + '//...:all-targets'")
}

func (b bazelClient) QueryAllSourceFileTargets() (m map[string]*BazelSourceFileTarget, err error) {
	targets, err := b.performBazelQuery("kind('source file', //...:all-targets)")
	if err != nil {
		return nil, err
	}
	return b.processBazelSourcefileTargets(targets, true)
}

func (b bazelClient) processBazelSourcefileTargets(targets []*Target,
	readSourcefileTargets bool) (map[string]*BazelSourceFileTarget, error) {
	var sourceTargets = make(map[string]*BazelSourceFileTarget)
	for _, target := range targets {
		var sourceFile = target.SourceFile
		if sourceFile != nil {
			var digest = sha256.New()
			digest.Write([]byte(*sourceFile.Name))
			for _, subinclude := range sourceFile.Subinclude {
				digest.Write([]byte(subinclude))
			}

			var workingDirectory = ""
			if readSourcefileTargets {
				workingDirectory = b.workingDirectory
			}

			var sourceFileTarget, err = NewBazelSourceFileTarget(
				sourceFile.GetName(),
				digest.Sum(nil),
				b.filesystem,
				workingDirectory,
			)

			if err != nil {
				return nil, err
			}

			sourceTargets[*sourceFileTarget.GetName()] = &sourceFileTarget
		}
	}
	return sourceTargets, nil
}

func (b bazelClient) performBazelQuery(query string) ([]*Target, error) {
	file, err := ioutil.TempFile("", ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(query)
	if err != nil {
		return nil, err
	}

	var cmd []string
	if b.verbose {
		cmd = append(cmd, "--bazelrc=/dev/null")
	}
	//cmd = append(cmd, b.startupOptions...)
	cmd = append(cmd, "--order_output=no")
	if b.keepGoing {
		cmd = append(cmd, "--keep_going")
	}
	//cmd = append(cmd, b.commandOptions...)
	cmd = append(cmd, "--query_file="+file.Name())

	result, err := b.bazel.Query(cmd...)
	if err != nil {
		return nil, err
	}
	return result.Target, nil
}

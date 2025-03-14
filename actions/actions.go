package actions

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// In our previous design, each platform should generate *_{OS}_{Arch}.go file
// Feb 12th, this design revoked, still keep the code.
// var currentSuffix = runtime.GOOS + "_" + runtime.GOARCH

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// envToString converts a env map to string
func envToString(envm map[string]string) string {
	var env []string

	for name, value := range envm {
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(env, "\n")
}

// Setenv sets the value of the Github Action environment variable named by the key.
func Setenv(envm map[string]string) {
	env, err := os.OpenFile(os.Getenv("GITHUB_ENV"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// should never happen,
	// it means current runtime is not Github actions if there's any errors
	must(err)

	env.WriteString(envToString(envm))

	// make sure we write it to the GITHUB_ENV
	env.Close()
}

// Setenv sets the value of the Github Action workflow output named by the key.
func SetOutput(envm map[string]string) {
	env, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	must(err)

	env.WriteString(envToString(envm))

	env.Close()
}

// Changes returns the changed files in current PR,
// which depends on ALL_CHANGED_FILES generated by tj-actions/changed-files action,
// if there's no content in ALL_CHANGED_FILES, it panic.
func Changes() []string {
	changes := os.Getenv("ALL_CHANGED_FILES")
	if changes == "" {
		panic("cannot find changes file!")
	}
	return strings.Fields(changes)
}

// Repository returns owner and repository name for the current repository
//
// Example: MeteorsLiu/llpkg, owner: MeteorsLiu, repo: llpkg
func Repository() (owner, repo string) {
	thisRepo := os.Getenv("GITHUB_REPOSITORY")
	if thisRepo == "" {
		panic("no github repo")
	}
	current := strings.Split(thisRepo, "/")
	return current[0], current[1]
}

// Token returns Github Token for current runner
func Token() string {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		panic("no GH_TOKEN")
	}
	return token
}

// CreateBranch creates a branch with the speficied name from the speficied tag.
// It returns an error if creating fail.
func CreateBranch(branchName, tag string) error {
	ret, err := exec.Command("git", "checkout", "-b", branchName, tag).CombinedOutput()
	if err != nil {
		return errors.New(string(ret))
	}
	ret, err = exec.Command("git", "push", "-u", "origin", branchName).CombinedOutput()
	if err != nil {
		return errors.New(string(ret))
	}
	return nil
}

// latestCommitMessageInPR returns the latest commit in PR using git
// so it's required to be at PR's ref.
// In Github Action, Checkout Action will switch to PR's ref automatically,
// it MUST NOT be used outside Github Action.
func latestCommitMessageInPR() string {
	// assume we're at PR's ref
	ret, err := exec.Command("git", "show", "-s", os.Getenv("GITHUB_SHA"), "--format", "%s").CombinedOutput()
	if err != nil {
		panic(string(ret))
	}
	return string(ret)
}

func mappedVersion() string {
	// get message via git
	message := latestCommitMessageInPR()

	// get the mapped version
	mappedVersion := regex(".*").FindString(message)

	if mappedVersion == "" {
		panic("invalid pr: no mapped version found")
	}
	return strings.TrimPrefix(mappedVersion, "Release-as: ")
}

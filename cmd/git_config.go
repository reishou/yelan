package cmd

import (
	"fmt"
	"github.com/reishou/yelan/util"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupGitConfig()
	},
}

func init() {
	gitCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", "--get", key)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// Remove trailing newline characters from the output
	value := strings.TrimSpace(string(output))
	return value, nil
}

func setGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	err := cmd.Run()
	return err
}

func isEmptyConfig(key string) bool {
	value, _ := getGitConfig(key)
	return len(value) == 0
}

func setupGitUser(name string, email string) error {
	err := setGitConfig("user.name", name)
	if err != nil {
		fmt.Println("failed to setup git user name: ", err)
		return err
	}

	err = setGitConfig("user.email", email)
	if err != nil {
		fmt.Println("failed to setup git user email: ", err)
		return err
	}

	return nil
}

func backupFile(sourcePath string) error {
	destinationPath := sourcePath + ".backup"
	cmd := exec.Command("mv", sourcePath, destinationPath)
	err := cmd.Run()
	return err
}

func setupGitConfig() {
	managed, _ := getGitConfig("dotfiles.managed")

	// if there is no user.email, we'll assume it's a new machine/setup and ask it
	if isEmptyConfig("user.email") {
		name := util.Read("What is your github author name? ")
		email := util.Read("What is your github author email? ")

		err := setupGitUser(name, email)
		if err != nil {
			fmt.Println("failed to setup git user email and name: ", err)
			return
		}
	} else if managed != "true" {
		// if user.email exists, let's check for dotfiles.managed config. If it is
		// not true, we'll back up the gitconfig file and set previous user.email and
		// user.name in the new one
		name, _ := getGitConfig("user.name")
		email, _ := getGitConfig("user.email")

		err := setupGitUser(name, email)
		if err != nil {
			fmt.Println("failed to setup git user email and name: ", err)
			return
		}

		err = backupFile("~/.gitconfig")
		if err != nil {
			fmt.Println("failed to backup git config: ", err)
			return
		}
		fmt.Println("moved ~/.gitconfig to ~/.gitconfig.backup")
	} else {
		// otherwise this gitconfig was already made by the dotfiles
		fmt.Println("already managed by dotfiles")
	}
	// todo: include the gitconfig.local file
	// finally make git knows this is a managed config already, preventing later
	// overrides by this script
	err := setGitConfig("dotfiles.managed", "true")
	if err != nil {
		fmt.Println("failed to setup git: ", err)
		return
	}

	name, _ := getGitConfig("user.name")
	email, _ := getGitConfig("user.email")
	fmt.Println("user.name: ", name)
	fmt.Println("user.email: ", email)
}

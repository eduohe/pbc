package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TSRAppLabs/pbc.git"

	"github.com/spf13/cobra"
)

func main() {
	pbc.InitDataDir()

	rootCmd.Execute()
}

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "pbc",
	}

	rootCmd.AddCommand(mkBuildCommand())
	rootCmd.AddCommand(mkProfileCommand())
	rootCmd.AddCommand(mkLintCommand())
}

func mkBuildCommand() *cobra.Command {
	var profileName string
	var name string

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "builds a pkpass",
		Long:  "builds a pkpass",
		Run: func(cmd *cobra.Command, args []string) {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			if profileName == "" {
				log.Fatal("Unable to compile without profile\n")
			}

			if name == "" {
				name = "pass.pkpass"
			}

			profile, err := pbc.GetProfile(profileName)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
				return
			}

			file, err := os.Create(name)

			if err != nil {

				fmt.Printf("Trying to create file: %v, %v", name, err)
				os.Exit(1)
			}

			err = pbc.Compile(root, profile, file)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)

			}
		},
	}
	buildCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Profile to use")
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "Resulting passbook file")
	return buildCmd
}

func mkProfileCommand() *cobra.Command {
	var name string
	var p12path string

	addProfile := func(cmd *cobra.Command, args []string) {
		if name != "" && p12path != "" {
			profile, err := pbc.CreateProfile(name, p12path)

			if err != nil {
				log.Fatal(err)
			}

			pbc.AddProfile(profile)
		}
	}

	profCmd := &cobra.Command{
		Use:   "profile",
		Short: "profile management",
		Long:  "manages a profile to add, rm, ls",
	}

	profAddCmd := &cobra.Command{
		Use:   "add",
		Short: "adds a profile",
		Run:   addProfile,
	}
	profAddCmd.Flags().StringVarP(&name, "profile", "p", "", "Name to give the profile")
	profAddCmd.Flags().StringVarP(&p12path, "cert", "c", "", "Cert to create profile with")

	profLsCmd := &cobra.Command{
		Use:   "ls",
		Short: "lists profiles",
		Long:  "lists profiles",
		Run: func(cmd *cobra.Command, args []string) {
			for _, prof := range pbc.ListProfiles() {
				fmt.Printf("\t%v\n", prof.Name)
			}
		},
	}

	profRmCmd := &cobra.Command{
		Use:   "rm",
		Short: "removes profiles",
		Long:  "removes all profile specified",
		Run: func(cmd *cobra.Command, args []string) {

			for _, arg := range args {
				pbc.DelProfile(arg)
			}

		},
	}

	profCmd.AddCommand(profAddCmd)
	profCmd.AddCommand(profLsCmd)
	profCmd.AddCommand(profRmCmd)

	return profCmd
}

func mkLintCommand() *cobra.Command {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "checks a pass for mistakes",
		Long:  "checks a pass for mistakes",
		Run: func(cmd *cobra.Command, args []string) {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			warn, err := pbc.LintPass(dir)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, msg := range warn {
				fmt.Println(msg)
			}

		},
	}

	return lintCmd
}

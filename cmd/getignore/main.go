package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SowTag/getignore/internal/getignore"
	"github.com/SowTag/getignore/internal/getignore/utils"
)

func main() {
	service := &getignore.GitignoreService{}

	var selectedLanguage *string

	if len(os.Args) > 1 {
		selectedLanguage = &os.Args[1]
	}

	if selectedLanguage == nil {
		ignores := unwrap(service.GetIgnores())

		fmt.Println("* Available languages:")
		utils.PrintDynamicTable(ignores)
		fmt.Println()
		printUsage()

		return
	}

	if _, err := os.Stat(".gitignore"); err == nil {
		panic("a .gitignore file already exists in the current directory")
	}

	ignore := unwrap(service.GetGitignoreContents(*selectedLanguage))

	unwrap(0, os.WriteFile(".gitignore", []byte(*ignore), 0644))

	fmt.Println("✓ gitignore template downloaded!")
}

func unwrap[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}

	return x
}

func printUsage() {
	programName := filepath.Base(os.Args[0])

	fmt.Printf("Usage: %s <language>\n", programName)
	fmt.Println("Note that <language> is case-insensitive but must be exact. For example, \"go\", \"GO\", etc. are valid selections of the Go gitignore while \"golang\" is not.")
	fmt.Println("Gitignores are fetched from", getignore.GitignoreRepoContentsURL)

}

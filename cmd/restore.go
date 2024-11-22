package cmd

import (
	"andreaG757/docker-back/pkg"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore the container",
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {

		} else {
			backupFiles, err := getDockerBackupFiles(storeDir)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			selectedFile, runc := restoreHuh(backupFiles)

			doRestore(selectedFile, runc)
			// fmt.Println(getBackedContainers())
		}
	},
}

func restoreHuh(backupFiles []string) ([]string, bool) {
	var selectedFiles []string
	var runc bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose the containers to store (press space or x to select)").
				Value(&selectedFiles).
				OptionsFunc(func() []huh.Option[string] {
					options := make([]huh.Option[string], len(backupFiles))

					for i, file := range backupFiles {
						filename := filepath.Base(file)
						options[i] = huh.Option[string]{
							Key:   filename,
							Value: file,
						}
					}
					return options
				}, nil),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to run in background the containers?").
				Value(&runc),
		),
	)

	form.Run()

	return selectedFiles, runc
}

func getDockerBackupFiles(storeDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(storeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tar") {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func doRestore(files []string, runc bool) {
	for _, filename := range files {
		parts := strings.Split(filename, "-")
		image := parts[2]
		name := strings.ReplaceAll(parts[3], ".tar", "")
		backupFileSave := storeDir + "/" + filename
		pkg.RunCommand("docker", "load", "-i", filename)
		if runc {
			pkg.RunCommand("docker", "run", "-it", "-d", "--name", name, image)
		}
		fmt.Println(backupFileSave)
	}
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}

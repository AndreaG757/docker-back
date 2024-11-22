package cmd

import (
	"andreaG757/docker-back/model"
	"andreaG757/docker-back/pkg"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.UserHomeDir()
	storeDir   = filepath.Join(homeDir, ".local/docker-back")
)

var storeCmd = &cobra.Command{
	Use:     "store",
	Short:   "Store the container",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CreateDirIfNotExist(storeDir)
		if len(args) > 0 {

		} else {
			c := storeHuh()
			if len(c) > 0 {
				doStore(c)
			} else {
				fmt.Println("No container was selected")
			}
		}
	},
}

func storeHuh() []model.DockerContainer {
	var selectedContainers []model.DockerContainer
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[model.DockerContainer]().
				Title("Choose the containers to store (press space or x to select)").
				Value(&selectedContainers).
				OptionsFunc(func() []huh.Option[model.DockerContainer] {
					containers, err := fetchContainerNames()
					if err != nil {
						fmt.Println("Error:", err)
						return nil
					}

					options := make([]huh.Option[model.DockerContainer], len(containers))
					for i, container := range containers {
						options[i] = huh.Option[model.DockerContainer]{
							Key: fmt.Sprintf(`[%s] Name: %s - ID: %s - Image: %s`,
								checkBackedContainer(container),
								container.Name,
								container.ID,
								container.Image),
							Value: container,
						}
					}

					return options
				}, nil),
		),
	)

	form.Run()

	return selectedContainers
}

func fetchContainerNames() ([]model.DockerContainer, error) {
	var containers []model.DockerContainer

	output, err := pkg.ReturnCommand("docker", "ps", "-a", "--format", "{{.ID}} {{.Names}} {{.Image}}")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		list := strings.Split(line, " ")

		containers = append(containers, model.DockerContainer{
			ID:    list[0],
			Name:  list[1],
			Image: list[2],
		})
	}
	return containers, nil
}

func checkBackedContainer(dc model.DockerContainer) string {
	filePath := filepath.Join(storeDir, backupName(dc))
	_, err := os.Stat(filePath)

	if err == nil {
		return "✔"
	} else {
		return "✘"
	}
}

func backupName(dc model.DockerContainer) string {
	return fmt.Sprintf("%s-%s-%s.tar", dc.Name, dc.Image, dc.ID)
}

func doStore(containers []model.DockerContainer) {
	for _, container := range containers {
		backupFileSave := storeDir + "/" + backupName(container)
		pkg.RunCommand("docker", "stop", container.ID)
		pkg.RunCommand("docker", "save", "-o", backupFileSave, container.Image)
		fmt.Println(backupFileSave)
	}
}

func init() {
	RootCmd.AddCommand(storeCmd)
}

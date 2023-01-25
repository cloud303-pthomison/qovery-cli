package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/qovery/qovery-cli/utils"
	"github.com/qovery/qovery-client-go"
	"github.com/spf13/cobra"
)

var applicationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an application",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Capture(cmd)

		tokenType, token, err := utils.GetAccessToken()
		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		client := utils.GetQoveryClient(tokenType, token)
		_, _, envId, err := getContextResourcesId(client)

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		if !utils.IsEnvironmentInATerminalState(envId, client) {
			utils.PrintlnError(fmt.Errorf("environment id '%s' is not in a terminal state. The request is not queued and you must wait "+
				"for the end of the current operation to run your command. Try again in a few moment", envId))
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		applications, _, err := client.ApplicationsApi.ListApplication(context.Background(), envId).Execute()

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		application := utils.FindByApplicationName(applications.GetResults(), applicationName)

		if application == nil {
			utils.PrintlnError(fmt.Errorf("application %s not found", applicationName))
			utils.PrintlnInfo("You can list all applications with: qovery application list")
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		var storage []qovery.ServiceStorageRequestStorageInner
		for _, s := range application.Storage {
			storage = append(storage, qovery.ServiceStorageRequestStorageInner{
				Id:         &s.Id,
				Type:       s.Type,
				Size:       s.Size,
				MountPoint: s.MountPoint,
			})
		}

		req := qovery.ApplicationEditRequest{
			Name:        application.Name,
			Description: application.Description.Get(),
			GitRepository: &qovery.ApplicationGitRepositoryRequest{
				Url:      *application.GitRepository.Url,
				Branch:   application.GitRepository.Branch,
				RootPath: application.GitRepository.RootPath,
			},
			BuildMode:           application.BuildMode,
			DockerfilePath:      application.DockerfilePath.Get(),
			Cpu:                 application.Cpu,
			Memory:              application.Memory,
			MinRunningInstances: application.MinRunningInstances,
			MaxRunningInstances: application.MaxRunningInstances,
			Healthcheck:         application.Healthcheck,
			AutoPreview:         application.AutoPreview,
			Ports:               application.Ports,
			Storage:             storage,
		}

		if applicationBranch != "" {
			req.GitRepository.Branch = &applicationBranch
		}

		_, _, err = client.ApplicationMainCallsApi.EditApplication(context.Background(), application.Id).ApplicationEditRequest(req).Execute()

		if err != nil {
			utils.PrintlnError(err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}

		utils.Println(fmt.Sprintf("Application %s updated!", pterm.FgBlue.Sprintf(applicationName)))
	},
}

func init() {
	applicationCmd.AddCommand(applicationUpdateCmd)
	applicationUpdateCmd.Flags().StringVarP(&organizationName, "organization", "", "", "Organization Name")
	applicationUpdateCmd.Flags().StringVarP(&projectName, "project", "", "", "Project Name")
	applicationUpdateCmd.Flags().StringVarP(&environmentName, "environment", "", "", "Environment Name")
	applicationUpdateCmd.Flags().StringVarP(&applicationName, "application", "n", "", "Application Name")
	applicationUpdateCmd.Flags().StringVarP(&applicationBranch, "branch", "", "", "Application Git Branch")

	_ = applicationUpdateCmd.MarkFlagRequired("application")
}

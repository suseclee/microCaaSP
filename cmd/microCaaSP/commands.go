package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/suseclee/microCaaSP/configs/constants"
	"github.com/suseclee/microCaaSP/pkg/microCaaSP"
)

func DeployCmd() *cobra.Command {
	microCaaSP := &microCaaSP.MicroCaaSP{}
	cmd := cobra.Command{
		Use:   "deploy",
		Short: "deploy microCaaSP cluster",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetString("version")
			microCaaSP.Init(version)
			microCaaSP.Deploy()
		},
	}

	cmd.Flags().String("version", constants.GetAvilableVersions()[0], "version to use ("+strings.Join(constants.GetAvilableVersions(), ",")+")")
	return &cmd
}

func Login() *cobra.Command {
	cmd := cobra.Command{
		Use:   "login",
		Short: "login microCaaSP cluster",
		Run: func(cmd *cobra.Command, args []string) {
			(&microCaaSP.MicroCaaSP{}).Login()
		},
	}
	return &cmd
}

func DestroyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "destroy",
		Short: "destroy microCaaSP cluster",
		Run: func(cmd *cobra.Command, args []string) {
			(&microCaaSP.MicroCaaSP{}).Destroy()
		},
	}
	return &cmd
}

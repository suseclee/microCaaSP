package main

import (
	"github.com/spf13/cobra"
	"github.com/suseclee/microCaaSP/pkg/microCaaSP"
)

func NewCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "deploy",
		Short: "deploy microCaaSP cluster",
		Run: func(cmd *cobra.Command, args []string) {
			(&microCaaSP.MicroCaaSP{}).Deploy()
		},
	}
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

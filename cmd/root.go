package cmd

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elewis787/boa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute is the command line applications entry function
func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "bcps",
		Long:    "bcps is a command line application that generates payment schedules for income sharing agreements",
		Example: "bcps",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			viper.AutomaticEnv()
			viper.SetEnvPrefix("bcps")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.SetHelpFunc(boa.HelpFunc)
	rootCmd.SetUsageFunc(boa.UsageFunc)

	return rootCmd.ExecuteContext(context.Background())
}

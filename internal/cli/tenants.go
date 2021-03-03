package cli

import (
	"fmt"

	"github.com/auth0/auth0-cli/internal/prompt"
	"github.com/spf13/cobra"
)

func tenantsCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenants",
		Short: "Manage configured tenants",
	}

	cmd.AddCommand(useTenantCmd(cli))
	return cmd
}

func useTenantCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use",
		Aliases: []string{"select"},
		Short:   "Set the active tenant",
		Long:    `auth0 tenants use <tenant>`,
		Args:    cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			prepareInteractivity(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var selectedTenant string
			if len(args) == 0 {
				tens, err := cli.listTenants()
				if err != nil {
					return fmt.Errorf("Unable to load tenants due to an unexpected error: %w", err)
				}

				tenNames := make([]string, len(tens))
				for i, t := range tens {
					tenNames[i] = t.Name
				}

				input := prompt.SelectInput("tenant", "Tenant:", "Tenant to activate", tenNames, true)
				if err := prompt.AskOne(input, &selectedTenant); err != nil {
					return fmt.Errorf("An unexpected error occurred: %w", err)
				}
			} else {
				requestedTenant := args[0]
				t, ok := cli.config.Tenants[requestedTenant]
				if !ok {
					return fmt.Errorf("Unable to find tenant %s; run `auth0 tenants use` to see your configured tenants or run `auth0 login` to configure a new tenant", requestedTenant)
				}
				selectedTenant = t.Name
			}

			cli.config.DefaultTenant = selectedTenant
			if err := cli.persistConfig(); err != nil {
				return fmt.Errorf("An error occurred while setting the default tenant: %w", err)
			}
			return nil
		},
	}

	return cmd
}
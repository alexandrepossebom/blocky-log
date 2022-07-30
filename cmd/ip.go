package cmd

import (
	"errors"
	"fmt"
	"net"

	"github.com/alexandrepossebom/blocky-log/utils"
	"github.com/spf13/cobra"
)

var eventTypeIP string

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show IP logs",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a ip argument")
		}
		if net.ParseIP(args[0]) != nil {
			return nil
		}
		return fmt.Errorf("invalid ip specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running its can take a while...")
		utils.ListEventsByIP(Host, args[0], time, eventTypeIP)
	},
}

func init() {
	ipCmd.Flags().StringVarP(&eventTypeIP, "event", "e", "all", "type of event to show")
	ipCmd.Flags().IntVarP(&time, "hours", "", 24, "hours to show")
	rootCmd.AddCommand(ipCmd)
}

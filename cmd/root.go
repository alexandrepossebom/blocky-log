package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexandrepossebom/blocky-log/config"
	"github.com/alexandrepossebom/blocky-log/utils"
	"github.com/spf13/cobra"
)

var time int
var eventType string
var Host string

var rootCmd = &cobra.Command{
	Use:   "blocky-log",
	Short: "Utility to show logs from blocky",
	Long:  `By default it's show blocked queries but you can change to show by IP or by event type`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running it's can take a while...")
		utils.ListEvents(Host, eventType, time)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	hosts := []string{}
	for _, host := range config.Get().Blockys {
		hosts = append(hosts, host.Name)
	}
	if len(hosts) == 0 {
		panic("No hosts found")
	} else if len(hosts) == 1 {
		rootCmd.PersistentFlags().StringVarP(&Host, "host", "", hosts[0], "host to use")
	} else {
		rootCmd.PersistentFlags().StringVarP(&Host, "host", "", hosts[0], "which host to use the valid options are: "+strings.Join(hosts, ", "))
	}
	rootCmd.Flags().IntVarP(&time, "hours", "", config.Get().Hours, "hours to show")
	rootCmd.Flags().StringVarP(&eventType, "event", "e", config.Get().DefaultEventType, "type of event to show, use 'all' to show all events")
}

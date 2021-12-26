package subcmd

import (
	"fmt"
	"strings"

	"github.com/laper32/regsm-console/src/cli/dpkg"
	"github.com/spf13/cobra"
)

func InitSearchCMD() *cobra.Command {
	var short []string
	search := &cobra.Command{
		Use: "search",
		Run: func(cmd *cobra.Command, args []string) {
			// This part may use a high large scale of hardware IO
			// Optimization required
			// Migrate to NoSQL(eg: MongoDB) required
			// Further consideration required
			allGameList := dpkg.AvailableGames()
			var foundGameList []dpkg.AvailableGame
			for _, content := range allGameList {
				for _, thisName := range args {
					if strings.Contains(content.Name, thisName) {
						foundGameList = append(foundGameList, content)
					}
				}
			}

			var result []dpkg.AvailableGame

			if len(short) > 0 {
				for _, content := range foundGameList {
					for _, thisShort := range short {
						if strings.Contains(content.Short, thisShort) {
							result = append(result, content)
						}
					}
				}
			} else {
				result = foundGameList
			}

			for _, v := range result {
				fmt.Printf("%v\n", v)
			}
		},
	}
	search.Flags().StringSliceVar(&short, "short", nil, "")

	return search
}

package subcmd

import (
	"fmt"
	"strings"

	"github.com/laper32/regsm-console/src/app/cli/dpkg"
	"github.com/spf13/cobra"
)

func InitSearchCMD() *cobra.Command {
	var short []string
	search := &cobra.Command{
		Use: "search",
		Run: func(cmd *cobra.Command, args []string) {
			// This idea comes from 'apt search'
			//
			// 	We want to search something, due to this is CLI, that we only
			// allow to search {name, short}, otherwise are not provided.
			// (It is not the CLI can do. If you want to do advanced search,
			// you should use something like ElasticSearch instead of typing in CLI.....)
			// 	It will search both 'name' and 'short', and return results accoring
			// your text.
			//
			// 	Also, this is list search, that you can search based on the list
			// of games what you want to find.
			//
			// Steps
			// 	1. Retrieve all available games.
			// 	2. Search by full name
			// 	3. Search by short, according the previous result on step 2.
			// 	Noting that if no 'short' then will return result directly.
			// 	4. Return result.

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

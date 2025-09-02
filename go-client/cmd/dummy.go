package cmd

import (
	"fmt"
	"go-client/lib/aiclient"
	"log"

	"github.com/spf13/cobra"
)

var dummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "Run dummy command",
	Run:   cmd_dummy,
}

func init() {
	rootCmd.AddCommand(dummyCmd)
}

func cmd_dummy(cmd *cobra.Command, args []string) {
	log.Println("Dummy command here")

	// appdebug.PrettyPrint(appconfig.AppCfg)

	// wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	// ctx := context.Background()
	// cr := cocktail.NewRepository(wvc, &ctx)

	// // cocktails, err := cr.GetListByNearText("Varmouth", 3)
	// cocktails, err := cr.GetByCocktailName("Cove")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// appdebug.PrettyPrint(cocktails)

	promptBuilder := aiclient.PromptBuilder()

	prompt := promptBuilder.
		WithContext("Prompt context").
		WithExamples("Prompt examples").
		WithInstructions("Prompt instructions").
		WithRole("Prompt role").
		WithTask("Prompt task").
		Get()

	fmt.Println(prompt)

}

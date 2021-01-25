package cmd

import (
	"github.com/spf13/cobra"
	"pagarme/transactions"
)

var boletoCmd = &cobra.Command{
	Use:   "boleto",
	Short: "Gerar boleto",
	Run: func(cmd *cobra.Command, args []string) {

		tb := transactions.TransactionBuilder{}

		amount, _ := cmd.Flags().GetFloat64("amount")
		tb.Amount(amount)

		name, _ := cmd.Flags().GetString("name")
		tb.Name(name)

		document, _ := cmd.Flags().GetString("document")
		tb.Document(document)

		tb.PaymentMethod(transactions.BOLETO)
		t := tb.Build()
		transactions.NewClient().Execute(t, transactions.BODY)
	},
}

func init() {
	rootCmd.AddCommand(boletoCmd)
	boletoCmd.Flags().Float64P("amount", "a", 0.0, "Amount value")
	boletoCmd.Flags().StringP("name", "n", "", "Name")
	boletoCmd.Flags().StringP("document", "d", "", "Document")
}

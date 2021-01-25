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
		tb.Amount(20.0)
		tb.Name("Leandro Greijal")
		tb.Country("BR")
		tb.TypeCustomer(transactions.INDIVIDUAL)
		tb.PaymentMethod(transactions.BOLETO)
		tb.Document("251.854.650-26")

		t := tb.Build()
		transactions.NewClient().Execute(t, transactions.BODY)

	},
}

func init() {
	rootCmd.AddCommand(boletoCmd)
}

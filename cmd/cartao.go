package cmd

import (
	"pagarme/transactions"

	"github.com/spf13/cobra"
)

var cartaoCmd = &cobra.Command{
	Use:   "cartao",
	Short: "Gerar cobramça cartão",
	Run: func(cmd *cobra.Command, args []string) {

		tb := transactions.TransactionBuilder{}
		tb.Amount(2.0)
		tb.Name("Leandro Greijal")
		tb.Country("BR")
		tb.TypeCustomer(transactions.INDIVIDUAL)
		tb.Document("251.854.650-26")
		tb.PaymentMethod(transactions.CREDIT_CARD)
		tb.CardNumber("41111111111111")
		tb.CardHolderName("Leandro")
		tb.CardExpirationDate("1028")
		tb.CardCVV("123")

		transaction := tb.Build()
		client := transactions.NewClient()
		publicKey := client.RecoverPublicKey()
		transaction.CreateCardHash(publicKey)

		client.Execute(transaction, transactions.BASIC_AUTH)
	},
}

func init() {
	rootCmd.AddCommand(cartaoCmd)
}

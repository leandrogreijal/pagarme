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

		amount, _ := cmd.Flags().GetFloat64("amount")
		tb.Amount(amount)

		name, _ := cmd.Flags().GetString("name")
		tb.Name(name)

		document, _ := cmd.Flags().GetString("document")
		tb.Document(document)

		country, _ := cmd.Flags().GetString("country")
		tb.Country(country)

		cardNumber, _ := cmd.Flags().GetString("cardNumber")
		tb.CardNumber(cardNumber)

		cardHolderName, _ := cmd.Flags().GetString("cardHolderName")
		tb.CardHolderName(cardHolderName)

		cardExpirationDate, _ := cmd.Flags().GetString("cardExpirationDate")
		tb.CardExpirationDate(cardExpirationDate)

		cardCVV, _ := cmd.Flags().GetString("cardCVV")
		tb.CardCVV(cardCVV)

		tb.PaymentMethod(transactions.CREDIT_CARD)

		transaction := tb.Build()
		client := transactions.NewClient()
		publicKey := client.RecoverPublicKey()
		transaction.CreateCardHash(publicKey)

		client.Execute(transaction, transactions.BASIC_AUTH)
	},
}

func init() {
	rootCmd.AddCommand(cartaoCmd)
	cartaoCmd.Flags().Float64P("amount", "a", 0.0, "Amount value")
	cartaoCmd.Flags().StringP("name", "n", "", "Name")
	cartaoCmd.Flags().StringP("document", "d", "", "Document")
	cartaoCmd.Flags().StringP("country", "C", "", "Country")
	cartaoCmd.Flags().StringP("cardNumber", "c", "", "Card Number")
	cartaoCmd.Flags().StringP("cardHolderName", "N", "", "Card Holder Name")
	cartaoCmd.Flags().StringP("cardExpirationDate", "e", "", "Card Expiration Date")
	cartaoCmd.Flags().StringP("cardCVV", "v", "", "Card CVV")
}

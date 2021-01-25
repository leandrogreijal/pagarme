package transactions

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteError500(t *testing.T) {

	tb := TransactionBuilder{}
	tb.Amount(2.0)
	tb.Name("Leandro Greijal")
	tb.Country("BR")
	tb.PaymentMethod(BOLETO)
	tb.Document("251.854.650-26")

	transaction := tb.Build()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}),
	)

	defer server.Close()

	client := client{server.Client(), server.URL}
	result, err := client.Execute(transaction, BASIC_AUTH)

	assertTest := assert.New(t)
	assertTest.Nil(result)
	assertTest.EqualError(err, "Mundipagg internal error. Path: /transactions")

}

func TestExecuteBoletoBasicAuth(t *testing.T) {

	tb := TransactionBuilder{}
	tb.Amount(2.0)
	tb.Name("Leandro Greijal")
	tb.Country("BR")
	tb.PaymentMethod(BOLETO)
	tb.Document("251.854.650-26")
	transaction := tb.Build()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/1/transactions" {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				io.WriteString(w, "{ \"object\": \"transaction\", \"status\": \"refused\"}")
			}
		}),
	)

	defer server.Close()

	client := client{server.Client(), server.URL}
	result, err := client.Execute(transaction, BASIC_AUTH)

	assertTest := assert.New(t)
	assertTest.Nil(err)
	assertTest.NotNil(result)
}

func TestTransactionMarshal(t *testing.T) {
	tb := TransactionBuilder{}
	tb.Amount(2.0)
	tb.Name("Leandro Greijal")
	tb.Country("BR")
	tb.PaymentMethod(BOLETO)
	tb.Document("251.854.650-26")

	transaction := tb.Build()
	json, _ := transaction.marshal()

	assertTest := assert.New(t)
	expectJson := "{\"api_key\":\"ak_test_qCS4GVwDKJhzbTn0Z3KIU4p4k79U17\",\"amount\":200,\"payment_method\":\"boleto\",\"customer\":{\"name\":\"Leandro Greijal\",\"country\":\"BR\",\"type\":\"individual\",\"documents\":[{\"type\":\"cpf\",\"number\":\"25185465026\"}]}}"

	assertTest.Equal(expectJson, string(json))

}

func TestTransactionBuild(t *testing.T) {

	tb := TransactionBuilder{}
	tb.Amount(2.0)
	tb.Name("Leandro Greijal")
	tb.Country("BR")
	tb.PaymentMethod(BOLETO)
	tb.Document("251.854.650-26")

	transactionTest := tb.Build()
	expect := transaction{}
	expect.ApiKey = "ak_test_qCS4GVwDKJhzbTn0Z3KIU4p4k79U17"
	expect.Amount = 200
	expect.Customer.Name = "Leandro Greijal"
	expect.Customer.Country = "BR"
	expect.Customer.CustomerType = INDIVIDUAL.String()
	expect.PaymentMethod = BOLETO.String()
	expect.Customer.Documents = append(expect.Customer.Documents, document{Number: "25185465026", DocumentType: CPF.String()})

	assertTest := assert.New(t)
	assertTest.Equal(expect, transactionTest)

}

func TestCreateDocumentCPF(t *testing.T) {

	tb := TransactionBuilder{}
	tb.Document("251.854.650-26")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("cpf", transactionTest.Customer.Documents[0].DocumentType)
	assertTest.Equal("individual", transactionTest.Customer.CustomerType)
	assertTest.Equal("25185465026", transactionTest.Customer.Documents[0].Number)

}

func TestCreateDocumentCNPJ(t *testing.T) {

	tb := TransactionBuilder{}
	tb.Document("30.516.297/0001-03")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("cnpj", transactionTest.Customer.Documents[0].DocumentType)
	assertTest.Equal("corporation", transactionTest.Customer.CustomerType)

}

func TestCreateDocumenInvalid(t *testing.T) {

	tb := TransactionBuilder{}
	_, err := tb.Document("kded;w;wd")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "Document.Number is invalid. Value: kded;w;wd")
}

func TestCountryInvalid(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.Country("o")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "Country is invalid. Value: o")
}

func TestCountry(t *testing.T) {
	tb := TransactionBuilder{}
	tb.Country("BR")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("BR", transactionTest.Customer.Country)
}

func TestNameInvalid(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.Name("")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "Name is invalid. Value: ")
}

func TestName(t *testing.T) {
	tb := TransactionBuilder{}
	tb.Name("Name")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("Name", transactionTest.Customer.Name)
}

func TestPaymentMethod(t *testing.T) {
	tb := TransactionBuilder{}
	tb.PaymentMethod(BOLETO)
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal(BOLETO.String(), transactionTest.PaymentMethod)
}

func TestCardCVV(t *testing.T) {
	tb := TransactionBuilder{}
	tb.CardCVV("123")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("123", transactionTest.CardCVV)
}

func TestCardCVVSize(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardCVV("1234")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardCVV is invalid. Value: 1234")

}

func TestCardCVVNoNumeric(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardCVV("a12")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardCVV is invalid. Value: a12")

}

func TestCardNumber(t *testing.T) {
	tb := TransactionBuilder{}
	tb.CardNumber("123")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("123", transactionTest.CardNumber)
}

func TestCardCardNumberSize(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardNumber("")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardNumber is invalid. Value: ")
}

func TestCardCardNumberNoNumeric(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardNumber("a12")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardNumber is invalid. Value: a12")
}

func TestCardExpirationDate(t *testing.T) {
	tb := TransactionBuilder{}
	tb.CardExpirationDate("0121")
	transactionTest := tb.Build()

	assertTest := assert.New(t)
	assertTest.Equal("0121", transactionTest.CardExpirationDate)
}

func TestCardExpirationDateSize(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardExpirationDate("")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardExpirationDate is invalid. Value: ")
}

func TestCardCardExpirationDateNoNumeric(t *testing.T) {
	tb := TransactionBuilder{}
	_, err := tb.CardExpirationDate("a12")

	assertTest := assert.New(t)
	assertTest.EqualError(err, "CardExpirationDate is invalid. Value: a12")
}

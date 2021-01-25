package transactions

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PaymentMethod int
type TypeCustomer int
type DocumentType int
type AuthenticationMethod int

const BASE_URL = "https://api.pagar.me/1"
const PATH_TRANSACTION = "/transactions"
const PATH_HASH = "/transactions/card_hash_key"
const API_KEY = "ak_test_qCS4GVwDKJhzbTn0Z3KIU4p4k79U17"

const (
	CREDIT_CARD PaymentMethod = iota
	BOLETO
)

const (
	INDIVIDUAL TypeCustomer = iota
	CORPORATION
)

const (
	CPF DocumentType = iota
	CNPJ
)

const (
	BODY AuthenticationMethod = iota
	PARAM
	BASIC_AUTH
)

func (p PaymentMethod) String() string {
	return [...]string{"credit_card", "boleto"}[p]
}

func (t TypeCustomer) String() string {
	return [...]string{"individual", "corporation"}[t]
}

func (d DocumentType) String() string {
	return [...]string{"cpf", "cnpj"}[d]
}

func NewTypeCustomer(value string) TypeCustomer {

	if strings.ToLower(value) == "corporation" {
		return CORPORATION
	}

	return INDIVIDUAL
}

type publicKey struct {
	Id        int    `json:"id"`
	PublicKey string `json:"public_key"`
	Ip        string `json:"ip"`
}

type transactionResponse struct {
	Status                string      `json:"status"`
	RefuseReason          string      `json:"refuse_reason"`
	StatusReason          string      `json:"status_reason"`
	AcquirerResponseCode  string      `json:"acquirer_response_code"`
	AcquirerName          string      `json:"acquirer_name"`
	AcquirerID            string      `json:"acquirer_id"`
	Tid                   int         `json:"tid"`
	Nsu                   int         `json:"nsu"`
	DateCreated           time.Time   `json:"date_created"`
	DateUpdated           time.Time   `json:"date_updated"`
	Amount                int         `json:"amount"`
	Installments          int         `json:"installments"`
	ID                    int         `json:"id"`
	CardHolderName        string      `json:"card_holder_name"`
	CardLastDigits        string      `json:"card_last_digits"`
	CardFirstDigits       string      `json:"card_first_digits"`
	CardBrand             string      `json:"card_brand"`
	CardPinMode           interface{} `json:"card_pin_mode"`
	CardMagstripeFallback bool        `json:"card_magstripe_fallback"`
	CvmPin                bool        `json:"cvm_pin"`
	PaymentMethod         string      `json:"payment_method"`
	CaptureMethod         string      `json:"capture_method"`
	BoletoURL             interface{} `json:"boleto_url"`
	BoletoBarcode         interface{} `json:"boleto_barcode"`
	BoletoExpirationDate  interface{} `json:"boleto_expiration_date"`
	Referer               string      `json:"referer"`
	IP                    string      `json:"ip"`
	Card                  struct {
		ID             string    `json:"id"`
		DateCreated    time.Time `json:"date_created"`
		DateUpdated    time.Time `json:"date_updated"`
		Brand          string    `json:"brand"`
		HolderName     string    `json:"holder_name"`
		FirstDigits    string    `json:"first_digits"`
		LastDigits     string    `json:"last_digits"`
		Country        string    `json:"country"`
		Fingerprint    string    `json:"fingerprint"`
		Valid          bool      `json:"valid"`
		ExpirationDate string    `json:"expiration_date"`
	} `json:"card"`
}

type TransactionI interface {
	marshal()
	CreateCardHash()
}

type transaction struct {
	ApiKey             string `json:"api_key,omitempty"`
	Amount             int64  `json:"amount"`
	CardHash           string `json:"card_hash,omitempty"`
	CardHolderName     string `json:"card_holder_name,omitempty"`
	CardExpirationDate string `json:"card_expiration_date,omitempty"`
	CardNumber         string `json:"card_number,omitempty"`
	CardCVV            string `json:"card_cvv,omitempty"`
	PaymentMethod      string `json:"payment_method,omitempty"`
	Customer           struct {
		ExternalId   string     `json:"number,omitempty"`
		Name         string     `json:"name,omitempty"`
		Country      string     `json:"country,omitempty"`
		CustomerType string     `json:"type,omitempty"`
		Documents    []document `json:"documents,omitempty"`
	} `json:"customer,omitempty"`
}

type document struct {
	DocumentType string `json:"type,omitempty"`
	Number       string `json:"number,omitempty"`
}

func (t *transaction) marshal() ([]byte, error) {
	return json.Marshal(t)
}

func (t *transaction) CreateCardHash(key publicKey) {
	rsaPublicKey := createRsaPublicKey(key.PublicKey)

	params := url.Values{}
	params.Add("card_number", t.CardNumber)
	params.Add("card_holder_name", t.CardHolderName)
	params.Add("card_expiration_date", t.CardExpirationDate)
	params.Add("card_cvv", t.CardCVV)
	querystring := params.Encode()

	pkcs1padding, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(querystring))

	if err != nil {
		log.Fatal("failed to EncryptPKCS1v15 " + err.Error())
	}
	pkcs1padding64 := b64.StdEncoding.EncodeToString(pkcs1padding)

	t.CardCVV = ""
	t.CardExpirationDate = ""
	t.CardHolderName = ""
	t.CardNumber = ""
	t.CardHash = strconv.Itoa(key.Id) + "_" + pkcs1padding64

}

type TransactionBuilderI interface {
	Build()
	Amount(value float64) *TransactionBuilder
	PaymentMethod(value PaymentMethod) *TransactionBuilder
	TypeCustomer(value TypeCustomer) *TransactionBuilder
	CardHolderName(value string) (*TransactionBuilder, error)
	Name(value string) (*TransactionBuilder, error)
	CardExpirationDate(value string) (*TransactionBuilder, error)
	CardNumber(value string) (*TransactionBuilder, error)
	CardCVV(value string) (*TransactionBuilder, error)
	Country(value string) (*TransactionBuilder, error)
	Document(value string) (*TransactionBuilder, error)
}

type TransactionBuilder struct {
	transaction transaction
}

func (b *TransactionBuilder) Build() transaction {
	b.transaction.ApiKey = API_KEY
	transactionFinal := b.transaction
	b.transaction = transaction{}
	return transactionFinal
}

func (b *TransactionBuilder) Amount(value float64) *TransactionBuilder {
	floatString := fmt.Sprintf("%.2f", value)
	floatString = strings.Replace(floatString, ".", "", -1)
	floatString = strings.Replace(floatString, ",", "", -1)
	amount, _ := strconv.ParseInt(floatString, 10, 64)
	b.transaction.Amount = amount
	return b
}

func (b *TransactionBuilder) CardHolderName(value string) (*TransactionBuilder, error) {

	if value == "" {
		return nil, &InvalidValueError{"CardHolderName", value}
	}

	b.transaction.CardHolderName = value
	return b, nil
}

func (b *TransactionBuilder) CardExpirationDate(value string) (*TransactionBuilder, error) {
	regex, _ := regexp.Compile("^\\d{4}$")

	if !regex.MatchString(value) {
		return nil, &InvalidValueError{"CardExpirationDate", value}
	}

	b.transaction.CardExpirationDate = value
	return b, nil
}

func (b *TransactionBuilder) CardNumber(value string) (*TransactionBuilder, error) {

	regex, _ := regexp.Compile("^\\d+$")

	if !regex.MatchString(value) {
		return nil, &InvalidValueError{"CardNumber", value}
	}

	b.transaction.CardNumber = value
	return b, nil
}

func (b *TransactionBuilder) CardCVV(value string) (*TransactionBuilder, error) {

	regex, _ := regexp.Compile("^\\d{3}$")

	if !regex.MatchString(value) {
		return nil, &InvalidValueError{"CardCVV", value}
	}

	b.transaction.CardCVV = value
	return b, nil
}

func (b *TransactionBuilder) PaymentMethod(value PaymentMethod) *TransactionBuilder {
	b.transaction.PaymentMethod = value.String()
	return b
}

func (b *TransactionBuilder) Name(value string) (*TransactionBuilder, error) {

	if strings.TrimSpace(value) == "" {
		return b, &InvalidValueError{"Name", value}
	}

	b.transaction.Customer.Name = value
	return b, nil
}

func (b *TransactionBuilder) Country(value string) (*TransactionBuilder, error) {

	regex, _ := regexp.Compile("^[a-zA-z]{2}$")

	if !regex.MatchString(value) {
		return b, &InvalidValueError{"Country", value}
	}

	b.transaction.Customer.Country = value
	return b, nil
}

func (b *TransactionBuilder) Document(value string) (*TransactionBuilder, error) {

	regexCPF, _ := regexp.Compile("[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}")
	regexCNPJ, _ := regexp.Compile("[0-9]{2}\\.[0-9]{3}\\.[0-9]{3}/[0-9]{4}-[0-9]{2}")

	if regexCPF.MatchString(value) {
		value = strings.Replace(value, ".", "", -1)
		value = strings.Replace(value, "-", "", -1)
		doc := document{DocumentType: CPF.String(), Number: value}
		b.transaction.Customer.Documents = append(b.transaction.Customer.Documents, doc)
		b.transaction.Customer.CustomerType = INDIVIDUAL.String()

		return b, nil
	}

	if regexCNPJ.MatchString(value) {

		value = strings.Replace(value, ".", "", -1)
		value = strings.Replace(value, "-", "", -1)
		value = strings.Replace(value, "/", "", -1)

		doc := document{DocumentType: CNPJ.String(), Number: value}
		b.transaction.Customer.Documents = append(b.transaction.Customer.Documents, doc)
		b.transaction.Customer.CustomerType = CORPORATION.String()

		return b, nil
	}
	return b, &InvalidValueError{"Document.Number", value}
}

type client struct {
	*http.Client
	url string
}

func NewClient() *client {
	return &client{
		new(http.Client),
		BASE_URL,
	}
}

func (c *client) Execute(transaction transaction, authenticationMethod AuthenticationMethod) (*transactionResponse, error) {

	jsonData, _ := transaction.marshal()
	reqUrl := c.url + PATH_TRANSACTION
	req, _ := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	switch authenticationMethod {
	case BODY:
	case PARAM:
		q := req.URL.Query()
		q.Add("api_key", API_KEY)
		req.URL.RawQuery = q.Encode()
	case BASIC_AUTH:
		req.SetBasicAuth(API_KEY, "x")
	}

	reqDump, _ := httputil.DumpRequest(req, true)
	log.Println("##### Request transaction #####")
	log.Println(string(reqDump))

	res, err := c.Do(req)

	if err != nil {
		log.Println("Stone response  error", err.Error())
		return nil, err
	}

	if res.StatusCode == 500 {
		err := InternalError{PATH_TRANSACTION}
		log.Println("Stone response  error", err.Error())
		return nil, &err
	}

	defer res.Body.Close()

	repDump, _ := httputil.DumpResponse(res, true)
	log.Println("##### Response  transaction ##### ")
	log.Println(string(repDump))

	result := transactionResponse{}
	json.NewDecoder(res.Body).Decode(result)

	return &result, nil
}

func (c *client) RecoverPublicKey() publicKey {
	log.Println("Recover Public Key")

	var body = []byte(`{"api_key":"` + API_KEY + `"}`)

	reqUrl := c.url + PATH_HASH
	req, _ := http.NewRequest(http.MethodGet, reqUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	reqDump, _ := httputil.DumpRequest(req, true)

	log.Println("##### Request PUBLIC KEY ######")
	log.Println(string(reqDump))

	res, errorRes := c.Do(req)

	if errorRes != nil {
		log.Fatal("Error request PUBLIC KEY", errorRes.Error())
	}

	resDump, _ := httputil.DumpResponse(res, true)
	log.Println("##### Response Stone PUBLIC KEY #####")
	log.Println(string(resDump))

	if res.StatusCode != 200 {
		log.Fatal("Error request PUBLIC KEY. Status: ", res.StatusCode)
	}

	publicKey := publicKey{}
	json.NewDecoder(res.Body).Decode(&publicKey)

	return publicKey
}

func createRsaPublicKey(value string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("failed to parse DER encoded public key: " + err.Error())
	}

	return pub.(*rsa.PublicKey)
}

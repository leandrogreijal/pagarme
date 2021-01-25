[![Build Status](https://pagar.me/static/logo_pagarme-f40e836118f75338095ebb5b461cd5ed.svg)](https://pagar.me/static/logo_pagarme-f40e836118f75338095ebb5b461cd5ed.svg)


### Installation

Pagar.me requires [GoLang](https://golang.org/) and

On MacOS you can install or upgrade to the latest released version with Homebrew:
```sh
$ brew install dep
$ brew upgrade dep
```
On Debian platforms you can install or upgrade to the latest version with apt-get:
```sh
$ sudo apt-get install go-dep
```
After run de comand in MAKEFILE
```sh
$  make build
```
### Run
```sh
$  ./bin/pagarme
```

```
Stone test

Usage:
  pagarme [command]

Available Commands:
  boleto      Gerar boleto
  cartao      Gerar cobramça cartão
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.pagarme.yaml)
  -h, --help            help for pagarme
  -t, --toggle          Help message for toggle

Use "pagarme [command] --help" for more information about a command.
```

##### Boleto

```sh
$  ./bin/pagarme boleto
```
```
Gerar boleto
Usage:
  pagarme boleto [flags]

Flags:
  -a, --amount float      Amount value
  -d, --document string   Document
  -h, --help              help for boleto
  -n, --name string       Name
```

Exemple:
```
  $  ./bin/pagarme boleto  --amount 33.00 --name Leandro --document 251.854.650-26
```

##### Credit card

```sh
$  ./bin/pagarme cartao
```
```
Gerar cobramça cartão

Usage:
  pagarme cartao [flags]

Flags:
  -a, --amount float                Amount value
  -v, --cardCVV string              Card CVV
  -e, --cardExpirationDate string   Card Expiration Date
  -N, --cardHolderName string       Card Holder Name
  -c, --cardNumber string           Card Number
  -C, --country string              Country
  -d, --document string             Document
  -h, --help                        help for cartao
  -n, --name string                 Name
```

Exemple:
```
  $  ./bin/pagarme cartao --amount 33.00 --name Leandro --document 251.854.650-26 --country br --cardNumber 41111111111111 --cardHolderName Leandro --cardExpirationDate 1028 --cardCVV 123
```

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Accounts struct {
	Accounts []Account `xml:"Account"`
}

type Account struct {
	UserId                 string                 `xml:"UserId"`
	Active                 bool                   `xml:"Active"`
	Code                   string                 `xml:"Code"`
	CreditLimit            float64                `xml:"CreditLimit"`
	RiskMultiplier         float64                `xml:"RiskMultiplier"`
	CollateralAllowed      bool                   `xml:"CollateralAllowed"`
	ShortSellAllowed       bool                   `xml:"ShortSellAllowed"`
	CreditAllowed          bool                   `xml:"CreditAllowed"`
	DefaultCurrency        string                 `xml:"DefaultCurrency"`
	DerivativeLevel        string                 `xml:"DerivativeLevel"`
	AllowedInstrumentTypes AllowedInstrumentTypes `xml:"AllowedInstrumentTypes"`
	SubAccounts            SubAccounts            `xml:"SubAccounts"`
	PropertiesNested       PropertiesNested       `xml:"PropertiesNested"`
}

type AllowedInstrumentTypes struct {
	AllowedInstrumentTypes []AllowedInstrumentType `xml:"AllowedInstrumentType"`
}

type AllowedInstrumentType struct {
	InstrumentType string `xml:"InstrumentType"`
}

type SubAccounts struct {
	SubAccounts []SubAccount `xml:"SubAccount"`
}

type SubAccount struct {
	SubAccountNumber       string                 `xml:"SubAccountNumber"`
	PortfolioMarketValue   float64                `xml:"PortfolioMarketValue"`
	Collateral             float64                `xml:"Collateral"`
	FutureBalance          float64                `xml:"FutureBalance"`
	VariationMargin        float64                `xml:"VariationMargin"`
	CashUsed               bool                   `xml:"CashUsed"`
	CollateralUsed         bool                   `xml:"CollateralUsed"`
	CurrencyAccounts       CurrencyAccounts       `xml:"CurrencyAccounts"`
	Positions              Positions              `xml:"Positions"`
	MarketReferencesNested MarketReferencesNested `xml:"MarketReferencesNested"`
}

type CurrencyAccounts struct {
	CurrencyAccounts []CurrencyAccount `xml:"CurrencyAccount"`
}

type CurrencyAccount struct {
	Currency            string  `xml:"Currency"`
	Balance             float64 `xml:"Balance"`
	Interest            float64 `xml:"Interest"`
	Margin              float64 `xml:"Margin"`
	CurrencyCreditLimit float64 `xml:"CurrencyCreditLimit"`
	ExternalMargin      float64 `xml:"ExternalMargin"`
}

type Positions struct {
	Positions []Position `xml:"Position"`
}

type Position struct {
	ShortName       string  `xml:"ShortName"`
	InstrumentId    string  `xml:"InstrumentId"`
	Volume          float64 `xml:"Volume"`
	Market          string  `xml:"Market"`
	Currency        string  `xml:"Currency"`
	DailyAmount     float64 `xml:"DailyAmount"`
	LoanAmount      float64 `xml:"LoanAmount"`
	MeanValue       float64 `xml:"MeanValue"`
	AverageValue    float64 `xml:"AverageValue"`
	TradingPrice    float64 `xml:"TradingPrice"`
	QuotingCurrency string  `xml:"QuotingCurrency"`
}

type MarketReferencesNested struct {
	MarketReferencesNested []MarketReference `xml:"MarketReference"`
}

type MarketReference struct {
	MarketName string `xml:"MarketName"`
	Reference  string `xml:"Reference"`
}

type PropertiesNested struct {
	PropertiesNested []PropertyNested `xml:"PropertyNested"`
}

type PropertyNested struct {
	Name           string         `xml:"Name"`
	PropertyValues PropertyValues `xml:"PropertyValues"`
}

type PropertyValues struct {
	PropertyValues []string `xml:"PropertyValue"`
}

func main() {
	xmlFile, err := os.Open("cmd/go_import/Accounts1K.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	defer func() {
		if err := xmlFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var accounts Accounts

	if err := xml.Unmarshal(byteValue, &accounts); err != nil {
		panic(err)
	}

	fmt.Println(accounts)
}

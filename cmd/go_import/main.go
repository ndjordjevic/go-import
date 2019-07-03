package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io"
	"log"
	"os"
	"time"
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

var (
	debug         = flag.Bool("debug", true, "enable debugging")
	password      = flag.String("password", "yourStrong(!)Password", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	user          = flag.String("user", "sa", "the database user")
)

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	if *debug {
		fmt.Printf("connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	stmt, err := conn.Prepare("select 1, 'abc'")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var somenumber int64
	var somechars string
	err = row.Scan(&somenumber, &somechars)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("somenumber:%d\n", somenumber)
	fmt.Printf("somechars:%s\n", somechars)

	xmlFile, err := os.Open("cmd/go_import/Accounts1M.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened Accounts1M.xml")
	defer func() {
		if err := xmlFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	decoder := xml.NewDecoder(xmlFile)
	start := time.Now()
	for {
		t, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}

			fmt.Println(tokenErr)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Space == "http://www.front.com" && t.Name.Local == "Account" {
				var account Account
				if err := decoder.DecodeElement(&account, &t); err != nil {
					fmt.Println(err)
				}

				//fmt.Println(account)
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Parsing took %s", elapsed)
}

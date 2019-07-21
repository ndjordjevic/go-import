package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"log"
	"os"
	"sync"
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
	database      = flag.String("database", "SOFTBROKER", "the database")
)

func main() {
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d", *server, *database, *user, *password, *port)
	if *debug {
		fmt.Printf("connString:%s\n", connString)
	}
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open db connection failed:", err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	xmlFile, err := os.Open("cmd/go_import/Accounts100K.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened Accounts100K.xml")
	defer func() {
		if err := xmlFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan *Account)

	stmt, err := db.Prepare("insert into dbo.Accounts (creation_time, modification_time, modification_type, user_id, trading_group_id, " +
		"credit_limit, short_sell_limit, order_value_limit, high_risk_collateral_factor, derivative_limit, risk_multiplier, " +
		"active, collateral_allowed, short_sell_allowed, credit_allowed, code, inactivation_comment, default_currency, " +
		"derivative_level, impersonate_cfd, gross_margin_calculation, cfd_account) " +
		" output inserted.id  values (?, ?, 'SYSTEM', NULL, NULL, ?, 0, 0, 0, 0, ?, ?, ?, ?, ?, ?, NULL, ?, ?, 0, 0, 0)")

	const noOfGoRoutines = 64
	var wg sync.WaitGroup
	wg.Add(100000)

	for g := 0; g < noOfGoRoutines; g++ {
		go func(g int) {
			for account := range ch {
				//fmt.Printf("Enter: %v\n", g)
				//println("Read from a channel code:", account.Code, "address:", &account)
				processAccount(account, stmt)
				//fmt.Printf("Exit: %v\n", g)
				wg.Done()
			}
		}(g)
	}

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var accounts Accounts

	start := time.Now()
	if err := xml.Unmarshal(byteValue, &accounts); err != nil {
		panic(err)
	}

	for i, _ := range accounts.Accounts {
		//processAccount(&a, db)
		//println(a.Code)
		ch <- &accounts.Accounts[i]
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Parsing and inserting into DB took %s", elapsed)
}

func processAccount(account *Account, stmt *sql.Stmt) {
	dt := time.Now()

	//println("Writing to db code: ", account.Code)

	_, err := stmt.Exec(dt, dt, account.CreditLimit, account.RiskMultiplier, account.Active, account.CollateralAllowed, account.ShortSellAllowed, account.CreditAllowed, account.Code, account.DefaultCurrency, account.DerivativeLevel)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

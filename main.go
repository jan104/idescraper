package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/go-resty/resty/v2"
	"github.com/jan104/idescraper/api"
	_ "github.com/lib/pq"
)

var GitCommit string = "undefined"
var BuildTime string = "No Time provided"

type IdeTokenGen struct {
	ACCESSTOKEN string `json:"access_token"`
	TOKENTYPE   string `json:"token_type"`
	EXPIRESEC   int    `json:"expires_in"`
	SCOPE       string `json:"scope"`
	JTI         string `json:"jti"`
}

func generateToken() (token string, err error) {
	fmt.Println("Generated a new token")
	apikey, apikeyfound := os.LookupEnv("APIKEY")
	secret, secretfound := os.LookupEnv("APISECRET")
	if !apikeyfound || !secretfound {
		fmt.Fprintln(os.Stderr, "Env not set")
		os.Exit(1)
	}
	var dump IdeTokenGen
	client := resty.New()
	resp, err := client.R().
		SetResult(&dump).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SetBasicAuth(apikey, secret).
		SetQueryParams(map[string]string{
			"grant_type": "client_credentials",
			"scope":      "read",
		}).
		Post("https://api.idealista.com/oauth/token")
	if err != nil {
		return "", err
	}
	fmt.Println("Success:", resp.IsSuccess(), "Code:", resp.StatusCode())
	return dump.ACCESSTOKEN, nil
}

func fetchIde(authtoken string, idx int) api.IdeResp {
	fmt.Println("Fetch idealista property API")
	var dump api.IdeResp
	filtermap := map[string]string{
		"center":       "28.1204,-16.7243",
		"distance":     "5000",
		"propertyType": "homes",
		"operation":    "sale",
		"country":      "es",
		"locale":       "en",
		"maxItems":     "50",
		"minPrice":     "100000.0",
		"maxPrice":     "850000.0",
		"order":        "modificationDate",
		"sort":         "asc",
		"numPage":      fmt.Sprint(idx),
	}
	fmt.Println("Using these filters:", filtermap)
	client := resty.New()
	resp, err := client.R().
		SetResult(&dump).
		SetAuthToken(authtoken).
		SetHeader("Content-Type", "multipart/form-data;").
		SetFormData(filtermap).
		Post("https://api.idealista.com/3.5/es/search")
	if err != nil {
		panic(err)
	}
	fmt.Println("Success:", resp.IsSuccess(), "Code:", resp.StatusCode())
	return dump
}

func main() {
	fmt.Println("Starting ideScraper routine")
	dialect := goqu.Dialect("postgres")
	token, err := generateToken()
	if err != nil {
		panic(err)
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBDBNAME"))
	pgDB, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 2; i++ {
		ideresp := fetchIde(token, i)
		ds := dialect.Insert("items").Rows(
			ideresp.ElementList,
		)
		insertSQL, _, _ := ds.ToSQL()
		rows, err := pgDB.Query(insertSQL)
		if err != nil {
			panic(err)
		}
		rows.Close()
	}
	fmt.Println("Finished writing to DB")
}

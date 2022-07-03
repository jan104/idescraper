package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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

var filtermap_circSmall = map[string]string{
	// klein: 28.095092,-16.720768 r=3725m
	"center":       "28.095092,-16.720768",
	"distance":     "3725",
	"propertyType": "homes",
	"operation":    "sale",
	"country":      "es",
	"locale":       "en",
	"maxItems":     "50",
	"minPrice":     "200000.0",
	"maxPrice":     "850000.0",
	"order":        "modificationDate",
	"sort":         "desc",
	//"numPage": fmt.Sprint(idx),
	"numPage": "1",
}

var filtermap_circBig = map[string]string{
	// gros: 27.625749,-16.639058 r=50.000m
	"center":       "27.625749,-16.639058",
	"distance":     "50000",
	"propertyType": "homes",
	"operation":    "sale",
	"country":      "es",
	"locale":       "en",
	"maxItems":     "50",
	"minPrice":     "200000.0",
	"maxPrice":     "850000.0",
	"order":        "modificationDate",
	"sort":         "desc",
	//"numPage": fmt.Sprint(idx),
	"numPage": "1",
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

func fetchIde(authtoken string, idx int, filtermap map[string]string) api.IdeResp {
	fmt.Println("Fetch idealista property API")
	var dump api.IdeResp
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
	fmt.Println("Success:", resp.IsSuccess(), "Code:", resp.StatusCode(), "Status:", resp.Status())
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
		ssmall := filtermap_circSmall
		ssmall["numPage"] = fmt.Sprint(i)
		ideresp := fetchIde(token, i, ssmall)
		ds := dialect.Insert("items").Rows(
			ideresp.ElementList,
		)
		insertSQL, _, _ := ds.ToSQL()
		rows, err := pgDB.Query(insertSQL)
		if err != nil {
			panic(err)
		}
		rows.Close()

		time.Sleep(2 * time.Second)

		bbig := filtermap_circBig
		bbig["numPage"] = fmt.Sprint(i)
		ideresp_big := fetchIde(token, i, bbig)
		ds = dialect.Insert("items").Rows(
			ideresp_big.ElementList,
		)
		insertSQL, _, _ = ds.ToSQL()
		rows, err = pgDB.Query(insertSQL)
		if err != nil {
			panic(err)
		}
		rows.Close()
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Finished writing to DB")
	os.Exit(0)
}

// klein: 28.095092,-16.720768 r=3725m
// gros: 27.625749,-16.639058 r=50.000m

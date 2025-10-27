package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/joho/godotenv"

	"HowAreYouFueling/api"
	"HowAreYouFueling/types"
	"HowAreYouFueling/ui"
	"HowAreYouFueling/util"
)

var (
	style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63"))

	purple             = lipgloss.Color("99")
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Align(lipgloss.Center)
	rowStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32")).Padding(0, 1)
	//postcodeInputStyle = style.Width(30)
)

var accessResponse types.AccessResponse

func main() {
	fmt.Println(headerStyle.Render("How Are You Fueling?"))

	// Load keys and secrets from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env file: ", err)
	}

	Auth_header := os.Getenv("AUTH_HEADER")
	apiKey := os.Getenv("API_KEY")
	// -------------------- End loading secrets --------------------

	// Check if previous token has expired
	t := time.Now()
	content, err := os.ReadFile("token_expiry.txt")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Token expiry file does not exist, getting new token...")
			content = []byte("1970-01-01 00:00:00")
		} else {
		fmt.Println("File reading error", err)
		return
		}
	}
	if string(content) > t.Format("2006-01-02 15:04:05") {
		fmt.Println("Token has not expired")
		accessResponse.AccessToken = os.Getenv("ACCESS_TOKEN")
		fmt.Println("Token: ", accessResponse.AccessToken)
	} else {
		fmt.Println("Getting new token...")
		// Get access token. ONLY NEEDS TO RUN WHEN PREVIOUS ONE EXPIRES
		urlToken := "https://api.onegov.nsw.gov.au/oauth/client_credential/accesstoken?grant_type=client_credentials"

		reqToken, _ := http.NewRequest("GET", urlToken, nil)

		reqToken.Header.Add("content-type", "application/json")
		reqToken.Header.Add("authorization", Auth_header)

		resToken, _ := http.DefaultClient.Do(reqToken)

		defer resToken.Body.Close()
		bodyToken, _ := ioutil.ReadAll(resToken.Body)

		if err := json.Unmarshal(bodyToken, &accessResponse); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		file, err := os.Create("token_expiry.txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		issuedAtStr := accessResponse.Issued
		expiresInStr := accessResponse.Expiry

		issuedAtMs, _ := strconv.ParseInt(issuedAtStr, 10, 64)
		expiresInSec, _ := strconv.ParseInt(expiresInStr, 10, 64)

		issuedAt := time.UnixMilli(issuedAtMs)

		expiry := issuedAt.Add(time.Duration(expiresInSec) * time.Second)

		l, err := file.WriteString(expiry.Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}

		envWriteErr := util.UpdateEnvValue(".env", "ACCESS_TOKEN", accessResponse.AccessToken)
		if envWriteErr != nil {
			fmt.Println("Error writing new access token")
		} else {
			fmt.Println(l, "Write success")
		}

		err = file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	p := tea.NewProgram(ui.InitialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	m := finalModel.(ui.Model)
	var selectedBrands []string
	for i := range m.Selected {
		selectedBrands = append(selectedBrands, m.Brands[i])
	}

	payload := types.Payload{
		FuelType:      "U91",
		Brand:         selectedBrands,
		NamedLocation: "2500",
		Reference: struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		}{
			Latitude:  "-33.4362551",
			Longitude: "151.2966549",
		},
		SortBy:        "Price",
		SortAscending: "True",
	}

	resp, err := api.GetFuelPrices(accessResponse.AccessToken, apiKey, payload)
	if err != nil {
		fmt.Println("Error fetching fuel prices:", err)
		return
	}

	var results [][]string
	for i := range resp.Prices {
		results = append(results, []string{
			resp.Stations[i].Name,
			resp.Stations[i].Address,
			resp.Prices[i].FuelType,
			fmt.Sprintf("$%.2f", resp.Prices[i].Price),
		})
	}
	resultTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			default:
				return rowStyle
			}
		}).
		Headers("NAME", "ADDRESS", "FUEL TYPE", "PRICE").
		Rows(results...)

	fmt.Println(resultTable)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

type GraphqlQuery struct {
	Query string `json:"query"`
}

type Contribution struct {
	ContributionCount int    `json:"contributionCount"`
	Date              string `json:"date"`
}

type Contributions struct {
	ContributionCalendar struct {
		Weeks []struct {
			ContributionDays []Contribution `json:"contributionDays"`
		} `json:"weeks"`
	} `json:"contributionCalendar"`
}

type User struct {
	ContributionsCollection Contributions `json:"contributionsCollection"`
}

type GraphqlResponseData struct {
	Data struct {
		User User `json:"user"`
	} `json:"data"`
}

func main() {
	username := os.Getenv("GITHUB_USERNAME")
	query := GraphqlQuery{
		Query: fmt.Sprintf(`
		{
			user(login: "%s") {
				contributionsCollection(from: "2023-01-01T00:00:00Z", to: "2023-12-31T23:59:59Z") {
					contributionCalendar {
						weeks {
							contributionDays {
								contributionCount
								date
							}
						}
					}
				}
			}
		}`, username),
	}

	b, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Authorization", "bearer "+os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data GraphqlResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Define colors
	colors := []*color.Color{
		color.New(color.FgWhite),  // White
		color.New(color.FgGreen),  // Green
		color.New(color.FgBlue),   // Blue
		color.New(color.FgYellow), // Yellow
		color.New(color.FgRed),    // Red
	}

	// Weekly contributions
	for _, week := range data.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		if len(week.ContributionDays) > 0 {
			weeklyContribution := 0
			for _, day := range week.ContributionDays {
				weeklyContribution += day.ContributionCount
			}

			contributionString := ""

			// Construct the string with colored blocks
			for j := 0; j < weeklyContribution; j++ {
				c := colors[(j/10)%len(colors)]
				contributionString += c.Sprint("■ ") // add colored blocks
			}

			fmt.Printf("%s: %s\n", week.ContributionDays[0].Date, contributionString)
		}
	}
}

/* pihole-stats
Copyright (C) 2020 Jeffrey Serio

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//"strconv"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

/*
Pi-hole stats for cli by Jeffrey Serio @hyperreal42 on Github/GitLab
WIP
TODO:
- Handle case in baseURL == "" and authorization == ""
*/

// Basic variables for Pihole instance ---
var (
	baseURL       = os.Getenv("PIHOLE_URL")
	urlSummary    = baseURL + "/api.php?summary"
	urlStatus     = baseURL + "/api.php?status"
	urlEnable     = baseURL + "/api.php?enable"
	urlDisable    = baseURL + "/api.php?disable"
	authorization = "&auth=" + os.Getenv("PIHOLE_AUTH")
)

// Colors ---
var (
	blue      = color.New(color.FgBlue).SprintFunc()
	green     = color.New(color.FgGreen).SprintFunc()
	red       = color.New(color.FgRed).SprintFunc()
	bold      = color.New(color.Bold).SprintFunc()
	underline = color.New(color.Underline).SprintFunc()
)

const usage = `
Pi-hole stats
-------------
This program prints Pi-hole stats to standard output

Usage:
pihole-stats            : Print stats from Pi-hole instance

pihole-stats [e | d]    : Enable (e) or disable (d) Pi-hole`

/*
Data structures ---
*/

// PiholeStats ---
type PiholeStats struct {
	Status                string `json:"status"`
	UniqueClients         string `json:"unique_clients"`
	ClientsEverSeen       string `json:"clients_ever_seen"`
	DomainsBeingBlocked   string `json:"domains_being_blocked"`
	AdsBlockedToday       string `json:"ads_blocked_today"`
	AdsPercentageToday    string `json:"ads_percentage_today"`
	DNSQueriesToday       string `json:"dns_queries_today"`
	QueriesCachedToday    string `json:"queries_cached"`
	QueriesForwardedToday string `json:"queries_forwarded"`
	UniqueDomainsToday    string `json:"unique_domains"`

	GravityLastUpdated struct {
		GravFileExists bool `json:"file_exists"`

		GravRelUp struct {
			Days    int `json:"days"`
			Hours   int `json:"hours"`
			Minutes int `json:"minutes"`
		} `json:"relative"`
	} `json:"gravity_last_updated"`
}

/*
 Helper functions ---
*/

// doRequest --- HTTP GET request to the API
// Returns byte array to be plugged into other functions
func doRequest(url string, auth string) ([]byte, error) {
	newURL := url + auth
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET request failed")
	}

	if newURL == "/api.php?summary&auth=" {
		err := errors.New("Pi-hole url or auth token not supplied")
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get HTTP response")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read HTTP response body")
	}
	return body, nil
}

// getSummary --- returns *PiholeStats instance
func getSummary(jsonKey []byte) (*PiholeStats, error) {

	data := &PiholeStats{}
	if err := json.Unmarshal(jsonKey, data); err != nil {
		return nil, err
	}
	return data, nil
}

func enablePihole() error {
	res, err := doRequest(urlEnable, authorization)
	if err != nil {
		return errors.Wrap(err, "Failed to enable Pi-hole")
	}
	fmt.Println(string(res))
	return nil
}

func disablePihole() error {
	res, err := doRequest(urlDisable, authorization)
	if err != nil {
		return errors.Wrap(err, "Failed to disable Pi-hole")
	}
	fmt.Println(string(res))
	return nil
}

func printContent() error {
	content, err := doRequest(urlSummary, authorization)
	if err != nil {
		return errors.Wrap(err, "Failed to run HTTP request for Pi-hole stats")
	}
	data, err := getSummary(content)
	if err != nil {
		return errors.Wrap(err, "Failed to get Pi-hole stats summary")
	}

	fmt.Printf("%s\n\n", bold(underline(red("Pi-hole Statistics"))))
	fmt.Printf("Pi-hole admin console: %s\n", baseURL)
	if data.Status == "enabled" {
		fmt.Printf("Status: %s\n", green("Enabled"))
	} else {
		fmt.Printf("Status: %s\n", red("Disabled"))
	}

	g := data.GravityLastUpdated
	if g.GravFileExists == true {
		fmt.Printf("Gravity last updated: %d days, %d hours, %d minutes\n", g.GravRelUp.Days, g.GravRelUp.Hours, g.GravRelUp.Minutes)
	} else {
		fmt.Println("Gravity has not been updated yet")
	}
	fmt.Printf("%s\n", blue("---"))

	dataMap := map[string]string{
		"Current unique clients":  data.UniqueClients,
		"Total clients ever seen": data.ClientsEverSeen,
		"Domains being blocked":   data.DomainsBeingBlocked,
		"Ads blocked today":       data.AdsBlockedToday,
		"Ads percentage today":    data.AdsPercentageToday,
		"DNS queries today":       data.DNSQueriesToday,
		"Queries cached today":    data.QueriesCachedToday,
		"Queries forwarded today": data.QueriesForwardedToday,
		"Unique domains today":    data.UniqueDomainsToday,
	}

	for i, j := range dataMap {
		fmt.Printf("%s: %s\n", i, j)
	}
	fmt.Printf("%s\n", blue("---"))

	return nil
}

func main() {
	if len(os.Args) < 2 {
		if err := printContent(); err != nil {
			log.Fatal(err)
		}
	} else {
		switch os.Args[1] {
		case "e":
			if err := enablePihole(); err != nil {
				log.Fatal(err)
			}
		case "d":
			if err := disablePihole(); err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println(usage)
		}
	}
}

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
	"strconv"

	"github.com/fatih/color"
)

/*
Pi-hole stats for cli by Jeffrey Serio @hyperreal42 on Github/GitLab
WIP
TODO:
*/

// Basic variables for Pihole instance ---
var (
	baseURL       = os.Getenv("PIHOLE_STATS_URL")
	urlSummary    = baseURL + "/api.php?summary"
	urlStatus     = baseURL + "/api.php?status"
	urlEnable     = baseURL + "/api.php?enable"
	urlDisable    = baseURL + "/api.php?disable"
	authorization = "&auth=" + os.Getenv("PIHOLE_STATS_AUTH")
)

// Colors ---
var (
	blue      = color.New(color.FgBlue).SprintFunc()
	green     = color.New(color.FgGreen).SprintFunc()
	red       = color.New(color.FgRed).SprintFunc()
	bold      = color.New(color.Bold).SprintFunc()
	underline = color.New(color.Underline).SprintFunc()
)

/*
Data structures ---
*/

// PiholeStats ---
type PiholeStats struct {
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
			Days    string `json:"days"`
			Hours   string `json:"hours"`
			Minutes string `json:"minutes"`
		} `json:"relative"`
	} `json:"gravity_last_updated"`
}

// piholeStatus -- Enabled or disabled
type piholeStatus struct {
	Status string `json:"status"`
}

/*
 Helper functions ---
*/

// errCheck --- check if err != nil
func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// doRequest --- HTTP GET request to the API
// Returns byte array to be plugged into other functions
func doRequest(u string, a string) []byte {
	newURL := u + a
	req, err := http.NewRequest("GET", newURL, nil)
	errCheck(err)

	res, err := http.DefaultClient.Do(req)
	errCheck(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	errCheck(err)
	return body
}

// getSummary --- returns *PiholeStats instance
func getSummary(jsonKey []byte) (*PiholeStats, error) {

	data := &PiholeStats{}
	if err := json.Unmarshal(jsonKey, data); err != nil {
		return nil, err
	}
	return data, nil
}

// getStatus --- returns *piholeStatus instance
func getStatus(jsonKey []byte) (*piholeStatus, error) {
	status := &piholeStatus{}
	if err := json.Unmarshal(jsonKey, status); err != nil {
		return nil, err
	}
	return status, nil
}

func enablePihole() {
	_ = doRequest(urlEnable, authorization)
}

func disablePihole() {
	_ = doRequest(urlDisable, authorization)
}

func getContent() {
	content := doRequest(urlSummary, authorization)
	data, err := getSummary(content)
	errCheck(err)

	statusReq := doRequest(urlStatus, authorization)
	status, err := getStatus(statusReq)
	errCheck(err)

	fmt.Printf("%s\n\n", bold(underline(red("Pi-hole Statistics"))))
	fmt.Printf("Pi-hole admin console: %s\n", baseURL)
	if status.Status == "enabled" {
		fmt.Printf("Status: %s\n", green("Enabled"))
	} else {
		fmt.Printf("Status: %s\n", red("Disabled"))
	}

	g := data.GravityLastUpdated
	if g.GravFileExists == true {
		gDays, _ := strconv.Atoi(g.GravRelUp.Days)
		gHours, _ := strconv.Atoi(g.GravRelUp.Hours)
		gMins, _ := strconv.Atoi(g.GravRelUp.Minutes)
		fmt.Printf("Gravity last updated: %d days, %d hours, %d minutes\n", gDays, gHours, gMins)
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
}

func printUsage() {
	fmt.Printf("USAGE:\n%s {e|d}\t", os.Args[0])
	fmt.Println("Enable or disable Pi-hole")
	fmt.Printf("%s\t", os.Args[0])
	fmt.Println("Get Pi-hole stats")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		getContent()
	} else {
		switch os.Args[1] {
		case "e":
			enablePihole()
		case "d":
			disablePihole()
		default:
			printUsage()
		}
	}
}

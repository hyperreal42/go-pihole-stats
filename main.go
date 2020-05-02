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
	"reflect"

	"github.com/fatih/color"
)

/*
Pi-hole stats for cli by Jeffrey Serio @hyperreal42 on Github/GitLab
WIP
TODO:
* Printf values with colored output in main()
* Implement command-line argument handling
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
	magenta   = color.New(color.FgMagenta).SprintFunc()
	bold      = color.New(color.Bold).SprintFunc()
	underline = color.New(color.Underline).SprintFunc()
)

/*
Data structures ---
*/

// PiholeStats ---
type PiholeStats struct {
	UniqueClients         string      `json:"unique_clients"`
	ClientsEverSeen       string      `json:"clients_ever_seen"`
	GravityLastUpdated    *gravLastUp `json:"gravity_last_updated"`
	DomainsBeingBlocked   string      `json:"domains_being_blocked"`
	AdsBlockedToday       string      `json:"ads_blocked_today"`
	AdsPercentageToday    string      `json:"ads_percentage_today"`
	DNSQueriesToday       string      `json:"dns_queries_today"`
	QueriesCachedToday    string      `json:"queries_cached"`
	QueriesForwardedToday string      `json:"queries_forwarded"`
	UniqueDomainsToday    string      `json:"unique_domains"`
}

// gravLastUp --- Get gravity last updated timestamp info
type gravLastUp struct {
	GravFileExists bool      `json:"file_exists"`
	GravRelUp      *relUnits `json:"relative"`
}

// relUnits --- Time units since t_0 relative to common human standards
type relUnits struct {
	Days    string `json:"days"`
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

// piholeStatus --- Enabled or disabled
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

// getRelUnits --- Returns *relUnits instance
func getRelUnits(jsonKey []byte) (*relUnits, error) {
	relUnits := &relUnits{}
	if err := json.Unmarshal(jsonKey, relUnits); err != nil {
		return nil, err
	}
	return relUnits, nil
}

// getGravUptime --- Returns *gravLastUp instance
func getGravUptime(jsonKey []byte) (*gravLastUp, error) {
	ru, err := getRelUnits(jsonKey)
	errCheck(err)
	gravLastUp := &gravLastUp{
		GravRelUp: ru,
	}
	if err := json.Unmarshal(jsonKey, gravLastUp); err != nil {
		return nil, err
	}
	return gravLastUp, nil
}

// getSummary --- returns *PiholeStats instance
func getSummary(jsonKey []byte) (*PiholeStats, error) {
	gravUp, err := getGravUptime(jsonKey)
	errCheck(err)
	data := &PiholeStats{
		GravityLastUpdated: gravUp,
	}
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

// toggleStatus --- enable/disable Pi-hole
func toggleStatus(url string) {
	statusReq := doRequest(urlStatus, authorization)
	status, err := getStatus(statusReq)
	errCheck(err)
	var req []byte

	if status.Status == "enabled" && url == urlDisable {
		req = doRequest(urlDisable, authorization)
	} else if status.Status == "disabled" && url == urlEnable {
		req = doRequest(urlEnable, authorization)
	}
	status, err = getStatus(req)
	fmt.Printf("Pi-hole status: %s\n", status.Status)
}

func enumerateContent(data interface{}) []interface{} {
	v := reflect.ValueOf(data)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}

func getContent() {
	content := doRequest(urlSummary, authorization)
	data, err := getSummary(content)
	errCheck(err)

	statusReq := doRequest(urlStatus, authorization)
	status, err := getStatus(statusReq)
	errCheck(err)

	fmt.Printf("%s\n", red("Pi-hole Statistics"))
	fmt.Printf("Pi-hole admin console: %s\n", baseURL)
	fmt.Printf("Status: %s\n", green(status.Status))
	fmt.Printf("%s\n", magenta("---"))

	fmt.Println(enumerateContent(*data))
}

func main() {
	/* content := doRequest(urlSummary, authorization)
	data, err := getSummary(content)
	errCheck(err)
	fmt.Println(data.UniqueClients)
	fmt.Println(data.GravityLastUpdated.GravFileExists)
	fmt.Println(*data.GravityLastUpdated.GravRelUp)
	statusReq := doRequest(urlStatus, authorization)
	status, err := getStatus(statusReq)
	errCheck(err)
	fmt.Println(status.Status) */

	/* 	args := os.Args
	   	if len(args) > 2 || args[1] == "help" {
	   		printUsage()
	   	} else if args[1] == "enable" || args[1] == "disable" {
	   		toggleStatus(args[1])
	   	} else {
	   		getContent()
	   	} */
	getContent()

}

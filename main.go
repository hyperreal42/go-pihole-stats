// pihole-stats
// Copyright (C) 2020 Jeffrey Serio

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// Pi-hole stats for cli by Jeffrey Serio @hyperreal42 on Github/GitLab
// WIP
// TODO:
// * Printf values with colored output in main()
// * Implement getStatus func
// * Implement enable/disable func

// Basic variables for Pihole instance
var (
	baseURL              = os.Getenv("PIHOLE_STATS_URL")
	urlSummary           = baseURL + "/api.php?summary"
	cookie        string = os.Getenv("PIHOLE_STATS_COOKIE")
	authorization string = os.Getenv("PIHOLE_STATS_AUTH")
)

// Colors
var (
	Blue      = color.New(color.FgBlue).SprintFunc()
	Green     = color.New(color.FgGreen).SprintFunc()
	Red       = color.New(color.FgRed).SprintFunc()
	Magenta   = color.New(color.FgMagenta).SprintFunc()
	Bold      = color.New(color.Bold).SprintFunc()
	Underline = color.New(color.Underline).SprintFunc()
)

// Data structures

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

type gravLastUp struct {
	GravFileExists bool      `json:"file_exists"`
	GravRelUp      *relUnits `json:"relative"`
}

type relUnits struct {
	Days    string `json:"days"`
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

// Helper functions
func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doRequest(urlSummary string, c string, a string) []byte {
	req, err := http.NewRequest("GET", urlSummary, nil)
	errCheck(err)

	req.Header.Add("cookie", c)
	req.Header.Add("authorization", a)

	res, err := http.DefaultClient.Do(req)
	errCheck(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	errCheck(err)
	return body
}

func getRelUnits(jsonKey []byte) (*relUnits, error) {
	relUnits := &relUnits{}
	if err := json.Unmarshal(jsonKey, relUnits); err != nil {
		return nil, err
	}
	return relUnits, nil
}

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

// func getStatus()
// func toggleEnable()

func main() {
	content := doRequest(urlSummary, cookie, authorization)
	data, err := getSummary(content)
	errCheck(err)
	fmt.Println(data.UniqueClients)
	fmt.Println(data.GravityLastUpdated.GravFileExists)
	fmt.Println(data.GravityLastUpdated.GravRelUp)
}

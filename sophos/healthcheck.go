package sophos

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type HealthCheck struct {
	TenantId   string
	TenantName string
	Endpoint   struct {
		Protection struct {
			Computer struct {
				Total        int `json:"total"`
				NotProtected int `json:"notFullyProtected"`
			} `json:"computer"`
			Server struct {
				Total        int `json:"total"`
				NotProtected int `json:"notFullyProtected"`
			} `json:"server"`
		} `json:"protection"`
		Policy struct {
			Computer struct {
				ThreatProtection struct {
					Total          int              `json:"total"`
					NotRecommended int              `json:"notOnRecommended"`
					Policies       []HealthPolicies `json:"policies,omitempty"`
				} `json:"threat-protection"`
			} `json:"computer"`
			Server struct {
				ThreatProtection struct {
					Total          int              `json:"total"`
					NotRecommended int              `json:"notOnRecommended"`
					Policies       []HealthPolicies `json:"policies,omitempty"`
				} `json:"server-threat-protection"`
			} `json:"server"`
		} `json:"policy"`
		Exclusions struct {
			Policy struct {
				Computer struct {
					Total         int                      `json:"total"`
					NumberOfRisks int                      `json:"numberOfSecurityRisks"`
					Exclusions    []HealthDeviceExclusions `json:"exclusions,omitempty"`
				} `json:"computer"`
				Server struct {
					Total         int                      `json:"total"`
					NumberOfRisks int                      `json:"numberOfSecurityRisks"`
					Exclusions    []HealthDeviceExclusions `json:"exclusions,omitempty"`
				} `json:"server"`
			} `json:"policy"`
			Global struct {
				NumberOfRisks   int                `json:"numberOfSecurityRisks"`
				LockedByAccount bool               `json:"lockedByManagingAccount"`
				Exclusions      []HealthExclusions `json:"scanningExclusions,omitempty"`
			} `json:"global"`
		} `json:"exclusions"`
		TamperProtection struct {
			Computer struct {
				Total    int `json:"total"`
				Disabled int `json:"disabled"`
			} `json:"computer"`
			Server struct {
				Total    int `json:"total"`
				Disabled int `json:"disabled"`
			} `json:"server"`
			Global bool `json:"global"`
		} `json:"tamperProtection"`
	} `json:"endpoint"`
}

type HealthPolicies struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	LockedByAccount bool   `json:"lockedByManagingAccount"`
	NotRecommended  int    `json:"notOnRecommended"`
}

type HealthDeviceExclusions struct {
	PolicyId           string             `json:"policyId"`
	PolicyName         string             `json:"policyName"`
	LockedByAccount    bool               `json:"lockedByManagingAccount"`
	ScanningExclusions []HealthExclusions `json:"scanningExclusions"`
}

type HealthExclusions struct {
	Id         string `json:"id"`
	Value      string `json:"value"`
	Type       string `json:"type"`
	ScanMode   string `json:"scanMode"`
	ReasonCode string `json:"reasonCode"`
}

var HealthChecks []HealthCheck
var timestamp time.Time

func GetHeathCheck(jwt string) {
	const rateLimit = 100 // Number of requests per minute
	timestamps := make([]time.Time, 0)

	timestamp = time.Now()
	total_tenants := len(Tenants)
	bar := progressbar.Default(int64(total_tenants), "Running health checks on tenants...")
	for _, tenant := range Tenants {
		now := time.Now()
		// Filter out timestamps older than 1 minute
		recentTimestamps := make([]time.Time, 0)
		for _, t := range timestamps {
			if now.Sub(t) < time.Minute {
				recentTimestamps = append(recentTimestamps, t)
			}
		}
		timestamps = recentTimestamps
		// If rate limit is reached, sleep for the necessary duration
		if len(timestamps) >= rateLimit {
			sleepDuration := time.Minute - now.Sub(timestamps[0])
			time.Sleep(sleepDuration)
			now = time.Now() // Update current time after sleep
		}
		// Add the current timestamp to the slice
		timestamps = append(timestamps, now)

		// Arrays start at 0, add + 1 for clenliness, otherwise it will say 0 of 1, 1 of 1, 2 of 1 etc.
		// I'd like this to be updated with a charmed UI, but for now this will do.
		bar.Add(1)
		// fmt.Printf("Processed %d of %d \n", i+1, len(Tenants))
		// Build string for URL
		uri := tenant.ApiHost
		url := uri + "/account-health-check/v1/health-check"
		// Create the Bearer Token
		var bearer = "Bearer " + jwt
		req, _ := http.NewRequest("GET", url, nil)
		// Add Headers to the request
		req.Header.Add("Authorization", bearer)
		req.Header.Add("X-Tenant-ID", tenant.Id)
		// Create a Client and do the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)

		}
		defer resp.Body.Close()
		// If everything is ok.
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
			}
			var healthcheck *HealthCheck
			err = json.Unmarshal(body, &healthcheck)
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
			}
			healthcheck.TenantId = tenant.Id
			healthcheck.TenantName = tenant.Name
			HealthChecks = append(HealthChecks, *healthcheck)

		} else if resp.StatusCode == 429 {
			// If we get a 429, we need to wait 2 seconds and try again. to get around rate limiting
			fmt.Println("Too many requests")
			time.Sleep(2 * time.Second)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
			}
			var healthcheck *HealthCheck
			err = json.Unmarshal(body, &healthcheck)
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
			}
			healthcheck.TenantId = tenant.Id
			healthcheck.TenantName = tenant.Name
			HealthChecks = append(HealthChecks, *healthcheck)
		} else {
			// Anything else, log it to a file.
			f, err := os.OpenFile(timestamp.Format("2006_01_02_15_04")+"_error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
			log.SetOutput(f)
			log.Println(tenant.Id, tenant.Name, resp.StatusCode, resp.Status)
		}
	}
}

func WriteCSV() {
	// Check to see if healthcheck already exists
	if CheckFileExists("healthcheck.csv") {
		fileName := timestamp.Format("2006_01_02_15_04") + "_healthcheck.csv"
		os.Rename("healthcheck.csv", fileName)
	}

	// Create a file
	csvFile, err := os.Create("healthcheck.csv")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// Create a writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()
	// Write the headers
	row := []string{
		"TenantId",
		"TenantName",
		"Endpoints_Total_Not_Protected",
		"Server_Total_Not_Protected",
		"Computer_Total_Policy_Issues",
		"Computer_Trouble_Policies",
		"Server_Total_Policy_Issues",
		"Server_Trouble_Policies",
		"Computer_Total_Tamper_Protection_Disabled",
		"Server_Total_Tamper_Protection_Disabled",
		"Global_Tamper_Protection_Enabled",
		"Computer_Total_High_Risk_Exclusions",
		"Computer_High_Risk_Exclusions_Name",
		"Server_Total_High_Risk_Exclusions",
		"Server_High_Risk_Exclusions_Name",
		"Global_Total_High_Risk_Exclusions",
	}
	if err := csvWriter.Write(row); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	// Write the data
	for _, healthcheck := range HealthChecks {
		// Handing of Nested Arrays
		var endpointTroublePolicy []string
		for _, items := range healthcheck.Endpoint.Policy.Computer.ThreatProtection.Policies {
			endpointTroublePolicy = append(endpointTroublePolicy, items.Name)
		}
		if len(endpointTroublePolicy) == 0 {
			endpointTroublePolicy = append(endpointTroublePolicy, "None")
		}
		// Handing of Nested Arrays
		var serverTroublePolicy []string
		for _, items := range healthcheck.Endpoint.Policy.Server.ThreatProtection.Policies {
			serverTroublePolicy = append(serverTroublePolicy, items.Name)
		}
		if len(serverTroublePolicy) == 0 {
			serverTroublePolicy = append(serverTroublePolicy, "None")
		}
		// Handing of Nested Arrays
		var computerExclusionName []string
		for _, items := range healthcheck.Endpoint.Exclusions.Policy.Computer.Exclusions {
			computerExclusionName = append(computerExclusionName, items.PolicyName)
		}
		if len(computerExclusionName) == 0 {
			computerExclusionName = append(computerExclusionName, "None")
		}
		// Handing of Nested Arrays
		var serverExclusionName []string
		for _, items := range healthcheck.Endpoint.Exclusions.Policy.Server.Exclusions {
			serverExclusionName = append(serverExclusionName, items.PolicyName)
		}
		if len(serverExclusionName) == 0 {
			serverExclusionName = append(serverExclusionName, "None")
		}
		// Build Row
		if healthcheck.Endpoint.Protection.Computer.NotProtected > 0 ||
			healthcheck.Endpoint.Protection.Server.NotProtected > 0 ||
			healthcheck.Endpoint.Policy.Computer.ThreatProtection.NotRecommended > 0 ||
			healthcheck.Endpoint.Policy.Server.ThreatProtection.NotRecommended > 0 ||
			healthcheck.Endpoint.TamperProtection.Computer.Disabled > 0 ||
			healthcheck.Endpoint.TamperProtection.Server.Disabled > 0 ||
			healthcheck.Endpoint.Exclusions.Policy.Computer.NumberOfRisks > 0 ||
			healthcheck.Endpoint.Exclusions.Policy.Server.NumberOfRisks > 0 ||
			healthcheck.Endpoint.Exclusions.Global.NumberOfRisks > 0 {
			row := []string{
				healthcheck.TenantId,
				healthcheck.TenantName,
				strconv.Itoa(healthcheck.Endpoint.Protection.Computer.NotProtected),
				strconv.Itoa(healthcheck.Endpoint.Protection.Server.NotProtected),
				strconv.Itoa(healthcheck.Endpoint.Policy.Computer.ThreatProtection.NotRecommended),
				strings.Join(endpointTroublePolicy, "/"),
				strconv.Itoa(healthcheck.Endpoint.Policy.Server.ThreatProtection.NotRecommended),
				strings.Join(serverTroublePolicy, "/"),
				strconv.Itoa(healthcheck.Endpoint.TamperProtection.Computer.Disabled),
				strconv.Itoa(healthcheck.Endpoint.TamperProtection.Server.Disabled),
				strconv.FormatBool(healthcheck.Endpoint.TamperProtection.Global),
				strconv.Itoa(healthcheck.Endpoint.Exclusions.Policy.Computer.NumberOfRisks),
				strings.Join(computerExclusionName, "/"),
				strconv.Itoa(healthcheck.Endpoint.Exclusions.Policy.Server.NumberOfRisks),
				strings.Join(serverExclusionName, "/"),
				strconv.Itoa(healthcheck.Endpoint.Exclusions.Global.NumberOfRisks),
			}
			err := csvWriter.Write(row)
			if err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	}
}

func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Removing for now, rate limiting prevents the use of concurrent calls.
// func StartHealthCheck(jwt string) {
// 	var wg sync.WaitGroup
// 	wg.Add(len(Tenants))
// 	for _, tenant := range Tenants {
// 		go GetHeathCheck(jwt, tenant.ApiHost, tenant.Id, tenant.Name, &wg)
// 	}
// 	wg.Wait()

// }

package sophos

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type TenantDetails struct {
	Items []struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		ShowAs        string `json:"showAs"`
		DataRegion    string `json:"dataRegion"`
		DataGeography string `json:"dataGeography"`
		BillingType   string `json:"billingType"`
		ApiHost       string `json:"apiHost"`
		Status        string `json:"status"`
	} `json:"items"`
}

type Tenant struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	ShowAs        string `json:"showAs"`
	DataRegion    string `json:"dataRegion"`
	DataGeography string `json:"dataGeography"`
	BillingType   string `json:"billingType"`
	ApiHost       string `json:"apiHost"`
	Status        string `json:"status"`
}

type TenantDetailsPages struct {
	Pages struct {
		CurrentPage int `json:"current"`
		Size        int `json:"size"`
		Total       int `json:"total"`
		Items       int `json:"items"`
		MaxSize     int `json:"maxSize"`
	} `json:"pages"`
}

var Tenants []Tenant

func GetTenants(jwt string, partnerid string, tenanttype string) {
	// This function is to get the list of total tenants, then in GetSubPageTenants we get all tenant info
	url := "https://api.central.sophos.com/" + tenanttype + "/v1/tenants?pageTotal=true"
	var bearer = "Bearer " + jwt
	var partner = partnerid
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("X-Partner-ID", partner)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	var tenantdetailpages TenantDetailsPages
	json.Unmarshal(body, &tenantdetailpages)
	page := tenantdetailpages.Pages.Total
	// Create a wait group
	var wg sync.WaitGroup
	wg.Add(page)
	for i := 1; i <= page; i++ {
		go GetSubPageTenants(jwt, partnerid, tenanttype, i, &wg)
	}
	wg.Wait()
}

func GetSubPageTenants(jwt string, partnerid string, tenanttype string, page int, wg *sync.WaitGroup) {
	// Build the URL
	url := "https://api.central.sophos.com/" + tenanttype + "/v1/tenants?page=" + strconv.Itoa(page)
	// Create the bearer Token
	var bearer = "Bearer " + jwt
	var partner = partnerid
	req, _ := http.NewRequest("GET", url, nil)
	// Add headers
	req.Header.Add("Authorization", bearer)
	req.Header.Add("X-Partner-ID", partner)
	// Create Client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	var tenantdetails *TenantDetails
	err = json.Unmarshal(body, &tenantdetails)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	for _, row := range tenantdetails.Items {
		if row.ApiHost != "" || row.Status == "Active" {
			Tenants = append(Tenants, row)
		}
	}
	wg.Done()
}

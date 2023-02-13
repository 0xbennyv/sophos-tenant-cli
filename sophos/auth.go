package sophos

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type OAuthJWT struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type WhoAmIDetails struct {
	Id      string `json:"id"`
	IdType  string `json:"idType"`
	ApiHost struct {
		Global string `json:"global"`
	} `json:"apiHosts"`
}

func OAuthToken() OAuthJWT {
	client_id := viper.Get("client_id").(string)
	client_secret := viper.Get("client_secret").(string)
	grant_type := "client_credentials"
	scope := "token"
	url := "https://id.sophos.com/api/v2/oauth2/token"
	payload := strings.NewReader("client_id=" + client_id + "&client_secret=" + client_secret + "&grant_type=" + grant_type + "&scope=" + scope)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", payload)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	var oauthjwt OAuthJWT
	err = json.Unmarshal(body, &oauthjwt)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	return oauthjwt
}

func Whoami(jwt string) WhoAmIDetails {
	url := "https://api.central.sophos.com/whoami/v1"
	var bearer = "Bearer " + jwt
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
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
	var whoamidetails WhoAmIDetails
	err = json.Unmarshal(body, &whoamidetails)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	return whoamidetails

}

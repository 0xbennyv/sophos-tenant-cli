package sophos

import (
	"fmt"
)

func Execute() {
	oauthjwt := OAuthToken()
	whoami := Whoami(oauthjwt.AccessToken)
	GetTenants(oauthjwt.AccessToken, whoami.Id, whoami.IdType)
	fmt.Println("Total Tenants: ", len(Tenants))
	GetHeathCheck(oauthjwt.AccessToken)
	WriteCSV()

}

package sophos

func Execute() {
	oauthjwt := OAuthToken()
	whoami := Whoami(oauthjwt.AccessToken)
	GetTenants(oauthjwt.AccessToken, whoami.Id, whoami.IdType)
	GetHeathCheck(oauthjwt.AccessToken)
	WriteCSV()

}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sophos_tenant_cli/cmd"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("")
		fmt.Println("    .:1010101010101010101:.")
		fmt.Println("   :101010101010101010101010")
		fmt.Println(" .1010101010101010101010101:")
		fmt.Println(" :1010101'::::::::::::::::'")
		fmt.Println(" 1010101:")
		fmt.Println(" :1010101`...............")
		fmt.Println("  10101010101010101010101'")
		fmt.Println("   `10101010101010101010101.")
		fmt.Println("     `:101010101010101010101,")
		fmt.Println("                    `01010101   +++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("                     :1010101   + Being first run time please set your  +")
		fmt.Println("  .::::::::::::::::::1010101:   + Client Id and Client Secret generated +")
		fmt.Println("  :1010101010101010101010101    + from your Enterprise Dashboard or     +")
		fmt.Println("  :10101010101010101010101'     + or Partner Dashboard                  +")
		fmt.Println("   `01010101010101010101'       +++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("")
		for {
			consoleReader := bufio.NewReader(os.Stdin)
			fmt.Print("Client ID: ")
			client_id, _ := consoleReader.ReadString('\n')
			client_id = strings.Replace(client_id, "\n", "", -1)
			client_id = strings.Replace(client_id, "\r", "", -1)
			client_id = strings.Replace(client_id, "\t", "", -1)
			fmt.Print("Client Secret: ")
			client_secret, _ := consoleReader.ReadString('\n')
			client_secret = strings.Replace(client_secret, "\n", "", -1)
			client_secret = strings.Replace(client_secret, "\r", "", -1)
			client_secret = strings.Replace(client_secret, "\t", "", -1)
			viper.SetDefault("client_id", client_id)
			viper.SetDefault("client_secret", client_secret)
			viper.WriteConfigAs("config.json")
			os.Exit(0)

		}
	}
	cmd.Execute()
}

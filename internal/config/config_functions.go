package config

import (
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/lithammer/shortuuid"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
)

func ReadConfig() (config structs.ModuleConfig) {
	hostname := os.Getenv("HOST_DOMAIN")
	token := viper.GetString("url-token")

	configElements := []structs.Element{
		{
			Label:            "Requirements",
			Type:             7,
			ExpectedJsonName: "",
			Rationale:        "An existing AWS account with Amazon GuardDuty activated.\n\nClick the Help icon for further information on how to configure this module.",
			Value:            "",
			PossibleValues:   nil,
			Required:         false,
		}, {
			Label:            "Elements Imported",
			Type:             7,
			ExpectedJsonName: "",
			Rationale:        "IP Addresses",
			Value:            "",
			PossibleValues:   nil,
			Required:         false,
		}, {
			Label:            "AWS Lambda Push URL",
			Type:             6,
			ExpectedJsonName: "",
			Rationale:        "Use the URL below in the AWS Lambda function. FQDN and port must resolvable and accessible from Internet.",
			Value:            fmt.Sprintf("https://%s:9000/ingress/awsgd/run?token=%s", hostname, token),
			PossibleValues:   nil,
			Required:         true,
		},
	}

	config.Fields = configElements

	return
}

func CreateAndSetUrlAuthToken() {
	if !viper.IsSet("url-token") {
		token := shortuuid.New()
		viper.Set("url-token", token)

		if err := viper.WriteConfig(); err != nil {
			log.Error("Error writing last run time to config file", err)
		}
	}
}

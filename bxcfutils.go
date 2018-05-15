package bxcfutils

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

type VCAPService struct {
	Name        string                    `json:"name"`
	Label       string                    `json:"label,omitempty"`
	Plan        string                    `json:"plan,omitempty"`
	Credentials MessageHubVcapCredentials `json:"credentials"`
}

type MessageHubVcapCredentials struct {
	InstanceID       string   `json:"instance_id"`
	MqLightLookupURL string   `json:"mqlight_lookup_url"`
	APIKey           string   `json:"api_key"`
	KafkaAdminURL    string   `json:"kafka_admin_url"`
	KafkaRestURL     string   `json:"kafka_rest_url"`
	KafkaBrokerSasl  []string `json:"kafka_brokers_sasl"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
}

func GetCurrEnvVCAPServices() (*cfenv.App, error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &appEnv.Services)
	return appEnv, nil
}

func GetBluemixMessageHubCredentials(name, plan string) (MessageHubVcapCredentials, error) {
	vcapServices := os.Getenv("VCAP_SERVICES")
	if len(vcapServices) == 0 {
		return MessageHubVcapCredentials{}, errors.New("vcapServices undefined")
	}
	var vcap map[string][]VCAPService
	err := json.Unmarshal([]byte(vcapServices), &vcap)
	if err != nil {
		return MessageHubVcapCredentials{}, errors.New("failed to parse vcapServices " + err.Error())
	}
	for vname, vservice := range vcap {
		if !strings.HasPrefix(vname, name) {
			continue
		}
		for i := range vservice {
			if len(plan) == 0 || plan == vservice[i].Plan {
				creds := MessageHubVcapCredentials{
					InstanceID:       vservice[i].Credentials.InstanceID,
					MqLightLookupURL: vservice[i].Credentials.MqLightLookupURL,
					APIKey:           vservice[i].Credentials.APIKey,
					KafkaAdminURL:    vservice[i].Credentials.KafkaAdminURL,
					KafkaRestURL:     vservice[i].Credentials.KafkaRestURL,
					KafkaBrokerSasl:  vservice[i].Credentials.KafkaBrokerSasl,
					User:             vservice[i].Credentials.User,
					Password:         vservice[i].Credentials.Password,
				}
				return creds, nil
			}
		}
	}
	return MessageHubVcapCredentials{}, errors.New("service instance not found in vcapServices")
}

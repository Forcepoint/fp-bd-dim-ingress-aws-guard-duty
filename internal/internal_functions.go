package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
)

func Register() {
	_, err := MakeRequest("POST", "register", buildMetadata())

	if err != nil {
		log.Error(err)
	}
}

func PushDataToController(data structs.IncomingItem) (*io.ReadCloser, int, error) {
	items := structs.ProcessedItemsWrapper{
		Items: []structs.ProcessedItem{{
			Source:      "Amazon GuardDuty",
			ServiceName: os.Getenv("MODULE_SVC_NAME"),
			Type:        "IP",
			Value:       data.RemoteIp,
		}},
	}

	log.Info(fmt.Sprintf("Processing %d item from %s", len(items.Items), os.Getenv("MODULE_SVC_NAME")))

	resp, err := MakeRequest("POST", "queue", items)

	if err != nil {
		log.Error(err)
		return nil, http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	return &resp.Body, resp.StatusCode, nil
}

func MakeRequest(httpMethod, internalEndpoint string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	registerUrl := fmt.Sprintf("http://%s:%s/internal/%s", os.Getenv("CONTROLLER_SVC_NAME"), os.Getenv("CONTROLLER_PORT"), internalEndpoint)

	req, err := http.NewRequest(httpMethod, registerUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	token := os.Getenv("INTERNAL_TOKEN")

	req.Header.Set("x-internal-token", token)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func buildMetadata() structs.ModuleMetadata {
	//TODO build these from yaml config?
	defaultEndpoint := structs.ModuleEndpoint{
		Secure:      false,
		Endpoint:    "/run",
		HttpMethods: []structs.HttpMethod{{"OPTIONS"}, {"POST"}},
	}

	healthEndpoint := structs.ModuleEndpoint{
		Secure:      true,
		Endpoint:    "/health",
		HttpMethods: []structs.HttpMethod{{"OPTIONS"}, {"GET"}},
	}

	testEndpoint := structs.ModuleEndpoint{
		Secure:      true,
		Endpoint:    "/config",
		HttpMethods: []structs.HttpMethod{{"OPTIONS"}, {"GET"}, {"POST"}},
	}

	desc := "Ingests intelligence received from security findings in Amazon GuardDuty."

	localPort := os.Getenv("LOCAL_PORT")
	moduleSvcName := os.Getenv("MODULE_SVC_NAME")

	return structs.ModuleMetadata{
		ModuleServiceName: moduleSvcName,
		ModuleDisplayName: "Amazon GuardDuty",
		ModuleDescription: desc,
		ModuleType:        "ingress",
		InboundRoute:      "/awsgd",
		InternalIP:        GetLocalIP(),
		InternalPort:      localPort,
		Configured:        true,
		Configurable:      false,
		IconURL:           os.Getenv("ICON_URL"),
		InternalEndpoints: []structs.ModuleEndpoint{defaultEndpoint, healthEndpoint, testEndpoint},
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error(err)
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

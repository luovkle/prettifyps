package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Container struct {
	Name       string `json:"Name"`
	State      string `json:"State"`
	Service    string `json:"Service"`
	Publishers []struct {
		URL           string `json:"URL"`
		TargetPort    int    `json:"TargetPort"`
		PublishedPort int    `json:"PublishedPort"`
		Protocol      string `json:"Protocol"`
	} `json:"Publishers"`
}

func main() {
	// Leer la entrada de la tubería
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Decodificar la entrada JSON
	var containers []Container
	err = json.Unmarshal(content, &containers)
	if err != nil {
		panic(err)
	}

	// Crear la tabla
	table := [][]string{{"SERVICE", "STATUS", "PORT"}}

	for i, container := range containers {
		service := container.Service
		state := container.State
		ports := []string{}

		for _, publisher := range container.Publishers {
			port := ""

			if publisher.URL != "" {
				port += publisher.URL + ":"
			}

			if publisher.PublishedPort != 0 {
				port += fmt.Sprintf("%d", publisher.TargetPort) + "->" + fmt.Sprintf("%d", publisher.PublishedPort)
			} else {
				port += fmt.Sprintf("%d", publisher.TargetPort)
			}

			port += "/" + publisher.Protocol
			ports = append(ports, port)
		}

		if i == len(containers)-1 && len(ports) > 0 && ports[0] == "/0" {
			ports = []string{ports[0]}
		}

		table = append(table, []string{service, state, strings.Join(ports, ", ")})
	}

	// Imprimir la tabla
	for _, row := range table {
		fmt.Printf("%-10s%-10s%-20s\n", row[0], row[1], row[2])
	}
}

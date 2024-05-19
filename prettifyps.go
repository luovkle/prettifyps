package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	// Leer la entrada de la tubería línea por línea
	reader := bufio.NewReader(os.Stdin)
	var containers []Container

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		var container Container
		err = json.Unmarshal([]byte(line), &container)
		if err != nil {
			panic(err)
		}

		containers = append(containers, container)
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
	// Encontrar la longitud máxima de cada columna
	maxServiceLen := 0
	for _, row := range table {
		if len(row[0]) > maxServiceLen {
			maxServiceLen = len(row[0])
		}
	}

	// Imprimir cada fila con formato adecuado
	for _, row := range table {
		fmt.Printf("%-*s %-10s %-20s\n", maxServiceLen+1, row[0], row[1], row[2])
	}
}

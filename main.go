package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	LocationMap    map[string]Location
	DistributorMap map[string][]distributor
)

func init() {
	DistributorMap = make(map[string][]distributor, 0)
	distributor1 := distributor{
		Name:    strings.ToLower("DISTRIBUTOR1"),
		Include: []string{"INDIA", "JAPAN"},
		Exclude: []string{"KARNATAKA-INDIA", "CHENNAI-TAMILNADU-INDIA"},
	}
	subDistributor := distributor{
		Name:    strings.ToLower("DISTRIBUTOR2"),
		Include: []string{"INDIA"},
		Exclude: []string{"TAMILNADU-INDIA"},
	}
	distributor1.SubDistributors = append(distributor1.SubDistributors, subDistributor)
	DistributorMap[distributor1.Name] = append(DistributorMap[distributor1.Name], distributor1)

	file, err := os.Open("./cities.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	LocationMap = make(map[string]Location)

	for _, row := range data {
		location := Location{
			City:    strings.ReplaceAll(strings.ToLower(row[3]), " ", "_"),
			State:   strings.ReplaceAll(strings.ToLower(row[4]), " ", "_"),
			Country: strings.ReplaceAll(strings.ToLower(row[5]), " ", "_"),
		}
		LocationMap[location.City] = location
		LocationMap[location.State] = location
		LocationMap[location.Country] = location

	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /distributor", PostDistributor)
	mux.HandleFunc("POST /distributor/{name}/sub-distributor", PostSubDistributor)
	mux.HandleFunc("POST /distributor/{name}/search", Search)

	fmt.Println("Starting server at port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

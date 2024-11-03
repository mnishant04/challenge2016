package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func PostDistributor(w http.ResponseWriter, r *http.Request) {

	distributorReq, err := UnmarshalJson[distributorRequest](r.Body)
	if err != nil {
		log.Printf("error while unMarshalling: %s", err)
		SendResponse(http.StatusInternalServerError, "Error while decoding body", nil, w)
		return
	}
	err = doesRegionExist(distributorReq)
	if err != nil {
		SendResponse(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	distributor, err := addPermission(distributorReq, distributorReq.Name, false)
	if err != nil {
		SendResponse(http.StatusInternalServerError, err.Error(), nil, w)
		return
	}
	DistributorMap[distributorReq.Name] = append(DistributorMap[distributorReq.Name], distributor)
	SendResponse(http.StatusCreated, "Success", "Distributor Added Successfully", w)
}

func PostSubDistributor(w http.ResponseWriter, r *http.Request) {
	parentDistributor := r.PathValue("name")
	dist, relation := checkDistributor(parentDistributor)
	if relation == "nil" {
		log.Printf("distributor: %s not found", parentDistributor)
		SendResponse(http.StatusNotFound, fmt.Sprintf("distributor: %s not found", parentDistributor), nil, w)
		return
	}
	subDistributorReq, err := UnmarshalJson[distributorRequest](r.Body)
	if err != nil {
		log.Printf("error while unMarshalling: %s", err)
		SendResponse(http.StatusInternalServerError, "Error while decoding body", nil, w)
		return
	}
	err = doesRegionExist(subDistributorReq)
	if err != nil {
		SendResponse(http.StatusBadRequest, err.Error(), nil, w)
		return
	}
	if relation == "child" {
		DistributorMap[parentDistributor] = append(DistributorMap[parentDistributor], dist)
	}

	subDistributor, err := addPermission(subDistributorReq, dist.Name, true)
	if err != nil {
		SendResponse(http.StatusForbidden, "Error while Adding Permission", nil, w)
		return
	}
	DistributorMap[parentDistributor] = append(DistributorMap[parentDistributor], subDistributor)

	SendResponse(http.StatusCreated, "Success", "Sub Distributor Added Successfully", w)

}

func Search(w http.ResponseWriter, r *http.Request) {
	distributorName := r.PathValue("name")
	distributor, relation := checkDistributor(distributorName)
	if relation == "nil" {	
		SendResponse(http.StatusInternalServerError, fmt.Sprintf("Distributor %s does not exit", distributorName), nil, w)
	}
	location, err := UnmarshalJson[Location](r.Body)
	if err != nil {
		log.Printf("error while unMarshalling: %s", err)
		SendResponse(http.StatusInternalServerError, "Error while decoding body", nil, w)
		return
	}
	if location.City != "" {
		isAuthorized := distributor.IsAuthorized(location.City, "city")
		if isAuthorized {
			SendResponse(http.StatusOK, "Success", "YES", w)
			return
		}
	}

	if location.State != "" {
		isAuthorized := distributor.IsAuthorized(location.State, "state")
		if isAuthorized {
			SendResponse(http.StatusOK, "Success", "YES", w)
			return
		}
	}

	if location.Country != "" {
		isAuthorized := distributor.IsAuthorized(location.Country, "country")
		if isAuthorized {
			SendResponse(http.StatusOK, "Success", "YES", w)
			return
		}
	}
	SendResponse(http.StatusOK, "Success", "NO", w)
}

func addPermission(distributorReq distributorRequest, parentDistributorName string, isSubDistributor bool) (distributor, error) {
	include_things := []string{}
	exclude_things := []string{}
	parent, _ := checkDistributor(parentDistributorName)
	for _, val := range distributorReq.Include {
		if isSubDistributor {
			authStatus := false
			if val.Country != "" {
				if parent.IsAuthorized(val.Country, "country") {
					authStatus = true
				}
			}
			if val.State != "" {
				if parent.IsAuthorized(val.State, "state") {
					authStatus = true
				}
			}
			if val.City != "" {
				if parent.IsAuthorized(val.City, "city") {
					authStatus = true
				}
			}
			if !authStatus {
				log.Println("You are not authorized")
				return distributor{}, errors.New("you are not authorized to give access")
			}
		}
		include_things = append(include_things, combineLocation(val))
	}
	for _, val := range distributorReq.Exclude {
		if isSubDistributor {
			authStatus := false
			if val.Country != "" {
				if parent.IsAuthorized(val.Country, "country") {
					authStatus = true
				}
			}
			if val.State != "" {
				if parent.IsAuthorized(val.State, "state") {
					authStatus = true
				}
			}
			if val.City != "" {
				if parent.IsAuthorized(val.City, "city") {
					authStatus = true
				}
			}
			if authStatus {
				return distributor{}, errors.New("you are not authorized to give access")
			}
		}
		exclude_things = append(exclude_things, combineLocation(val))
	}
	dist := distributor{
		Name:            parentDistributorName,
		Include:         include_things,
		Exclude:         exclude_things,
		SubDistributors: []distributor{},
	}
	return dist, nil

}

func combineLocation(location Location) string {
	include_perm := ""
	if (location.Country != "") || (location.State != "") || (location.City != "") {
		if location.Country != "" {
			include_perm += strings.ReplaceAll(strings.ToLower(location.Country), " ", "_")
		}
		if location.State != "" {
			include_perm = "_" + include_perm
			include_perm = strings.ReplaceAll(strings.ToLower(location.State), " ", "_") + include_perm
		}
		if location.City != "" {
			include_perm = "_" + include_perm
			include_perm = strings.ReplaceAll(strings.ToLower(location.City), " ", "_") + include_perm
		}
	}
	return include_perm
}

func (d distributor) IsAuthorized(region, loc string) bool {
	location, found := LocationMap[strings.ToLower(region)]
	new_reg1 := ""
	new_reg2 := ""
	new_reg3 := ""

	fmt.Println("location", location, "region", region, d)
	if !found {
		fmt.Println("Location not found")
		return false
	} else {
		if loc == "city" {
			new_reg1 = location.City + "_" + location.State + "_" + location.Country
			new_reg2 = location.State + "_" + location.Country
			new_reg3 = location.Country
		} else if loc == "state" {
			new_reg1 = location.State + "_" + location.Country
			new_reg2 = location.Country
			new_reg3 = location.Country
		} else if loc == "country" {
			new_reg1 = location.Country
			new_reg2 = location.Country
			new_reg3 = location.Country
		}
	}

	for _, i := range d.Include {
		if strings.HasPrefix(new_reg1, i) || strings.HasPrefix(new_reg2, i) || strings.HasPrefix(new_reg3, i) {
			return true
		}
	}

	for _, e := range d.Exclude {
		if strings.HasPrefix(new_reg1, e) || strings.HasPrefix(new_reg2, e) || strings.HasPrefix(new_reg3, e) {
			return false
		}
	}

	return false
}

func checkDistributor(parentDistributor string) (distributor, string) {

	for key, value := range DistributorMap {
		if key == parentDistributor {
			return value[0], "parent"
		} else {
			for _, childDistributor := range value {
				if childDistributor.Name == parentDistributor {
					return childDistributor, "child"
				}
			}
		}
	}
	return distributor{}, "nil"
}

func doesRegionExist(distributorReq distributorRequest) error {
	for _, location_in := range distributorReq.Include {
		if location_in.City != "" {
			if _, ok := LocationMap[strings.ToLower(location_in.City)]; !ok {
				return fmt.Errorf("region %s does not exist", location_in.City)
			}
		}
		if location_in.State != "" {
			if _, ok := LocationMap[strings.ToLower(location_in.State)]; !ok {
				return fmt.Errorf("region %s does not exist", location_in.State)
			}
		}
		if location_in.Country != "" {
			if _, ok := LocationMap[strings.ToLower(location_in.Country)]; !ok {
				return fmt.Errorf("region %s does not exist", location_in.Country)
			}
		}
	}
	for _, location_ex := range distributorReq.Exclude {
		if location_ex.City != "" {
			if _, ok := LocationMap[strings.ToLower(location_ex.City)]; !ok {
				return fmt.Errorf("region %s does not exist", location_ex.City)
			}
		}
		if location_ex.State != "" {
			if _, ok := LocationMap[strings.ToLower(location_ex.State)]; !ok {
				return fmt.Errorf("region %s does not exist", location_ex.State)
			}
		}
		if location_ex.Country != "" {
			if _, ok := LocationMap[strings.ToLower(location_ex.Country)]; !ok {
				return fmt.Errorf("region %s does not exist", location_ex.Country)
			}
		}
	}
	return nil
}

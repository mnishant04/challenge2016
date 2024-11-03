package main

type distributor struct {
	Name            string
	Include         []string
	Exclude         []string
	SubDistributors []distributor
}

type Location struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

type distributorRequest struct {
	Name    string     `json:"name"`
	Include []Location `json:"include"`
	Exclude []Location `json:"exclude,omitempty"`
}

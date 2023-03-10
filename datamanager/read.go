package datamanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Data struct {
	Results []struct {
		First   string `json:"first"`
		Last    string `json:"last"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Created string `json:"created"`
		Balance string `json:"balance"`
	} `json:"results"`
}

func readData(location string) ([][]string, error) {
	//Read the JSON data from location
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		err = fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
		return nil, err
	}

	//Convert the JSON to data structure
	var data Data
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		err = errors.New("can't get data from empty record")
		return nil, err
	}

	//Store the obtained data into a 2D slice
	result := make([][]string, len(data.Results))
	for i, v := range data.Results {
		result[i] = append(result[i], v.First, v.Last, v.Email, v.Address, v.Created, v.Balance)
	}
	return result, nil
}

// GetData reads data from a location and receives the number of records that should be read.
// It calls the readData() func as many times it needs to ensure that it gets the number of lines
// specified in the input. It returns a 2D slice of strings with the read records and an error if there
// is one.
func GetData(location string, noOfRecords int) ([][]string, error) {
	var finalData [][]string
	//Call the ReadData() as many times as needed to get the number of records specified
	for noOfRecords > len(finalData) {
		data, err := readData(location)
		if err != nil {
			return nil, err
		}
		finalData = append(finalData, data...)
	}
	finalData = finalData[:noOfRecords]

	return finalData, nil
}

package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pamaehara/guardiao/domain/api"
	log "github.com/sirupsen/logrus"
)

func SalesUnits(reqAPI api.RequestApi, source string, ownCredential bool, capture bool, branchID string) (api.ReportUnitsResponse, error) {
	client := http.DefaultClient

	log.Info("##################################")
	log.Info("#       REPORT UNITS             #")
	log.Info("##################################")

	req, err := http.NewRequest("GET", getReportUnitsURL(reqAPI, ownCredential, source, capture, branchID), nil)
	for _, header := range reqAPI.GetHeaders() {
		req.Header.Add(header.Name, header.Value)
	}

	log.Trace("Performing search: ", req)

	resp, err := client.Do(req)

	log.Info(resp)

	if err != nil {
		log.Info("Error perfoming search", err)
		return api.ReportUnitsResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Trace("Response: ", string(body))

	if err != nil {
		log.Info("Error perfoming search", err)
		return api.ReportUnitsResponse{}, err
	}

	flightsResponse := api.ReportUnitsResponse{}
	json.Unmarshal(body, &flightsResponse)
	flightsResponse.CorrelationID = resp.Header["Gtw-Transaction-Id"][0]

	return flightsResponse, nil
}

func ReportSales(reqAPI api.RequestApi, source string, agencyID string, branchID string, capture bool, officeID string, reportDate string) (api.ReportSalesResponse, error) {
	client := http.DefaultClient

	log.Info("##################################")
	log.Info("#             SEARCH             #")
	log.Info("##################################")

	req, err := http.NewRequest("GET", getReportSalesURL(reqAPI, agencyID, branchID, capture, officeID, reportDate, source), nil)
	for _, header := range reqAPI.GetHeaders() {
		req.Header.Add(header.Name, header.Value)
	}

	log.Trace("Performing search: ", req)
	resp, err := client.Do(req)

	log.Info(resp)

	if err != nil {
		log.Info("Error perfoming search", err)
		return api.ReportSalesResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Trace("Response: ", string(body))

	if err != nil {
		log.Info("Error perfoming search", err)
		return api.ReportSalesResponse{}, err
	}

	flightsResponse := api.ReportSalesResponse{}
	json.Unmarshal(body, &flightsResponse)

	return flightsResponse, nil
}

func getReportUnitsURL(a api.RequestApi, ownCredential bool, source string, capture bool, branchID string) string {
	return fmt.Sprintf("%s/units?branch-Id=%s&source=%s", a.GetURL(), branchID, source)
}

func getReportSalesURL(a api.RequestApi, agencyID string, branchID string, capture bool, officeID string, reportDate string, source string) string {
	return fmt.Sprintf("%s/sales?branch-Id=%s&office-Id=%s&reportDate=%s&source=%s", a.GetURL(), branchID, officeID, reportDate, source)
}

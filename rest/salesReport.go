package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pamaehara/guardiao/config"
	"github.com/pamaehara/guardiao/domain/api"
	"github.com/pamaehara/guardiao/domain/model"
	"github.com/pamaehara/guardiao/integration"
)

func SalesReportHandler(res http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	fmt.Println("params: ", params)

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// bodyString := string(bodyBytes)

	var reqAPI api.RequestApi
	json.Unmarshal(bodyBytes, &reqAPI)
	fmt.Println(reqAPI)

	reportUnitsRS, err := integration.SalesUnits(reqAPI, "LAT", false, false, "")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(reportUnitsRS)
	jsonValue, _ := json.Marshal(reportUnitsRS)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(jsonValue))

}

func LocHandler(res http.ResponseWriter, req *http.Request) {
	loc := strings.TrimPrefix(req.URL.Path, "/salesReport/")
	switch {
	case req.Method == "GET":
		getReserva(res, req, loc)
	default:
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Foi mal!")
	}
}

func getReserva(res http.ResponseWriter, req *http.Request, loc string) {
	conn := config.GetConn()
	defer conn.Close()

	var reserva model.Reserva
	conn.Table("reserva").Where("loc = ?", loc).First(&reserva)

	json, _ := json.Marshal(reserva)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(json))
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"git.reservafacil.tur.br/gateway/guardiao/conexao"
	"git.reservafacil.tur.br/gateway/guardiao/domain"
)

// SalesReportHandler responde ao serviço de consulta de reservas
func SalesReportHandler(res http.ResponseWriter, req *http.Request) {
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
	conn := conexao.GetConn()
	defer conn.Close()

	var reserva domain.Reserva
	conn.Table("reserva").Where("loc = ?", loc).First(&reserva)

	json, _ := json.Marshal(reserva)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(json))
}

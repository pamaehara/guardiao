package main

import (
	"fmt"

	"git.reservafacil.tur.br/gateway/guardiao/conexao"
	"git.reservafacil.tur.br/gateway/guardiao/domain"
)

func main() {
	db := conexao.GetConn()
	defer db.Close()

	var reserva domain.Reserva

	db.Table("reserva").Where("loc = ?", "FPFSHK").First(&reserva)
	fmt.Println(reserva)
}

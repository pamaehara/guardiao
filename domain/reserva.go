package domain

import (
	"time"
)

// Reserva representa a tabela de reserva do GWAPI
type Reserva struct {
	Id                int       `json:"id" gorm:"column:id"`
	CodigoCia         string    `json:"codigoCia" gorm:"column:codigo_cia"`
	Credencial_id     int       `json:"credencialId"`
	Data_atualizacao  time.Time `json:"dataAtualizacao"`
	Data_criacao      time.Time `json:"dataCriacao"`
	Loc               string    `json:"loc"`
	Office_id_emissao string    `json:"officeIdEmissao"`
	Office_id_reserva string    `json:"officeIdReserva"`
	Produto           string    `json:"produto"`
	Sist_emis         string    `json:"sistEmis"`
	Status            string    `json:"status"`
	Pedido_id         int64     `json:"pedidoId"`
	Arquivada         bool      `json:"arquivada"`
}

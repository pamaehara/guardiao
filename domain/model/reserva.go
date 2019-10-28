package model

import (
	"time"
)

// Reserva representa a tabela de reserva do GWAPI
type Reserva struct {
	Id                int       `json:"id" gorm:"column:id"`
	CodigoCia         string    `json:"codigoCia" gorm:"column:codigo_cia"`
	Credencial_id     int       `json:"credencialId" gorm:"column:credencial_id"`
	Data_atualizacao  time.Time `json:"dataAtualizacao" gorm:"column:data_atualizacao"`
	Data_criacao      time.Time `json:"dataCriacao" gorm:"column:data_criacao"`
	Loc               string    `json:"loc" gorm:"column:loc"`
	Office_id_emissao string    `json:"officeIdEmissao" gorm:"column:office_id_emissao"`
	Office_id_reserva string    `json:"officeIdReserva" gorm:"column:office_id_reserva"`
	Produto           string    `json:"produto" gorm:"column:produto"`
	Sist_emis         string    `json:"sistEmis" gorm:"column:sist_emis"`
	Status            string    `json:"status" gorm:"column:status"`
	Pedido_id         int64     `json:"pedidoId" gorm:"column:pedido_id"`
	Arquivada         bool      `json:"arquivada" gorm:"column:arquivada"`
}

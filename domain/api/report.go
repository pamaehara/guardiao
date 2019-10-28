package api

import (
	"fmt"
	"strings"
)

type ReportUnitsResponse struct {
	Units []struct {
		BranchID  string `json:"branchId"`
		UnitCodes []struct {
			Code    string `json:"code"`
			Source  string `json:"source"`
			Capture bool   `json:"capture"`
		} `json:"unitCodes"`
	} `json:"units"`
	CorrelationID string
}

type ReportSalesResponse struct {
	ReportDate string `json:"reportDate"`
	Sales      []struct {
		Iata        string `json:"iata"`
		SalesUnit   string `json:"salesUnit"`
		UserAgentID string `json:"userAgentId"`
		Source      string `json:"source"`
		Locator     string `json:"locator"`
		Ticket      []struct {
			TicketToken string `json:"ticketToken"`
			Number      string `json:"number"`
			Status      string `json:"status"`
			Pax         struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Gender    string `json:"gender"`
			} `json:"pax"`
			TransactionDate string `json:"transactionDate"`
			Payments        []struct {
				Type  string  `json:"type"`
				Value float64 `json:"value"`
			} `json:"payments"`
		} `json:"ticket"`
	} `json:"sales"`
	CorrelationID string
}

func (s ReportUnitsResponse) Print() string {
	var message []string

	for _, unit := range s.Units {
		message = append(message, fmt.Sprintf("BranchID %s\n", unit.BranchID))
		for _, unitCode := range unit.UnitCodes {
			message = append(message, fmt.Sprintf("\tCode %s - Source %s - Capture %v", unitCode.Code, unitCode.Source, unitCode.Capture))
		}
	}

	return strings.Join(message, "")
}

func (s ReportSalesResponse) Print() string {
	var message []string

	message = append(message, s.ReportDate)

	for _, sale := range s.Sales {
		message = append(message, fmt.Sprintf("%s | %s | %s | %s | %s\n", sale.Locator, sale.Source, sale.Iata, sale.SalesUnit, sale.UserAgentID))
		for _, tkt := range sale.Ticket {
			message = append(message, fmt.Sprintf("\t%s | %s | %s | %s | %s\n", tkt.Number, tkt.Status, fmt.Sprint(tkt.Pax.Gender, " ", tkt.Pax.FirstName, " ", tkt.Pax.LastName), tkt.TransactionDate, tkt.TicketToken))
		}
		message = append(message, fmt.Sprintln())
	}

	return strings.Join(message, "")
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agionfriddo/pdf-poc/internal/data"
	pdfgen "github.com/agionfriddo/pdf-poc/internal/pdf"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		payroll, err := data.FetchPayroll(ctx)
		if err != nil {
			log.Printf("error fetching payroll: %v", err)
			http.Error(w, "failed to fetch data", http.StatusBadGateway)
			return
		}

		method := r.URL.Query().Get("method")
		var pdfBytes []byte
		switch method {
		case "html":
			if payroll.Country == "CA" {
				pdfBytes, err = pdfgen.GenerateCanadaPayrollPDFViaHTML(payroll)
			} else {
				pdfBytes, err = pdfgen.GeneratePayrollPDFViaHTML(payroll)
			}
		case "playwright":
			if payroll.Country == "CA" {
				pdfBytes, err = pdfgen.GenerateCanadaPayrollPDFViaPlaywright(payroll)
			} else {
				pdfBytes, err = pdfgen.GeneratePayrollPDFViaPlaywright(payroll)
			}
		default:
			if payroll.Country == "CA" {
				pdfBytes, err = pdfgen.GenerateCanadaPayrollPDF(payroll)
			} else {
				pdfBytes, err = pdfgen.GeneratePayrollPDF(payroll)
			}
		}
		if err != nil {
			log.Printf("error generating pdf: %v", err)
			http.Error(w, "failed to generate pdf", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		filename := fmt.Sprintf("paystub_%s.pdf", payroll.Employee.ID)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdfBytes)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}


package models

import "time"

// Laptop representa un equipo disponible en el inventario.
type Laptop struct {
    ID           int       `db:"id" json:"id"`
    Name         string    `db:"name" json:"name"`
    Brand        string    `db:"brand" json:"brand"`
    SerialNumber string    `db:"serial_number" json:"serial_number"`
    Status       string    `db:"status" json:"status"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// LaptopRequest contiene los datos para registrar una laptop.
type LaptopRequest struct {
    Name         string `json:"name" binding:"required"`
    Brand        string `json:"brand" binding:"required"`
    SerialNumber string `json:"serial_number" binding:"required"`
}

// LaptopLoanRequest identifica la laptop que se desea pedir.
type LaptopLoanRequest struct {
    LaptopID int `json:"laptop_id" binding:"required"`
}

// LaptopLoan representa un préstamo de laptop.
type LaptopLoan struct {
    ID         int        `db:"id" json:"id"`
    UserID     int        `db:"user_id" json:"user_id"`
    LaptopID   int        `db:"laptop_id" json:"laptop_id"`
    LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
    DueAt      time.Time  `db:"due_at" json:"due_at"`
    ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
    Status     string     `db:"status" json:"status"`
}

// LaptopLoanDetail combina el préstamo con la información de la laptop.
type LaptopLoanDetail struct {
    LaptopLoan
    LaptopName  string `db:"laptop_name" json:"laptop_name"`
    LaptopBrand string `db:"laptop_brand" json:"laptop_brand"`
}

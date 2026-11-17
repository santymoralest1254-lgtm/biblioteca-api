package repository

import (
	"database/sql"
	"errors"
	"time"

	"biblioteca-api/internal/models"

	"github.com/jmoiron/sqlx"
)

var (
	ErrLaptopNotFound        = errors.New("laptop no encontrada")
	ErrLaptopNotAvailable    = errors.New("laptop no disponible")
	ErrActiveLaptopLoan      = errors.New("el usuario ya tiene un préstamo activo")
	ErrLaptopLoanNotFound    = errors.New("préstamo no encontrado")
	ErrLaptopLoanForbidden   = errors.New("préstamo no pertenece al usuario")
	ErrLaptopAlreadyReturned = errors.New("préstamo ya devuelto")
)

type LaptopRepository struct {
	db *sqlx.DB
}

func NewLaptopRepository(db *sqlx.DB) *LaptopRepository {
	return &LaptopRepository{db: db}
}

// Create registra una laptop en el inventario.
func (r *LaptopRepository) Create(req models.LaptopRequest) (*models.Laptop, error) {
	result, err := r.db.Exec(`
        INSERT INTO laptops (name, brand, serial_number)
        VALUES (?, ?, ?)`,
		req.Name, req.Brand, req.SerialNumber,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(int(id))
}

// FindByID consulta una laptop por su identificador.
func (r *LaptopRepository) FindByID(id int) (*models.Laptop, error) {
	var laptop models.Laptop

	err := r.db.Get(&laptop, `
        SELECT id, name, brand, serial_number, status, created_at
        FROM laptops
        WHERE id = ?`, id,
	)
	if err != nil {
		return nil, err
	}

	return &laptop, nil
}

// FindAll devuelve todas las laptops.
func (r *LaptopRepository) FindAll() ([]models.Laptop, error) {
	laptops := []models.Laptop{}

	err := r.db.Select(&laptops, `
        SELECT id, name, brand, serial_number, status, created_at
        FROM laptops
        ORDER BY id DESC`,
	)

	return laptops, err
}

// FindAvailable devuelve solamente las laptops disponibles.
func (r *LaptopRepository) FindAvailable() ([]models.Laptop, error) {
	laptops := []models.Laptop{}

	err := r.db.Select(&laptops, `
        SELECT id, name, brand, serial_number, status, created_at
        FROM laptops
        WHERE status = 'available'
        ORDER BY id DESC`,
	)

	return laptops, err
}

// SetStatus actualiza el estado de una laptop.
func (r *LaptopRepository) SetStatus(id int, status string) error {
	_, err := r.db.Exec(`
        UPDATE laptops
        SET status = ?
        WHERE id = ?`, status, id,
	)

	return err
}

// CreateLoan registra un préstamo de laptop.
func (r *LaptopRepository) CreateLoan(userID, laptopID int, dueAt time.Time) (*models.LaptopLoan, error) {
	result, err := r.db.Exec(`
        INSERT INTO laptop_loans (user_id, laptop_id, due_at)
        VALUES (?, ?, ?)`,
		userID, laptopID, dueAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindLoanByID(int(id))
}

// FindLoanByID consulta un préstamo.
func (r *LaptopRepository) FindLoanByID(id int) (*models.LaptopLoan, error) {
	var loan models.LaptopLoan

	err := r.db.Get(&loan, `
        SELECT id, user_id, laptop_id, loaned_at, due_at, returned_at, status
        FROM laptop_loans
        WHERE id = ?`, id,
	)
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// HasActiveLoan comprueba si el usuario ya tiene esa laptop prestada.
func (r *LaptopRepository) HasActiveLoan(userID, laptopID int) (bool, error) {
	var count int

	err := r.db.Get(&count, `
        SELECT COUNT(*)
        FROM laptop_loans
        WHERE user_id = ? AND laptop_id = ? AND status = 'active'`,
		userID, laptopID,
	)

	return count > 0, err
}

// MarkLoanReturned registra la devolución del préstamo.
func (r *LaptopRepository) MarkLoanReturned(id int) error {
	result, err := r.db.Exec(`
        UPDATE laptop_loans
        SET status = 'returned', returned_at = CURRENT_TIMESTAMP
        WHERE id = ? AND status = 'active'`,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindLoansByUser devuelve los préstamos de un usuario.
func (r *LaptopRepository) FindLoansByUser(userID int) ([]models.LaptopLoanDetail, error) {
	loans := []models.LaptopLoanDetail{}

	err := r.db.Select(&loans, `
        SELECT ll.id, ll.user_id, ll.laptop_id, ll.loaned_at,
               ll.due_at, ll.returned_at, ll.status,
               l.name AS laptop_name, l.brand AS laptop_brand
        FROM laptop_loans ll
        INNER JOIN laptops l ON l.id = ll.laptop_id
        WHERE ll.user_id = ?
        ORDER BY ll.loaned_at DESC`,
		userID,
	)

	return loans, err
}

// BorrowAtomically registra un préstamo y actualiza el inventario
// dentro de una sola transacción.
func (r *LaptopRepository) BorrowAtomically(
	userID, laptopID int,
	dueAt time.Time,
) (*models.LaptopLoan, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Bloquea la laptop para evitar préstamos simultáneos.
	var status string

	err = tx.Get(&status, `
        SELECT status
        FROM laptops
        WHERE id = ?
        FOR UPDATE`,
		laptopID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLaptopNotFound
	}
	if err != nil {
		return nil, err
	}

	if status != "available" {
		return nil, ErrLaptopNotAvailable
	}

	// Comprueba si el usuario ya tiene un préstamo activo.
	var count int

	err = tx.Get(&count, `
        SELECT COUNT(*)
        FROM laptop_loans
        WHERE user_id = ?
        AND laptop_id = ?
        AND status = 'active'`,
		userID, laptopID,
	)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ErrActiveLaptopLoan
	}

	// Registra el préstamo.
	result, err := tx.Exec(`
        INSERT INTO laptop_loans (user_id, laptop_id, due_at)
        VALUES (?, ?, ?)`,
		userID, laptopID, dueAt,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Actualiza el inventario.
	_, err = tx.Exec(`
        UPDATE laptops
        SET status = 'loaned'
        WHERE id = ?`,
		laptopID,
	)

	if err != nil {
		return nil, err
	}

	var loan models.LaptopLoan

	err = tx.Get(&loan, `
        SELECT id, user_id, laptop_id, loaned_at,
               due_at, returned_at, status
        FROM laptop_loans
        WHERE id = ?`,
		id,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &loan, nil
}

// ReturnAtomically registra la devolución y libera la laptop
// dentro de una sola transacción.
func (r *LaptopRepository) ReturnAtomically(
	userID, loanID int,
) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var loan models.LaptopLoan

	// Bloquea el préstamo durante la devolución.
	err = tx.Get(&loan, `
        SELECT id, user_id, laptop_id, loaned_at,
               due_at, returned_at, status
        FROM laptop_loans
        WHERE id = ?
        FOR UPDATE`,
		loanID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrLaptopLoanNotFound
	}
	if err != nil {
		return err
	}

	if loan.UserID != userID {
		return ErrLaptopLoanForbidden
	}

	if loan.Status == "returned" {
		return ErrLaptopAlreadyReturned
	}

	// Actualiza el préstamo.
	_, err = tx.Exec(`
        UPDATE laptop_loans
        SET status = 'returned',
            returned_at = CURRENT_TIMESTAMP
        WHERE id = ?`,
		loanID,
	)

	if err != nil {
		return err
	}

	// Devuelve la laptop al inventario.
	_, err = tx.Exec(`
        UPDATE laptops
        SET status = 'available'
        WHERE id = ?`,
		loan.LaptopID,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

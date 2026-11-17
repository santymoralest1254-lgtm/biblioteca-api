package service

import (
	"errors"
	"time"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

var (
	ErrLaptopNotAvailable    = errors.New("esta laptop no está disponible")
	ErrLaptopAlreadyBorrowed = errors.New("ya tienes un préstamo activo de esta laptop")
	ErrLaptopLoanNotFound    = errors.New("préstamo de laptop no encontrado")
	ErrLaptopAlreadyReturned = errors.New("esta laptop ya fue devuelta")
	ErrLaptopLoanForbidden   = errors.New("este préstamo no pertenece al usuario")
)

const laptopLoanDurationDays = 10

type LaptopService struct {
	laptopRepo *repository.LaptopRepository
}

func NewLaptopService(laptopRepo *repository.LaptopRepository) *LaptopService {
	return &LaptopService{laptopRepo: laptopRepo}
}

// Create registra una nueva laptop.
func (s *LaptopService) Create(req models.LaptopRequest) (*models.Laptop, error) {
	return s.laptopRepo.Create(req)
}

// List devuelve todas las laptops.
func (s *LaptopService) List() ([]models.Laptop, error) {
	return s.laptopRepo.FindAll()
}

// ListAvailable devuelve solamente las laptops disponibles.
func (s *LaptopService) ListAvailable() ([]models.Laptop, error) {
	return s.laptopRepo.FindAvailable()
}

// Borrow registra el préstamo de una laptop.
func (s *LaptopService) Borrow(userID int, req models.LaptopLoanRequest) (*models.LaptopLoan, error) {
	dueAt := time.Now().AddDate(0, 0, laptopLoanDurationDays)

	loan, err := s.laptopRepo.BorrowAtomically(
		userID,
		req.LaptopID,
		dueAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLaptopNotAvailable):
			return nil, ErrLaptopNotAvailable
		case errors.Is(err, repository.ErrActiveLaptopLoan):
			return nil, ErrLaptopAlreadyBorrowed
		}
		return nil, err
	}

	return loan, nil
}

// Return registra la devolución de una laptop.
func (s *LaptopService) Return(userID, loanID int) error {
	err := s.laptopRepo.ReturnAtomically(userID, loanID)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLaptopLoanNotFound):
			return ErrLaptopLoanNotFound
		case errors.Is(err, repository.ErrLaptopLoanForbidden):
			return ErrLaptopLoanForbidden
		case errors.Is(err, repository.ErrLaptopAlreadyReturned):
			return ErrLaptopAlreadyReturned
		}
		return err
	}

	return nil
}

// ListByUser devuelve los préstamos de un usuario.
func (s *LaptopService) ListByUser(userID int) ([]models.LaptopLoanDetail, error) {
	return s.laptopRepo.FindLoansByUser(userID)
}

package model

import "time"

// User — kutubxona foydalanuvchisi (a'zo yoki admin).
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`    // "-" → JSON javobga hech qachon chiqmaydi
	Role         string    `json:"role"` // "member" yoki "admin"
	CreatedAt    time.Time `json:"created_at"`
}

// Book — kutubxonadagi kitob.
type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Year        int       `json:"year"`
	TotalCopies int       `json:"total_copies"` // kutubxonada jami nechta nusxa bor
	CreatedAt   time.Time `json:"created_at"`
}

// Loan — ijara: foydalanuvchi kitobni oldi.
type Loan struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id"` // FK → books.id
	UserID     int        `json:"user_id"` // FK → users.id
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueAt      time.Time  `json:"due_at"`      // qaytarish muddati
	ReturnedAt *time.Time `json:"returned_at"` // nil → hali qaytarilmagan
}

// LoanView — ro'yxatda ko'rsatish uchun: ijara + kitob nomi + foydalanuvchi ismi
type LoanView struct {
	ID         int        `json:"id"`
	BookTitle  string     `json:"book_title"`
	UserName   string     `json:"user_name"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueAt      time.Time  `json:"due_at"`
	ReturnedAt *time.Time `json:"returned_at"`
}

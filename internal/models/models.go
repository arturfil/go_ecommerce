package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// DBMode is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Struct to hold all models
type Models struct {
	DB DBModel
}

// Models Initializer
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Meeting Model
type Meeting struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	MeetingDate time.Time `json:"meeting_date"`
	Image       string    `json:"image"`
	IsRecurring bool      `json:"is_recurring"`
	PlanID      string    `json:"plan_id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// Order Struct
type Order struct {
	ID            int         `json:"id"`
	MeetingID     int         `json:"meeting_id"`
	TransactionID int         `json:"transaction_id"`
	CustomerID    int         `json:"customer_id"`
	StatusID      int         `json:"status_id"`
	Quantity      int         `json:"quantity"`
	Amount        int         `json:"amount"`
	Transaction   Transaction `json:"transaction"`
	Meeting       Meeting     `json:"meeting"`
	Customer      Customer    `json:"customer"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

// Statuses struct
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TODO: MISSING CHARGE ONCE FEATURE TO SOLVE, plan_id returns null!
func (m *DBModel) GetMeeting(id int) (Meeting, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var meeting Meeting

	query := `
        SELECT id, name, price, description, is_recurring, plan_id, meeting_date FROM meetings where id = ?
    `

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&meeting.ID,
		&meeting.Name,
		&meeting.Price,
		&meeting.Description,
		&meeting.IsRecurring,
		&meeting.PlanID,
		&meeting.MeetingDate,
	)
	if err != nil {
		return meeting, err
	}

	return meeting, nil
}

// InsertTransaction inserts a new txn, and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO transactions (
            amount,
            currency,
            last_four,
            bank_return_code,
            transaction_status_id,
            expiry_month,
            expiry_year,
            payment_intent,
            payment_method,
            created_at,
            updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := m.DB.ExecContext(
		ctx,
		query,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// InsertOrder inserts a new txn, and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO orders (
             meeting_id,
             transaction_id,
             customer_id,
             status_id,
             quantity,
             amount,
             created_at,
             updated_at
         ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := m.DB.ExecContext(
		ctx,
		query,
		order.MeetingID,
		order.TransactionID,
		order.CustomerID,
		order.StatusID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// InsertOrder inserts a new txn, and returns its id
func (m *DBModel) InsertCustomer(customer Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO customers
            (first_name, last_name, email, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?) 
    `

	result, err := m.DB.ExecContext(
		ctx,
		query,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var u User

	query := `
        SELECT id, first_name, last_name, email, password, created_at, updated_at
        FROM users
        WHERE email = ?
    `

	row := m.DB.QueryRowContext(
		ctx,
		query,
		email,
	)

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(
		ctx,
		"SELECT id, password FROM users WHERE email = ?",
		email,
	)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DBModel) UpdatePasswordForUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        UPDATE users set password = ? where id = ?
    `

	_, err := m.DB.ExecContext(ctx, query, hash, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Getting a section of an object
func (m *DBModel) GetAllOrdersPaginated(pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    offset := (page - 1) * pageSize

	var orders []*Order

	query := `
        
        SELECT
            o.id,
            m.id,
            m.name,
            o.transaction_id ,
            o.customer_id ,
            o.status_id ,
            o.quantity ,
            o.amount,
            c.id,
            c.first_name,
            c.last_name,
            c.email,
            t.id,
            t.currency,
            t.amount,
            t.last_four,
            t.expiry_month,
            t.expiry_year,
            t.payment_intent,
            t.bank_return_code,
            o.created_at,
            o.updated_at
        FROM orders o 

        LEFT JOIN meetings m ON o.meeting_id = m.id
        LEFT JOIN transactions t ON t.id = o.transaction_id
        LEFT JOIN customers c ON o.customer_id = c.id

        WHERE m.is_recurring = 0

        ORDER BY o.created_at DESC
        LIMIT ? OFFSET ?;
    `

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err 
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.Meeting.ID,
			&o.Meeting.Name,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
			&o.Transaction.ID,
			&o.Transaction.Currency,
			&o.Transaction.Amount,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.CreatedAt,
			&o.UpdatedAt,
		)

		if err != nil {
			return nil, 0, 0, err
		}

		orders = append(orders, &o)
	}

    query = `
        SELECT COUNT(o.id)
        FROM orders o
        LEFT JOIN meetings m on (o.meeting_id = m.id) 
        where
        m.is_recurring = 0
    `

    var totalRecords int
    countRow := m.DB.QueryRowContext(ctx, query)
    err = countRow.Scan(&totalRecords)
    if err != nil {
        return nil, 0, 0, err
    }

    lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

func (m *DBModel) GetAllOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `
        
        SELECT
            o.id,
            m.id,
            m.name,
            o.transaction_id ,
            o.customer_id ,
            o.status_id ,
            o.quantity ,
            o.amount,
            c.id,
            c.first_name,
            c.last_name,
            c.email,
            t.id,
            t.currency,
            t.amount,
            t.last_four,
            t.expiry_month,
            t.expiry_year,
            t.payment_intent,
            t.bank_return_code,
            o.created_at,
            o.updated_at
        FROM orders o 

        LEFT JOIN meetings m ON o.meeting_id = m.id
        LEFT JOIN transactions t ON t.id = o.transaction_id
        LEFT JOIN customers c ON o.customer_id = c.id

        WHERE m.is_recurring = 0

        ORDER BY o.created_at DESC;
    `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.Meeting.ID,
			&o.Meeting.Name,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
			&o.Transaction.ID,
			&o.Transaction.Currency,
			&o.Transaction.Amount,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.CreatedAt,
			&o.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	return orders, nil
}

func (m *DBModel) GetAllSubscriptionsPaginated(pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    offset := (page - 1) * pageSize

	var orders []*Order

	query := `
        
        SELECT
            o.id,
            m.id,
            m.name,
            o.transaction_id ,
            o.customer_id ,
            o.status_id ,
            o.quantity ,
            o.amount,
            c.id,
            c.first_name,
            c.last_name,
            c.email,
            t.id,
            t.currency,
            t.amount,
            t.last_four,
            t.expiry_month,
            t.expiry_year,
            t.payment_intent,
            t.bank_return_code,
            o.created_at,
            o.updated_at
        FROM orders o 

        LEFT JOIN meetings m ON o.meeting_id = m.id
        LEFT JOIN transactions t ON t.id = o.transaction_id
        LEFT JOIN customers c ON o.customer_id = c.id

        WHERE m.is_recurring = 1

        ORDER BY o.created_at DESC
        LIMIT ? OFFSET ?;
    `

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.Meeting.ID,
			&o.Meeting.Name,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
			&o.Transaction.ID,
			&o.Transaction.Currency,
			&o.Transaction.Amount,
			&o.Transaction.LastFour,
			&o.Transaction.ExpiryMonth,
			&o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent,
			&o.Transaction.BankReturnCode,
			&o.CreatedAt,
			&o.UpdatedAt,
		)

		if err != nil {
			return nil, 0, 0, err
		}

		orders = append(orders, &o)
	}

    query = `
        SELECT COUNT(o.id)
        FROM orders o
        LEFT JOIN meetings m on (o.meeting_id = m.id)
        WHERE
        m.is_recurring = 0
    `

    var totalRecords int
    countRow := m.DB.QueryRowContext(ctx, query)
    err = countRow.Scan(&totalRecords)
    if err != nil {
        return nil, 0, 0, err
    }

    lastPage := totalRecords / pageSize
    return orders, lastPage, totalRecords, nil
}

func (m *DBModel) GetOrderById(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
        
        SELECT
            o.id,
            m.id,
            m.name,
            o.transaction_id ,
            o.customer_id ,
            o.status_id ,
            o.quantity ,
            o.amount,
            c.id,
            c.first_name,
            c.last_name,
            c.email,
            t.id,
            t.currency,
            t.amount,
            t.last_four,
            t.expiry_month,
            t.expiry_year,
            t.payment_intent,
            t.bank_return_code,
            o.created_at,
            o.updated_at
        FROM orders o 

        LEFT JOIN meetings m ON o.meeting_id = m.id
        LEFT JOIN transactions t ON t.id = o.transaction_id
        LEFT JOIN customers c ON o.customer_id = c.id

        WHERE o.id = ? 
    `

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&o.ID,
		&o.Meeting.ID,
		&o.Meeting.Name,
		&o.TransactionID,
		&o.CustomerID,
		&o.StatusID,
		&o.Quantity,
		&o.Amount,
		&o.Customer.ID,
		&o.Customer.FirstName,
		&o.Customer.LastName,
		&o.Customer.Email,
		&o.Transaction.ID,
		&o.Transaction.Currency,
		&o.Transaction.Amount,
		&o.Transaction.LastFour,
		&o.Transaction.ExpiryMonth,
		&o.Transaction.ExpiryYear,
		&o.Transaction.PaymentIntent,
		&o.Transaction.BankReturnCode,
		&o.CreatedAt,
		&o.UpdatedAt,
	)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *DBModel) UpdateOrderStatus(id, statusID int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    query := "update orders set status_id = ? where id = ?"
    _, err := m.DB.ExecContext(ctx, query, statusID, id)
    if err != nil {
        return err
    }

    return nil
}

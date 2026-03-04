package repository

import (
	"database/sql"
	"github.com/hxseqwe/korochki-est/internal/model"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(app *model.Application) error {
	query := `INSERT INTO applications (user_id, course_name, start_date, payment_method, status) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	return r.db.QueryRow(query, app.UserID, app.CourseName, app.StartDate,
		app.PaymentMethod, app.Status).Scan(&app.ID, &app.CreatedAt)
}

func (r *ApplicationRepository) GetByUserID(userID int) ([]*model.Application, error) {
	rows, err := r.db.Query(`SELECT id, user_id, course_name, start_date, payment_method, 
                             status, review, created_at, updated_at 
                             FROM applications WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*model.Application
	for rows.Next() {
		app := &model.Application{}
		err := rows.Scan(&app.ID, &app.UserID, &app.CourseName, &app.StartDate,
			&app.PaymentMethod, &app.Status, &app.Review, &app.CreatedAt, &app.UpdatedAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *ApplicationRepository) GetAll() ([]*model.Application, error) {
	rows, err := r.db.Query(`SELECT a.id, a.user_id, a.course_name, a.start_date, 
                             a.payment_method, a.status, a.review, a.created_at, a.updated_at,
                             u.login, u.full_name, u.phone, u.email
                             FROM applications a
                             JOIN users u ON a.user_id = u.id
                             ORDER BY a.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*model.Application
	for rows.Next() {
		app := &model.Application{}
		user := &model.User{}
		err := rows.Scan(&app.ID, &app.UserID, &app.CourseName, &app.StartDate,
			&app.PaymentMethod, &app.Status, &app.Review, &app.CreatedAt, &app.UpdatedAt,
			&user.Login, &user.FullName, &user.Phone, &user.Email)
		if err != nil {
			return nil, err
		}
		app.User = user
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *ApplicationRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE applications SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *ApplicationRepository) UpdateApplication(id int, courseName string, startDate string, paymentMethod string) error {
	query := `UPDATE applications SET course_name = $1, start_date = $2, payment_method = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := r.db.Exec(query, courseName, startDate, paymentMethod, id)
	return err
}

func (r *ApplicationRepository) DeleteApplication(id int) error {
	query := `DELETE FROM applications WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ApplicationRepository) AddReview(id int, review string) error {
	query := `UPDATE applications SET review = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, review, id)
	return err
}

func (r *ApplicationRepository) GetByID(id int) (*model.Application, error) {
	app := &model.Application{}
	query := `SELECT id, user_id, course_name, start_date, payment_method, 
              status, review, created_at, updated_at 
              FROM applications WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&app.ID, &app.UserID, &app.CourseName,
		&app.StartDate, &app.PaymentMethod, &app.Status, &app.Review, &app.CreatedAt, &app.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return app, nil
}

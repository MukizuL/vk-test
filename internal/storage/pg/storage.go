package pg

import (
	"context"
	"fmt"

	"github.com/MukizuL/vk-test/internal/dto"
	filters2 "github.com/MukizuL/vk-test/internal/filters"
	"github.com/MukizuL/vk-test/internal/models"
	"github.com/google/uuid"
)

// CreateNewUser Creates a new user with given login and password. Returns userID and an error.
func (s *Storage) CreateNewUser(ctx context.Context, login, passwordHash string) (string, error) {
	userID := uuid.New()

	_, err := s.conn.Exec(ctx, `INSERT INTO users (id, login, passwordHash) VALUES ($1, $2, $3)`, userID, login, passwordHash)
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}

// GetUserByLogin Fetches user from database and stores all data in User struct. Returns User and an error.
func (s *Storage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, passwordHash FROM users WHERE login = $1`, login).
		Scan(&user.ID, &user.Login, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID Fetches user from database and stores all data in User struct. Returns User and an error.
func (s *Storage) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, passwordHash FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.Login, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) CreateAd(ctx context.Context, login string, req *dto.CreateAdRequest) (*models.Ad, error) {
	var ad models.Ad

	query := `
		INSERT INTO ads (user_login, title, description, image, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_login, title, description, image, price, created_at`

	args := []any{login, req.Title, req.Description, req.ImageURL, req.Price}

	err := s.conn.QueryRow(ctx, query, args...).
		Scan(
			&ad.ID,
			&ad.UserLogin,
			&ad.Title,
			&ad.Description,
			&ad.ImageURL,
			&ad.Price,
			&ad.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

func (s *Storage) GetAds(ctx context.Context, filters filters2.Filters) ([]models.Ad, filters2.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, user_login, title, description, image, price, created_at
		FROM ads
		WHERE price >= $1 AND price <= $2
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())

	fmt.Printf("Executing query: %s\nWith params: min=%v, max=%v, limit=%v, offset=%v\n",
		query, filters.GetMin(), filters.GetMax(), filters.Limit(), filters.Offset())

	rows, err := s.conn.Query(ctx, query, filters.GetMin()*100, filters.GetMax()*100, filters.Limit(), filters.Offset())
	if err != nil {
		return nil, filters2.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var ads []models.Ad

	for rows.Next() {
		var ad models.Ad

		err := rows.Scan(
			&totalRecords,
			&ad.ID,
			&ad.UserLogin,
			&ad.Title,
			&ad.Description,
			&ad.ImageURL,
			&ad.Price,
			&ad.CreatedAt,
		)

		if err != nil {
			return nil, filters2.Metadata{}, err
		}

		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		return nil, filters2.Metadata{}, err
	}

	metadata := filters2.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return ads, metadata, nil
}

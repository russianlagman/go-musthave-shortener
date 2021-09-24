package sqlstore

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"log"
	"strconv"
)

// store.Store interface implementation
var _ store.Store = (*Store)(nil)

func (s *Store) ReadURL(id string) (string, error) {
	rawID, err := s.toRawID(id)
	if err != nil {
		return "", fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}

	q := `
SELECT original_url FROM urls WHERE id=$1`
	var url string
	err = s.db.QueryRow(q, rawID).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("read url query: %w", err)
	}

	return url, err
}

func (s *Store) WriteURL(url string, uid string) (string, error) {
	if err := store.ValidateURL(url); err != nil {
		return "", err
	}

	q := `
INSERT INTO urls (uid, original_url)
VALUES ($1, $2)
RETURNING id`
	var rawID uint64
	err := s.db.QueryRow(q, uid, url).Scan(&rawID)
	if err != nil {
		return "", fmt.Errorf("write url query: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.baseURL, s.fromRawID(rawID)), nil
}

func (s *Store) ReadAllURLs(uid string) []store.Record {
	var result []store.Record

	q := `
SELECT id, original_url FROM urls WHERE uid=$1`
	rows, err := s.db.Query(q, uid)
	if err != nil {
		return result
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var (
			rawID       uint64
			originalURL string
		)
		if err := rows.Scan(&rawID, &originalURL); err != nil {
			log.Printf("scan error: %v", err)
			break
		}
		result = append(result, store.Record{
			ID:          s.fromRawID(rawID),
			OriginalURL: originalURL,
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, s.fromRawID(rawID)),
		})
	}

	return result
}

// fromRawID converts sql id to short string id
func (s *Store) fromRawID(id uint64) string {
	return strconv.FormatUint(id, s.base)
}

// toRawID converts short string id to sql id
func (s *Store) toRawID(id string) (uint64, error) {
	rawID, err := strconv.ParseUint(id, s.base, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}
	return rawID, nil
}

// idFromInt64 converts sql id to short string id
func (s *Store) idFromInt64(id int64) string {
	return strconv.FormatInt(id, s.base)
}

// idToInt64 converts short string id to sql id
func (s *Store) idToInt64(id string) (int64, error) {
	rawID, err := strconv.ParseInt(id, s.base, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}
	return rawID, nil
}

// shortUrl returns short url of the id
func (s *Store) shortUrl(id string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, id)
}

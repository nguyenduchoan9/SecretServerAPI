package repo

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/nguyenduchoan9/coderschool.go/assignment.2/goapi/src/db"
	encrypt "github.com/nguyenduchoan9/coderschool.go/assignment.2/goapi/src/encryption"
	"github.com/nguyenduchoan9/coderschool.go/assignment.2/goapi/src/model"
)

func AddSecret(originalSecret string, expireAfterViews int, expireAfter int) (*model.HashKey, *model.AppError) {
	createdAtParams := time.Now()
	expireAtParams := createdAtParams.Add(time.Duration(expireAfter) * time.Minute)
	encryptedText, err := encrypt.Encrypt(originalSecret)
	if err != nil {
		return nil, &model.AppError{err, "The issue occured while encrypted content", http.StatusInternalServerError}
	}
	hashKey := GenerateHashBy(originalSecret)

	sqlStatement := `INSERT INTO hashkey (secret_text, hash, created_at, expires_at, remaining_views)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, secret_text, hash, created_at, expires_at, remaining_views`
	var id, remainingViews int
	var secretText, hash string
	var createdAt, expireAt time.Time
	err = db.GetDB().QueryRow(sqlStatement, encryptedText, hashKey, createdAtParams, expireAtParams, expireAfterViews).Scan(&id, &secretText, &hash, &createdAt, &expireAt, &remainingViews)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, &model.AppError{err, "The secret already exists.", http.StatusNotFound}
		}
		return nil, &model.AppError{err, "The issue occured while creating secret.", http.StatusNotFound}

	}

	return &model.HashKey{
		ID:             id,
		SecretText:     originalSecret,
		Hash:           hash,
		CreatedAt:      createdAt,
		ExpireAt:       expireAt,
		RemainingViews: remainingViews,
	}, nil
}

func GetSecretBy(hashKey string) (*model.HashKey, *model.AppError) {
	sqlStatement := `SELECT id, secret_text, hash, created_at, expires_at, remaining_views FROM hashkey WHERE hash=$1`
	var id, remainingViews int
	var encryptedSecretText, hash string
	var createdAt, expireAt time.Time
	row := db.GetDB().QueryRow(sqlStatement, hashKey)
	switch error := row.Scan(&id, &encryptedSecretText, &hash, &createdAt, &expireAt, &remainingViews); error {
	case sql.ErrNoRows:
		return nil, &model.AppError{error, "The hash does not exist!", http.StatusBadRequest}
	case nil:
		fmt.Printf("id=%d secret_text=%s hash=%s created_at=%s expires_at=%s remaining_views=%d\n",
			id, encryptedSecretText, hash, createdAt.String(), expireAt.String(), remainingViews)
	default:
		return nil, &model.AppError{error, "The sever has an issue.", http.StatusInternalServerError}
	}

	if remainingViews > 0 {
		updateRemainingViews(id, remainingViews-1)
	} else {
		return nil, &model.AppError{nil, "Exceed the number of accessing key", http.StatusUnauthorized}
	}

	secretText, err := encrypt.Decrypt(encryptedSecretText)
	if err != nil {
		panic(err)
	}
	return &model.HashKey{
		ID:             id,
		SecretText:     secretText,
		Hash:           hash,
		CreatedAt:      createdAt,
		ExpireAt:       expireAt,
		RemainingViews: remainingViews,
	}, nil
}

func updateRemainingViews(id, remainingViews int) {
	sqlStatement := `UPDATE hashkey SET remaining_views = $1 WHERE id = $2`
	_, err := db.GetDB().Exec(sqlStatement, remainingViews, id)
	if err != nil {
		panic(err)
	}
}

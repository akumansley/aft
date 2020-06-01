package auth

import (
	"awans.org/aft/internal/db"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"hash"
)

var (
	ErrInvalid = fmt.Errorf("%w: invalid token", ErrAuth)
)

// TODO add a timestamp
func TokenForUser(appDB db.DB, user db.Record) (string, error) {
	mac, err := getOrCreateMac(appDB)
	if err != nil {
		return "", err
	}

	bytes, err := user.Id().MarshalBinary()
	if err != nil {
		return "", err
	}
	bytesMac := mac.Sum(bytes)
	binaryToken := append(bytes, bytesMac...)
	token := base64.StdEncoding.EncodeToString(binaryToken)
	return token, nil
}

func UserForToken(appDB db.DB, b64Token string) (db.Record, error) {
	binaryToken, err := base64.StdEncoding.DecodeString(b64Token)
	if err != nil {
		return nil, err
	}
	uuidBytes := binaryToken[:16]
	providedMacBytes := binaryToken[16:]

	mac, err := getOrCreateMac(appDB)
	computedMacBytes := mac.Sum(uuidBytes)

	if !hmac.Equal(providedMacBytes, computedMacBytes) {
		return nil, ErrInvalid
	}
	id, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return nil, err
	}

	tx := appDB.NewTx()
	user, err := tx.FindOne(UserModel.Name, db.Eq("id", id))
	if err != nil {
		return nil, ErrInvalid
	}
	return user, nil
}

func getOrCreateMac(appDB db.DB) (hash.Hash, error) {
	tx := appDB.NewTx()

	rec, err := tx.FindOne(AuthKeyModel.Name, db.Eq("active", true))
	if errors.Is(db.ErrNotFound, err) {
		rec, err = createAuthKey()
		rwtx := appDB.NewRWTx()
		rwtx.Insert(rec)
		rwtx.Commit()
	}
	b64KeyIf := rec.Get("key")
	b64Key := b64KeyIf.(string)
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	return mac, nil
}

func createAuthKey() (db.Record, error) {
	akStore := db.RecordForModel(AuthKeyModel)
	// 128-bit key
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#Session_ID_Length
	c := 16
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	b64Key := base64.StdEncoding.EncodeToString(b)

	akStore.Set("id", uuid.New())
	akStore.Set("active", true)
	akStore.Set("key", b64Key)

	return akStore, nil
}

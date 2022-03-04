package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3"
)

type GetUserOpts struct {
	ID    uint64 `sql:"col:id"`
	Email string `sql:"col:email"`
}

type User struct {
	ID           uint64    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAti   int64     `db:"created_ati" json:"createdAti"`
	DeltedAti    null.Int  `db:"deleted_ati" json:"-"`
	Wallets      []*Wallet `json:"-"`
}

func (user *User) HashPwd() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return nil
}

func (user *User) CmpPwd(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pwd))
}

type Token struct {
	AccessToken string    `json:"accessToken"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

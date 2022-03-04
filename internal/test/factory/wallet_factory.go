package factory

import (
	"cro_test/internal/model"
	"time"

	"github.com/Pallinder/go-randomdata"
	gofactory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/genutil"
)

type UserFactory struct {
	*gofactory.Factory
}

func (us UserFactory) HasWallets(walletF *WalletFactory, num int32) *UserFactory {
	ass := walletF.ToAssociation().ReferField("ID").ForeignField("UserID").ForeignKey("user_id")
	return &UserFactory{
		us.HasMany("Wallets", ass, num),
	}
}

var User = &UserFactory{gofactory.New(
	&model.User{},
	attr.Uint("ID", genutil.SeqUint(1, 1)),
	attr.Str("Email", randomdata.Email),
	attr.Str("PasswordHash", genutil.RandAlph(10)),
	attr.Int("CreatedAti", func() int { return int(time.Now().Unix()) }),
).Table("users"),
}

type WalletFactory struct {
	*gofactory.Factory
}

func (wf WalletFactory) Amount(a string) *WalletFactory {
	return &WalletFactory{
		wf.Attrs(
			attr.Str("Amount", genutil.FixStr(a)),
		),
	}
}

func (wf WalletFactory) Currency(a string) *WalletFactory {
	return &WalletFactory{
		wf.Attrs(
			attr.Str("Currency", genutil.FixStr(a)),
		),
	}
}

var Wallet = &WalletFactory{gofactory.New(
	&model.Wallet{},
	attr.Uint("ID", genutil.SeqUint(1, 1)),
	attr.Str("SerialNumber", genutil.RandUUID()),
	attr.Uint("UserID", genutil.SeqUint(1, 1)),
	attr.Str("Currency", genutil.SeqStrSet("twd", "usd", "xxx", "ooo", "aaa")),
	attr.Str("Amount", genutil.FixStr("0")),
	attr.Int("CreatedAti", func() int { return int(time.Now().Unix()) }),
	attr.Int("UpdatedAti", func() int { return int(time.Now().Unix()) }),
).Table("wallets")}

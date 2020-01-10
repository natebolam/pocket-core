package types

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewValidator(t *testing.T) {
	type args struct {
		addr          sdk.ValAddress
		consPubKey    crypto.PubKey
		tokensToStake sdk.Int
		chains        []string
		serviceURL    string
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want Validator
	}{
		{"defaultValidator", args{sdk.ValAddress(pub.Address()), pub, sdk.ZeroInt(), []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, "google.com"},
			Validator{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.ZeroInt(),
				Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
				ServiceURL:              "google.com",
				UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: bonded
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValidator(tt.args.addr, tt.args.consPubKey, tt.args.chains, tt.args.serviceURL, tt.args.tokensToStake); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ABCIValidatorUpdate(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   abci.ValidatorUpdate
	}{
		{"testingABCIValidatorUpdate Unbonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub),
			Power:  int64(0),
		}},
		{"testingABCIValidatorUpdate Bonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub),
			Power:  int64(0),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ABCIValidatorUpdate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ABCIValidatorUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ABCIValidatorUpdateZero(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   abci.ValidatorUpdate
	}{
		{"testingABCIValidatorUpdate Unbonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.OneInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub),
			Power:  int64(0),
		}},
		{"testingABCIValidatorUpdate Bonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.OneInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub),
			Power:  int64(0),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ABCIValidatorUpdateZero(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ABCIValidatorUpdateZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_AddStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Default Add Token Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.NewInt(100)},
			Validator{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.NewInt(100),
				UnstakingCompletionTime: time.Time{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.AddStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ConsAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / 0 power", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
		{"Default Test / 1 power", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.NewInt(1000000),
			UnstakingCompletionTime: time.Time{},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsensusPower(); got != tt.want {
				t.Errorf("ConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_Equals(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		v2 Validator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Default Test Equal", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, true},
		{"Default Test Not Equal", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.Equals(tt.args.v2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ValAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ValAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsAddr(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ConsAddress(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsPubKey(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   crypto.PubKey
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, pub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsPubKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetConsensusPower(); got != tt.want {
				t.Errorf("GetConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.BondStatus
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Bonded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Int
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ZeroInt()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsJailed(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsJailed(); got != tt.want {
				t.Errorf("IsJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsStaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unbonding false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsStaked(); got != tt.want {
				t.Errorf("IsStaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonding false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonded true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaked(); got != tt.want {
				t.Errorf("IsUnstaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaking(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / bonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unbonding true", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unbonded false", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaking(); got != tt.want {
				t.Errorf("IsUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_PotentialConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / potential power 0", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.PotentialConsensusPower(); got != tt.want {
				t.Errorf("PotentialConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_RemoveStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Remove 0 tokens having 0 tokens ", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.ZeroInt()}, Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Remove 99 tokens having 100 tokens ", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.NewInt(100),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.NewInt(99)}, Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.OneInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.RemoveStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_UpdateStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	type args struct {
		newStatus sdk.BondStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Test Bonded -> Unbonding", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unbonding}, Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unbonding -> Unbonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonding,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unbonded}, Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unbonded -> Bonded", fields{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unbonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Bonded}, Validator{
			Address:                 sdk.ValAddress(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.UpdateStatus(tt.args.newStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
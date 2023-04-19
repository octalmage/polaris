// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package config

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pkg.berachain.dev/polaris/eth/accounts"
)

const (
	// PolarBech32Prefix defines the Bech32 prefix of the local polaris chain.
	PolarBech32Prefix = "polar"
)

// We want to ensure that the config is only being setup conce.
var initConfig sync.Once

// SetupCosmosConfig sets up the Cosmos SDK configuration to be compatible with the semantics of etheruem.
func SetupCosmosConfig() {
	SetupCosmosConfigWith(PolarBech32Prefix, "bera", "abera")
}

// SetupCosmosConfigWith sets up the Cosmos SDK configuration to be compatible with the semantics of etheruem.
func SetupCosmosConfigWith(bech32Prefix, baseDenom, attoDenom string) {
	initConfig.Do(func() {
		// set the address prefixes
		config := sdk.GetConfig()
		SetBech32Prefixes(bech32Prefix, config)
		SetBip44CoinType(config)
		RegisterDenoms(baseDenom, attoDenom)
		config.Seal()
	})
}

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(prefix string, config *sdk.Config) {
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	var (
		bech32PrefixAccAddr = prefix
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
		bech32PrefixAccPub = prefix + sdk.PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
		bech32PrefixValAddr = prefix + sdk.PrefixValidator + sdk.PrefixOperator
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
		bech32PrefixValPub = prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
		bech32PrefixConsAddr = prefix + sdk.PrefixValidator + sdk.PrefixConsensus
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
		bech32PrefixConsPub = prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
	)
	config.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(accounts.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms(baseDenom, attoDenom string) {
	if err := sdk.RegisterDenom(baseDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(attoDenom, sdk.NewDecWithPrec(1, accounts.EtherDecimals)); err != nil {
		panic(err)
	}
}

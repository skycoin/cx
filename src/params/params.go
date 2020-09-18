package params

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

var (
	// MainNetDistribution Skycoin mainnet coin distribution parameters
	MainNetDistribution = Distribution{
		MaxCoinSupply:        [100000000],
		InitialUnlockedCount: 2,
		UnlockAddressRate:    5,
		UnlockTimeInterval:   31536000,
		Addresses: []string{
			"TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX",
			"2PBcLADETphmqWV7sujRZdh3UcabssgKAEB",
		},
	}

	// UserVerifyTxn transaction verification parameters for user-created transactions
	UserVerifyTxn = VerifyTxn{
		// BurnFactor can be overriden with `USER_BURN_FACTOR` env var
		BurnFactor: 100,
		// MaxTransactionSize can be overriden with `USER_MAX_TXN_SIZE` env var
		MaxTransactionSize: 134217728, // in bytes
		// MaxDropletPrecision can be overriden with `USER_MAX_DECIMALS` env var
		MaxDropletPrecision: 3,
	}
)

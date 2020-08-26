package wallet

import "strings"

// CXPrefix is the prefix for cx coin names.
const CXPrefix = "cx_"

// HasCXPrefix returns true if string has 'cx:' prefix.
func HasCXPrefix(s string) bool {
	return strings.HasPrefix(s, CXPrefix)
}

// IsCXCoin returns true if coin type represents that of a cx coin.
func IsCXCoin(coin CoinType) bool {
	return HasCXPrefix(string(coin))
}

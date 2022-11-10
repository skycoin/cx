package loader

import (
	"encoding/json"
	"fmt"
)

type PackageList struct {
	Packages []string
}

func (pl PackageList) MarshalBinary() ([]byte, error) {
	return json.Marshal(pl)
}

func (pl *PackageList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pl)
}

// Encode a package and put it in the specified package list
func (packageList *PackageList) appendPackage(newPackage *Package, database string) error {
	hash, err := newPackage.getHash()
	if err != nil {
		return err
	}
	packageList.Packages = append(packageList.Packages, fmt.Sprintf("%x", hash[:]))
	err = newPackage.saveToDatabase(hash, database)
	if err != nil {
		return err
	}
	return nil
}

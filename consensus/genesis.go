package consensus

import v1 "github.com/attestantio/go-eth2-client/api/v1"

type GenesisResponseJSON struct {
	Data *v1.Genesis `json:"data"`
}

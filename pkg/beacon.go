package pkg

import (
	"context"
	eth2client "github.com/attestantio/go-eth2-client"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (c *ConsensusClient) getClient() (eth2client.Service, error) {
	ctx := context.Background()
	clientAPI, err := http.New(ctx, http.WithAddress(c.APIEndpoint), http.WithLogLevel(zerolog.WarnLevel))
	return clientAPI, err
}

func (c *ConsensusClient) GetBeaconBlockHeader(slot string) (*v1.BeaconBlockHeader, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create eth2client")
	}
	bbh, err := client.(eth2client.BeaconBlockHeadersProvider).BeaconBlockHeader(context.Background(), slot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the beacon block header")
	}
	return bbh, nil
}

func (c *ConsensusClient) GetFinalityCheckpoint(slot string) (*v1.Finality, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create eth2client")
	}
	fcp, err := client.(eth2client.FinalityProvider).Finality(context.Background(), slot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get finality checkpoint")
	}
	return fcp, nil
}

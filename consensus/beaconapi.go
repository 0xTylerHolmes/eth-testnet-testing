package consensus

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

//TODO: Temporary library while PRs are being submitted to go-eth2-client

const (
	TIMEOUT_SECONDS = 10 * time.Second //TODO: make modular with client
)

type RequestError struct {
	// Method the GET/POST
	Method string
	// StatusCode of the response
	StatusCode int
	// Endpoint where we sent the request
	Endpoint string
	// Data included in the response
	Data []byte
}

func (r RequestError) Error() string {
	return fmt.Sprintf("%s error with status-code %d from endpoint %s with %s", r.Method, r.StatusCode, r.Endpoint, r.Data)
}

// postRequest takes the beaconEndpoint, an api, and json data and performs a post request and returns the response
// note this is a temporary workaround
func postRequest(beaconEndpoint string, api string, jsonData []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", beaconEndpoint+api, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create post request")
	}
	// set the content-type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ethereum-testnet-tools/0.1") //TODO: really necessary?

	client := &http.Client{
		Timeout: TIMEOUT_SECONDS,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform POST")
	}
	defer resp.Body.Close() // do we need to handle this error?

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the POST response")
	}

	if resp.StatusCode/100 == 2 {
		// request was successful
		return body, nil
	}
	return body, RequestError{
		Method:     "POST",
		StatusCode: resp.StatusCode,
		Endpoint:   beaconEndpoint,
		Data:       body,
	}
}

// getRequest submits a GET request to the beacon endpoint. We always prefer json
// note this is a temporary work around
func getRequest(beaconEndpoint string, jsonData []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", beaconEndpoint+string(jsonData[:]), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get request")
	}
	// set the content-type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ethereum-testnet-tools/0.1") //TODO: really necessary?

	client := &http.Client{
		Timeout: TIMEOUT_SECONDS,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform GET")
	}
	defer resp.Body.Close() // do we need to handle this error?

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the GET response")
	}

	if resp.StatusCode/100 == 2 {
		// request was successful
		return body, nil
	}
	return body, RequestError{
		Method:     "GET",
		StatusCode: resp.StatusCode,
		Endpoint:   beaconEndpoint,
		Data:       body,
	}
}

// GetBeaconBlockHeader get beacon block header from beacon api
func GetBeaconBlockHeader(beaconEndpoint string, slot string) (*v1.BeaconBlockHeader, error) {
	var beaconBlockHeaderResponse BeaconBlockHeaderResponseJSON
	resp, err := getRequest(beaconEndpoint, []byte(fmt.Sprintf("/eth/v1/beacon/headers/%s", slot)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &beaconBlockHeaderResponse)
	if err != nil {
		return nil, err
	}
	return beaconBlockHeaderResponse.Data, nil
}

// GetCheckpoints get checkpoints from beacon api
func GetCheckpoints(beaconEndpoint string, slot string) (*v1.Finality, error) {
	var checkpointResponse CheckpointResponseJSON
	resp, err := getRequest(beaconEndpoint, []byte(fmt.Sprintf("/eth/v1/beacon/states/%s/finality_checkpoints", slot)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &checkpointResponse)
	if err != nil {
		return nil, err
	}
	return checkpointResponse.Data, nil
}

// PostAttestation post an attestation to beacon api pool
func PostAttestation(beaconEndpoint string, attestation phase0.Attestation) error {
	var attestations []*phase0.Attestation
	attestations = append(attestations, &attestation)
	jsonData, err := json.Marshal(attestations)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json for attestation")
	}
	_, err = postRequest(beaconEndpoint, "/eth/v1/beacon/pool/attestations", jsonData)
	return err
}

// PostAttestorSlashing post an attester slashing to beacon api attester slashings pool
func PostAttestorSlashing(beaconEndpoint string, slashing phase0.AttesterSlashing) error {
	slashingData := AttestorSlashingJSON{
		Attestation1: slashing.Attestation1,
		Attestation2: slashing.Attestation2,
	}
	jsonData, err := json.Marshal(slashingData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json for attestor slashing")
	}
	_, err = postRequest(beaconEndpoint, "/eth/v1/beacon/pool/attester_slashings", jsonData)
	return err
}

// PostProposerSlashing post a proposer slashing to the beacon api proposer slashings pool
func PostProposerSlashing(beaconEndpoint string, slashing phase0.ProposerSlashing) error {
	slashingData := ProposerSlashingJSON{
		SignedHeader1: slashing.SignedHeader1,
		SignedHeader2: slashing.SignedHeader2,
	}
	jsonData, err := json.Marshal(slashingData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json for proposer slashing")
	}
	_, err = postRequest(beaconEndpoint, "/eth/v1/beacon/pool/proposer_slashings", jsonData)
	return err
}

// PostSignedVoluntaryExit post a signed voluntary exit to the beacon api voluntary exits pool
func PostSignedVoluntaryExit(beaconEndpoint string, exit phase0.SignedVoluntaryExit) error {
	exitData := SignedVoluntaryExitJSON{
		Message:   exit.Message,
		Signature: exit.Signature.String(),
	}
	jsonData, err := json.Marshal(exitData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json for voluntary exit")
	}
	_, err = postRequest(beaconEndpoint, "/eth/v1/beacon/pool/voluntary_exits", jsonData)
	return err
}

// PostSignedBLSToExecutionChange post a signed BLS to execution change message to the beacon api BLS to execution change pool
func PostSignedBLSToExecutionChange(beaconEndpoint string, change capella.SignedBLSToExecutionChange) error {
	var blsToExecutionChanges []SignedBLSToExecutionChangeJSON
	blsToExecutionData := SignedBLSToExecutionChangeJSON{
		Message:   change.Message,
		Signature: change.Signature.String(),
	}
	blsToExecutionChanges = append(blsToExecutionChanges, blsToExecutionData)
	jsonData, err := json.Marshal(blsToExecutionChanges)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json for signed bls to execution change")
	}
	_, err = postRequest(beaconEndpoint, "/eth/v1/beacon/pool/bls_to_execution_changes", jsonData)
	return err
}

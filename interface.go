package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"

	"github.com/libsv/go-bt/v2/bscript"
)

// Example demo implementation of a service provider
type opnsServiceProvider struct {
	// Extend your dependencies or custom values
}

type Ordinal struct {
	Outpoint string `json:"outpoint"`
	Height   uint32 `json:"height"`
	Idx      string `json:"idx"`
	Owner    string `json:"owner"`
}

// GetPaymailByAlias is a demo implementation of this interface
func (d *opnsServiceProvider) GetAddressStringByAlias(_ context.Context, alias, domain string) (string, error) {
	query := map[string]interface{}{
		"opns": map[string]interface{}{
			"status": 1,
			"domain": alias,
		},
	}
	buf := bytes.NewBuffer([]byte{})
	address := ""
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return address, err
	}

	if resp, err := http.Post("https://ordinals.gorillapool.io/api/txos/search/unspent", "application/json", buf); err != nil {
		return address, err
	} else {
		var ordinals []*Ordinal
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&ordinals); err != nil {
			return address, err
		}
		if len(ordinals) == 0 {
			return address, fmt.Errorf("not-found")
		}
		address = ordinals[0].Owner
	}
	return address, nil
}

// GetPaymailByAlias is a demo implementation of this interface
func (d *opnsServiceProvider) GetPaymailByAlias(ctx context.Context, alias, domain string,
	_ *server.RequestMetadata,
) (*paymail.AddressInformation, error) {
	if add, err := d.GetAddressStringByAlias(ctx, alias, domain); err != nil {
		return nil, err
	} else {
		return &paymail.AddressInformation{
			Alias:       alias,
			Domain:      domain,
			LastAddress: add,
		}, nil
	}
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (d *opnsServiceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string,
	senderValidation bool, _ *server.RequestMetadata,
) (*paymail.ResolutionPayload, error) {
	// Generate a new destination / output for the basic address resolution
	if add, err := d.GetAddressStringByAlias(ctx, alias, domain); err != nil {
		return nil, err
	} else if p2pkh, err := bscript.NewP2PKHFromAddress(add); err != nil {
		return nil, err
	} else {
		response := &paymail.ResolutionPayload{
			Output: hex.EncodeToString(*p2pkh),
		}
		// if senderValidation {
		// 	if response.Signature, err = bitcoin.SignMessage(
		// 		p.PrivateKey, response.Output, false,
		// 	); err != nil {
		// 		return nil, errors.New("invalid signature: " + err.Error())
		// 	}
		// }
		return response, nil
	}

	// return DemoCreateAddressResolutionResponse(ctx, alias, domain, senderValidation)
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (d *opnsServiceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64, _ *server.RequestMetadata,
) (*paymail.PaymentDestinationPayload, error) {
	// Generate a new destination for the p2p request
	output := &paymail.PaymentOutput{
		Satoshis: satoshis,
	}
	if add, err := d.GetAddressStringByAlias(ctx, alias, domain); err != nil {
		return nil, err
	} else if p2pkh, err := bscript.NewP2PKHFromAddress(add); err != nil {
		return nil, err
	} else {
		output.Script = hex.EncodeToString(*p2pkh)
		output.Address = add
		// Create the response
		return &paymail.PaymentDestinationPayload{
			Outputs:   []*paymail.PaymentOutput{output},
			Reference: "1234567890", // todo: this should be unique per request
		}, nil
	}
}

// RecordTransaction is a demo implementation of this interface
func (d *opnsServiceProvider) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction, _ *server.RequestMetadata,
) (*paymail.P2PTransactionPayload, error) {
	// Record the tx into your datastore layer
	return DemoRecordTransaction(ctx, p2pTx)
}

// VerifyMerkleRoots is a demo implementation of this interface
func (d *opnsServiceProvider) VerifyMerkleRoots(ctx context.Context, merkleProofs []*spv.MerkleRootConfirmationRequestItem) error {
	// Verify the Merkle roots
	return nil
}

func (d *opnsServiceProvider) AddContact(
	ctx context.Context,
	requesterPaymail string,
	contact *paymail.PikeContactRequestPayload,
) error {
	return nil
}

func (d *opnsServiceProvider) CreatePikeOutputResponse(
	ctx context.Context,
	alias, domain, senderPubKey string,
	satoshis uint64,
	metaData *server.RequestMetadata,
) (*paymail.PikePaymentOutputsResponse, error) {
	return nil, nil
}

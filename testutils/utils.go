///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package testutils

import (
	"crypto/rand"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"testing"
)

func LoadPublicKeyTesting(i interface{}) (*rsa.PublicKey, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	privKey, err := LoadPrivateKeyTesting(i)
	if err != nil {
		return nil, errors.Errorf("Could not load private key: %v", err)
	}

	return privKey.GetPublic(), nil
}

func LoadPrivateKeyTesting(i interface{}) (*rsa.PrivateKey, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		return nil, errors.Errorf("Could not load public key: %v", err)
	}

	return privKey, nil

}

func LoadEllipticPublicKey() (*eddsa.PrivateKey, error) {

	ecKey, err := eddsa.NewKeypair(rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Failed to generate new keypair: %v", err)
	}

	return ecKey, nil

}

// Utility function which signs a round info message
func SignRoundInfoRsa(ri *pb.RoundInfo, i interface{}) error {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		return errors.Errorf("Could not load public key: %v", err)
	}

	err = signature.Sign(ri, privKey)
	if err != nil {
		return errors.Errorf("Could not sign round info: %+v", err)
	}
	return nil
}

func SignRoundInfoEddsa(ri *pb.RoundInfo, key *eddsa.PrivateKey, i interface{}) error {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoEddsa is restricted to testing only. Got %T", i)
	}
	err := signature.SignEddsa(ri, key)
	if err != nil {
		return errors.Errorf("Could not sign round info: %+v", err)
	}
	return nil

}

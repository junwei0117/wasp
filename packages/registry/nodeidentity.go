package registry

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/plugins/database"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/key"
)

// NodeIdentityProvider is a subset of the registry interface
// providing access to the persistent node identity information.
type NodeIdentityProvider interface {
	GetNodeIdentity() (*key.Pair, error)
	GetNodePublicKey() (kyber.Point, error)
}

// GetNodeIdentity implements NodeIdentityProvider.
func (r *Impl) GetNodeIdentity() (*key.Pair, error) {
	var err error
	var pair *key.Pair
	var dbKey []byte = dbKeyForNodeIdentity()
	var exists bool
	var data []byte
	partition := database.GetRegistryPartition()
	exists, err = partition.Has(dbKey)
	if !exists {
		pair = key.NewKeyPair(r.keySuite)
		if data, err = keyPairToBytes(pair); err != nil {
			return nil, err
		}
		partition.Set(dbKey, data)
		r.log.Info("Node identity key pair generated.")
		return pair, nil
	}
	if data, err = partition.Get(dbKey); err != nil {
		return nil, err
	}
	if pair, err = keyPairFromBytes(data); err != nil {
		return nil, err
	}
	return pair, nil
}

// GetNodePublicKey implements NodeIdentityProvider.
func (r *Impl) GetNodePublicKey() (kyber.Point, error) {
	var err error
	var pair *key.Pair
	if pair, err = r.GetNodeIdentity(); err != nil {
		return nil, err
	}
	return pair.Public, nil
}

func dbKeyForNodeIdentity() []byte {
	return database.MakeKey(database.ObjectTypeNodeIdentity)
}

func keyPairToBytes(pair *key.Pair) ([]byte, error) {
	var err error
	var w bytes.Buffer
	if err = util.WriteMarshaled(&w, pair.Private); err != nil {
		return nil, err
	}
	if err = util.WriteMarshaled(&w, pair.Public); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func keyPairFromBytes(buf []byte) (*key.Pair, error) {
	var err error
	r := bytes.NewReader(buf)
	pair := key.Pair{}
	if err = util.ReadMarshaled(r, pair.Private); err != nil {
		return nil, err
	}
	if err = util.ReadMarshaled(r, pair.Public); err != nil {
		return nil, err
	}
	return &pair, nil
}
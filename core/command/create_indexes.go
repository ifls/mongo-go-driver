// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package command

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/result"
	"github.com/mongodb/mongo-go-driver/core/session"
	"github.com/mongodb/mongo-go-driver/core/wiremessage"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

// CreateIndexes represents the createIndexes command.
//
// The createIndexes command creates indexes for a namespace.
type CreateIndexes struct {
	NS           Namespace
	Indexes      bson.Arr
	Opts         []bson.Elem
	WriteConcern *writeconcern.WriteConcern
	Clock        *session.ClusterClock
	Session      *session.Client

	result result.CreateIndexes
	err    error
}

// Encode will encode this command into a wire message for the given server description.
func (ci *CreateIndexes) Encode(desc description.SelectedServer) (wiremessage.WireMessage, error) {
	cmd, err := ci.encode(desc)
	if err != nil {
		return nil, err
	}

	return cmd.Encode(desc)
}

func (ci *CreateIndexes) encode(desc description.SelectedServer) (*Write, error) {
	cmd := bson.Doc{
		{"createIndexes", bson.String(ci.NS.Collection)},
		{"indexes", bson.Array(ci.Indexes)},
	}
	cmd = append(cmd, ci.Opts...)

	return &Write{
		Clock:        ci.Clock,
		DB:           ci.NS.DB,
		Command:      cmd,
		WriteConcern: ci.WriteConcern,
		Session:      ci.Session,
	}, nil
}

// Decode will decode the wire message using the provided server description. Errors during decoding
// are deferred until either the Result or Err methods are called.
func (ci *CreateIndexes) Decode(desc description.SelectedServer, wm wiremessage.WireMessage) *CreateIndexes {
	rdr, err := (&Write{}).Decode(desc, wm).Result()
	if err != nil {
		ci.err = err
		return ci
	}

	return ci.decode(desc, rdr)
}

func (ci *CreateIndexes) decode(desc description.SelectedServer, rdr bson.Raw) *CreateIndexes {
	ci.err = bson.Unmarshal(rdr, &ci.result)
	return ci
}

// Result returns the result of a decoded wire message and server description.
func (ci *CreateIndexes) Result() (result.CreateIndexes, error) {
	if ci.err != nil {
		return result.CreateIndexes{}, ci.err
	}
	return ci.result, nil
}

// Err returns the error set on this command.
func (ci *CreateIndexes) Err() error { return ci.err }

// RoundTrip handles the execution of this command using the provided wiremessage.ReadWriter.
func (ci *CreateIndexes) RoundTrip(ctx context.Context, desc description.SelectedServer, rw wiremessage.ReadWriter) (result.CreateIndexes, error) {
	cmd, err := ci.encode(desc)
	if err != nil {
		return result.CreateIndexes{}, err
	}

	rdr, err := cmd.RoundTrip(ctx, desc, rw)
	if err != nil {
		return result.CreateIndexes{}, err
	}

	return ci.decode(desc, rdr).Result()
}

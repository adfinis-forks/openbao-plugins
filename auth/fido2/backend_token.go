// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fido2

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/openbao/openbao/sdk/v2/framework"
	"github.com/openbao/openbao/sdk/v2/helper/roottoken"
	"github.com/openbao/openbao/sdk/v2/logical"
)

func (b *backend) pathSelfServiceEnroll(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entity, err := b.System().EntityInfo(req.EntityID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, logical.CodedError(http.StatusNotFound, "entity not found")
	}

	aliasIdx := slices.IndexFunc(entity.Aliases, func(alias *logical.Alias) bool {
		return alias.MountAccessor == req.MountAccessor
	})

	if aliasIdx == -1 {
		return nil, fmt.Errorf("no alias available for self-service, please ask your authentication admin to create an alias for mount-accessor %q and entity %q", req.MountAccessor, req.EntityID)
	}

	otp, err := roottoken.GenerateOTP(16)
	if err != nil {
		return nil, err
	}

	token := fmt.Sprintf("v1.%s.%s", otp, entity.Aliases[aliasIdx].Name)
	tokenStoragePath := escapedPathJoin("v1", "enrollment-tokens", entity.Aliases[aliasIdx].Name, token)
	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   tokenStoragePath,
		Value: []byte{},
	})
	if err != nil {
		return nil, err
	}

	resp := b.Secret(SecretTypeEnrollmentToken).Response(map[string]any{
		"token": token,
	}, map[string]any{
		"storage_path": tokenStoragePath,
	})
	resp.Secret.Renewable = false
	resp.Secret.TTL = 60 * time.Second

	return resp, nil
}

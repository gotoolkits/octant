/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dashConfigFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/internal/kubeconfig/fake"
)

func Test_kubeContextGenerator(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	kc := &kubeconfig.KubeConfig{
		CurrentContext: "current-context",
	}

	loader := fake.NewMockLoader(controller)
	loader.EXPECT().
		Load("/path").
		Return(kc, nil)

	configLoaderFuncOpt := func(x *ContextsGenerator) {
		x.ConfigLoader = loader
	}

	dashConfig := dashConfigFake.NewMockDash(controller)
	dashConfig.EXPECT().KubeConfigPath().Return("/path")
	dashConfig.EXPECT().ContextName().Return("")

	kgc := NewContextsGenerator(dashConfig, configLoaderFuncOpt)

	assert.Equal(t, "kubeConfig", kgc.Name())

	ctx := context.Background()
	e, err := kgc.Event(ctx)
	require.NoError(t, err)

	assert.Equal(t, event.EventTypeKubeConfig, e.Type)

	resp := kubeContextsResponse{
		CurrentContext: kc.CurrentContext,
	}

	assert.Equal(t, resp, e.Data)
}

//go:build unix

// SPDX-FileCopyrightText: 2022 k0s authors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"

	apcli "github.com/k0sproject/k0s/pkg/autopilot/client"
	apcont "github.com/k0sproject/k0s/pkg/autopilot/controller"
	aproot "github.com/k0sproject/k0s/pkg/autopilot/controller/root"
	"github.com/k0sproject/k0s/pkg/component/manager"
	"github.com/k0sproject/k0s/pkg/config"

	"github.com/k0sproject/k0s/pkg/kubernetes"
	"github.com/sirupsen/logrus"
)

var _ manager.Component = (*Autopilot)(nil)

type Autopilot struct {
	K0sVars            *config.CfgVars
	KubeletExtraArgs   string
	AdminClientFactory kubernetes.ClientFactoryInterface
	Workloads          bool
}

func (a *Autopilot) Init(ctx context.Context) error {
	return nil
}

func (a *Autopilot) Start(ctx context.Context) error {
	log := logrus.WithFields(logrus.Fields{"component": "autopilot"})

	autopilotClientFactory := &apcli.ClientFactory{
		ClientFactoryInterface: a.AdminClientFactory,
	}

	autopilotRoot, err := apcont.NewRootController(aproot.RootConfig{
		InvocationID:        a.K0sVars.InvocationID,
		KubeConfig:          a.K0sVars.AdminKubeConfigPath,
		K0sDataDir:          a.K0sVars.DataDir,
		KubeletExtraArgs:    a.KubeletExtraArgs,
		Mode:                "controller",
		ManagerPort:         8899,
		MetricsBindAddr:     "0",
		HealthProbeBindAddr: "0",
	}, logrus.WithFields(logrus.Fields{"component": "autopilot"}), a.Workloads, a.AdminClientFactory, autopilotClientFactory)
	if err != nil {
		return fmt.Errorf("failed to create autopilot controller: %w", err)
	}

	go func() {
		if err := autopilotRoot.Run(ctx); err != nil {
			log.WithError(err).Error("Error while running autopilot")
		}
	}()

	return nil
}

func (a *Autopilot) Stop() error {
	return nil
}

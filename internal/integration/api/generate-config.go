// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build integration_api

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/siderolabs/talos/internal/integration/base"
	"github.com/siderolabs/talos/pkg/images"
	machineapi "github.com/siderolabs/talos/pkg/machinery/api/machine"
	"github.com/siderolabs/talos/pkg/machinery/client"
	clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// GenerateConfigSuite ...
type GenerateConfigSuite struct {
	base.K8sSuite

	ctx       context.Context //nolint:containedctx
	ctxCancel context.CancelFunc
}

// SuiteName ...
func (suite *GenerateConfigSuite) SuiteName() string {
	return "api.GenerateConfigSuite"
}

// SetupTest ...
func (suite *GenerateConfigSuite) SetupTest() {
	// make sure we abort at some point in time, but give enough room for Recovers
	suite.ctx, suite.ctxCancel = context.WithTimeout(context.Background(), 30*time.Minute)
}

// TearDownTest ...
func (suite *GenerateConfigSuite) TearDownTest() {
	if suite.ctxCancel != nil {
		suite.ctxCancel()
	}
}

// TestGenerate verifies the generated config API.
func (suite *GenerateConfigSuite) TestGenerate() {
	request := &machineapi.GenerateConfigurationRequest{
		ConfigVersion: "v1alpha1",
		MachineConfig: &machineapi.MachineConfig{
			Type: machineapi.MachineConfig_MachineType(machine.TypeInit),
			NetworkConfig: &machineapi.NetworkConfig{
				Hostname: "testhost",
			},
			InstallConfig: &machineapi.InstallConfig{
				InstallDisk:  "/dev/sdb",
				InstallImage: images.DefaultInstallerImage,
			},
			KubernetesVersion: constants.DefaultKubernetesVersion,
		},
		ClusterConfig: &machineapi.ClusterConfig{
			Name: "talos-default",
			ControlPlane: &machineapi.ControlPlaneConfig{
				Endpoint: "http://localhost",
			},
			ClusterNetwork: &machineapi.ClusterNetworkConfig{
				DnsDomain: "cluster.test",
				CniConfig: &machineapi.CNIConfig{
					Name: constants.CustomCNI,
					Urls: []string{
						"https://docs.projectcalico.org/archive/v3.20/manifests/canal.yaml",
					},
				},
			},
		},
	}

	node := suite.RandomDiscoveredNodeInternalIP(machine.TypeControlPlane)
	ctx := client.WithNodes(suite.ctx, node)

	reply, err := suite.Client.GenerateConfiguration(
		ctx,
		request,
	)

	suite.Require().NoError(err)

	data := reply.Messages[0].GetData()

	config, err := configloader.NewFromBytes(data[0])

	suite.Require().NoError(err)

	suite.Require().EqualValues(request.MachineConfig.Type, config.Machine().Type())
	suite.Require().EqualValues(request.ClusterConfig.Name, config.Cluster().Name())
	suite.Require().EqualValues(request.ClusterConfig.ControlPlane.Endpoint, config.Cluster().Endpoint().String())
	suite.Require().EqualValues(request.ClusterConfig.ClusterNetwork.DnsDomain, config.Cluster().Network().DNSDomain())
	suite.Require().EqualValues(
		request.ClusterConfig.ClusterNetwork.CniConfig.Name,
		config.Cluster().Network().CNI().Name(),
	)
	suite.Require().EqualValues(
		request.ClusterConfig.ClusterNetwork.CniConfig.Urls,
		config.Cluster().Network().CNI().URLs(),
	)
	suite.Require().EqualValues(
		fmt.Sprintf("%s:v%s", constants.KubeletImage, request.MachineConfig.KubernetesVersion),
		config.Machine().Kubelet().Image(),
	)
	suite.Require().EqualValues(request.MachineConfig.InstallConfig.InstallDisk, config.Machine().Install().Disk())
	suite.Require().EqualValues(request.MachineConfig.InstallConfig.InstallImage, config.Machine().Install().Image())
	suite.Require().EqualValues(request.MachineConfig.NetworkConfig.Hostname, config.Machine().Network().Hostname())
	suite.Require().EqualValues(request.MachineConfig.NetworkConfig.Hostname, config.Machine().Network().Hostname())

	talosconfig, err := clientconfig.FromBytes(reply.Messages[0].Talosconfig)

	suite.Require().NoError(err)

	context := talosconfig.Contexts[request.ClusterConfig.Name]

	suite.Require().NotNil(context)

	suite.Require().NotEmpty(context.CA)
	suite.Require().NotEmpty(context.Crt)
	suite.Require().NotEmpty(context.Key)
	suite.Require().Greater(len(context.Endpoints), 0)
	suite.Require().EqualValues(context.Endpoints[0], config.Cluster().Endpoint().Hostname())

	// now generate control plane join config
	request.MachineConfig.Type = machineapi.MachineConfig_MachineType(machine.TypeControlPlane)

	reply, err = suite.Client.GenerateConfiguration(
		ctx,
		request,
	)

	suite.Require().NoError(err)

	data = reply.Messages[0].GetData()

	joinedConfig, err := configloader.NewFromBytes(data[0])

	suite.Require().NoError(err)

	suite.Require().EqualValues(request.MachineConfig.Type, joinedConfig.Machine().Type())
	suite.Require().EqualValues(request.ClusterConfig.Name, joinedConfig.Cluster().Name())
	suite.Require().EqualValues(request.ClusterConfig.ControlPlane.Endpoint, joinedConfig.Cluster().Endpoint().String())
	suite.Require().EqualValues(
		request.ClusterConfig.ClusterNetwork.DnsDomain,
		joinedConfig.Cluster().Network().DNSDomain(),
	)
	suite.Require().EqualValues(
		fmt.Sprintf("%s:v%s", constants.KubeletImage, request.MachineConfig.KubernetesVersion),
		joinedConfig.Machine().Kubelet().Image(),
	)
	suite.Require().EqualValues(request.MachineConfig.InstallConfig.InstallDisk, config.Machine().Install().Disk())
	suite.Require().EqualValues(
		request.MachineConfig.InstallConfig.InstallImage,
		joinedConfig.Machine().Install().Image(),
	)
	suite.Require().EqualValues(
		request.MachineConfig.NetworkConfig.Hostname,
		joinedConfig.Machine().Network().Hostname(),
	)

	suite.Require().EqualValues(config.Machine().Security().IssuingCA(), joinedConfig.Machine().Security().IssuingCA())
	suite.Require().EqualValues(config.Machine().Security().Token(), joinedConfig.Machine().Security().Token())
	suite.Require().EqualValues(
		config.Cluster().AESCBCEncryptionSecret(),
		joinedConfig.Cluster().AESCBCEncryptionSecret(),
	)
	suite.Require().EqualValues(config.Cluster().IssuingCA(), joinedConfig.Cluster().IssuingCA())
	suite.Require().EqualValues(config.Cluster().Token(), joinedConfig.Cluster().Token())
	suite.Require().EqualValues(config.Cluster().Etcd().CA(), config.Cluster().Etcd().CA())
}

func init() {
	allSuites = append(allSuites, new(GenerateConfigSuite))
}

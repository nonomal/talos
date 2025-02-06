// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package docker

import (
	"context"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"github.com/siderolabs/talos/pkg/provision"
)

//nolint:gocyclo
func (p *provisioner) Reflect(ctx context.Context, clusterName, stateDirectory string) (provision.Cluster, error) {
	res := &result{
		clusterInfo: provision.ClusterInfo{
			ClusterName: clusterName,
		},
		statePath: stateDirectory,
	}

	// find network assuming network name == cluster name
	networks, err := p.listNetworks(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	if len(networks) > 0 {
		network := networks[0]

		var cidr netip.Prefix

		cidr, err = netip.ParsePrefix(network.IPAM.Config[0].Subnet)
		if err != nil {
			return nil, err
		}

		res.clusterInfo.Network.Name = network.Name
		res.clusterInfo.Network.CIDRs = []netip.Prefix{cidr}
		res.clusterInfo.Network.GatewayAddrs = []netip.Addr{}

		var addr netip.Addr

		if addr, err = netip.ParseAddr(network.IPAM.Config[0].Gateway); err == nil {
			res.clusterInfo.Network.GatewayAddrs = append(res.clusterInfo.Network.GatewayAddrs, addr)
		}

		mtuStr, ok := network.Options["com.docker.network.driver.mtu"]
		if !ok {
			mtuStr = network.Options["mtu"] // Use the podman version of the option
		}

		res.clusterInfo.Network.MTU, err = strconv.Atoi(mtuStr)
		if err != nil {
			return nil, err
		}
	}

	// find nodes (containers)
	nodes, err := p.listNodes(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	var mappedKubernetesPort string

	for _, node := range nodes {
		t, err := machine.ParseType(node.Labels["talos.type"])
		if err != nil {
			return nil, err
		}

		container, err := p.client.ContainerInspect(ctx, node.ID)
		if err != nil {
			return nil, err
		}

		var ips []netip.Addr

		if network, ok := node.NetworkSettings.Networks[res.clusterInfo.Network.Name]; ok {
			ip := network.IPAddress
			if ip == "" && network.IPAMConfig != nil {
				ip = network.IPAMConfig.IPv4Address
			}

			addr, err := netip.ParseAddr(ip)
			if err != nil {
				return nil, err
			}

			ips = append(ips, addr)
		}

		for port, portBinding := range container.HostConfig.PortBindings {
			if port.Int() == constants.DefaultControlPlanePort {
				for _, binding := range portBinding {
					mappedKubernetesPort = binding.HostPort
				}
			}
		}

		res.clusterInfo.Nodes = append(res.clusterInfo.Nodes,
			provision.NodeInfo{
				ID:   node.ID,
				Name: strings.TrimLeft(node.Names[0], "/"),
				Type: t,

				IPs: ips,

				NanoCPUs: container.HostConfig.Resources.NanoCPUs,
				Memory:   container.HostConfig.Resources.Memory,
			})
	}

	if mappedKubernetesPort != "" {
		res.clusterInfo.KubernetesEndpoint = "https://" + net.JoinHostPort("127.0.0.1", mappedKubernetesPort)
	}

	return res, nil
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/mdlayher/ethtool"
	"github.com/siderolabs/gen/optional"
	"go.uber.org/zap"

	"github.com/siderolabs/talos/internal/app/machined/pkg/controllers/runtime"
	"github.com/siderolabs/talos/pkg/machinery/resources/network"
)

// EthernetSpecController reports Ethernet link statuses.
type EthernetSpecController struct{}

// Name implements controller.Controller interface.
func (ctrl *EthernetSpecController) Name() string {
	return "network.EthernetSpecController"
}

// Inputs implements controller.Controller interface.
func (ctrl *EthernetSpecController) Inputs() []controller.Input {
	return []controller.Input{
		{
			Namespace: network.NamespaceName,
			Type:      network.EthernetSpecType,
			Kind:      controller.InputWeak,
		},
	}
}

// Outputs implements controller.Controller interface.
func (ctrl *EthernetSpecController) Outputs() []controller.Output {
	return nil
}

// Run implements controller.Controller interface.
func (ctrl *EthernetSpecController) Run(ctx context.Context, r controller.Runtime, logger *zap.Logger) error {
	// wait for udevd to be healthy, which implies that all link renames are done
	if err := runtime.WaitForDevicesReady(ctx, r,
		[]controller.Input{
			{
				Namespace: network.NamespaceName,
				Type:      network.EthernetSpecType,
				Kind:      controller.InputWeak,
			},
		},
	); err != nil {
		return err
	}

	ethClient, err := ethtool.New()
	if err != nil {
		logger.Warn("error dialing ethtool socket", zap.Error(err))

		return nil
	}

	defer ethClient.Close() //nolint:errcheck

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.EventCh():
		}

		specs, err := safe.ReaderListAll[*network.EthernetSpec](ctx, r)
		if err != nil {
			return fmt.Errorf("error reading EthernetSpec resources: %w", err)
		}

		var errs error

		for spec := range specs.All() {
			if err = ctrl.apply(ethClient, spec); err != nil {
				errs = errors.Join(errs, fmt.Errorf("error configuring %q: %w", spec.Metadata().ID(), err))
			}
		}

		if errs != nil {
			return fmt.Errorf("failed to reconcile Ethernet specs: %w", errs)
		}

		r.ResetRestartBackoff()
	}
}

func optionalFromPtr[T any](ptr *T) optional.Optional[T] {
	if ptr == nil {
		return optional.None[T]()
	}

	return optional.Some(*ptr)
}

func (ctrl *EthernetSpecController) apply(
	ethClient *ethtool.Client,
	spec *network.EthernetSpec,
) error {
	ringSpec := spec.TypedSpec().Rings

	var zeroRingSpec network.EthernetRingsSpec

	if ringSpec != zeroRingSpec {
		if err := ethClient.SetRings(ethtool.Rings{
			Interface: ethtool.Interface{
				Name: spec.Metadata().ID(),
			},
			RX:           optionalFromPtr(ringSpec.RX),
			RXMini:       optionalFromPtr(ringSpec.RXMini),
			RXJumbo:      optionalFromPtr(ringSpec.RXJumbo),
			TX:           optionalFromPtr(ringSpec.TX),
			RXBufLen:     optionalFromPtr(ringSpec.RXBufLen),
			CQESize:      optionalFromPtr(ringSpec.CQESize),
			TXPush:       optionalFromPtr(ringSpec.TXPush),
			RXPush:       optionalFromPtr(ringSpec.RXPush),
			TXPushBufLen: optionalFromPtr(ringSpec.TXPushBufLen),
			TCPDataSplit: optionalFromPtr(ringSpec.TCPDataSplit),
		}); err != nil {
			return fmt.Errorf("error updating rings: %w", err)
		}
	}

	return nil
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qemu

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/hashicorp/go-getter/v2"

	"github.com/siderolabs/talos/pkg/machinery/constants"
	"github.com/siderolabs/talos/pkg/provision"
)

func (p *provisioner) preflightChecks(ctx context.Context, request provision.ClusterRequest, options provision.Options, arch Arch) error {
	checkContext := preflightCheckContext{
		request: request,
		options: options,
		arch:    arch,
	}

	for _, check := range []func(ctx context.Context) error{
		checkContext.verifyRoot,
		checkContext.checkKVM,
		checkContext.qemuExecutable,
		checkContext.checkFlashImages,
		checkContext.swtpmExecutable,
		checkContext.cniDirectories,
		checkContext.cniBundle,
		checkContext.checkIptables,
	} {
		if err := check(ctx); err != nil {
			return err
		}
	}

	return nil
}

type preflightCheckContext struct {
	request provision.ClusterRequest
	options provision.Options
	arch    Arch
}

func (check *preflightCheckContext) verifyRoot(context.Context) error {
	if os.Geteuid() != 0 {
		return errors.New("error: please run as root user (CNI requirement), we recommend running with `sudo -E`")
	}

	return nil
}

func (check *preflightCheckContext) checkKVM(context.Context) error {
	f, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("error opening /dev/kvm, please make sure KVM support is enabled in Linux kernel: %w", err)
	}

	return f.Close()
}

func (check *preflightCheckContext) qemuExecutable(context.Context) error {
	if check.arch.QemuExecutable() == "" {
		return fmt.Errorf("QEMU executable (qemu-system-%s or qemu-kvm) not found, please install QEMU with package manager", check.arch.QemuArch())
	}

	return nil
}

func (check *preflightCheckContext) checkFlashImages(context.Context) error {
	for _, flashImage := range check.arch.PFlash(check.options.UEFIEnabled, check.options.ExtraUEFISearchPaths) {
		if len(flashImage.SourcePaths) == 0 {
			continue
		}

		found := false

		for _, path := range flashImage.SourcePaths {
			_, err := os.Stat(path)
			if err == nil {
				found = true

				break
			}
		}

		if !found {
			return fmt.Errorf("the required flash image was not found in any of the expected paths for (%q), "+
				"please install it with the package manager or specify --extra-uefi-search-paths", flashImage.SourcePaths)
		}
	}

	return nil
}

func (check *preflightCheckContext) swtpmExecutable(context.Context) error {
	if check.options.TPM2Enabled {
		if _, err := exec.LookPath("swtpm"); err != nil {
			return fmt.Errorf("swtpm not found in PATH, please install swtpm-tools with the package manager: %w", err)
		}
	}

	return nil
}

func (check *preflightCheckContext) cniDirectories(context.Context) error {
	cniDirs := slices.Clone(check.request.Network.CNI.BinPath)
	cniDirs = append(cniDirs, check.request.Network.CNI.CacheDir, check.request.Network.CNI.ConfDir)

	for _, cniDir := range cniDirs {
		st, err := os.Stat(cniDir)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("error checking CNI directory %q: %w", cniDir, err)
			}

			fmt.Fprintf(check.options.LogWriter, "creating %q\n", cniDir)

			err = os.MkdirAll(cniDir, 0o777)
			if err != nil {
				return err
			}

			continue
		}

		if !st.IsDir() {
			return fmt.Errorf("CNI path %q exists, but it's not a directory", cniDir)
		}
	}

	return nil
}

func (check *preflightCheckContext) cniBundle(ctx context.Context) error {
	var missing bool

	requiredCNIPlugins := []string{"bridge", "firewall", "static", "tc-redirect-tap"}

	for _, cniPlugin := range requiredCNIPlugins {
		missing = true

		for _, binPath := range check.request.Network.CNI.BinPath {
			_, err := os.Stat(filepath.Join(binPath, cniPlugin))
			if err == nil {
				missing = false

				break
			}
		}

		if missing {
			break
		}
	}

	if !missing {
		return nil
	}

	if check.request.Network.CNI.BundleURL == "" {
		return fmt.Errorf("error: required CNI plugins %q were not found in %q", requiredCNIPlugins, check.request.Network.CNI.BinPath)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	client := getter.Client{}
	src := strings.ReplaceAll(check.request.Network.CNI.BundleURL, constants.ArchVariable, runtime.GOARCH)
	dst := check.request.Network.CNI.BinPath[0]

	fmt.Fprintf(check.options.LogWriter, "downloading CNI bundle from %q to %q\n", src, dst)

	_, err = client.Get(ctx, &getter.Request{
		Src:     src,
		Dst:     dst,
		Pwd:     pwd,
		GetMode: getter.ModeDir,
	})

	return err
}

func (check *preflightCheckContext) checkIptables(ctx context.Context) error {
	_, err := iptables.New()
	if err != nil {
		return fmt.Errorf("error accessing iptables: %w", err)
	}

	return nil
}

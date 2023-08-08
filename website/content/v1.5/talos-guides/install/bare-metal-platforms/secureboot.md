---
title: "SecureBoot"
description: "Booting Talos in SecureBoot mode on UEFI platforms."
---

Talos now supports booting on UEFI systems in SecureBoot mode.
When combined with TPM-based disk encryption, this provides [Trusted Boot](https://0pointer.net/blog/brave-new-trusted-boot-world.html) experience.

> Note: SecureBoot is not supported on x86 platforms in BIOS mode.

The implementation is using [systemd-boot](https://www.freedesktop.org/wiki/Software/systemd/systemd-boot/) as a boot menu implementation, while the
Talos kernel, initramfs and cmdline arguments are combined into the [Unified Kernel Image](https://uapi-group.org/specifications/specs/unified_kernel_image/) (UKI) format.
UEFI firmware loads the `systemd-boot` bootloader, which then loads the UKI image.
Both `systemd-boot` and Talos `UKI` image are signed with the key, which is enrolled into the UEFI firmware.

As Talos Linux is fully contained in the UKI image, the full operating system is verified and booted by the UEFI firmware.

> Note: There is no support at the moment to upgrade non-UKI (GRUB-based) Talos installation to use UKI/SecureBoot, so a fresh installation is required.

## Acquiring the SecureBoot Assets

Sidero Labs is going to provide Talos images signed with the Sidero Labs Secureboot key in August 2023, but until then, you can use the following instructions to build your own.

## SecureBoot with Custom Keys

### Generating the Keys

Talos requires two set of keys to be used for the SecureBoot process:

* SecureBoot key is used to sign the boot assets and it is enrolled into the UEFI firmware.
* PCR Signing Key is used to sign the TPM policy, which is used to seal the disk encryption key.

The same key might be used for both, but it is recommended to use separate keys for each purpose.

Talos provides a utility to generate the keys, but existing PKI infrastructure can be used as well:

```shell
$ talosctl gen secureboot uki --common-name "SecureBoot Key"
writing _out/uki-signing-cert.pem
writing _out/uki-signing-key.pem
```

The generated certificate and private key are written to disk in PEM-encoded format (RSA 4096-bit key).

PCR signing key can be generated with:

```shell
$ talosctl gen secureboot pcr
writing _out/pcr-signing-key.pem
writing _out/pcr-signing-public-key.pem
```

Two files containing the private and public key are written to disk in PEM-encoded format (RSA 2048-bit key).

Optionally, UEFI automatic key enrollment database can be generated using the `_out/uki-signing-*` files as input:

```shell
$ talosctl gen secureboot database
writing _out/db.auth
writing _out/KEK.auth
writing _out/PK.auth
```

These files can be used to enroll the keys into the UEFI firmware automatically when booting from a SecureBoot ISO while UEFI firmware is in the setup mode.

### Generating the SecureBoot Assets

Once the keys are generated, they can be used to sign the Talos boot assets to generate required ISO images, PXE boot assets, disk images, installer containers, etc.
In this guide we will generate a SecureBoot ISO image and an installer image.

```shell
$ docker run --rm -t -v $PWD/_out:/secureboot:ro -v $PWD/_out:/out ghcr.io/siderolabs/imager:{{< release >}} secureboot-iso
profile ready:
arch: amd64
platform: metal
secureboot: true
version: {{< release >}}
input:
  kernel:
    path: /usr/install/amd64/vmlinuz
  initramfs:
    path: /usr/install/amd64/initramfs.xz
  sdStub:
    path: /usr/install/amd64/systemd-stub.efi
  sdBoot:
    path: /usr/install/amd64/systemd-boot.efi
  baseInstaller:
    imageRef: ghcr.io/siderolabs/installer:v1.5.0-alpha.3-35-ge0f383598-dirty
  secureboot:
    signingKeyPath: /secureboot/uki-signing-key.pem
    signingCertPath: /secureboot/uki-signing-cert.pem
    pcrSigningKeyPath: /secureboot/pcr-signing-key.pem
    pcrPublicKeyPath: /secureboot/pcr-signing-public-key.pem
    platformKeyPath: /secureboot/PK.auth
    keyExchangeKeyPath: /secureboot/KEK.auth
    signatureKeyPath: /secureboot/db.auth
output:
  kind: iso
  outFormat: raw
skipped initramfs rebuild (no system extensions)
kernel command line: talos.platform=metal console=ttyS0 console=tty0 init_on_alloc=1 slab_nomerge pti=on consoleblank=0 nvme_core.io_timeout=4294967295 printk.devkmsg=on ima_template=ima-ng ima_appraise=fix ima_hash=sha512 lockdown=confidentiality
UKI ready
ISO ready
output asset path: /out/metal-amd64-secureboot.iso
```

Next, the installer image should be generated to install Talos to disk on a SecureBoot-enabled system:

```shell
$ docker run --rm -t -v $PWD/_out:/secureboot:ro -v $PWD/_out:/out ghcr.io/siderolabs/imager:{{< release >}} secureboot-installer
profile ready:
arch: amd64
platform: metal
secureboot: true
version: {{< release >}}
input:
  kernel:
    path: /usr/install/amd64/vmlinuz
  initramfs:
    path: /usr/install/amd64/initramfs.xz
  sdStub:
    path: /usr/install/amd64/systemd-stub.efi
  sdBoot:
    path: /usr/install/amd64/systemd-boot.efi
  baseInstaller:
    imageRef: ghcr.io/siderolabs/installer:{{< release >}}
  secureboot:
    signingKeyPath: /secureboot/uki-signing-key.pem
    signingCertPath: /secureboot/uki-signing-cert.pem
    pcrSigningKeyPath: /secureboot/pcr-signing-key.pem
    pcrPublicKeyPath: /secureboot/pcr-signing-public-key.pem
    platformKeyPath: /secureboot/PK.auth
    keyExchangeKeyPath: /secureboot/KEK.auth
    signatureKeyPath: /secureboot/db.auth
output:
  kind: installer
  outFormat: raw
skipped initramfs rebuild (no system extensions)
kernel command line: talos.platform=metal console=ttyS0 console=tty0 init_on_alloc=1 slab_nomerge pti=on consoleblank=0 nvme_core.io_timeout=4294967295 printk.devkmsg=on ima_template=ima-ng ima_appraise=fix ima_hash=sha512 lockdown=confidentiality
UKI ready
installer container image ready
output asset path: /out/metal-amd64-secureboot-installer.tar
```

The generated container image should be pushed to some container registry which Talos can access during the installation, e.g.:

```shell
crane push _out/metal-amd64-secureboot-installer.tar ghcr.io/<user>/installer-amd64-secureboot:{{< release >}}
```

The generated ISO and installer images might be further customized with system extensions, extra kernel command line arguments, etc.

## Booting Talos Linux in SecureBoot Mode

In this guide we will use the generated ISO image to boot Talos Linux in SecureBoot mode, followed by submitting machine configuration to the machine in maintenance mode.
We will use one the ways to generate and submit machine configuration to the node, please refer to the [Production Notes]({{< relref "../../../introduction/prodnotes" >}}) for the full guide.

First, make sure SecureBoot is enabled in the UEFI firmware.
For the first boot, the UEFI firmware should be in the setup mode, so that the keys can be enrolled into the UEFI firmware automatically.

> Note: There are other ways to enroll the keys into the UEFI firmware, but this is out of scope of this guide.

Once Talos is running in maintenance mode, verify that secure boot is enabled:

```shell
$ talosctl -n <IP> get securitystate --insecure
NODE   NAMESPACE   TYPE            ID              VERSION   SECUREBOOT
       runtime     SecurityState   securitystate   1         true
```

Now we will generate the machine configuration for the node supplying the `installer` container image we pushed earlier, and applying the patch to enable TPM-based disk encryption (requires TPM 2.0):

```yaml
# tpm-disk-encryption.yaml
machine:
  systemDiskEncryption:
    ephemeral:
      provider: luks2
      keys:
        - slot: 0
          tpm: {}
    state:
      provider: luks2
      keys:
        - slot: 0
          tpm: {}
```

Generate machine configuration:

```shell
talosctl gen config <cluster-name> https://<endpoint>:6443 --install-image=ghcr.io/<user>/installer-amd64-secureboot:{{< release >}} --install-disk=/dev/sda --config-patch @tpm-disk-encryption.yaml
```

Apply machine configuration to the node:

```shell
talosctl -n <IP> apply-config --insecure -f controlplane.yaml
```

Talos will perform the installation to the disk and reboot the node.
Please make sure that the ISO image is not attached to the node anymore, otherwise the node will boot from the ISO image again.

Once the node is rebooted, verify that the node is running in secure boot mode:

```shell
talosctl -n <IP> --talosconfig=talosconfig get securitystate
```

## Upgrading Talos Linux

Any change to the boot asset (kernel, initramfs, kernel command line) requires the UKI to be regenerated and the installer image to be rebuilt.
Follow the steps above to generate new installer image updating the boot assets: use new Talos version, add a system extension, or modify the kernel command line.
Once the new `installer` image is pushed to the registry, [upgrade]({{< relref "../../upgrading-talos" >}}) the node using the new installer image.

It is important to preserve the UKI signing key and the PCR signing key, otherwise the node will not be able to boot with the new UKI and unlock the encrypted partitions.

## Disk Encryption with TPM

When encrypting the disk partition for the first time, Talos Linux generates a random disk encryption key and seals (encrypts) it with the TPM device.
The TPM unlock policy is configured to trust the expected policy signed by the PCR signing key.
This way TPM unlocking doesn't depend on the exact [PCR measurements](https://uapi-group.org/specifications/specs/linux_tpm_pcr_registry/), but rather on the expected policy signed by the PCR signing key and the state of SecureBoot (PCR 7 measurement, including secureboot status and the list of enrolled keys).

When the UKI image is generated, the UKI is measured and expected measurements are combined into TPM unlock policy and signed with the PCR signing key.
During the boot process, `systemd-stub` component of the UKI performs measurements of the UKI sections into the TPM device.
Talos Linux during the boot appends to the PCR register the measurements of the boot phases, and once the boot reaches the point of mounting the encrypted disk partition,
the expected signed policy from the UKI is matched against measured values to unlock the TPM, and TPM unseals the disk encryption key which is then used to unlock the disk partition.

During the upgrade, as long as the new UKI is contains PCR policy signed with the same PCR signing key, and SecureBoot state has not changed the disk partition will be unlocked successfully.

Disk encryption is also tied to the state of PCR register 7, so that it unlocks only if SecureBoot is enabled and the set of enrolled keys hasn't changed.

## Other Boot Options

Unified Kernel Image (UKI) is a UEFI-bootable image which can be booted directly from the UEFI firmware skipping the `systemd-boot` bootloader.
In network boot mode, the UKI can be used directly as well, as it contains the full set of boot assets required to boot Talos Linux.

When SecureBoot is enabled, the UKI image ignores any kernel command line arguments passed to it, but rather uses the kernel command line arguments embedded into the UKI image itself.
If kernel command line arguments need to be changed, the UKI image needs to be rebuilt with the new kernel command line arguments.
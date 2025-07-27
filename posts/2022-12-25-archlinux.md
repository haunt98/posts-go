# Install Arch Linux

Install Arch Linux is thing I always want to do for my laptop/PC since I had my laptop in ninth grade.

This is not a guide for everyone, this is just save for myself in a future and for anyone who want to walk in my shoes.

## [Installation guide](https://wiki.archlinux.org/index.php/Installation_guide)

### Pre-installation

Check disks carefully:

```sh
lsblk
```

[USB flash installation medium](https://wiki.archlinux.org/index.php/USB_flash_installation_medium)

#### Verify the boot mode

Check UEFI mode:

```sh
cat /sys/firmware/efi/fw_platform_size
# 64 or 32 is UEFI
# No such file or directory is BIOS
```

#### Connect to the internet

For wifi, use [iwd](https://wiki.archlinux.org/index.php/Iwd).

#### Partition the disks

[GPT fdisk](https://wiki.archlinux.org/index.php/GPT_fdisk):

```sh
cgdisk /dev/sdx
```

- [Partition scheme](https://wiki.archlinux.org/index.php/Partitioning#Partition_scheme)
- [EFI system partition](https://wiki.archlinux.org/title/EFI_system_partition)

UEFI/GPT layout:

| Mount point | Partition                             | Partition type                            | Suggested size | gdisk code |
| ----------- | ------------------------------------- | ----------------------------------------- | -------------- | ---------- |
| `/mnt/efi`  | `/dev/efi_system_partition`           | EFI System Partition                      | 512 MiB        | EF00       |
| `/mnt/boot` | `/dev/extended_boot_loader_partition` | Extended Boot Loader Partition (XBOOTLDR) | 1 GiB          | EA00       |
| `/mnt`      | `/dev/root_partition`                 | Root Partition                            |                | 8300       |

Why not `/boot/efi`? See
[Lennart Poettering comment](https://github.com/systemd/systemd/pull/3757#issuecomment-234290236).

BIOS/GPT layout:

| Mount point | Partition             | Partition type      | Suggested size | gdisk code |
| ----------- | --------------------- | ------------------- | -------------- | ---------- |
|             |                       | BIOS boot partition | 1 MiB          | EF02       |
| `/mnt`      | `/dev/root_partition` | Root Partition      |                | 8300       |

Format:

```sh
# efi
mkfs.fat -F32 /dev/efi_system_partition

# boot
mkfs.fat -F32 /dev/extended_boot_loader_partition

# root
mkfs.ext4 -L ROOT /dev/root_partition

# root with btrfs (optional)
mkfs.btrfs -L ROOT /dev/root_partition
```

Mount:

```sh
# root
mount /dev/root_partition /mnt

# root with btrfs (optional)
mount -o compress=zstd /dev/root_partition /mnt

# efi
mount --mkdir /dev/efi_system_partition /mnt/efi

# boot
mount --mkdir /dev/extended_boot_loader_partition /mnt/boot
```

### Installation

Please check [Mirrors](https://wiki.archlinux.org/title/Mirrors) if you have slow Internet.

```sh
pacstrap -K /mnt base linux linux-firmware

# AMD (optional)
pacstrap -K /mnt amd-ucode

# Intel (optional)
pacstrap -K /mnt intel-ucode

# Btrfs (optional)
pacstrap -K /mnt btrfs-progs

# Text editor
pacstrap -K /mnt neovim
```

### Configure

#### [fstab](https://wiki.archlinux.org/index.php/Fstab)

```sh
genfstab -U /mnt >> /mnt/etc/fstab
```

#### Chroot

```sh
arch-chroot /mnt
```

#### Time

```sh
# Change Region/City to your location
ln -sf /usr/share/zoneinfo/Region/City /etc/localtime

hwclock --systohc
```

#### Localization

Edit `/etc/locale.gen` then uncomment `# en_US.UTF-8 UTF-8` by removing `#` at the beginning.

Generate locales:

```sh
locale-gen
```

Edit `/etc/locale.conf`:

```txt
LANG=en_US.UTF-8
```

#### Network configuration

Edit `/etc/hostname`:

```txt
myhostname
```

#### Initramfs

Edit `/etc/mkinitcpio.conf`:

```txt
# https://wiki.archlinux.org/title/mkinitcpio#Common_hooks
# Replace base udev with systemd
#
HOOKS=(systemd ... )
```

```sh
mkinitcpio -P
```

#### Root password

```sh
passwd
```

#### [NetworkManager](https://wiki.archlinux.org/title/NetworkManager)

```sh
pacman -Syu networkmanager dhcpcd iwd
systemctl enable NetworkManager.service
systemctl enable systemd-resolved.service
```

Edit `/etc/NetworkManager/conf.d/wifi_backend.conf`:

```txt
[device]
wifi.backend=iwd
```

Edit `/etc/NetworkManager/conf.d/dns.conf`:

```txt
[main]
dns=systemd-resolved
```

#### [Bluetooth](https://wiki.archlinux.org/title/Bluetooth)

```sh
pacman -Syu bluez
systemctl enable bluetooth.service
```

#### Clock

Use [systemd-timesyncd](https://wiki.archlinux.org/title/Systemd-timesyncd)

```sh
timedatectl set-ntp true

timedatectl status
```

#### Boot loader

Use [systemd-boot](https://wiki.archlinux.org/index.php/Systemd-boot)

Install using XBOOTLDR:

```sh
bootctl --esp-path=/efi --boot-path=/boot install

systemctl enable systemd-boot-update.service
```

[Label partition](https://wiki.archlinux.org/index.php/persistent_block_device_naming#by-label)

Edit `/efi/loader/loader.conf`:

```txt
default	archlinux.conf
timeout 4
editor no
console-mode max
```

Edit `/boot/loader/entries/archlinux.conf`:

```txt
title Arch Linux
linux /vmlinuz-linux

# Intel (optional)
initrd /intel-ucode.img

# AMD (optional)
initrd /amd-ucode.img

initrd /initramfs-linux.img

# Kernel parameters (optional)
#
# Acer Nitro AN515-45
# https://wiki.archlinux.org/title/backlight#Kernel_command-line_options
# acpi_backlight=vendor
#
# NVIDIA
# https://wiki.archlinux.org/title/NVIDIA#DRM_kernel_mode_setting
# nvidia-drm.modeset=1
#
options root="LABEL=ROOT" rw quiet loglevel=3 nowatchdog module_blacklist=iTCO_wdt,sp5100_tco ipv6.disable=1 init_on_alloc=1 init_on_free=1 page_alloc.shuffle=1
```

## [General recommendations](https://wiki.archlinux.org/index.php/General_recommendations)

Always remember to check **dependencies** when install packages.

### System administration

[Sudo](https://wiki.archlinux.org/index.php/sudo):

```sh
pacman -Syu sudo zsh

EDITOR=nvim visudo
# Uncomment group wheel by removing % at the beginning of %wheel ...

# Add user if don't want to use systemd-homed
useradd -m -G wheel -c "The Joker" joker

# Or using zsh (optional)
useradd -m -G wheel -s /usr/bin/zsh -c "The Joker" joker

# Set password
passwd joker
```

- [systemd-homed (optional if no useradd before)](https://wiki.archlinux.org/index.php/Systemd-homed):
- [Home Directories](https://systemd.io/HOME_DIRECTORY/)

```sh
systemctl enable systemd-homed.service

homectl create joker --real-name="The Joker" --member-of=wheel

# Using zsh (optional)
homectl update joker --shell=/usr/bin/zsh
```

**Note**: Can not run `homectl` when install Arch Linux. Should run on the first boot.

### Desktop Environment

GPU drivers:

- [Intel graphics](https://wiki.archlinux.org/title/Intel_graphics)
- [NVIDIA](https://wiki.archlinux.org/title/NVIDIA)
- [AMDGPU](https://wiki.archlinux.org/title/AMDGPU)

#### [KDE](https://wiki.archlinux.org/title/KDE)

See [KDE Distributions/Packaging Recommendations](https://community.kde.org/Distributions/Packaging_Recommendations)

```sh
pacman -Syu plasma-desktop

# Login manager
pacman -Syu sddm
```

#### Worth trying

- [COSMIC](https://wiki.archlinux.org/title/COSMIC)
- [Pantheon](https://wiki.archlinux.org/title/Pantheon)

## [List of applications](https://wiki.archlinux.org/index.php/List_of_applications)

### [pacman](https://wiki.archlinux.org/index.php/pacman)

Uncomment in `/etc/pacman.conf`:

```txt
# Misc options
Color
ParallelDownloads
```

```sh
systemctl enable paccache.timer
```

### [Pipewire](https://wiki.archlinux.org/title/PipeWire)

```sh
pacman -Syu pipewire wireplumber \
	pipewire-alsa pipewire-pulse
```

See
[Advanced Linux Sound Architecture/Firmware](https://wiki.archlinux.org/title/Advanced_Linux_Sound_Architecture#Firmware)

### [Flatpak](https://wiki.archlinux.org/title/Flatpak)

```sh
pacman -Syu flatpak
```

## [Improving performance](https://wiki.archlinux.org/index.php/improving_performance)

- https://wiki.archlinux.org/index.php/swap#Swap_file
- https://wiki.archlinux.org/index.php/swap#Swappiness
- https://wiki.archlinux.org/index.php/Systemd/Journal#Journal_size_limit
- https://wiki.archlinux.org/index.php/Core_dump#Disabling_automatic_core_dumps
- https://wiki.archlinux.org/title/Ext4#Enabling_fast_commit_in_existing_filesystems
- https://wiki.archlinux.org/index.php/Solid_state_drive#Periodic_TRIM
- https://wiki.archlinux.org/index.php/Silent_boot
- https://wiki.archlinux.org/title/Improving_performance#Watchdogs
- https://wiki.archlinux.org/title/Sysctl#Enable_TCP_Fast_Open
- [Fast commits for ext4](https://lwn.net/Articles/842385/)
- [TCP Fast Open: expediting web services](https://lwn.net/Articles/508865/)
- [The search for the correct amount of split-lock misery](https://lwn.net/Articles/911219/)

Edit `/etc/systemd/journald.conf.d/00-journal-size.conf` then restart:

```txt
[Journal]
SystemMaxUse=50M
```

Edit `/etc/systemd/coredump.conf.d/custom.conf` then restart:

```txt
[Coredump]
Storage=none
ProcessSizeMax=0
```

Enable ext4 fast commit:

```sh
tune2fs -O fast_commit /dev/partition
```

Periodic TRIM:

```sh
systemctl enable fstrim.timer
```

Edit `/etc/sysctl.d/99-sysctl.conf`:

```txt
# Enable TCP Fast Open
net.ipv4.tcp_fastopen = 3

kernel.split_lock_mitigate = 0
```

## [Security](https://wiki.archlinux.org/title/Security)

- https://wiki.archlinux.org/title/IPv6#Disable_IPv6
- [add init_on_alloc/init_on_free boot options](https://lwn.net/Articles/791380/)
- [mm: Randomize free memory](https://lwn.net/Articles/776228/)
- [mm: introduce Designated Movable Blocks](https://lwn.net/Articles/925941/)

## Hardware dependent

- https://wiki.archlinux.org/title/Laptop
- https://wiki.archlinux.org/title/ASUS_Linux
- https://wiki.archlinux.org/title/PRIME

## Experiment

Do it at your own risk!!!

- https://wiki.archlinux.org/title/Pacman/Pacnew_and_Pacsave
- https://wiki.archlinux.org/title/Unified_kernel_image
- https://wiki.archlinux.org/title/Unified_Extensible_Firmware_Interface
- https://github.com/Foxboron/sbctl
- https://github.com/AdnanHodzic/auto-cpufreq
- https://github.com/nbfc-linux/nbfc-linux
- https://github.com/erpalma/throttled

## Maintenance

See [pacman/Tips and tricks](https://wiki.archlinux.org/title/Pacman/Tips_and_tricks)

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
ls /sys/firmware/efi/efivars
```

#### Connect to the internet

For wifi, use [iwd](https://wiki.archlinux.org/index.php/Iwd).

#### Partition the disks

[GPT fdisk](https://wiki.archlinux.org/index.php/GPT_fdisk):

```sh
cgdisk /dev/sdx
```

[Partition scheme](https://wiki.archlinux.org/index.php/Partitioning#Partition_scheme)

UEFI/GPT layout:

| Mount point | Partition                             | Partition type                 | Suggested size |
| ----------- | ------------------------------------- | ------------------------------ | -------------- |
| `/mnt/efi`  | `/dev/efi_system_partition`           | EFI System Partition           | 512 MiB        |
| `/mnt/boot` | `/dev/extended_boot_loader_partition` | Extended Boot Loader Partition | 1 GiB          |
| `/mnt`      | `/dev/root_partition`                 | Root Partition                 |                |

Why not `/boot/efi`?
See [Lennart Poettering comment](https://github.com/systemd/systemd/pull/3757#issuecomment-234290236).

BIOS/GPT layout:

| Mount point | Partition             | Partition type      | Suggested size |
| ----------- | --------------------- | ------------------- | -------------- |
|             |                       | BIOS boot partition | 1 MiB          |
| `/mnt`      | `/dev/root_partition` | Root Partition      |                |

LVM:

```sh
# Create physical volumes
pvcreate /dev/sdaX

# Create volume groups
vgcreate RootGroup /dev/sdaX /dev/sdaY

# Create logical volumes
lvcreate -l +100%FREE RootGroup -n rootvol
```

Format:

```sh
# efi
mkfs.fat -F32 /dev/efi_system_partition

# boot
mkfs.fat -F32 /dev/extended_boot_loader_partition

# root
mkfs.ext4 -L ROOT /dev/root_partition

# root with btrfs
mkfs.btrfs -L ROOT /dev/root_partition

# root on lvm
mkfs.ext4 /dev/RootGroup/rootvol
```

Mount:

```sh
# root
mount /dev/root_partition /mnt

# root with btrfs
mount -o compress=zstd /dev/root_partition /mnt

# root on lvm
mount /dev/RootGroup/rootvol /mnt

# efi
mount --mkdir /dev/efi_system_partition /mnt/efi

# boot
mount --mkdir /dev/extended_boot_loader_partition /mnt/boot
```

### Installation

```sh
pacstrap -K /mnt base linux linux-firmware

# AMD
pacstrap -K /mnt amd-ucode

# Intel
pacstrap -K /mnt intel-ucode

# Btrfs
pacstrap -K /mnt btrfs-progs

# LVM
pacstrap -K /mnt lvm2

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

#### Time zone

```sh
ln -sf /usr/share/zoneinfo/Region/City /etc/localtime

hwclock --systohc
```

#### Localization:

Edit `/etc/locale.gen`:

```txt
# Uncomment en_US.UTF-8 UTF-8
```

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
# LVM
# https://wiki.archlinux.org/title/Install_Arch_Linux_on_LVM#Adding_mkinitcpio_hooks
HOOKS=(base udev ... block lvm2 filesystems)

# https://wiki.archlinux.org/title/mkinitcpio#Common_hooks
# Replace udev with systemd
```

```sh
mkinitcpio -P
```

#### Root password

```sh
passwd
```

#### Addition

##### [NetworkManager (WIP)](https://wiki.archlinux.org/title/NetworkManager)

```sh
pacman -Syu networkmanager dhcpcd iwd
systemctl enable NetworkManager.service
systemctl enable systemd-resolved.service
```

Edit `/etc/NetworkManager/conf.d/dns.conf`:

```txt
[main]
dns=systemd-resolved
```

Edit `/etc/NetworkManager/conf.d/dhcp-client.conf`:

```txt
[main]
dhcp=dhcpcd
```

Edit `/etc/NetworkManager/conf.d/wifi_backend.conf`:

```txt
[device]
wifi.backend=iwd
```

See [dhcpcd](https://wiki.archlinux.org/title/Dhcpcd)

Append `/etc/dhcpcd.conf`

```txt
noarp
nohook resolv.conf
```

##### Bluetooth

```sh
pacman -Syu bluez
systemctl enable bluetooth.service
```

##### Clock

```sh
timedatectl set-ntp true
```

#### Boot loader

##### [systemd-boot](https://wiki.archlinux.org/index.php/Systemd-boot)

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

# Intel
initrd /intel-ucode.img

# AMD
initrd /amd-ucode.img

initrd /initramfs-linux.img

# Kernel parameters
#
# Acer Nitro AN515-45
# https://wiki.archlinux.org/title/backlight#Kernel_command-line_options
# acpi_backlight=vendor
#
# NVIDIA
# https://wiki.archlinux.org/title/NVIDIA#DRM_kernel_mode_setting
# nvidia-drm.modeset=1
options root="LABEL=ROOT" rw
```

## [General recommendations](https://wiki.archlinux.org/index.php/General_recommendations)

Always remember to check **dependencies** when install packages.

### System administration

[Sudo](https://wiki.archlinux.org/index.php/sudo):

```sh
pacman -Syu sudo

EDITOR=nvim visudo
# Uncomment group wheel

# Add user if don't want to use systemd-homed
useradd -m -G wheel -c "The Joker" joker

# Or using zsh
useradd -m -G wheel -s /usr/bin/zsh -c "The Joker" joker

# Set password
passwd joker
```

[doas (WIP)](https://wiki.archlinux.org/title/Doas)

```sh
pacman -Syu opendoas

chown -c root:root /etc/doas.conf
chmod -c 0400 /etc/doas.conf
```

Edit `/etc/doas.conf` (must end with a newline):

```txt
permit persist :wheel
```

[systemd-homed (WIP)](https://wiki.archlinux.org/index.php/Systemd-homed):

```sh
systemctl enable systemd-homed.service

homectl create joker --real-name="The Joker" --member-of=wheel

# Using zsh
homectl update joker --shell=/usr/bin/zsh
```

**Note**:
Can not run `homectl` when install Arch Linux.
Should run on the first boot.

### Desktop Environment

Install [Xorg](https://wiki.archlinux.org/index.php/Xorg):

```sh
pacman -Syu xorg-server
```

#### [GNOME](https://wiki.archlinux.org/index.php/GNOME)

```sh
pacman -Syu gnome-shell \
	gnome-control-center gnome-system-monitor power-profiles-daemon \
	gnome-tweaks gnome-backgrounds gnome-screenshot gnome-keyring gnome-logs gnome-firmware \
	gnome-console gnome-text-editor \
	nautilus xdg-user-dirs-gtk file-roller evince eog

# Login manager
pacman -Syu gdm
systemctl enable gdm.service
```

## [List of applications](https://wiki.archlinux.org/index.php/List_of_applications)

### [pacman](https://wiki.archlinux.org/index.php/pacman)

Uncomment in `/etc/pacman.conf`:

```txt
# Misc options
Color
ParallelDownloads
```

### [Pipewire (WIP)](https://wiki.archlinux.org/title/PipeWire)

```sh
pacman -Syu pipewire wireplumber \
	pipewire-alsa pipewire-pulse \
	gst-plugin-pipewire pipewire-v4l2
```

See [Advanced Linux Sound Architecture](https://wiki.archlinux.org/title/Advanced_Linux_Sound_Architecture)

```sh
pacman -Syu sof-firmware
```

### [Flatpak (WIP)](https://wiki.archlinux.org/title/Flatpak)

```sh
pacman -Syu flatpak
```

## [Improving performance](https://wiki.archlinux.org/index.php/improving_performance)

https://wiki.archlinux.org/index.php/swap#Swap_file

https://wiki.archlinux.org/index.php/swap#Swappiness

https://wiki.archlinux.org/index.php/Systemd/Journal#Journal_size_limit

https://wiki.archlinux.org/index.php/Core_dump#Disabling_automatic_core_dumps

https://wiki.archlinux.org/index.php/Solid_state_drive#Periodic_TRIM

https://wiki.archlinux.org/index.php/Silent_boot

https://wiki.archlinux.org/title/Improving_performance#Watchdogs

## In the end

This guide is updated regularly I promise.

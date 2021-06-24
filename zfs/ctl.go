// Copyright 2021 lack
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zfs

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type zfsctl struct {
	cmd string
}

func ZFSCtl(cmd string) *zfsctl {
	zfsctl := &zfsctl{cmd: cmd}
	return zfsctl
}

// CreateFileSystem creates zfs filesystem
// 	zfs create [-p] [-o property=value] ... <filesystem>
func (z *zfsctl) CreateFileSystem(ctx context.Context, name string, properties map[string]string) *execute {
	args := []string{"create", "-p"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// CreateVolume creates zfs block volume
// 	zfs create [-ps] [-b blocksize] [-o property=value] ... -V <size> <volume>
func (z *zfsctl) CreateVolume(ctx context.Context, name string, block int64, properties map[string]string, size string) *execute {
	args := []string{"create", "-ps"}
	if block > 0 {
		args = append(args, fmt.Sprintf("-b %d", block))
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, "-V", size, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// DestroyFileSystemOrVolume destroy zfs filesystem or volume
// 	zfs destroy [-fnpRrv] <filesystem|volume>
func (z *zfsctl) DestroyFileSystemOrVolume(ctx context.Context, name, options string) *execute {
	args := []string{"destroy"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// DestroySnapshot destroy zfs snapshot
// 	zfs destroy [-dnpRrv] <filesystem|volume>@<snap>[%<snap>][,...]
func (z *zfsctl) DestroySnapshot(ctx context.Context, name, options string) *execute {
	args := []string{"destroy"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// DestroyBookmark destroy zfs bookmark that belongs to filesystem or volume
// 	zfs destroy <filesystem|volume>#<bookmark>
func (z *zfsctl) DestroyBookmark(ctx context.Context, name string) *execute {
	args := []string{"destroy", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Snapshot create snapshot that belongs to volume or filesystem
// 	zfs snapshot|snap [-r] [-o property=value] ... <filesystem|volume>@<snap> ...
func (z *zfsctl) Snapshot(ctx context.Context, name string, properties map[string]string) *execute {
	args := []string{"snapshot", "-r"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Rollback uses old data from snapshot
// 	zfs rollback [-rRf] <snapshot>
func (z *zfsctl) Rollback(ctx context.Context, options string, snap string) *execute {
	args := []string{"rollback"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, snap)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Clone creates a volume or filesystem from snapshot
// 	zfs clone [-p] [-o property=value] ... <snapshot> <filesystem|volume>
func (z *zfsctl) Clone(ctx context.Context, name string, properties map[string]string, source string) *execute {
	args := []string{"clone", "-p"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, source, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Promote
// 	zfs promote <clone-filesystem>
func (z *zfsctl) Promote(ctx context.Context, name string) *execute {
	args := []string{"promote", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Rename renames volume, filesystem or snapshot
// 	zfs rename [-f] <filesystem|volume|snapshot> <filesystem|volume|snapshot>
func (z *zfsctl) Rename(ctx context.Context, name, newName string, force bool) *execute {
	args := []string{"rename"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name, newName)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// RenameFileSystemOrVolume renames volume or filesystem
// 	zfs rename [-f] -p <filesystem|volume> <filesystem|volume>
func (z *zfsctl) RenameFileSystemOrVolume(ctx context.Context, name, newName string, force bool) *execute {
	args := []string{"rename"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, "-p")
	args = append(args, name, newName)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// RenameSnapshot renames snapshot:
// 	zfs rename -r <snapshot> <snapshot>
func (z *zfsctl) RenameSnapshot(ctx context.Context, name, newName string) *execute {
	args := []string{"rename", "-r"}
	args = append(args, name, newName)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Bookmark setup bookmark for snapshot:
// 	zfs bookmark <snapshot> <bookmark>
func (z *zfsctl) Bookmark(ctx context.Context, snapshot, bookmark string) *execute {
	args := []string{"bookmark", snapshot, bookmark}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// List get all datasets:
// 	zfs list [-Hp] [-r|-d max] [-o property[,...]] [-s property]...
//            [-S property]... [-t type[,...]] [filesystem|volume|snapshot] ...
func (z *zfsctl) List(ctx context.Context, name, options, max string, oProperties []string, sProperty, SProperty, t string) *execute {
	args := []string{"list"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(max) > 0 {
		args = append(args, max)
	}
	if oProperties != nil {
		o := "-o "
		for _, p := range oProperties {
			o += p + ","
		}
		args = append(args, strings.TrimSuffix(o, ","))
	}
	if len(sProperty) > 0 {
		args = append(args, sProperty)
	}
	if len(SProperty) > 0 {
		args = append(args, SProperty)
	}
	if len(t) > 0 {
		args = append(args, "-t "+t)
	}
	if len(name) > 0 {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Set setup dataset property:
// 	zfs set <property=value> ... <filesystem|volume|snapshot> ...
func (z *zfsctl) Set(ctx context.Context, name string, properties map[string]string) *execute {
	args := []string{"set"}
	if properties != nil {
		kv := ""
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Get gets dataset's property:
// 	zfs get [-rHp] [-d max] [-o "all" | field[,...]]
//            [-t type[,...]] [-s source[,...]]
//            <"all" | property[,...]> [filesystem|volume|snapshot|bookmark] ...
func (z *zfsctl) Get(ctx context.Context, name, options, max string, out []string, t, s string, properties ...string) *execute {
	args := []string{"get"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(max) > 0 {
		args = append(args, "-d "+max)
	}
	if out != nil {
		o := "-o "
		for _, p := range properties {
			o += p + ","
		}
		args = append(args, strings.TrimSuffix(o, ","))
	}
	if len(t) > 0 {
		args = append(args, "-t "+t)
	}
	if len(s) > 0 {
		args = append(args, "-s ", s)
	}
	if properties != nil {
		o := ""
		for _, p := range properties {
			o += p + ","
		}
		o = strings.TrimSuffix(o, ",")
		args = append(args, o)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Inherit :
// 	zfs inherit [-rS] <property> <filesystem|volume|snapshot> ...
func (z *zfsctl) Inherit(ctx context.Context, name string, options string, property string) *execute {
	args := []string{"inherit"}
	if len(options) != 0 {
		args = append(args, options)
	}
	args = append(args, property, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Upgrade zfs upgrade
// 	zfs upgrade [-v]
func (z *zfsctl) Upgrade(ctx context.Context, v bool) *execute {
	args := []string{"upgrade"}
	if v {
		args = append(args, "-v")
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// UpgradeFileSystem
// 	zfs upgrade [-r] [-V version] <-a | filesystem ...>
func (z *zfsctl) UpgradeFileSystem(ctx context.Context, name, version string, r, all bool) *execute {
	args := []string{"upgrade"}
	if r {
		args = append(args, "-r")
	}
	if len(version) > 0 {
		args = append(args, "-V "+version)
	}
	if all {
		args = append(args, "-a")
	} else {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Userspace
// 	zfs userspace [-Hinp] [-o field[,...]] [-s field] ...
//            [-S field] ... [-t type[,...]] <filesystem|snapshot>
func (z *zfsctl) Userspace(ctx context.Context, name, options string, fields []string, sField, SField, t string) *execute {
	args := []string{"userspace"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if fields != nil {
		o := "-o "
		for _, field := range fields {
			o += field + ","
		}
		args = append(args, strings.TrimSuffix(o, ","))
	}
	if len(sField) > 0 {
		args = append(args, "-s "+sField)
	}
	if len(SField) > 0 {
		args = append(args, "-S "+SField)
	}
	if len(t) > 0 {
		args = append(args, "-t "+t)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// GroupSpace
// 	zfs groupspace [-Hinp] [-o field[,...]] [-s field] ...
//            [-S field] ... [-t type[,...]] <filesystem|snapshot>
func (z *zfsctl) GroupSpace(ctx context.Context, name, options string, fields []string, sField, SField, t string) *execute {
	args := []string{"groupspace"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if fields != nil {
		o := "-o "
		for _, field := range fields {
			o += field + ","
		}
		o = strings.TrimSuffix(o, ",")
		args = append(args, o)
	}
	if len(sField) > 0 {
		args = append(args, "-s "+sField)
	}
	if len(SField) > 0 {
		args = append(args, "-S "+SField)
	}
	if len(t) > 0 {
		args = append(args, "-t "+t)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Mount mounts all endpoint
// 	zfs mount
func (z *zfsctl) Mount(ctx context.Context) *execute {
	args := []string{"mount"}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// MountFileSystem
// 	zfs mount [-vO] [-o opts] <-a | filesystem>
func (z *zfsctl) MountFileSystem(ctx context.Context, name, options, opts string, all bool) *execute {
	args := []string{"mount"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(opts) > 0 {
		args = append(args, opts)
	}
	if all {
		args = append(args, "-a")
	} else {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Umount
// 	zfs unmount [-f] <-a | filesystem|mountpoint>
func (z *zfsctl) Umount(ctx context.Context, name string, force, all bool) *execute {
	args := []string{"umount"}
	if force {
		args = append(args, "-f")
	}
	if all {
		args = append(args, "-a")
	} else {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Share
// 	zfs share <-a [nfs|smb] | filesystem>
func (z *zfsctl) Share(ctx context.Context, name string, all bool, kind string) *execute {
	args := []string{"share"}
	if all {
		args = append(args, "-a")
		if len(kind) > 0 {
			args = append(args, kind)
		}
	} else {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Unshare
// 	zfs unshare <-a [nfs|smb] | filesystem|mountpoint>
func (z *zfsctl) Unshare(ctx context.Context, name string, all bool, kind string) *execute {
	args := []string{"unshare"}
	if all {
		args = append(args, "-a")
		if len(kind) > 0 {
			args = append(args, kind)
		}
	} else {
		args = append(args, name)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// SendSnapshot Examples:
// 	zfs send [-DnPpRvLec] [-[i|I] snapshot] <snapshot>
func (z *zfsctl) SendSnapshot(ctx context.Context, name, options string, i string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(i) > 0 {
		args = append(args, "-i "+i)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// SendAndRecv Examples:
// 	zfs send tank/sla@snap | ssh storage@192.168.2.120 -i test "sudo /usr/sbin/zfs receive -Fu dup/sla@snap"
func (z *zfsctl) SendAndRecv(ctx context.Context, source, target, user, host string) *execute {
	args := []string{"send", source}
	if len(user) != 0 && len(host) != 0 {
		args = append(args, fmt.Sprintf(`| ssh %s@%s -i %s "sudo %s receive -Fu %s"`, user, host, user, z.cmd, target))
	} else {
		args = append(args, fmt.Sprintf(`| %s receive -Fu %s`, z.cmd, target))
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// IncrementSendAndRecv Examples:
// 	zfs send -i tank/sla@snap1 tank/sla@snap2 | ssh storage@192.168.2.120 -i test "sudo /usr/sbin/zfs receive -Fu dup/sla@snap2"
func (z *zfsctl) IncrementSendAndRecv(ctx context.Context, source, lastSnapshot, target, user, host string) *execute {
	args := []string{"send", "-i", lastSnapshot, source}
	if len(user) != 0 && len(host) != 0 {
		args = append(args, fmt.Sprintf(`| ssh %s@%s -i %s "sudo %s receive -Fu %s"`, user, host, user, z.cmd, target))
	} else {
		args = append(args, fmt.Sprintf(`| %s receive -Fu %s`, z.cmd, target))
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Send Examples:
// 	zfs send [-Lec] [-i snapshot|bookmark] <filesystem|volume|snapshot>
func (z *zfsctl) Send(ctx context.Context, name, options string, i string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(i) > 0 {
		args = append(args, "-i "+i)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// SendToken Examples:
// 	zfs send [-nvPe] -t <receive_resume_token>
func (z *zfsctl) SendToken(ctx context.Context, name, options string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, "-t", name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Receive Examples:
// 	zfs receive [-vnsFu] [-o <property>=<value>] ... [-x <property>] ...
//            <filesystem|volume|snapshot>
func (z *zfsctl) Receive(ctx context.Context, name, options string, properties map[string]string, xProperty string) *execute {
	args := []string{"receive"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, strings.TrimSuffix(kv, " "))
	}
	if len(xProperty) > 0 {
		args = append(args, "-x "+xProperty)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// ReceiveFileSystem Examples:
// 	zfs receive [-vnsFu] [-o <property>=<value>] ... [-x <property>] ...
//            <filesystem|volume|snapshot>
func (z *zfsctl) ReceiveFileSystem(ctx context.Context, name, options string, properties map[string]string, xProperty string, de string) *execute {
	args := []string{"receive"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, strings.TrimSuffix(kv, " "))
	}
	if len(xProperty) > 0 {
		args = append(args, "-x "+xProperty)
	}
	switch de {
	case "-d":
		args = append(args, "-d")
	case "-e":
		args = append(args, "-e")
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// ReceiveAll Examples:
// 	zfs receive -A <filesystem|volume>
func (z *zfsctl) ReceiveAll(ctx context.Context, name string) *execute {
	args := []string{"receive", "-A", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Allow1 Examples:
// 	zfs allow <filesystem|volume>
func (z *zfsctl) Allow1(ctx context.Context, name string) *execute {
	args := []string{"allow", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Allow2 Examples:
// 	zfs allow [-ldug] <"everyone"|user|group>[,...] <perm|@setname>[,...]
//            <filesystem|volume>
func (z *zfsctl) Allow2(ctx context.Context, name, options, authority, perm string) *execute {
	args := []string{"allow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, authority, perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Allow3 Examples:
// 	zfs allow [-ld] -e <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow3(ctx context.Context, name, options, perm string) *execute {
	args := []string{"allow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Allow4 Examples:
// 	zfs allow -c <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow4(ctx context.Context, name, perm string) *execute {
	args := []string{"allow", "-c", perm, name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Allow5 Examples:
// 	zfs allow -s @setname <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow5(ctx context.Context, name, setname, perm string) *execute {
	args := []string{"allow", "-s", setname, perm, name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Unallow1 Examples:
// 	zfs unallow [-rldug] <"everyone"|user|group>[,...]
//            [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow1(ctx context.Context, name, options, authority, perm string) *execute {
	args := []string{"unallow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, authority, perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Unallow2 Examples:
// 	zfs unallow [-rld] -e [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow2(ctx context.Context, name, options, perm string) *execute {
	args := []string{"unallow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Unallow3 Examples:
// 	zfs unallow [-r] -c [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow3(ctx context.Context, name string, r bool, perm string) *execute {
	args := []string{"unallow"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, "-c", perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Unallow4 Examples:
// 	zfs unallow [-r] -s @setname [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow4(ctx context.Context, name string, r bool, setmame, perm string) *execute {
	args := []string{"unallow"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, "-s", setmame, perm, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Hold Examples:
// 	zfs hold [-r] <tag> <snapshot> ...
func (z *zfsctl) Hold(ctx context.Context, name string, r bool, tag string) *execute {
	args := []string{"hold"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, tag, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Holds Examples:
// 	zfs holds [-r] <snapshot> ...
func (z *zfsctl) Holds(ctx context.Context, name string, r bool) *execute {
	args := []string{"holds"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Release Examples:
// 	zfs  release [-r] <tag> <snapshot> ...
func (z *zfsctl) Release(ctx context.Context, name string, r bool, tag string) *execute {
	args := []string{"release"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, tag, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Diff Examples:
// 	zfs diff [-FHt] <snapshot> [snapshot|filesystem]
func (z *zfsctl) Diff(ctx context.Context, name, options string, target string) *execute {
	args := []string{"diff"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	if len(target) > 0 {
		args = append(args, target)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

type zpoolctl struct {
	cmd string
}

func ZPoolCtl(cmd string) *zpoolctl {
	z := &zpoolctl{cmd: cmd}
	return z
}

// Create creates new pool
// Examples:
// 	zpool create [-fnd] [-o property=value] ...
//            [-O file-system-property=value] ...
//            [-m mountpoint] [-R root] <pool> <vdev>
func (z *zpoolctl) Create(ctx context.Context, name, options, point, root string, properties map[string]string, raid string, devs ...string) *execute {
	args := []string{"create"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	if len(point) > 0 {
		args = append(args, "-m "+point)
	}
	if len(root) > 0 {
		args = append(args, "-R "+root)
	}
	args = append(args, name, raid)
	for _, dev := range devs {
		args = append(args, dev)
	}

	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Destroy delete a pool
// Examples:
// 	zpool destroy [-f] <pool>
func (z *zpoolctl) Destroy(ctx context.Context, name string, force bool) *execute {
	args := []string{"destroy"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Add Examples:
//	zpool add [-fgLnP] [-o property=value] <pool> <vdev>
func (z *zpoolctl) Add(ctx context.Context, name string, options string, devs ...string) *execute {
	args := []string{"add"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Remove Examples:
//	zpool remove <pool> <device>
func (z *zpoolctl) Remove(ctx context.Context, name string, devs ...string) *execute {
	args := []string{"remove", name}
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// LabelClear Examples:
//	zpool labelclear [-f] <vdev>
func (z *zpoolctl) LabelClear(ctx context.Context, device string, force bool) *execute {
	args := []string{"labelclear"}
	if force {
		args = append(args, "-f", device)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// List Examples:
//	zpool list [-gHLpPv] [-o property[,...]] [-T d|u] [pool]
func (z *zpoolctl) List(ctx context.Context, name, options string, properties []string, t string) *execute {
	args := []string{"list"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if properties != nil {
		kv := "-o "
		for _, v := range properties {
			kv += v + ","
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	if len(t) > 0 {
		args = append(args, "-T "+t)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Iostat Examples:
// 	zpool iostat [[[-c [script1,script2,...][-lq]]|[-rw]] [-T d | u] [-ghHLpPvy]
//            [[pool ...]|[pool vdev ...]|[vdev ...]] [interval [count]]
func (z *zpoolctl) Iostat(ctx context.Context, name, scripts, t, options string, devs ...string) *execute {
	args := []string{"iostat"}
	if len(scripts) > 0 {
		args = append(args, scripts)
	}
	if len(t) > 0 {
		args = append(args, t)
	}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Status Examples:
// 	zpool status [-c [script1,script2,...]] [-gLPvxD][-T d|u] [pool] ... [interval [count]]
func (z *zpoolctl) Status(ctx context.Context, scripts, options, t string, name string) *execute {
	args := []string{"status"}
	if len(scripts) > 0 {
		args = append(args, scripts)
	}
	if len(t) > 0 {
		args = append(args, t)
	}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Online Examples:
//	zpool online <pool> <device> ...
func (z *zpoolctl) Online(ctx context.Context, name string, devs ...string) *execute {
	args := []string{"online"}
	args = append(args, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Offline Examples:
//	zpool offline [-f] [-t] <pool> <device> ...
func (z *zpoolctl) Offline(ctx context.Context, name string, force, t bool, devs ...string) *execute {
	args := []string{"offline"}
	args = append(args, name)
	if force {
		args = append(args, "-f")
	}
	if t {
		args = append(args, "-t")
	}
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Clear Examples:
//	zpool clear [-nF] <pool> [device]
func (z *zpoolctl) Clear(ctx context.Context, name, options string, devs ...string) *execute {
	args := []string{"clear"}
	args = append(args, options, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Reopen Examples:
//	zpool reopen <pool>
func (z *zpoolctl) Reopen(ctx context.Context, name string) *execute {
	args := []string{"reopen", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Attach Examples:
//	zpool attach [-f] [-o property=value] <pool> <device> <new-device>
func (z *zpoolctl) Attach(ctx context.Context, name string, force bool, properties map[string]string, dev, newDev string) *execute {
	args := []string{"attach"}
	if force {
		args = append(args, "-f")
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s,", k, v)
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	args = append(args, name, dev, newDev)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Detach Examples:
//	zpool detach <pool> <device>
func (z *zpoolctl) Detach(ctx context.Context, name, dev string) *execute {
	args := []string{"attach", name, dev}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Replace Examples:
//	zpool replace [-f] [-o property=value] <pool> <device> [new-device]
func (z *zpoolctl) Replace(ctx context.Context, name string, force bool, properties map[string]string, dev string, newDev string) *execute {
	args := []string{"replace"}
	if force {
		args = append(args, "-f")
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s,", k, v)
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	args = append(args, name, dev)
	if len(newDev) > 0 {
		args = append(args, newDev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Split Examples:
//	zpool split [-gLnP] [-R altroot] [-o mntopts]
//            [-o property=value] <pool> <newpool> [<device> ...]
func (z *zpoolctl) Split(ctx context.Context, name, options, altRoot, mntopts string, properties map[string]string, newName string, devs ...string) *execute {
	args := []string{"split"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(altRoot) > 0 {
		args = append(args, altRoot)
	}
	if len(mntopts) > 0 {
		args = append(args, mntopts)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s,", k, v)
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	args = append(args, name, newName)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Examples:
//	zpool scrub [-s | -p] <pool> ...
func (z *zpoolctl) scrub(ctx context.Context, name, option string) *execute {
	args := []string{"scrub"}
	if len(option) > 0 {
		args = append(args, option)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Import1 Examples:
//	zpool import [-d dir] [-D]
func (z *zpoolctl) Import1(ctx context.Context, dir string, d bool) *execute {
	args := []string{"import"}
	if len(dir) > 0 {
		args = append(args, dir)
	}
	if d {
		args = append(args, "-D")
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Import2 Examples:
//	zpool import [-d dir] [-D]
func (z *zpoolctl) Import2(ctx context.Context, dir string, d bool) *execute {
	args := []string{"import"}
	if len(dir) > 0 {
		args = append(args, dir)
	}
	if d {
		args = append(args, "-D")
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Import3 Examples:
//	zpool import [-d dir | -c cachefile] [-F [-n]] <pool | id>
func (z *zpoolctl) Import3(ctx context.Context, name, dir, file string, force, n bool) *execute {
	args := []string{"import"}
	switch {
	case len(dir) > 0:
		args = append(args, "-d "+dir)
	case len(file) > 0:
		args = append(args, "-c "+file)
	}
	if force {
		args = append(args, "-f")
		if n {
			args = append(args, "-n")
		}
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Import4 Examples:
//	zpool import [-o mntopts] [-o property=value] ...
//            [-d dir | -c cachefile] [-D] [-f] [-m] [-N] [-R root] [-F [-n]] -a
func (z *zpoolctl) Import4(
	ctx context.Context,
	mntopts string, properties map[string]string,
	dir, file string, d, m, N bool, root string, force, n bool) *execute {
	args := []string{"import"}
	if len(mntopts) > 0 {
		args = append(args, "-o "+mntopts)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s,", k, v)
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	switch {
	case len(dir) > 0:
		args = append(args, "-d "+dir)
	case len(file) > 0:
		args = append(args, "-c "+file)
	}
	if d {
		args = append(args, "-d")
	}
	if m {
		args = append(args, "-m")
	}
	if N {
		args = append(args, "-N")
	}
	if len(root) > 0 {
		args = append(args, "-R "+root)
	}
	if force {
		args = append(args, "-f")
		if n {
			args = append(args, "-n")
		}
	}
	args = append(args, "-a")
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Import5 Examples:
//	zpool import [-o mntopts] [-o property=value] ...
//            [-d dir | -c cachefile] [-D] [-f] [-m] [-N] [-R root] [-F [-n]]
//            <pool | id> [newpool]
func (z *zpoolctl) Import5(
	ctx context.Context,
	name, mntopts string, properties map[string]string,
	dir, file string, d, m, N bool, root string, force, n bool, newPool string) *execute {
	args := []string{"import"}
	if len(mntopts) > 0 {
		args = append(args, "-o "+mntopts)
	}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s,", k, v)
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	switch {
	case len(dir) > 0:
		args = append(args, "-d "+dir)
	case len(file) > 0:
		args = append(args, "-c "+file)
	}
	if d {
		args = append(args, "-d")
	}
	if m {
		args = append(args, "-m")
	}
	if N {
		args = append(args, "-N")
	}
	if len(root) > 0 {
		args = append(args, "-R "+root)
	}
	if force {
		args = append(args, "-f")
		if n {
			args = append(args, "-n")
		}
	}
	args = append(args, name)
	if len(newPool) > 0 {
		args = append(args, newPool)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Export Examples:
//	zpool zpoexport [-af] <pool> ...
func (z *zpoolctl) Export(ctx context.Context, name string, options string) *execute {
	args := []string{"export"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Upgrade1 Examples:
//	zpool upgrade [-v]
func (z *zpoolctl) Upgrade1(ctx context.Context, v bool) *execute {
	args := []string{"upgrade"}
	if v {
		args = append(args, "-v")
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Upgrade2 Examples:
//	zpool upgrade [-V version] <-a | pool ...>
func (z *zpoolctl) Upgrade2(ctx context.Context, name, version string) *execute {
	args := []string{"upgrade"}
	if len(version) > 0 {
		args = append(args, "-V "+version)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Reguid Examples:
//	zpool reguid <pool>
func (z *zpoolctl) Reguid(ctx context.Context, name string) *execute {
	args := []string{"reguid", name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// History Examples:
//	zpool history [-il] [<pool>]
func (z *zpoolctl) History(ctx context.Context, options string, pool string) *execute {
	args := []string{"history"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(pool) > 0 {
		args = append(args, pool)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Event Examples:
//	zpool events [-vHfc]
func (z *zpoolctl) Event(ctx context.Context, options string) *execute {
	args := []string{"event"}
	if len(options) > 0 {
		args = append(args, options)
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Get Examples:
//	zpool get [-Hp] [-o "all" | field[,...]] <"all" | property[,...]> <pool> ...
func (z *zpoolctl) Get(ctx context.Context, name, options string, out []string, properties ...string) *execute {
	args := []string{"get"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(options) > 0 {
		args = append(args, options)
	}
	if out != nil && len(out) > 0 {
		options := "-o "
		for _, o := range out {
			options += o + ","
		}
		options = strings.TrimSuffix(options, ",")
		args = append(args, options)
	}
	if properties != nil {
		kv := ""
		for _, p := range properties {
			kv += p + ","
		}
		kv = strings.TrimSuffix(kv, ",")
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Set Examples:
//	zpool set <property=value> <pool>
func (z *zpoolctl) Set(ctx context.Context, name, k, v string) *execute {
	args := []string{"set", k + "=" + v, name}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

// Sync Examples:
//	zpool sync [pool] ...
func (z *zpoolctl) Sync(ctx context.Context, names ...string) *execute {
	args := []string{"sync"}
	if names != nil {
		for _, name := range names {
			args = append(args, name)
		}
	}
	return &execute{ctx: ctx, name: z.cmd, args: args}
}

type execute struct {
	ctx  context.Context
	name string
	args []string
}

func (e *execute) Commit() string {
	return fmt.Sprintf(`%s %s`, e.name, strings.Join(e.args, " "))
}

func (e *execute) Exec() ([]byte, error) {
	if e.ctx == nil {
		e.ctx = context.Background()
	}
	return execution(e.ctx, fmt.Sprintf(`%s %s`, e.name, strings.Join(e.args, " ")))
}

func (e *execute) Bash() *exec.Cmd {
	if e.ctx == nil {
		e.ctx = context.Background()
	}
	shell := fmt.Sprintf(`%s %s`, e.name, strings.Join(e.args, " "))
	return exec.CommandContext(e.ctx, "/bin/bash", "-c", shell)
}

func execution(ctx context.Context, shell string) ([]byte, error) {
	data, err := exec.CommandContext(ctx, "/bin/sh", "-c", shell).CombinedOutput()
	if bytes.HasSuffix(data, []byte("\n")) {
		data = bytes.TrimSuffix(data, []byte("\n"))
	}
	if err != nil {
		return nil, fmt.Errorf("%v:%v", err, string(data))
	}
	return data, nil
}

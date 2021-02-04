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
	"fmt"
	"os/exec"
	"strings"
)

type zfsctl struct {
	cmd string
}

func ZfsCtl(cmd string) *zfsctl {
	zfsctl := &zfsctl{cmd: cmd}
	return zfsctl
}

// Examples:
// 	zfs create [-p] [-o property=value] ... <filesystem>
func (z *zfsctl) CreateFileSystem(name string, properties map[string]string) *execute {
	args := []string{"create", "-p"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs create [-ps] [-b blocksize] [-o property=value] ... -V <size> <volume>
func (z *zfsctl) CreateVolume(name string, block int64, properties map[string]string, size string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs destroy [-fnpRrv] <filesystem|volume>
func (z *zfsctl) DestroyFileSystemOrVolume(name, options string) *execute {
	args := []string{"destroy"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs destroy [-dnpRrv] <filesystem|volume>@<snap>[%<snap>][,...]
func (z *zfsctl) DestroySnapshot(name, options string) *execute {
	args := []string{"destroy"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs destroy <filesystem|volume>#<bookmark>
func (z *zfsctl) DestroyBookmark(name string) *execute {
	args := []string{"destroy", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs snapshot|snap [-r] [-o property=value] ... <filesystem|volume>@<snap> ...
func (z *zfsctl) Snapshot(name string, properties map[string]string) *execute {
	args := []string{"snapshot", "-r"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs rollback [-rRf] <snapshot>
func (z *zfsctl) Rollback(options string, snap string) *execute {
	args := []string{"rollback"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, snap)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs clone [-p] [-o property=value] ... <snapshot> <filesystem|volume>
func (z *zfsctl) Clone(name string, properties map[string]string, source string) *execute {
	args := []string{"clone", "-p"}
	if properties != nil {
		kv := "-o "
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, source, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs promote <clone-filesystem>
func (z *zfsctl) Promote(name string) *execute {
	args := []string{"promote", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs rename [-f] <filesystem|volume|snapshot> <filesystem|volume|snapshot>
func (z *zfsctl) Rename(name, newName string, force bool) *execute {
	args := []string{"rename"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name, newName)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs rename [-f] -p <filesystem|volume> <filesystem|volume>
func (z *zfsctl) RenameFileSystemOrVolume(name, newName string, force bool) *execute {
	args := []string{"rename"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, "-p")
	args = append(args, name, newName)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs rename -r <snapshot> <snapshot>
func (z *zfsctl) RenameSnapshot(name, newName string) *execute {
	args := []string{"rename", "-r"}
	args = append(args, name, newName)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs bookmark <snapshot> <bookmark>
func (z *zfsctl) Bookmark(snapshot, bookmark string) *execute {
	args := []string{"bookmark", snapshot, bookmark}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs list [-Hp] [-r|-d max] [-o property[,...]] [-s property]...
//            [-S property]... [-t type[,...]] [filesystem|volume|snapshot] ...
func (z *zfsctl) List(name, options, max string, oProperties []string, sProperty, SProperty, t string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs set <property=value> ... <filesystem|volume|snapshot> ...
func (z *zfsctl) Set(name string, properties map[string]string) *execute {
	args := []string{"set"}
	if properties != nil {
		kv := ""
		for k, v := range properties {
			kv += fmt.Sprintf("%s=%s ", k, v)
		}
		args = append(args, kv)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs get [-rHp] [-d max] [-o "all" | field[,...]]
//            [-t type[,...]] [-s source[,...]]
//            <"all" | property[,...]> [filesystem|volume|snapshot|bookmark] ...
func (z *zfsctl) Get(name, options, max string, out []string, t, s string, properties ...string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs inherit [-rS] <property> <filesystem|volume|snapshot> ...
func (z *zfsctl) Inherit(name string, options string, property string) *execute {
	args := []string{"inherit"}
	if len(options) != 0 {
		args = append(args, options)
	}
	args = append(args, property, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs upgrade [-v]
func (z *zfsctl) Upgrade(v bool) *execute {
	args := []string{"upgrade"}
	if v {
		args = append(args, "-v")
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs upgrade [-r] [-V version] <-a | filesystem ...>
func (z *zfsctl) UpgradeFileSystem(name, version string, r, all bool) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs userspace [-Hinp] [-o field[,...]] [-s field] ...
//            [-S field] ... [-t type[,...]] <filesystem|snapshot>
func (z *zfsctl) Userspace(name, options string, fields []string, sField, SField, t string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs groupspace [-Hinp] [-o field[,...]] [-s field] ...
//            [-S field] ... [-t type[,...]] <filesystem|snapshot>
func (z *zfsctl) Groupspace(name, options string, fields []string, sField, SField, t string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs mount
func (z *zfsctl) Mount() *execute {
	args := []string{"mount"}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs mount [-vO] [-o opts] <-a | filesystem>
func (z *zfsctl) MountFileSystem(name, options, opts string, all bool) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unmount [-f] <-a | filesystem|mountpoint>
func (z *zfsctl) Umount(name string, force, all bool) *execute {
	args := []string{"umount"}
	if force {
		args = append(args, "-f")
	}
	if all {
		args = append(args, "-a")
	} else {
		args = append(args, name)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs share <-a [nfs|smb] | filesystem>
func (z *zfsctl) Share(name string, all bool, kind string) *execute {
	args := []string{"share"}
	if all {
		args = append(args, "-a")
		if len(kind) > 0 {
			args = append(args, kind)
		}
	} else {
		args = append(args, name)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unshare <-a [nfs|smb] | filesystem|mountpoint>
func (z *zfsctl) Unshare(name string, all bool, kind string) *execute {
	args := []string{"unshare"}
	if all {
		args = append(args, "-a")
		if len(kind) > 0 {
			args = append(args, kind)
		}
	} else {
		args = append(args, name)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs send [-DnPpRvLec] [-[i|I] snapshot] <snapshot>
func (z *zfsctl) SendSnapshot(name, options string, i string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(i) > 0 {
		args = append(args, "-i "+i)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs send tank/sla@snap | ssh storage@192.168.2.120 -i test "sudo /usr/sbin/zfs receive -Fu dup/sla@snap"
func (z *zfsctl) SendAndRecv(source, target, user, host string) *execute {
	args := []string{"send", source}
	if len(user) != 0 && len(host) != 0 {
		args = append(args, fmt.Sprintf(`| ssh %s@%s -i %s "sudo %s receive -Fu %s"`, user, host, user, z.cmd, target))
	} else {
		args = append(args, fmt.Sprintf(`| %s receive -Fu %s`, z.cmd, target))
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs send -i tank/sla@snap1 tank/sla@snap2 | ssh storage@192.168.2.120 -i test "sudo /usr/sbin/zfs receive -Fu dup/sla@snap2"
func (z *zfsctl) IncrementSendAndRecv(source, lastSnapshot, target, user, host string) *execute {
	args := []string{"send", "-i", lastSnapshot, source}
	if len(user) != 0 && len(host) != 0 {
		args = append(args, fmt.Sprintf(`| ssh %s@%s -i %s "sudo %s receive -Fu %s"`, user, host, user, z.cmd, target))
	} else {
		args = append(args, fmt.Sprintf(`| %s receive -Fu %s`, z.cmd, target))
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs send [-Lec] [-i snapshot|bookmark] <filesystem|volume|snapshot>
func (z *zfsctl) Send(name, options string, i string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(i) > 0 {
		args = append(args, "-i "+i)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs send [-nvPe] -t <receive_resume_token>
func (z *zfsctl) SendToken(name, options string) *execute {
	args := []string{"send"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, "-t", name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs receive [-vnsFu] [-o <property>=<value>] ... [-x <property>] ...
//            <filesystem|volume|snapshot>
func (z *zfsctl) Receive(name, options string, properties map[string]string, xProperty string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs receive [-vnsFu] [-o <property>=<value>] ... [-x <property>] ...
//            <filesystem|volume|snapshot>
func (z *zfsctl) ReceiveFileSystem(name, options string, properties map[string]string, xProperty string, de string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs receive -A <filesystem|volume>
func (z *zfsctl) ReceiveAll(name string) *execute {
	args := []string{"receive", "-A", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs allow <filesystem|volume>
func (z *zfsctl) Allow1(name string) *execute {
	args := []string{"allow", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs allow [-ldug] <"everyone"|user|group>[,...] <perm|@setname>[,...]
//            <filesystem|volume>
func (z *zfsctl) Allow2(name, options, authority, perm string) *execute {
	args := []string{"allow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, authority, perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs allow [-ld] -e <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow3(name, options, perm string) *execute {
	args := []string{"allow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs allow -c <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow4(name, perm string) *execute {
	args := []string{"allow", "-c", perm, name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs allow -s @setname <perm|@setname>[,...] <filesystem|volume>
func (z *zfsctl) Allow5(name, setname, perm string) *execute {
	args := []string{"allow", "-s", setname, perm, name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unallow [-rldug] <"everyone"|user|group>[,...]
//            [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow1(name, options, authority, perm string) *execute {
	args := []string{"unallow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, authority, perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unallow [-rld] -e [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow2(name, options, perm string) *execute {
	args := []string{"unallow"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unallow [-r] -c [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow3(name string, r bool, perm string) *execute {
	args := []string{"unallow"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, "-c", perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs unallow [-r] -s @setname [<perm|@setname>[,...]] <filesystem|volume>
func (z *zfsctl) Unallow4(name string, r bool, setmame, perm string) *execute {
	args := []string{"unallow"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, "-s", setmame, perm, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs hold [-r] <tag> <snapshot> ...
func (z *zfsctl) Hold(name string, r bool, tag string) *execute {
	args := []string{"hold"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, tag, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs holds [-r] <snapshot> ...
func (z *zfsctl) Holds(name string, r bool) *execute {
	args := []string{"holds"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs  release [-r] <tag> <snapshot> ...
func (z *zfsctl) Release(name string, r bool, tag string) *execute {
	args := []string{"release"}
	if r {
		args = append(args, "-r")
	}
	args = append(args, tag, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zfs diff [-FHt] <snapshot> [snapshot|filesystem]
func (z *zfsctl) Diff(name, options string, target string) *execute {
	args := []string{"diff"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	if len(target) > 0 {
		args = append(args, target)
	}
	return &execute{name: z.cmd, args: args}
}

type zpoolctl struct {
	cmd string
}

func ZPoolCtl(cmd string) *zpoolctl {
	z := &zpoolctl{cmd: cmd}
	return z
}

// Examples:
// 	zpool create [-fnd] [-o property=value] ...
//            [-O file-system-property=value] ...
//            [-m mountpoint] [-R root] <pool> <vdev>
func (z *zpoolctl) Create(name, options, point, root string, properties map[string]string, raid string, devs ...string) *execute {
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

	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zpool destroy [-f] <pool>
func (z *zpoolctl) Destroy(name string, force bool) *execute {
	args := []string{"destroy"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool add [-fgLnP] [-o property=value] <pool> <vdev>
func (z *zpoolctl) Add(name string, options string, devs ...string) *execute {
	args := []string{"add"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool remove <pool> <device>
func (z *zpoolctl) Remove(name string, devs ...string) *execute {
	args := []string{"remove", name}
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool labelclear [-f] <vdev>
func (z *zpoolctl) LabelClear(device string, force bool) *execute {
	args := []string{"labelclear"}
	if force {
		args = append(args, "-f", device)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool list [-gHLpPv] [-o property[,...]] [-T d|u] [pool]
func (z *zpoolctl) List(name, options string, properties []string, t string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zpool iostat [[[-c [script1,script2,...][-lq]]|[-rw]] [-T d | u] [-ghHLpPvy]
//            [[pool ...]|[pool vdev ...]|[vdev ...]] [interval [count]]
func (z *zpoolctl) Iostat(name, scripts, t, options string, devs ...string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
// 	zpool  status [-c [script1,script2,...]] [-gLPvxD][-T d|u] [pool] ... [interval [count]]
func (z *zpoolctl) Status(scripts, options, t string, name string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool online <pool> <device> ...
func (z *zpoolctl) Online(name string, devs ...string) *execute {
	args := []string{"online"}
	args = append(args, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool offline [-f] [-t] <pool> <device> ...
func (z *zpoolctl) Offline(name string, force, t bool, devs ...string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool clear [-nF] <pool> [device]
func (z *zpoolctl) Clear(name, options string, devs ...string) *execute {
	args := []string{"clear"}
	args = append(args, options, name)
	for _, dev := range devs {
		args = append(args, dev)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool reopen <pool>
func (z *zpoolctl) Reopen(name string) *execute {
	args := []string{"reopen", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool attach [-f] [-o property=value] <pool> <device> <new-device>
func (z *zpoolctl) Attach(name string, force bool, properties map[string]string, dev, newDev string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool detach <pool> <device>
func (z *zpoolctl) Detach(name, dev string) *execute {
	args := []string{"attach", name, dev}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool replace [-f] [-o property=value] <pool> <device> [new-device]
func (z *zpoolctl) Replace(name string, force bool, properties map[string]string, dev string, newDev string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool split [-gLnP] [-R altroot] [-o mntopts]
//            [-o property=value] <pool> <newpool> [<device> ...]
func (z *zpoolctl) Split(name, options, altRoot, mntopts string, properties map[string]string, newName string, devs ...string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool scrub [-s | -p] <pool> ...
func (z *zpoolctl) scrub(name, option string) *execute {
	args := []string{"scrub"}
	if len(option) > 0 {
		args = append(args, option)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool import [-d dir] [-D]
func (z *zpoolctl) Import1(dir string, d bool) *execute {
	args := []string{"import"}
	if len(dir) > 0 {
		args = append(args, dir)
	}
	if d {
		args = append(args, "-D")
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool import [-d dir] [-D]
func (z *zpoolctl) Import2(dir string, d bool) *execute {
	args := []string{"import"}
	if len(dir) > 0 {
		args = append(args, dir)
	}
	if d {
		args = append(args, "-D")
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool import [-d dir | -c cachefile] [-F [-n]] <pool | id>
func (z *zpoolctl) Import3(name, dir, file string, force, n bool) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool import [-o mntopts] [-o property=value] ...
//            [-d dir | -c cachefile] [-D] [-f] [-m] [-N] [-R root] [-F [-n]] -a
func (z *zpoolctl) Import4(
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool import [-o mntopts] [-o property=value] ...
//            [-d dir | -c cachefile] [-D] [-f] [-m] [-N] [-R root] [-F [-n]]
//            <pool | id> [newpool]
func (z *zpoolctl) Import5(
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool zpoexport [-af] <pool> ...
func (z *zpoolctl) Export(name string, options string) *execute {
	args := []string{"export"}
	if len(options) > 0 {
		args = append(args, options)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool upgrade [-v]
func (z *zpoolctl) Upgrade1(v bool) *execute {
	args := []string{"upgrade"}
	if v {
		args = append(args, "-v")
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool upgrade [-V version] <-a | pool ...>
func (z *zpoolctl) Upgrade2(name, version string) *execute {
	args := []string{"upgrade"}
	if len(version) > 0 {
		args = append(args, "-V "+version)
	}
	args = append(args, name)
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool reguid <pool>
func (z *zpoolctl) Reguid(name string) *execute {
	args := []string{"reguid", name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool history [-il] [<pool>]
func (z *zpoolctl) History(options string, pool string) *execute {
	args := []string{"history"}
	if len(options) > 0 {
		args = append(args, options)
	}
	if len(pool) > 0 {
		args = append(args, pool)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool events [-vHfc]
func (z *zpoolctl) Event(options string) *execute {
	args := []string{"event"}
	if len(options) > 0 {
		args = append(args, options)
	}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool get [-Hp] [-o "all" | field[,...]] <"all" | property[,...]> <pool> ...
func (z *zpoolctl) Get(name, options string, out []string, properties ...string) *execute {
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
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool set <property=value> <pool>
func (z *zpoolctl) Set(name, k, v string) *execute {
	args := []string{"set", k + "=" + v, name}
	return &execute{name: z.cmd, args: args}
}

// Examples:
//	zpool sync [pool] ...
func (z *zpoolctl) Sync(names ...string) *execute {
	args := []string{"sync"}
	if names != nil {
		for _, name := range names {
			args = append(args, name)
		}
	}
	return &execute{name: z.cmd, args: args}
}

type execute struct {
	name string
	args []string
}

func (e *execute) Commit() string {
	return fmt.Sprintf(`%s %s`, e.name, strings.Join(e.args, " "))
}

func (e *execute) Exec() ([]byte, error) {
	return execution(fmt.Sprintf(`%s %s`, e.name, strings.Join(e.args, " ")))
}

func execution(shell string) ([]byte, error) {
	data, err := exec.Command("/bin/sh", "-c", shell).CombinedOutput()
	if bytes.HasSuffix(data, []byte("\n")) {
		data = bytes.TrimSuffix(data, []byte("\n"))
	}
	if err != nil {
		return nil, fmt.Errorf("%v:%v", err, string(data))
	}
	return data, nil
}

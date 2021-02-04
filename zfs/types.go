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

const (
	// Pool, Volume, FileSystem, Snapshot 的 lb:status 属性
	// ready 表示正在被使用中
	Ready = "ready"
	// recycling 表示等待删除，每隔一段时间会进行一次垃圾回收
	// lb:status 为 recycling 且不再被其他资源依赖时就会被删除
	Recycling = "recycling"
)

const (
	// lb:display 的属性
	// show 表示资源可见
	Show = "show"
	// hidden 表示资源不可见，一般作为其他资源的辅助使用，不存储数据
	Hidden = "hidden"
)

const (
	STATUS       = "lb:status"
	DISPLAY      = "lb:display"
	RELATIONSHIP = "lb:relationship"
	SHADOW       = "lb:shadow"
)

type System struct {
	Pools map[string]*Pool `json:"pools"`

	FileSystems map[string]*Filesystem `json:"fileSystems"`

	Volumes map[string]*Volume `json:"volumes"`

	Snapshots map[string]*Snapshot `json:"snapshot"`
}

type Pool struct {
	Name string `json:"name"`
	// zfs pool 容量
	Size string `json:"size,omitempty" zfs:"size"`

	Capacity string `json:"capacity,omitempty" zfs:"capacity"`

	Altroot string `json:"altroot,omitempty" zfs:"altroot"`

	Health string `json:"health,omitempty" zfs:"health"`

	Guid string `json:"guid,omitempty" zfs:"guid"`

	Version string `json:"version,omitempty" zfs:"version"`

	Bootfs string `json:"bootfs,omitempty" zfs:"bootfs"`

	Delegation string `json:"delegation,omitempty" zfs:"delegation"`

	AutoReplace string `json:"autoReplace,omitempty" zfs:"authreplace"`

	Cachefile string `json:"cachefile,omitempty" zfs:"cachefile"`

	FailMode string `json:"failMode,omitempty" zfs:"failmode"`

	Listsnapshots string `json:"listsnapshots,omitempty" zfs:"listsnapshots"`

	AutoExpand string `json:"autoExpand,omitempty" zfs:"autoexpand"`

	Dedupditto string `json:"dedupditto,omitempty" zfs:"dedupditto"`

	Dedupratio string `json:"dedupratio,omitempty" zfs:"dedupratio"`

	// zfs pool 可用空间
	Free string `json:"free,omitempty" zfs:"free"`

	// zfs pool 已分配空间
	Allocated string `json:"allocated,omitempty" zfs:"allocated"`

	Readonly string `json:"readonly,omitempty" zfs:"readonly"`

	Ashift string `json:"ashift,omitempty" zfs:"ashift"`

	Comment string `json:"comment,omitempty" zfs:"comment"`

	Expandsize string `json:"expandsize,omitempty" zfs:"expandsize"`

	Freeing string `json:"freeing,omitempty" zfs:"freeing"`

	Fragmentation string `json:"fragmentation,omitempty" zfs:"fragmentation"`

	Leaked string `json:"leaked,omitempty" zfs:"leaked"`

	Multihost string `json:"multihost,omitempty" zfs:"multihost"`

	FeatureAsyncDestroy string `json:"featureAsyncDestroy,omitempty" zfs:"feature@async_destroy"`

	FeatureEmptyBpobj string `json:"featureEmptyBpobj,omitempty" zfs:"feature@empty_bpobj"`

	FeatureLz4Compress string `json:"featureLz4Compress,omitempty" zfs:"feature@lz4_compress"`

	FeatureMultiVdevCrashDump string `json:"featureMultiVdevCrashDump,omitempty" zfs:"feature@multi_vdev_crash_dump"`

	FeatureSpacemapHistogram string `json:"featureSpacemapHistogram,omitempty" zfs:"feature@spacemap_histogram"`

	FeatureEnabledTxg string `json:"featureEnabledTxg,omitempty" zfs:"feature@enabled_txg"`

	FeatureHoleBirth string `json:"featureHoleBirth,omitempty" zfs:"feature@hole_birth"`

	FeatureExtensibleDataset string `json:"featureExtensibleDataset,omitempty" zfs:"feature@extensible_dataset"`

	FeatureEmbeddedData string `json:"featureEmbeddedData,omitempty" zfs:"feature@embedded_data"`

	FeatureBookmarks string `json:"featureBookmarks,omitempty" zfs:"feature@bookmarks"`

	FeatureFilesystemLimits string `json:"featureFilesystemLimits,omitempty" zfs:"feature@filesystem_limits"`

	FeatureLargeBlocks string `json:"featureLargeBlocks,omitempty" zfs:"feature@large_blocks"`

	FeatureLargeDnode string `json:"featureLargeDnode,omitempty" zfs:"feature@large_dnode"`

	FeatureSha512 string `json:"featureSha512,omitempty" zfs:"feature@sha512"`

	FeatureSkein string `json:"featureSkein,omitempty" zfs:"feature@skein"`

	FeatureEdonr string `json:"featureEdonr,omitempty" zfs:"feature@edonr "`

	FeatureUserobjAccounting string `json:"featureUserobjAccounting,omitempty" zfs:"feature@userobj_accounting"`

	// tank filesystem
	Type string `json:"type,omitempty" zfs:"type"`

	Creation string `json:"creation,omitempty" zfs:"creation"`

	Used string `json:"used,omitempty" zfs:"used"`

	Available string `json:"available,omitempty" zfs:"available"`

	Referenced string `json:"referenced,omitempty" zfs:"referenced"`

	Compressratio string `json:"compressratio,omitempty" zfs:"compressratio"`

	Mounted string `json:"mounted,omitempty" zfs:"mounted"`

	Quota string `json:"quota,omitempty" zfs:"quota"`

	Reservation string `json:"reservation,omitempty" zfs:"reservation"`

	Recordsize string `json:"recordsize,omitempty" zfs:"recordsize"`

	Mountpoint string `json:"mountpoint,omitempty" zfs:"mountpoint"`

	Sharenfs string `json:"sharenfs,omitempty" zfs:"sharenfs"`

	Checksum string `json:"checksum,omitempty" zfs:"checksum"`

	Compression string `json:"compression,omitempty" zfs:"compression"`

	Atime string `json:"atime,omitempty" zfs:"atime"`

	Devices string `json:"devices,omitempty" zfs:"devices"`

	Exec string `json:"exec,omitempty" zfs:"exec"`

	Setuid string `json:"setuid,omitempty" zfs:"setuid"`

	Zoned string `json:"zoned,omitempty" zfs:"zoned"`

	Snapdir string `json:"snapdir,omitempty" zfs:"snapdir"`

	Aclinherit string `json:"aclinherit,omitempty" zfs:"aclinherit"`

	Createtxg string `json:"createtxg,omitempty" zfs:"createtxg"`

	Canmount string `json:"canmount,omitempty" zfs:"canmount"`

	Xattr string `json:"xattr,omitempty" zfs:"xattr"`

	Copies string `json:"copies,omitempty" zfs:"copies"`

	Utf8only string `json:"utf8Only,omitempty" zfs:"utf8only"`

	Normalization string `json:"normalization,omitempty" zfs:"normalization"`

	Casesensitivity string `json:"casesensitivity,omitempty" zfs:"casesensitivity"`

	Vscan string `json:"vscan,omitempty" zfs:"vscan"`

	Nbmand string `json:"nbmand,omitempty" zfs:"nbmand"`

	Sharesmb string `json:"sharesmb,omitempty" zfs:"sharesmb"`

	Refquota string `json:"refquota,omitempty" zfs:"refquota"`

	Refreservation string `json:"refreservation,omitempty" zfs:"refreservation"`

	Primarycache string `json:"primarycache,omitempty" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache,omitempty" zfs:"secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots,omitempty" zfs:"usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset,omitempty" zfs:"usedbydataset"`

	Usedbychildren string `json:"usedbychildren,omitempty" zfs:"usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation,omitempty" zfs:"usedbyrefreservation"`

	Logbias string `json:"logbias,omitempty" zfs:"logbias"`

	Dedup string `json:"dedup,omitempty" zfs:"dedup"`

	Mlslabel string `json:"mlslabel,omitempty" zfs:"mlslabel"`

	Sync string `json:"sync,omitempty" zfs:"sync"`

	Dnodesize string `json:"dnodesize,omitempty" zfs:"dnodesize"`

	Refcompressratio string `json:"refcompressratio,omitempty" zfs:"refcompressratio"`

	Written string `json:"written,omitempty" zfs:"written"`

	Logicalused string `json:"logicalused,omitempty" zfs:"logicalused"`

	Logicalreferenced string `json:"logicalreferenced,omitempty" zfs:"logicalreferenced"`

	Volmode string `json:"volmode,omitempty" zfs:"volmode"`

	FilesystemLimit string `json:"filesystemLimit,omitempty" zfs:"filesystem_limit"`

	SnapshotLimit string `json:"snapshotLimit,omitempty" zfs:"snapshot_limit"`

	FilesystemCount string `json:"filesystemCount,omitempty" zfs:"filesystem_count"`

	SnapshotCount string `json:"snapshotCount,omitempty" zfs:"snapshot_count"`

	Snapdev string `json:"snapdev,omitempty" zfs:"snapdev"`

	Acltype string `json:"acltype,omitempty" zfs:"acltype"`

	Context string `json:"context,omitempty" zfs:"context"`

	Fscontext string `json:"fscontext,omitempty" zfs:"fscontext"`

	Defcontext string `json:"defcontext,omitempty" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext,omitempty" zfs:"rootcontext"`

	Relatime string `json:"relatime,omitempty" zfs:"relatime"`

	RedundantMetadata string `json:"redundantMetadata,omitempty" zfs:"redundant_metadata"`

	Overlay string `json:"overlay,omitempty" zfs:"overlay"`

	Encryption string `json:"encryption,omitempty" zfs:"encryption"`

	Keylocation string `json:"keylocation,omitempty" zfs:"keylocation"`

	Keyformat string `json:"keyformat,omitempty" zfs:"keyformat"`

	Pbkdf2iters string `json:"pbkdf2Iters,omitempty" zfs:"pbkdf2iters"`

	SpecialSmallBlocks string `json:"specialSmallBlocks,omitempty" zfs:"special_small_blocks"`
	// 自定义属性，用于垃圾回收
	Status string `json:"status,omitempty" zfs:"lb:status"`
	// Pool 下的 Volume, FileSystem, Snapshot, Clone
	Devs map[string]string `json:"devs,omitempty"`
}

type Filesystem struct {
	Pool string `json:"pool"`

	Name string `json:"name"`
	// 克隆该 Volume 的 Snapshot 名称
	Source string `json:"source,omitempty"`

	Type string `json:"type,omitempty" zfs:"type"`

	Creation string `json:"creation,omitempty" zfs:"creation"`

	Used string `json:"used,omitempty" zfs:"used"`

	Available string `json:"available,omitempty" zfs:"available"`

	Referenced string `json:"referenced,omitempty" zfs:"referenced"`

	Compressratio string `json:"compressratio,omitempty" zfs:"compressratio"`

	Mounted string `json:"mounted,omitempty" zfs:"mounted"`

	Quota string `json:"quota,omitempty" zfs:"quota"`

	Reservation string `json:"reservation,omitempty" zfs:"reservation"`

	Recordsize string `json:"recordsize,omitempty" zfs:"recordsize"`

	Mountpoint string `json:"mountpoint,omitempty" zfs:"mountpoint"`

	Sharenfs string `json:"sharenfs,omitempty" zfs:"sharenfs"`

	Checksum string `json:"checksum,omitempty" zfs:"checksum"`

	Compression string `json:"compression,omitempty" zfs:"compression"`

	Atime string `json:"atime,omitempty" zfs:"atime"`

	Devices string `json:"devices,omitempty" zfs:"devices"`

	Exec string `json:"exec,omitempty" zfs:"exec"`

	Setuid string `json:"setuid,omitempty" zfs:"setuid"`

	Readonly string `json:"readonly,omitempty" zfs:"readonly"`

	Zoned string `json:"zoned,omitempty" zfs:"zoned"`

	Snapdir string `json:"snapdir,omitempty" zfs:"snapdir"`

	Aclinherit string `json:"aclinherit,omitempty" zfs:"aclinherit"`

	Createtxg string `json:"createtxg,omitempty" zfs:"createtxg"`

	Canmount string `json:"canmount,omitempty" zfs:"canmount"`

	Xattr string `json:"xattr,omitempty" zfs:"xattr"`

	Copies string `json:"copies,omitempty" zfs:"copies"`

	Version string `json:"version,omitempty" zfs:"version"`

	Utf8only string `json:"utf8Only,omitempty" zfs:"utf8only"`

	Normalization string `json:"normalization,omitempty" zfs:"normalization"`

	Casesensitivity string `json:"casesensitivity,omitempty" zfs:"casesensitivity"`

	Vscan string `json:"vscan,omitempty" zfs:"vscan"`

	Nbmand string `json:"nbmand,omitempty" zfs:"nbmand"`

	Sharesmb string `json:"sharesmb,omitempty" zfs:"sharesmb"`

	Refquota string `json:"refquota,omitempty" zfs:"refquota"`

	Refreservation string `json:"refreservation,omitempty" zfs:"refreservation"`

	Guid string `json:"guid,omitempty" zfs:"guid"`

	Primarycache string `json:"primarycache,omitempty" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache,omitempty" zfs:"secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots,omitempty" zfs:"usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset,omitempty" zfs:"usedbydataset"`

	Usedbychildren string `json:"usedbychildren,omitempty" zfs:"usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation,omitempty" zfs:"usedbyrefreservation"`

	Logbias string `json:"logbias,omitempty" zfs:"logbias"`

	Dedup string `json:"dedup,omitempty" zfs:"dedup"`

	Mlslabel string `json:"mlslabel,omitempty" zfs:"mlslabel"`

	Sync string `json:"sync,omitempty" zfs:"sync"`

	Dnodesize string `json:"dnodesize,omitempty" zfs:"dnodesize"`

	Refcompressratio string `json:"refcompressratio,omitempty" zfs:"refcompressratio"`

	Written string `json:"written,omitempty" zfs:"written"`

	Logicalused string `json:"logicalused,omitempty" zfs:"logicalused"`

	Logicalreferenced string `json:"logicalreferenced,omitempty" zfs:"logicalreferenced"`

	Volmode string `json:"volmode,omitempty" zfs:"volmode"`

	FilesystemLimit string `json:"filesystemLimit,omitempty" zfs:"filesystem_limit"`

	SnapshotLimit string `json:"snapshotLimit,omitempty" zfs:"snapshot_limit"`

	FilesystemCount string `json:"filesystemCount,omitempty" zfs:"filesystem_count"`

	SnapshotCount string `json:"snapshotCount,omitempty" zfs:"snapshot_count"`

	Snapdev string `json:"snapdev,omitempty" zfs:"snapdev"`

	Acltype string `json:"acltype,omitempty" zfs:"acltype"`

	Context string `json:"context,omitempty" zfs:"context"`

	Fscontext string `json:"fscontext,omitempty" zfs:"fscontext"`

	Defcontext string `json:"defcontext,omitempty" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext,omitempty" zfs:"rootcontext"`

	Relatime string `json:"relatime,omitempty" zfs:"relatime"`

	RedundantMetadata string `json:"redundantMetadata,omitempty" zfs:"redundant_metadata"`

	Overlay string `json:"overlay,omitempty" zfs:"overlay"`
	// 自定义属性,是否可见
	Display string `json:"display,omitempty" zfs:"lb:display"`
	// Shadow, FileSystem 在创建时会默认生成一份快照
	// 防止快照发送时由于只剩下一份快照而无法进行增量发送
	Shadow string `json:"shadow,omitempty" zfs:"lb:shadow"`
	// 自定义属性
	Status string `json:"status,omitempty" zfs:"lb:status"`
	// FileSystem 下的 Snapshot
	Snapshots map[string]string `json:"snapshots,omitempty"`
}

type Volume struct {
	Pool string `json:"pool"`

	Name string `json:"name"`
	// 克隆该 Volume 的 Snapshot 名称
	Source string `json:"source,omitempty"`

	Type string `json:"type,omitempty" zfs:"type"`

	Creation string `json:"creation,omitempty" zfs:"creation"`

	Used string `json:"used,omitempty" zfs:"used"`

	Available string `json:"available,omitempty" zfs:"available"`

	Referenced string `json:"referenced,omitempty" zfs:"referenced"`

	Compressratio string `json:"compressratio,omitempty" zfs:"compressratio"`

	Reservation string `json:"reservation,omitempty" zfs:"reservation"`

	Volsize string `json:"volsize,omitempty" zfs:"volsize"`

	Volblocksize string `json:"volblocksize,omitempty" zfs:"volblocksize"`

	Checksum string `json:"checksum,omitempty" zfs:"checksum"`

	Compression string `json:"compression,omitempty" zfs:"compression"`

	Readonly string `json:"readonly,omitempty" zfs:"readonly"`

	Copies string `json:"copies,omitempty" zfs:"copies"`

	Refreservation string `json:"refreservation,omitempty" zfs:"refreservation"`

	Guid string `json:"guid,omitempty" zfs:"guid"`

	Primarycache string `json:"primarycache,omitempty" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache,omitempty" zfs:"secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots,omitempty" zfs:"usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset,omitempty" zfs:"usedbydataset"`

	Usedbychildren string `json:"usedbychildren,omitempty" zfs:"usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation,omitempty" zfs:"usedbyrefreservation"`

	Logbias string `json:"logbias,omitempty" zfs:"logbias"`

	Dedup string `json:"dedup,omitempty" zfs:"dedup"`

	Mlslabel string `json:"mlslabel,omitempty" zfs:"mlslabel"`

	Sync string `json:"sync,omitempty" zfs:"sync"`

	Refcompressratio string `json:"refcompressratio,omitempty" zfs:"refcompressratio"`

	Written string `json:"written,omitempty" zfs:"written"`

	Logicalused string `json:"logicalused,omitempty" zfs:"logicalused"`

	Logicalreferenced string `json:"logicalreferenced,omitempty" zfs:"logicalreferenced"`

	Volmode string `json:"volmode,omitempty" zfs:"volmode"`

	FilesystemLimit string `json:"filesystemLimit,omitempty" zfs:"filesystem_limit"`

	SnapshotLimit string `json:"snapshotLimit,omitempty" zfs:"snapshot_limit"`

	SnapshotCount string `json:"snapshotCount,omitempty" zfs:"snapshot_count"`

	Snapdev string `json:"snapdev,omitempty" zfs:"snapdev"`

	Acltype string `json:"acltype,omitempty" zfs:"acltype"`

	Context string `json:"context,omitempty" zfs:"context"`

	Fscontext string `json:"fscontext,omitempty" zfs:"fscontext"`

	Defcontext string `json:"defcontext,omitempty" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext,omitempty" zfs:"rootcontext"`

	RedundantMetadata string `json:"redundantMetadata,omitempty" zfs:"redundant_metadata"`
	// 自定义属性,是否可见
	Display string `json:"display,omitempty" zfs:"lb:display"`
	// Shadow, Volume 在创建时会默认生成一份快照
	// 防止快照发送时由于只剩下一份快照而无法进行增量发送
	Shadow string `json:"shadow,omitempty" zfs:"lb:shadow"`
	// 自定义属性
	Status string `json:"status,omitempty" zfs:"lb:status"`
	// FileSystem 下的 Snapshot
	Snapshots map[string]string `json:"snapshots,omitempty"`
}

type Snapshot struct {
	Pool string `json:"pool"`

	Parent string `json:"parent"`

	Name string `json:"name"`

	Type string `json:"type,omitempty" zfs:"type"`

	Creation string `json:"creation,omitempty" zfs:"creation"`

	Used string `json:"used,omitempty" zfs:"used"`

	Referenced string `json:"referenced,omitempty" zfs:"referenced"`

	Compressratio string `json:"compressratio,omitempty" zfs:"compressratio"`

	Devices string `json:"devices,omitempty" zfs:"devices"`

	Exec string `json:"exec,omitempty" zfs:"exec"`

	Setuid string `json:"setuid,omitempty" zfs:"setuid"`

	Createtxg string `json:"createtxg,omitempty" zfs:"createtxg"`

	Version string `json:"version,omitempty" zfs:"version"`

	Utf8only string `json:"utf8Only,omitempty" zfs:"utf8only"`

	Normalization string `json:"normalization,omitempty" zfs:"normalization"`

	Casesensitivity string `json:"casesensitivity,omitempty" zfs:"casesensitivity"`

	Vscan string `json:"vscan,omitempty" zfs:"vscan"`

	Nbmand string `json:"nbmand,omitempty" zfs:"nbmand"`

	Guid string `json:"guid,omitempty" zfs:"guid"`

	Primarycache string `json:"primarycache,omitempty" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache,omitempty" zfs:"secondarycache"`

	DeferDestroy string `json:"deferDestroy,omitempty" zfs:"defer_destroy"`

	Userrefs string `json:"userrefs,omitempty" zfs:"userrefs"`

	Mlslabel string `json:"mlslabel,omitempty" zfs:"mlslabel"`

	Refcompressratio string `json:"refcompressratio,omitempty" zfs:"refcompressratio"`

	Written string `json:"written,omitempty" zfs:"written"`

	Logicalreferenced string `json:"logicalreferenced,omitempty" zfs:"logicalreferenced"`

	Acltype string `json:"acltype,omitempty" zfs:"acltype"`

	Context string `json:"context,omitempty" zfs:"context"`

	Fscontext string `json:"fscontext,omitempty" zfs:"fscontext"`

	Defcontext string `json:"defcontext,omitempty" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext,omitempty" zfs:"rootcontext"`
	// 自定义属性,是否可见
	Display string `json:"display,omitempty" zfs:"lb:display"`
	// Snapshot 下的 Clone
	Clones map[string]string `json:"clones,omitempty" zfs:"clones"`
	// 自定义属性
	Status string `json:"status,omitempty" zfs:"lb:status"`
	// 自定义属性, snapshot send 后关联的 filesystem 或是 volume
	// 存在这个属性时，删除 snapshot 时同时删除 filesystem 或是 volume
	RelationShip string `json:"relationShip,omitempty" zfs:"lb:relationship"`
}

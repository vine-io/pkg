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
	// Pool, Volume, Snapshot lb:status property
	Ready = "ready"
	Recycling = "recycling"
)

const (
	STATUS = "lb:status"
)

type VolumeKind string

const (
	VKVolume     VolumeKind = "volume"
	VKFileSystem VolumeKind = "filesystem"
)

type PoolTypes struct {
	Name string `json:"name"`
	// zfs pool 容量
	Size string `json:"size" zfs:"size"`

	Capacity string `json:"capacity" zfs:"capacity"`

	AltRoot string `json:"altRoot" zfs:"altroot"`

	Health string `json:"health" zfs:"health"`

	Guid string `json:"guid" zfs:"guid"`

	Version string `json:"version" zfs:"version"`

	BootFs string `json:"bootFs" zfs:"bootfs"`

	Delegation string `json:"delegation" zfs:"delegation"`

	AutoReplace string `json:"autoReplace" zfs:"authreplace"`

	CacheFile string `json:"cacheFile" zfs:"cachefile"`

	FailMode string `json:"failMode" zfs:"failmode"`

	ListSnapshots string `json:"listSnapshots" zfs:"listsnapshots"`

	AutoExpand string `json:"autoExpand" zfs:"autoexpand"`

	DedupDitto string `json:"dedupDitto" zfs:"dedupditto"`

	DedupRatio string `json:"dedupRatio" zfs:"dedupratio"`

	// zfs pool 可用空间
	Free string `json:"free" zfs:"free"`

	// zfs pool 已分配空间
	Allocated string `json:"allocated" zfs:"allocated"`

	Readonly string `json:"readonly" zfs:"readonly"`

	AShift string `json:"aShift" zfs:"ashift"`

	Comment string `json:"comment" zfs:"comment"`

	ExpandSize string `json:"expandSize" zfs:"expandsize"`

	Freeing string `json:"freeing" zfs:"freeing"`

	Fragmentation string `json:"fragmentation" zfs:"fragmentation"`

	Leaked string `json:"leaked" zfs:"leaked"`

	Multihost string `json:"multihost" zfs:"multihost"`

	FeatureAsyncDestroy string `json:"featureAsyncDestroy" zfs:"feature@async_destroy"`

	FeatureEmptyBpobj string `json:"featureEmptyBpobj" zfs:"feature@empty_bpobj"`

	FeatureLz4Compress string `json:"featureLz4Compress" zfs:"feature@lz4_compress"`

	FeatureMultiVdevCrashDump string `json:"featureMultiVdevCrashDump" zfs:"feature@multi_vdev_crash_dump"`

	FeatureSpaceMapHistogram string `json:"featureSpaceMapHistogram" zfs:"feature@spacemap_histogram"`

	FeatureEnabledTxg string `json:"featureEnabledTxg" zfs:"feature@enabled_txg"`

	FeatureHoleBirth string `json:"featureHoleBirth" zfs:"feature@hole_birth"`

	FeatureExtensibleDataset string `json:"featureExtensibleDataset" zfs:"feature@extensible_dataset"`

	FeatureEmbeddedData string `json:"featureEmbeddedData" zfs:"feature@embedded_data"`

	FeatureBookmarks string `json:"featureBookmarks" zfs:"feature@bookmarks"`

	FeatureFilesystemLimits string `json:"featureFilesystemLimits" zfs:"feature@filesystem_limits"`

	FeatureLargeBlocks string `json:"featureLargeBlocks" zfs:"feature@large_blocks"`

	FeatureLargeDnode string `json:"featureLargeDnode" zfs:"feature@large_dnode"`

	FeatureSha512 string `json:"featureSha512" zfs:"feature@sha512"`

	FeatureSkein string `json:"featureSkein" zfs:"feature@skein"`

	FeatureEdonr string `json:"featureEdonr" zfs:"feature@edonr "`

	FeatureUserobjAccounting string `json:"featureUserobjAccounting" zfs:"feature@userobj_accounting"`

	// tank filesystem
	Type string `json:"type" zfs:"type"`

	Creation string `json:"creation" zfs:"creation"`

	Used string `json:"used" zfs:"used"`

	Available string `json:"available" zfs:"available"`

	Referenced string `json:"referenced" zfs:"referenced"`

	CompressRatio string `json:"compressRatio" zfs:"compressratio"`

	Mounted string `json:"mounted" zfs:"mounted"`

	Quota string `json:"quota" zfs:"quota"`

	Reservation string `json:"reservation" zfs:"reservation"`

	RecordSize string `json:"recordSize" zfs:"recordsize"`

	MountPoint string `json:"mountPoint" zfs:"mountpoint"`

	SharenFs string `json:"sharenFs" zfs:"sharenfs"`

	Checksum string `json:"checksum" zfs:"checksum"`

	Compression string `json:"compression" zfs:"compression"`

	ATime string `json:"aTime" zfs:"atime"`

	Devices string `json:"devices" zfs:"devices"`

	Exec string `json:"exec" zfs:"exec"`

	SetUid string `json:"setUid" zfs:"setuid"`

	Zoned string `json:"zoned" zfs:"zoned"`

	SnapDir string `json:"snapDir" zfs:"snapdir"`

	AclinHerit string `json:"aclinHerit" zfs:"aclinherit"`

	Createtxg string `json:"createtxg" zfs:"createtxg"`

	CanMount string `json:"canMount" zfs:"canmount"`

	Xattr string `json:"xattr" zfs:"xattr"`

	Copies string `json:"copies" zfs:"copies"`

	Utf8only string `json:"utf8Only" zfs:"utf8only"`

	Normalization string `json:"normalization" zfs:"normalization"`

	Casesensitivity string `json:"casesensitivity" zfs:"casesensitivity"`

	Vscan string `json:"vscan" zfs:"vscan"`

	Nbmand string `json:"nbmand" zfs:"nbmand"`

	Sharesmb string `json:"sharesmb" zfs:"sharesmb"`

	Refquota string `json:"refquota" zfs:"refquota"`

	Refreservation string `json:"refreservation" zfs:"refreservation"`

	Primarycache string `json:"primarycache" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots" zfs:"usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset" zfs:"usedbydataset"`

	Usedbychildren string `json:"usedbychildren" zfs:"usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation" zfs:"usedbyrefreservation"`

	Logbias string `json:"logbias" zfs:"logbias"`

	Dedup string `json:"dedup" zfs:"dedup"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel"`

	Sync string `json:"sync" zfs:"sync"`

	Dnodesize string `json:"dnodesize" zfs:"dnodesize"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio"`

	Written string `json:"written" zfs:"written"`

	Logicalused string `json:"logicalused" zfs:"logicalused"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced"`

	Volmode string `json:"volmode" zfs:"volmode"`

	FilesystemLimit string `json:"filesystemLimit" zfs:"filesystem_limit"`

	SnapshotLimit string `json:"snapshotLimit" zfs:"snapshot_limit"`

	FilesystemCount string `json:"filesystemCount" zfs:"filesystem_count"`

	SnapshotCount string `json:"snapshotCount" zfs:"snapshot_count"`

	Snapdev string `json:"snapdev" zfs:"snapdev"`

	Acltype string `json:"acltype" zfs:"acltype"`

	Context string `json:"context" zfs:"context"`

	Fscontext string `json:"fscontext" zfs:"fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext"`

	Relatime string `json:"relatime" zfs:"relatime"`

	RedundantMetadata string `json:"redundantMetadata" zfs:"redundant_metadata"`

	Overlay string `json:"overlay" zfs:"overlay"`

	Encryption string `json:"encryption" zfs:"encryption"`

	Keylocation string `json:"keylocation" zfs:"keylocation"`

	Keyformat string `json:"keyformat" zfs:"keyformat"`

	Pbkdf2iters string `json:"pbkdf2Iters" zfs:"pbkdf2iters"`

	SpecialSmallBlocks string `json:"specialSmallBlocks" zfs:"special_small_blocks"`
	// 自定义属性，用于垃圾回收
	Status string `json:"status" zfs:"lb:status"`
}

type VolumeTypes struct {
	Name string `json:"name"`

	Type string `json:"type" zfs:"type"`

	Creation string `json:"creation" zfs:"creation"`

	Used string `json:"used" zfs:"used"`

	Available string `json:"available" zfs:"available"`

	Referenced string `json:"referenced" zfs:"referenced"`

	Compressratio string `json:"compressratio" zfs:"compressratio"`

	Reservation string `json:"reservation" zfs:"reservation"`

	Volsize string `json:"volsize" zfs:"volsize"`

	Volblocksize string `json:"volblocksize" zfs:"volblocksize"`

	Checksum string `json:"checksum" zfs:"checksum"`

	Compression string `json:"compression" zfs:"compression"`

	Readonly string `json:"readonly" zfs:"readonly"`

	Copies string `json:"copies" zfs:"copies"`

	Refreservation string `json:"refreservation" zfs:"refreservation"`

	Guid string `json:"guid" zfs:"guid"`

	Primarycache string `json:"primarycache" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots" zfs:"usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset" zfs:"usedbydataset"`

	Usedbychildren string `json:"usedbychildren" zfs:"usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation" zfs:"usedbyrefreservation"`

	Logbias string `json:"logbias" zfs:"logbias"`

	Dedup string `json:"dedup" zfs:"dedup"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel"`

	Sync string `json:"sync" zfs:"sync"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio"`

	Written string `json:"written" zfs:"written"`

	Logicalused string `json:"logicalused" zfs:"logicalused"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced"`

	Volmode string `json:"volmode" zfs:"volmode"`

	FilesystemLimit string `json:"filesystemLimit" zfs:"filesystem_limit"`

	SnapshotLimit string `json:"snapshotLimit" zfs:"snapshot_limit"`

	SnapshotCount string `json:"snapshotCount" zfs:"snapshot_count"`

	Snapdev string `json:"snapdev" zfs:"snapdev"`

	Acltype string `json:"acltype" zfs:"acltype"`

	Context string `json:"context" zfs:"context"`

	Fscontext string `json:"fscontext" zfs:"fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext"`

	RedundantMetadata string `json:"redundantMetadata" zfs:"redundant_metadata"`
	// 自定义属性
	Status string `json:"status" zfs:"lb:status"`
}

type SnapshotTypes struct {
	Pool string `json:"pool"`

	Parent string `json:"parent"`

	Name string `json:"name"`

	Type string `json:"type" zfs:"type"`

	Creation string `json:"creation" zfs:"creation"`

	Used string `json:"used" zfs:"used"`

	Referenced string `json:"referenced" zfs:"referenced"`

	Compressratio string `json:"compressratio" zfs:"compressratio"`

	Devices string `json:"devices" zfs:"devices"`

	Exec string `json:"exec" zfs:"exec"`

	Setuid string `json:"setuid" zfs:"setuid"`

	Createtxg string `json:"createtxg" zfs:"createtxg"`

	Version string `json:"version" zfs:"version"`

	Utf8only string `json:"utf8Only" zfs:"utf8only"`

	Normalization string `json:"normalization" zfs:"normalization"`

	Casesensitivity string `json:"casesensitivity" zfs:"casesensitivity"`

	Vscan string `json:"vscan" zfs:"vscan"`

	Nbmand string `json:"nbmand" zfs:"nbmand"`

	Guid string `json:"guid" zfs:"guid"`

	Primarycache string `json:"primarycache" zfs:"primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache"`

	DeferDestroy string `json:"deferDestroy" zfs:"defer_destroy"`

	Userrefs string `json:"userrefs" zfs:"userrefs"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio"`

	Written string `json:"written" zfs:"written"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced"`

	Acltype string `json:"acltype" zfs:"acltype"`

	Context string `json:"context" zfs:"context"`

	Fscontext string `json:"fscontext" zfs:"fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext"`
	// 自定义属性
	Status string `json:"status" zfs:"lb:status"`
}

// +gogo:deepcopy-gen=true
type Pool struct {
	Name string `json:"name"`

	Types PoolTypes `json:"types"`

	Volumes map[string]*Volume `json:"volumes"`

	Snapshots map[string]*Snapshot `json:"snapshots"`

	Clones map[string]*Volume `json:"clones"`
}

// +gogo:deepcopy-gen=true
type Volume struct {
	Name string `json:"name"`
	// if this field is not empty, volume is a clone.
	// this field empty snapshot's name.
	Parent string `json:"parent"`

	Kind VolumeKind `json:"kind"`

	Types VolumeTypes `json:"types"`

	Snapshots map[string]*Snapshot `json:"snapshots"`
}

// +gogo:deepcopy-gen=true
type Snapshot struct {
	Name string `json:"name"`
	// volume name
	Source string `json:"source"`

	Types SnapshotTypes `json:"types"`

	Clones map[string]*Volume `json:"clones"`
}
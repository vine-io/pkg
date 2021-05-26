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

// Pool
// +gogo:deepcopy-gen=true
type Pool struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// zfs pool 容量
	Size string `json:"size" zfs:"size" protobuf:"bytes,2,opt,name=size"`

	Capacity string `json:"capacity" zfs:"capacity" protobuf:"bytes,3,opt,name=capacity"`

	AltRoot string `json:"altRoot" zfs:"altroot" protobuf:"bytes,4,opt,name=altRoot"`

	Health string `json:"health" zfs:"health" protobuf:"bytes,5,opt,name=health"`

	Guid string `json:"guid" zfs:"guid" protobuf:"bytes,6,opt,name=guid"`

	Version string `json:"version" zfs:"version" protobuf:"bytes,7,opt,name=version"`

	BootFs string `json:"bootFs" zfs:"bootfs" protobuf:"bytes,8,opt,name=bootFs"`

	Delegation string `json:"delegation" zfs:"delegation" protobuf:"bytes,9,opt,name=delegation"`

	AutoReplace string `json:"autoReplace" zfs:"authreplace" protobuf:"bytes,10,opt,name=autoReplace"`

	CacheFile string `json:"cacheFile" zfs:"cachefile" protobuf:"bytes,11,opt,name=cacheFile"`

	FailMode string `json:"failMode" zfs:"failmode" protobuf:"bytes,12,opt,name=failMode"`

	ListSnapshots string `json:"listSnapshots" zfs:"listsnapshots" protobuf:"bytes,13,opt,name=listSnapshots"`

	AutoExpand string `json:"autoExpand" zfs:"autoexpand" protobuf:"bytes,14,opt,name=autoExpand"`

	DedupDitto string `json:"dedupDitto" zfs:"dedupditto" protobuf:"bytes,15,opt,name=dedupDitto"`

	DedupRatio string `json:"dedupRatio" zfs:"dedupratio" protobuf:"bytes,16,opt,name=dedupRatio"`

	// zfs pool 可用空间
	Free string `json:"free" zfs:"free" protobuf:"bytes,17,opt,name=free"`

	// zfs pool 已分配空间
	Allocated string `json:"allocated" zfs:"allocated" protobuf:"bytes,18,opt,name=allocated"`

	Readonly string `json:"readonly" zfs:"readonly" protobuf:"bytes,19,opt,name=readonly"`

	AShift string `json:"aShift" zfs:"ashift" protobuf:"bytes,20,opt,name=aShift"`

	Comment string `json:"comment" zfs:"comment" protobuf:"bytes,21,opt,name=comment"`

	ExpandSize string `json:"expandSize" zfs:"expandsize" protobuf:"bytes,22,opt,name=expandSize"`

	Freeing string `json:"freeing" zfs:"freeing" protobuf:"bytes,23,opt,name=freeing"`

	Fragmentation string `json:"fragmentation" zfs:"fragmentation" protobuf:"bytes,24,opt,name=fragmentation"`

	Leaked string `json:"leaked" zfs:"leaked" protobuf:"bytes,25,opt,name=leaked"`

	Multihost string `json:"multihost" zfs:"multihost" protobuf:"bytes,26,opt,name=multihost"`

	FeatureAsyncDestroy string `json:"featureAsyncDestroy" zfs:"feature@async_destroy" protobuf:"bytes,27,opt,name=featureAsyncDestroy"`

	FeatureEmptyBpobj string `json:"featureEmptyBpobj" zfs:"feature@empty_bpobj" protobuf:"bytes,28,opt,name=featureEmptyBpobj"`

	FeatureLz4Compress string `json:"featureLz4Compress" zfs:"feature@lz4_compress" protobuf:"bytes,29,opt,name=featureLz4Compress"`

	FeatureMultiVdevCrashDump string `json:"featureMultiVdevCrashDump" zfs:"feature@multi_vdev_crash_dump" protobuf:"bytes,30,opt,name=featureMultiVdevCrashDump"`

	FeatureSpaceMapHistogram string `json:"featureSpaceMapHistogram" zfs:"feature@spacemap_histogram" protobuf:"bytes,31,opt,name=featureSpaceMapHistogram"`

	FeatureEnabledTxg string `json:"featureEnabledTxg" zfs:"feature@enabled_txg" protobuf:"bytes,32,opt,name=featureEnabledTxg"`

	FeatureHoleBirth string `json:"featureHoleBirth" zfs:"feature@hole_birth" protobuf:"bytes,33,opt,name=featureHoleBirth"`

	FeatureExtensibleDataset string `json:"featureExtensibleDataset" zfs:"feature@extensible_dataset" protobuf:"bytes,34,opt,name=featureExtensibleDataset"`

	FeatureEmbeddedData string `json:"featureEmbeddedData" zfs:"feature@embedded_data" protobuf:"bytes,35,opt,name=featureEmbeddedData"`

	FeatureBookmarks string `json:"featureBookmarks" zfs:"feature@bookmarks" protobuf:"bytes,36,opt,name=featureBookmarks"`

	FeatureFilesystemLimits string `json:"featureFilesystemLimits" zfs:"feature@filesystem_limits" protobuf:"bytes,37,opt,name=featureFilesystemLimits"`

	FeatureLargeBlocks string `json:"featureLargeBlocks" zfs:"feature@large_blocks" protobuf:"bytes,38,opt,name=featureLargeBlocks"`

	FeatureLargeDnode string `json:"featureLargeDnode" zfs:"feature@large_dnode" protobuf:"bytes,39,opt,name=featureLargeDnode"`

	FeatureSha512 string `json:"featureSha512" zfs:"feature@sha512" protobuf:"bytes,40,opt,name=featureSha512"`

	FeatureSkein string `json:"featureSkein" zfs:"feature@skein" protobuf:"bytes,41,opt,name=featureSkein"`

	FeatureEdonr string `json:"featureEdonr" zfs:"feature@edonr " protobuf:"bytes,42,opt,name=featureEdonr"`

	FeatureUserobjAccounting string `json:"featureUserobjAccounting" zfs:"feature@userobj_accounting" protobuf:"bytes,43,opt,name=featureUserobjAccounting"`

	// tank filesystem
	Type string `json:"type" zfs:"type" protobuf:"bytes,44,opt,name=type"`

	Creation string `json:"creation" zfs:"creation" protobuf:"bytes,45,opt,name=creation"`

	Used string `json:"used" zfs:"used" protobuf:"bytes,46,opt,name=used"`

	Available string `json:"available" zfs:"available" protobuf:"bytes,47,opt,name=available"`

	Referenced string `json:"referenced" zfs:"referenced" protobuf:"bytes,48,opt,name=referenced"`

	CompressRatio string `json:"compressRatio" zfs:"compressratio" protobuf:"bytes,49,opt,name=compressRatio"`

	Mounted string `json:"mounted" zfs:"mounted" protobuf:"bytes,50,opt,name=mounted"`

	Quota string `json:"quota" zfs:"quota" protobuf:"bytes,51,opt,name=quota"`

	Reservation string `json:"reservation" zfs:"reservation" protobuf:"bytes,52,opt,name=reservation"`

	RecordSize string `json:"recordSize" zfs:"recordsize" protobuf:"bytes,53,opt,name=recordSize"`

	MountPoint string `json:"mountPoint" zfs:"mountpoint" protobuf:"bytes,54,opt,name=mountPoint"`

	SharenFs string `json:"sharenFs" zfs:"sharenfs" protobuf:"bytes,55,opt,name=sharenFs"`

	Checksum string `json:"checksum" zfs:"checksum" protobuf:"bytes,56,opt,name=checksum"`

	Compression string `json:"compression" zfs:"compression" protobuf:"bytes,57,opt,name=compression"`

	ATime string `json:"aTime" zfs:"atime" protobuf:"bytes,58,opt,name=aTime"`

	Devices string `json:"devices" zfs:"devices" protobuf:"bytes,59,opt,name=devices"`

	Exec string `json:"exec" zfs:"exec" protobuf:"bytes,60,opt,name=exec"`

	SetUid string `json:"setUid" zfs:"setuid" protobuf:"bytes,61,opt,name=setUid"`

	Zoned string `json:"zoned" zfs:"zoned" protobuf:"bytes,62,opt,name=zoned"`

	SnapDir string `json:"snapDir" zfs:"snapdir" protobuf:"bytes,63,opt,name=snapDir"`

	AclinHerit string `json:"aclinHerit" zfs:"aclinherit" protobuf:"bytes,64,opt,name=aclinHerit"`

	Createtxg string `json:"createtxg" zfs:"createtxg" protobuf:"bytes,65,opt,name=createtxg"`

	CanMount string `json:"canMount" zfs:"canmount" protobuf:"bytes,66,opt,name=canMount"`

	Xattr string `json:"xattr" zfs:"xattr" protobuf:"bytes,67,opt,name=xattr"`

	Copies string `json:"copies" zfs:"copies" protobuf:"bytes,68,opt,name=copies"`

	Utf8only string `json:"utf8Only" zfs:"utf8only" protobuf:"bytes,69,opt,name=utf8Only"`

	Normalization string `json:"normalization" zfs:"normalization" protobuf:"bytes,70,opt,name=normalization"`

	Casesensitivity string `json:"casesensitivity" zfs:"casesensitivity" protobuf:"bytes,71,opt,name=casesensitivity"`

	Vscan string `json:"vscan" zfs:"vscan" protobuf:"bytes,72,opt,name=vscan"`

	Nbmand string `json:"nbmand" zfs:"nbmand" protobuf:"bytes,73,opt,name=nbmand"`

	Sharesmb string `json:"sharesmb" zfs:"sharesmb" protobuf:"bytes,74,opt,name=sharesmb"`

	Refquota string `json:"refquota" zfs:"refquota" protobuf:"bytes,75,opt,name=refquota"`

	Refreservation string `json:"refreservation" zfs:"refreservation" protobuf:"bytes,76,opt,name=refreservation"`

	Primarycache string `json:"primarycache" zfs:"primarycache" protobuf:"bytes,77,opt,name=primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache" protobuf:"bytes,78,opt,name=secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots" zfs:"usedbysnapshots" protobuf:"bytes,79,opt,name=usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset" zfs:"usedbydataset" protobuf:"bytes,80,opt,name=usedbydataset"`

	Usedbychildren string `json:"usedbychildren" zfs:"usedbychildren" protobuf:"bytes,81,opt,name=usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation" zfs:"usedbyrefreservation" protobuf:"bytes,82,opt,name=usedbyrefreservation"`

	Logbias string `json:"logbias" zfs:"logbias" protobuf:"bytes,83,opt,name=logbias"`

	Dedup string `json:"dedup" zfs:"dedup" protobuf:"bytes,84,opt,name=dedup"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel" protobuf:"bytes,85,opt,name=mlslabel"`

	Sync string `json:"sync" zfs:"sync" protobuf:"bytes,86,opt,name=sync"`

	Dnodesize string `json:"dnodesize" zfs:"dnodesize" protobuf:"bytes,87,opt,name=dnodesize"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio" protobuf:"bytes,88,opt,name=refcompressratio"`

	Written string `json:"written" zfs:"written" protobuf:"bytes,89,opt,name=written"`

	Logicalused string `json:"logicalused" zfs:"logicalused" protobuf:"bytes,90,opt,name=logicalused"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced" protobuf:"bytes,91,opt,name=logicalreferenced"`

	Volmode string `json:"volmode" zfs:"volmode" protobuf:"bytes,92,opt,name=volmode"`

	FilesystemLimit string `json:"filesystemLimit" zfs:"filesystem_limit" protobuf:"bytes,93,opt,name=filesystemLimit"`

	SnapshotLimit string `json:"snapshotLimit" zfs:"snapshot_limit" protobuf:"bytes,94,opt,name=snapshotLimit"`

	FilesystemCount string `json:"filesystemCount" zfs:"filesystem_count" protobuf:"bytes,95,opt,name=filesystemCount"`

	SnapshotCount string `json:"snapshotCount" zfs:"snapshot_count" protobuf:"bytes,96,opt,name=snapshotCount"`

	Snapdev string `json:"snapdev" zfs:"snapdev" protobuf:"bytes,97,opt,name=snapdev"`

	Acltype string `json:"acltype" zfs:"acltype" protobuf:"bytes,98,opt,name=acltype"`

	Context string `json:"context" zfs:"context" protobuf:"bytes,99,opt,name=context"`

	Fscontext string `json:"fscontext" zfs:"fscontext" protobuf:"bytes,100,opt,name=fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext" protobuf:"bytes,101,opt,name=defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext" protobuf:"bytes,102,opt,name=rootcontext"`

	Relatime string `json:"relatime" zfs:"relatime" protobuf:"bytes,103,opt,name=relatime"`

	RedundantMetadata string `json:"redundantMetadata" zfs:"redundant_metadata" protobuf:"bytes,104,opt,name=redundantMetadata"`

	Overlay string `json:"overlay" zfs:"overlay" protobuf:"bytes,105,opt,name=overlay"`

	Encryption string `json:"encryption" zfs:"encryption" protobuf:"bytes,106,opt,name=encryption"`

	Keylocation string `json:"keylocation" zfs:"keylocation" protobuf:"bytes,107,opt,name=keylocation"`

	Keyformat string `json:"keyformat" zfs:"keyformat" protobuf:"bytes,108,opt,name=keyformat"`

	Pbkdf2iters string `json:"pbkdf2Iters" zfs:"pbkdf2iters" protobuf:"bytes,109,opt,name=pbkdf2Iters"`

	SpecialSmallBlocks string `json:"specialSmallBlocks" zfs:"special_small_blocks" protobuf:"bytes,110,opt,name=specialSmallBlocks"`
}

// Volume
// +gogo:deepcopy-gen=true
type Volume struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	Source string `json:"source" protobuf:"bytes,2,opt,name=source"`

	Type string `json:"type" zfs:"type" protobuf:"bytes,3,opt,name=type"`

	Creation string `json:"creation" zfs:"creation" protobuf:"bytes,4,opt,name=creation"`

	Used string `json:"used" zfs:"used" protobuf:"bytes,5,opt,name=used"`

	Available string `json:"available" zfs:"available" protobuf:"bytes,6,opt,name=available"`

	Referenced string `json:"referenced" zfs:"referenced" protobuf:"bytes,7,opt,name=referenced"`

	Compressratio string `json:"compressratio" zfs:"compressratio" protobuf:"bytes,8,opt,name=compressratio"`

	Reservation string `json:"reservation" zfs:"reservation" protobuf:"bytes,9,opt,name=reservation"`

	Volsize string `json:"volsize" zfs:"volsize" protobuf:"bytes,10,opt,name=volsize"`

	Volblocksize string `json:"volblocksize" zfs:"volblocksize" protobuf:"bytes,11,opt,name=volblocksize"`

	Checksum string `json:"checksum" zfs:"checksum" protobuf:"bytes,12,opt,name=checksum"`

	Compression string `json:"compression" zfs:"compression" protobuf:"bytes,13,opt,name=compression"`

	Readonly string `json:"readonly" zfs:"readonly" protobuf:"bytes,14,opt,name=readonly"`

	Copies string `json:"copies" zfs:"copies" protobuf:"bytes,15,opt,name=copies"`

	Refreservation string `json:"refreservation" zfs:"refreservation" protobuf:"bytes,16,opt,name=refreservation"`

	Guid string `json:"guid" zfs:"guid" protobuf:"bytes,17,opt,name=guid"`

	Primarycache string `json:"primarycache" zfs:"primarycache" protobuf:"bytes,18,opt,name=primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache" protobuf:"bytes,19,opt,name=secondarycache"`

	Usedbysnapshots string `json:"usedbysnapshots" zfs:"usedbysnapshots" protobuf:"bytes,20,opt,name=usedbysnapshots"`

	Usedbydataset string `json:"usedbydataset" zfs:"usedbydataset" protobuf:"bytes,21,opt,name=usedbydataset"`

	Usedbychildren string `json:"usedbychildren" zfs:"usedbychildren" protobuf:"bytes,22,opt,name=usedbychildren"`

	Usedbyrefreservation string `json:"usedbyrefreservation" zfs:"usedbyrefreservation" protobuf:"bytes,23,opt,name=usedbyrefreservation"`

	Logbias string `json:"logbias" zfs:"logbias" protobuf:"bytes,24,opt,name=logbias"`

	Dedup string `json:"dedup" zfs:"dedup" protobuf:"bytes,25,opt,name=dedup"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel" protobuf:"bytes,26,opt,name=mlslabel"`

	Sync string `json:"sync" zfs:"sync" protobuf:"bytes,27,opt,name=sync"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio" protobuf:"bytes,28,opt,name=refcompressratio"`

	Written string `json:"written" zfs:"written" protobuf:"bytes,29,opt,name=written"`

	Logicalused string `json:"logicalused" zfs:"logicalused" protobuf:"bytes,30,opt,name=logicalused"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced" protobuf:"bytes,31,opt,name=logicalreferenced"`

	Volmode string `json:"volmode" zfs:"volmode" protobuf:"bytes,32,opt,name=volmode"`

	FilesystemLimit string `json:"filesystemLimit" zfs:"filesystem_limit" protobuf:"bytes,33,opt,name=filesystemLimit"`

	SnapshotLimit string `json:"snapshotLimit" zfs:"snapshot_limit" protobuf:"bytes,34,opt,name=snapshotLimit"`

	SnapshotCount string `json:"snapshotCount" zfs:"snapshot_count" protobuf:"bytes,35,opt,name=snapshotCount"`

	Snapdev string `json:"snapdev" zfs:"snapdev" protobuf:"bytes,36,opt,name=snapdev"`

	Acltype string `json:"acltype" zfs:"acltype" protobuf:"bytes,37,opt,name=acltype"`

	Context string `json:"context" zfs:"context" protobuf:"bytes,38,opt,name=context"`

	Fscontext string `json:"fscontext" zfs:"fscontext" protobuf:"bytes,39,opt,name=fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext" protobuf:"bytes,40,opt,name=defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext" protobuf:"bytes,41,opt,name=rootcontext"`

	RedundantMetadata string `json:"redundantMetadata" zfs:"redundant_metadata" protobuf:"bytes,42,opt,name=redundantMetadata"`
}

// Snapshot
// +gogo:deepcopy-gen=true
type Snapshot struct {
	Parent string `json:"parent" protobuf:"bytes,1,opt,name=parent"`

	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`

	Type string `json:"type" zfs:"type" protobuf:"bytes,3,opt,name=type"`

	Creation string `json:"creation" zfs:"creation" protobuf:"bytes,4,opt,name=creation"`

	Used string `json:"used" zfs:"used" protobuf:"bytes,5,opt,name=used"`

	Referenced string `json:"referenced" zfs:"referenced" protobuf:"bytes,6,opt,name=referenced"`

	Compressratio string `json:"compressratio" zfs:"compressratio" protobuf:"bytes,7,opt,name=compressratio"`

	Devices string `json:"devices" zfs:"devices" protobuf:"bytes,8,opt,name=devices"`

	Exec string `json:"exec" zfs:"exec" protobuf:"bytes,9,opt,name=exec"`

	Setuid string `json:"setuid" zfs:"setuid" protobuf:"bytes,10,opt,name=setuid"`

	Createtxg string `json:"createtxg" zfs:"createtxg" protobuf:"bytes,11,opt,name=createtxg"`

	Version string `json:"version" zfs:"version" protobuf:"bytes,12,opt,name=version"`

	Utf8only string `json:"utf8Only" zfs:"utf8only" protobuf:"bytes,13,opt,name=utf8Only"`

	Normalization string `json:"normalization" zfs:"normalization" protobuf:"bytes,14,opt,name=normalization"`

	Casesensitivity string `json:"casesensitivity" zfs:"casesensitivity" protobuf:"bytes,15,opt,name=casesensitivity"`

	Vscan string `json:"vscan" zfs:"vscan" protobuf:"bytes,16,opt,name=vscan"`

	Nbmand string `json:"nbmand" zfs:"nbmand" protobuf:"bytes,17,opt,name=nbmand"`

	Guid string `json:"guid" zfs:"guid" protobuf:"bytes,18,opt,name=guid"`

	Primarycache string `json:"primarycache" zfs:"primarycache" protobuf:"bytes,19,opt,name=primarycache"`

	Secondarycache string `json:"secondarycache" zfs:"secondarycache" protobuf:"bytes,20,opt,name=secondarycache"`

	DeferDestroy string `json:"deferDestroy" zfs:"defer_destroy" protobuf:"bytes,21,opt,name=deferDestroy"`

	Userrefs string `json:"userrefs" zfs:"userrefs" protobuf:"bytes,22,opt,name=userrefs"`

	Mlslabel string `json:"mlslabel" zfs:"mlslabel" protobuf:"bytes,23,opt,name=mlslabel"`

	Refcompressratio string `json:"refcompressratio" zfs:"refcompressratio" protobuf:"bytes,24,opt,name=refcompressratio"`

	Written string `json:"written" zfs:"written" protobuf:"bytes,25,opt,name=written"`

	Logicalreferenced string `json:"logicalreferenced" zfs:"logicalreferenced" protobuf:"bytes,26,opt,name=logicalreferenced"`

	Acltype string `json:"acltype" zfs:"acltype" protobuf:"bytes,27,opt,name=acltype"`

	Context string `json:"context" zfs:"context" protobuf:"bytes,28,opt,name=context"`

	Fscontext string `json:"fscontext" zfs:"fscontext" protobuf:"bytes,29,opt,name=fscontext"`

	Defcontext string `json:"defcontext" zfs:"defcontext" protobuf:"bytes,30,opt,name=defcontext"`

	Rootcontext string `json:"rootcontext" zfs:"rootcontext" protobuf:"bytes,31,opt,name=rootcontext"`
}

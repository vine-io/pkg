package zfs
//
//import (
//	"bufio"
//	"bytes"
//	"fmt"
//	"io"
//	"os/exec"
//	"path"
//	"reflect"
//	"runtime"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//
//	utilzfs "lb.io/util/zfs"
//)
//
//var once sync.Once
//
//var adm = &ZFSadm{
//	mux:    sync.RWMutex{},
//	System: &System{},
//}
//
//type ZFSadm struct {
//	mux sync.RWMutex
//
//	zpool string
//	zfs   string
//
//	System *System
//}
//
//func Default() (*ZFSadm, error) {
//	if runtime.GOOS != "linux" {
//		return nil, fmt.Errorf("zfs must work in linux")
//	}
//	adm.mux.Lock()
//	adm.lazy()
//	adm.mux.Unlock()
//	if len(adm.zfs) == 0 {
//		return nil, fmt.Errorf("could't find command 'zfs', zfs is not installed")
//	}
//	if len(adm.zpool) == 0 {
//		return nil, fmt.Errorf("could't find command 'zpool', zfs is not installed")
//	}
//	go adm.GarbageCollecting()
//	return adm, nil
//}
//
//func (z *ZFSadm) lazy() {
//	once.Do(func() {
//		data, err := exec.Command("sh", "-c", `lsmod | grep zfs`).CombinedOutput()
//		if err != nil || string(data) == "" {
//			return
//		}
//
//		// 内核加载 zfs 模块
//		_, _ = exec.Command("modprobe", "zfs").CombinedOutput()
//
//		data, err = exec.Command("which", "zpool").CombinedOutput()
//		if err != nil || string(data) == "" {
//			return
//		}
//		adm.zpool = strings.TrimSuffix(string(data), "\n")
//		data, err = exec.Command("which", "zfs").CombinedOutput()
//		if err != nil || string(data) == "" {
//			return
//		}
//		adm.zfs = strings.TrimSuffix(string(data), "\n")
//
//		adm.System.Pools = map[string]*Pool{}
//		adm.System.FileSystems = map[string]*Filesystem{}
//		adm.System.Volumes = map[string]*Volume{}
//		adm.System.Snapshots = map[string]*Snapshot{}
//
//		// 获取所有的 pool
//		// zpool list -Hp -o name
//		shell := utilzfs.ZPool(adm.zpool).
//			List("", "-Hp", []string{"name"}, "")
//		klog.V(4).Infoln(shell.Commit())
//		out, _ := shell.Exec()
//		if len(out) == 0 {
//			return
//		}
//		names := strings.Split(string(out), "\n")
//		for _, name := range names {
//			adm.System.Pools[name] = &Pool{Name: name, Devs: map[string]string{}}
//		}
//
//		// 解析所有的pool
//		for name, pool := range adm.System.Pools {
//			shell = utilzfs.ZPool(adm.zpool).
//				Get(name, "-Hp", nil, "all")
//			klog.V(4).Infoln(shell.Commit())
//			out, _ = shell.Exec()
//			setValue(out, pool)
//
//			shell = utilzfs.Zfs(adm.zfs).
//				Get(name, "-Hp", "", nil, "", "", "all")
//			klog.V(4).Infoln(shell.Commit())
//			out, _ = shell.Exec()
//			setValue(out, pool)
//		}
//
//		// 获取所有的 filesystem
//		shell = utilzfs.Zfs(adm.zfs).
//			List("", "-Hp", "", []string{"name", DISPLAY, "origin"}, "", "", "filesystem")
//		klog.V(4).Infoln(shell.Commit())
//		out, _ = shell.Exec()
//
//		lines := strings.Split(string(out), "\n")
//		for _, line := range lines {
//			parts := strings.Split(line, "\t")
//			name := strings.TrimSpace(parts[0])
//			if !strings.Contains(name, "/") {
//				continue
//			}
//			if v := strings.TrimSpace(parts[1]); v == Hidden {
//				continue
//			}
//			pool := strings.Split(name, "/")[0]
//			var source string
//			if v := strings.TrimSpace(parts[2]); v != "-" {
//				source = v
//			}
//			// clones 为 - 表示为文件系统
//			adm.System.Pools[pool].Devs[name] = name
//			adm.System.FileSystems[name] = &Filesystem{Pool: pool, Name: name, Source: source, Snapshots: map[string]string{}}
//		}
//
//		// 解析所有的filesystem
//		for name, filesystem := range adm.System.FileSystems {
//			shell = utilzfs.Zfs(adm.zfs).
//				Get(name, "-Hp", "", nil, "", "", "all")
//			klog.V(4).Infoln(shell.Commit())
//			out, _ = shell.Exec()
//			setValue(out, filesystem)
//		}
//
//		// 获取所有的 volume
//		shell = utilzfs.Zfs(adm.zfs).
//			List("", "-Hp", "", []string{"name", DISPLAY, "origin"}, "", "", "volume")
//		klog.V(4).Infoln(shell.Commit())
//		out, _ = shell.Exec()
//
//		lines = strings.Split(string(out), "\n")
//		for _, line := range lines {
//			parts := strings.Split(line, "\t")
//			name := strings.TrimSpace(parts[0])
//			if !strings.Contains(name, "/") {
//				continue
//			}
//			if v := strings.TrimSpace(parts[1]); v == Hidden {
//				continue
//			}
//			pool := strings.Split(name, "/")[0]
//			var source string
//			if v := strings.TrimSpace(parts[2]); v != "-" {
//				source = v
//			}
//			adm.System.Pools[pool].Devs[name] = name
//			adm.System.Volumes[name] = &Volume{Pool: pool, Name: name, Source: source, Snapshots: map[string]string{}}
//		}
//
//		// 解析所有的volume
//		for name, volume := range adm.System.Volumes {
//			shell = utilzfs.Zfs(adm.zfs).
//				Get(name, "-Hp", "", nil, "", "", "all")
//			klog.V(4).Infoln(shell.Commit())
//			out, _ = shell.Exec()
//			setValue(out, volume)
//		}
//
//		// 获取所有的 snapshot
//		shell = utilzfs.Zfs(adm.zfs).
//			List("", "-Hp", "", []string{"name", DISPLAY}, "", "", "snapshot")
//		klog.V(4).Infoln(shell.Commit())
//		out, _ = shell.Exec()
//
//		lines = strings.Split(string(out), "\n")
//		for _, line := range lines {
//			parts := strings.Split(line, "\t")
//			name := strings.TrimSpace(parts[0])
//			if !strings.Contains(name, "/") {
//				continue
//			}
//			if v := strings.TrimSpace(parts[1]); v == Hidden {
//				continue
//			}
//			pool := strings.Split(name, "/")[0]
//			adm.System.Pools[pool].Devs[name] = name
//			parent := strings.Split(name, "@")[0]
//			if fileSystem, ok := adm.System.FileSystems[parent]; ok {
//				fileSystem.Snapshots[name] = name
//			}
//			if volume, ok := adm.System.Volumes[parent]; ok {
//				volume.Snapshots[name] = name
//			}
//			adm.System.Snapshots[name] = &Snapshot{Pool: pool, Name: name, Parent: parent, Clones: map[string]string{}}
//		}
//
//		// 解析所有的snapshot
//		for name, snapshot := range adm.System.Snapshots {
//			shell = utilzfs.Zfs(adm.zfs).
//				Get(name, "-Hp", "", nil, "", "", "all")
//			klog.V(4).Infoln(shell.Commit())
//			out, _ = shell.Exec()
//			setSnapshot(out, snapshot)
//		}
//	})
//}
//
//func (z *ZFSadm) GarbageCollecting() {
//	timer := time.NewTicker(time.Second * 15)
//	for range timer.C {
//		pools, _ := z.GetPools()
//		for _, pool := range pools {
//			flag := false
//			z.mux.RLock()
//			if pool.Status == Recycling && len(pool.Devs) == 0 {
//				flag = true
//			}
//			z.mux.RUnlock()
//			if flag {
//				klog.V(4).Infof("detected the pool to be recovered and started to destroy it")
//				z.destroyPool(pool.Name)
//			}
//		}
//		fileSystems, _ := z.GetFileSystems()
//		for _, fileSystem := range fileSystems {
//			flag := false
//			z.mux.RLock()
//			if fileSystem.Status == Recycling && len(fileSystem.Snapshots) == 0 {
//				flag = true
//			}
//			z.mux.RUnlock()
//			if flag {
//				klog.V(4).Infof("detected the fileSystem to be recovered and started to destroy it")
//				if len(fileSystem.Shadow) != 0 {
//					execute := utilzfs.Zfs(z.zfs).DestroySnapshot(fileSystem.Shadow, "-dr")
//					klog.V(4).Infoln(execute.Commit())
//					_, _ = execute.Exec()
//				}
//				z.destroyFileSystem(fileSystem)
//			}
//		}
//		volumes, _ := z.GetVolumes()
//		for _, volume := range volumes {
//			flag := false
//			z.mux.RLock()
//			if volume.Status == Recycling && len(volume.Snapshots) == 0 {
//				flag = true
//			}
//			z.mux.RUnlock()
//			if flag {
//				klog.V(4).Infof("detected the volume to be recovered and started to destroy it")
//				if len(volume.Shadow) != 0 {
//					execute := utilzfs.Zfs(z.zfs).DestroySnapshot(volume.Shadow, "-dr")
//					klog.V(4).Infoln(execute.Commit())
//					_, _ = execute.Exec()
//				}
//				z.destroyVolume(volume)
//			}
//		}
//		snapshots, _ := z.GetSnapshots()
//		for _, snapshot := range snapshots {
//			flag := false
//			z.mux.RLock()
//			if snapshot.Status == Recycling && len(snapshot.Clones) == 0 {
//				flag = true
//			}
//			z.mux.RUnlock()
//			if flag {
//				klog.V(4).Infof("detected the snapshot to be recovered and started to destroy it")
//				z.destroySnapshot(snapshot)
//				if len(snapshot.RelationShip) != 0 {
//					execute := utilzfs.Zfs(z.zfs).DestroyFileSystemOrVolume(snapshot.RelationShip, "")
//					klog.V(4).Infoln(execute.Commit())
//					_, _ = execute.Exec()
//				}
//			}
//		}
//	}
//}
//
//// return zfs pool:
//func (z *ZFSadm) QueryPoolIO(name string) (string, string, string, string, error) {
//	var (
//		readIO  = "0"
//		writeIO = "0"
//		read    = "0"
//		write   = "0"
//	)
//
//	execute := utilzfs.ZPool(z.zpool).Iostat(name, "", "", "-Hp")
//	data, _ := execute.Exec()
//	line := strings.Split(strings.TrimSuffix(string(data), "\n"), "\t")
//	if len(line) > 6 {
//		readIO = line[3]
//		writeIO = line[4]
//		read = line[5]
//		write = line[6]
//	}
//
//	return readIO, writeIO, read, write, nil
//}
//
//func (z *ZFSadm) CreatePool(name, compression, raid string, quota float64, devices []string) (*Pool, error) {
//
//	if _, err := z.GetPool(name); err == nil {
//		return nil, fmt.Errorf("pool '%s' is already exists", name)
//	}
//
//	pool := &Pool{Name: name, Devs: map[string]string{}}
//	execute := utilzfs.ZPool(z.zpool).Create(name, "", "", "", nil, raid, devices...)
//	klog.V(4).Infoln(execute.Commit())
//	data, err := execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("create pool '%s': %v", name, err)
//	}
//
//	execute = utilzfs.ZPool(z.zpool).Get(name, "-Hp", nil, "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ = execute.Exec()
//	setValue(data, pool)
//
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ = execute.Exec()
//	setValue(data, pool)
//
//	poolSize, _ := strconv.ParseFloat(pool.Size, 64)
//	_quota := int64(poolSize * quota / 100)
//	properties := map[string]string{
//		"compression": compression,
//		"quota":       strconv.FormatInt(_quota, 10),
//		STATUS:        Ready,
//	}
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	data, err = execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("create pool '%s': %v", name, err)
//	}
//	z.mux.Lock()
//	z.System.Pools[name] = pool
//	z.mux.Unlock()
//	return pool, nil
//}
//
//func (z *ZFSadm) ExpensePool(name string, devices ...string) (*Pool, error) {
//
//	pool, err := z.GetPool(name)
//	if err != nil {
//		return nil, err
//	}
//
//	execute := utilzfs.ZPool(z.zpool).Add(name, "", devices...)
//	klog.V(4).Infoln(execute.Commit())
//	data, err := execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("expense pool '%s': %v", name, err)
//	}
//
//	execute = utilzfs.ZPool(z.zpool).Get(name, "-Hp", nil, "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ = execute.Exec()
//
//	z.mux.Lock()
//	setValue(data, pool)
//	z.System.Pools[name] = pool
//	z.mux.Unlock()
//
//	return pool, nil
//}
//
//func (z *ZFSadm) GetPools() ([]*Pool, error) {
//	z.mux.RLock()
//	defer z.mux.RUnlock()
//
//	pools := make([]*Pool, 0)
//	for _, v := range z.System.Pools {
//		pools = append(pools, v)
//	}
//	return pools, nil
//}
//
//func (z *ZFSadm) GetPool(name string) (*Pool, error) {
//	z.mux.Lock()
//	defer z.mux.Unlock()
//
//	pool, ok := z.System.Pools[name]
//	if !ok {
//		return nil, fmt.Errorf("pool '%s' is not exists", name)
//	}
//	execute := utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	data, _ := execute.Exec()
//	setValue(data, pool)
//	return pool, nil
//}
//
//func (z *ZFSadm) DeletePool(name string) (*Pool, error) {
//	pool, err := z.GetPool(name)
//	if err != nil {
//		return nil, fmt.Errorf("pool '%s' is not exists: %v", name, err)
//	}
//
//	if len(pool.Devs) != 0 {
//		return nil, fmt.Errorf("pool '%s' has accoicated volume or fileSystem", name)
//	}
//
//	// 不直接删除，修改 lb:status 属性，等待后端回收
//	properties := map[string]string{STATUS: Recycling}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, err = execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("delete pool '%s': %v", name, err)
//	}
//
//	z.mux.Lock()
//	pool.Status = Recycling
//	z.mux.Unlock()
//
//	return pool, nil
//}
//
//func (z *ZFSadm) destroyPool(name string) {
//	execute := utilzfs.ZPool(z.zpool).Destroy(name, false)
//	klog.V(4).Infoln(execute.Commit())
//	_, err := execute.Exec()
//	if err != nil {
//		klog.Errorf("destroy pool '%s': %v", name, err)
//		return
//	}
//	z.mux.Lock()
//	delete(z.System.Pools, name)
//	z.mux.Unlock()
//}
//
//func (z *ZFSadm) GetFileSystems() ([]*Filesystem, error) {
//	z.mux.RLock()
//	defer z.mux.RUnlock()
//
//	filesystems := make([]*Filesystem, 0)
//	for _, v := range z.System.FileSystems {
//		filesystems = append(filesystems, v)
//	}
//	return filesystems, nil
//}
//
//func (z *ZFSadm) GetFileSystem(name string) (*Filesystem, error) {
//	z.mux.Lock()
//	defer z.mux.Unlock()
//
//	v, ok := z.System.FileSystems[name]
//	if !ok {
//		return nil, fmt.Errorf("filesystem '%s' is not exists", name)
//	}
//	execute := utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	data, _ := execute.Exec()
//	setValue(data, v)
//	return v, nil
//}
//
//func (z *ZFSadm) CreateFileSystem(name string, properties map[string]string) (*Filesystem, error) {
//
//	pool := strings.SplitN(name, "/", 2)[0]
//	if _, err := z.GetPool(pool); err != nil {
//		return nil, err
//	}
//	if _, err := z.GetFileSystem(name); err == nil {
//		return nil, fmt.Errorf("fileSystem '%s' is already exists", name)
//	}
//
//	execute := utilzfs.Zfs(z.zfs).CreateFileSystem(name, nil)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create filesystem '%s': %v", name, err)
//	}
//
//	snapshot := name + "@" + "shadow"
//	execute = utilzfs.Zfs(z.zfs).Snapshot(snapshot, map[string]string{DISPLAY: Hidden})
//	klog.V(4).Infoln(execute.Commit())
//	_, _ = execute.Exec()
//
//	if properties == nil {
//		properties = map[string]string{}
//	}
//	properties[STATUS] = Ready
//	properties[SHADOW] = snapshot
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create filesystem '%s': %v", name, err)
//	}
//
//	filesystem := &Filesystem{Pool: pool, Name: name, Snapshots: map[string]string{}}
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//
//	z.mux.Lock()
//	setValue(data, filesystem)
//	z.System.FileSystems[name] = filesystem
//	z.System.Pools[pool].Devs[name] = name
//	z.mux.Unlock()
//	return filesystem, nil
//}
//
//func (z *ZFSadm) ShareFileSystem(name string, ips []string, mode string) error {
//	// zfs set sharenfs='rw=@192.168.2.0/24,rw=@192.168.221.0/24,all_squash,insecure' tank/test
//	if _, err := z.GetFileSystem(name); err != nil {
//		return err
//	}
//
//	properties := map[string]string{"sync": "disabled"}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, _ = execute.Exec()
//
//	args := ""
//	for _, ip := range ips {
//		args += fmt.Sprintf(`%s=@%s,`, mode, ip)
//	}
//	args += "no_root_squash,insecure"
//	properties = map[string]string{"sharenfs": fmt.Sprintf(`'%s'`, args)}
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return fmt.Errorf("failed to share filesystem '%s' => %v: %v", name, ips, err)
//	}
//
//	return nil
//}
//
//func (z *ZFSadm) UnShareFileSystem(name string) error {
//	// zfs set sharenfs=off tank/test
//	if _, err := z.GetFileSystem(name); err != nil {
//		return err
//	}
//
//	properties := map[string]string{"sharenfs": "off"}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return fmt.Errorf("failed to unshare filesystem '%s': %v", name, err)
//	}
//	return nil
//}
//
//// zfs 中删除 filesystem 时，如果 filesystem 存在 snapshot，强制删除时
//// 会删除 snapshot 及其关联的 clone
//func (z *ZFSadm) DeleteFileSystem(name string, force bool) (*Filesystem, error) {
//
//	fileSystem, err := z.GetFileSystem(name)
//	if err != nil {
//		return nil, err
//	}
//
//	// 不直接删除，修改 lb:status 属性，等待后端回收
//	properties := map[string]string{STATUS: Recycling}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, err = execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("destroy filesystem '%s': %v", name, err)
//	}
//
//	z.mux.Lock()
//	fileSystem.Status = Recycling
//	z.mux.Unlock()
//
//	return fileSystem, nil
//}
//
//func (z *ZFSadm) destroyFileSystem(fs *Filesystem) {
//	options := "-rf"
//	execute := utilzfs.Zfs(z.zfs).DestroyFileSystemOrVolume(fs.Name, options)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		klog.Errorf("failed to destroy filesystem '%s': %v", fs.Name, err)
//		return
//	}
//	z.mux.Lock()
//	delete(z.System.FileSystems, fs.Name)
//	delete(z.System.Pools[fs.Pool].Devs, fs.Name)
//	if len(fs.Source) != 0 {
//		delete(z.System.Snapshots[fs.Source].Clones, fs.Name)
//	}
//	z.mux.Unlock()
//}
//
//func (z *ZFSadm) GetVolumes() ([]*Volume, error) {
//	z.mux.RLock()
//	defer z.mux.RUnlock()
//
//	volumes := make([]*Volume, 0)
//	for _, v := range z.System.Volumes {
//		volumes = append(volumes, v)
//	}
//	return volumes, nil
//}
//
//func (z *ZFSadm) GetVolume(name string) (*Volume, error) {
//	z.mux.Lock()
//	defer z.mux.Unlock()
//
//	volume, ok := z.System.Volumes[name]
//	if !ok {
//		return nil, fmt.Errorf("volume '%s' is not exists", name)
//	}
//	execute := utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	data, _ := execute.Exec()
//	setValue(data, volume)
//	return volume, nil
//}
//
//func (z *ZFSadm) CreateVolume(name string, properties map[string]string, size int64) (*Volume, error) {
//
//	pool := strings.SplitN(name, "/", 2)[0]
//	if _, err := z.GetPool(pool); err != nil {
//		return nil, err
//	}
//	if _, err := z.GetVolume(name); err == nil {
//		return nil, fmt.Errorf("volume '%s' is alreay exists", name)
//	}
//
//	execute := utilzfs.Zfs(z.zfs).CreateVolume(name, 4096, nil, strconv.FormatInt(size, 10))
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create volume '%s': %v", name, err)
//	}
//
//	snapshot := name + "@" + "shadow"
//	execute = utilzfs.Zfs(z.zfs).Snapshot(snapshot, map[string]string{DISPLAY: Hidden})
//	klog.V(4).Infoln(execute.Commit())
//	_, _ = execute.Exec()
//
//	if properties == nil {
//		properties = map[string]string{}
//	}
//	properties[STATUS] = Ready
//	properties[SHADOW] = snapshot
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create volume '%s': %v", name, err)
//	}
//
//	volume := &Volume{Pool: pool, Name: name, Snapshots: map[string]string{}}
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//	z.mux.Lock()
//	setValue(data, volume)
//	z.System.Volumes[name] = volume
//	z.System.Pools[pool].Devs[name] = name
//	z.mux.Unlock()
//	return volume, nil
//}
//
//func (z *ZFSadm) DeleteVolume(name string, force bool) (*Volume, error) {
//	volume, err := z.GetVolume(name)
//	if err != nil {
//		return nil, fmt.Errorf("volume '%s' is not exists", name)
//	}
//
//	// 不直接删除，修改 lb:status 属性，等待后端回收
//	properties := map[string]string{STATUS: Recycling}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, err = execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("destroy volume '%s': %v", name, err)
//	}
//
//	z.mux.Lock()
//	volume.Status = Recycling
//	z.mux.Unlock()
//
//	return volume, nil
//}
//
//func (z *ZFSadm) destroyVolume(vol *Volume) {
//	options := "-rf"
//	execute := utilzfs.Zfs(z.zfs).DestroyFileSystemOrVolume(vol.Name, options)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		klog.Errorf("failed to destroy volume '%s': %v", vol.Name, err)
//		return
//	}
//	z.mux.Lock()
//	delete(z.System.Volumes, vol.Name)
//	delete(z.System.Pools[vol.Pool].Devs, vol.Name)
//	if len(vol.Source) != 0 {
//		delete(z.System.Snapshots[vol.Source].Clones, vol.Name)
//	}
//	z.mux.Unlock()
//}
//
//func (z *ZFSadm) GetSnapshots() ([]*Snapshot, error) {
//	z.mux.RLock()
//	defer z.mux.RUnlock()
//
//	snapshots := make([]*Snapshot, 0)
//	for _, snap := range z.System.Snapshots {
//		snapshots = append(snapshots, snap)
//	}
//	return snapshots, nil
//}
//
//func (z *ZFSadm) GetSnapshot(name string) (*Snapshot, error) {
//	z.mux.RLock()
//	defer z.mux.RUnlock()
//
//	for k, v := range z.System.Snapshots {
//		if k == name {
//			return v, nil
//		}
//	}
//	return nil, fmt.Errorf("snapshot '%s' is not exsits", name)
//}
//
//func (z *ZFSadm) CreateSnapshot(parent, name string, properties map[string]string) (*Snapshot, error) {
//
//	if !strings.Contains(name, "@") {
//		return nil, fmt.Errorf("missing '@' in snapshot name")
//	}
//
//	var fileSystem *Filesystem
//	var volume *Volume
//	var err error
//	pool := strings.SplitN(parent, "/", 2)[0]
//	if fileSystem, err = z.GetFileSystem(parent); err != nil {
//		if volume, err = z.GetVolume(parent); err != nil {
//			return nil, fmt.Errorf("parent '%s' is not exists", parent)
//		}
//	}
//
//	if _, err := z.GetSnapshot(name); err == nil {
//		return nil, fmt.Errorf("snapshot '%s' is already exists", name)
//	}
//
//	execute := utilzfs.Zfs(z.zfs).Snapshot(name, nil)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create snapshot '%s' => '%s': %v", parent, name, err)
//	}
//
//	if properties == nil {
//		properties = map[string]string{}
//	}
//	properties[STATUS] = Ready
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create snapshot '%s': %v", name, err)
//	}
//
//	snapshot := &Snapshot{Pool: pool, Parent: strings.SplitN(parent, "@", 2)[0], Name: name, Clones: map[string]string{}}
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//	z.mux.Lock()
//	setSnapshot(data, snapshot)
//	z.System.Snapshots[name] = snapshot
//	z.System.Pools[pool].Devs[name] = name
//
//	if fileSystem != nil {
//		fileSystem.Snapshots[name] = name
//	} else if volume != nil {
//		volume.Snapshots[name] = name
//	}
//	z.mux.Unlock()
//	return snapshot, nil
//}
//
//func (z *ZFSadm) DeleteSnapshot(name string, force bool) (*Snapshot, error) {
//
//	snapshot, err := z.GetSnapshot(name)
//	if err != nil {
//		return nil, err
//	}
//
//	// 不直接删除，修改 lb:status 属性，等待后端回收
//	properties := map[string]string{STATUS: Recycling}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, err = execute.Exec()
//	if err != nil {
//		return nil, fmt.Errorf("destroy snapshot '%s': %v", name, err)
//	}
//
//	z.mux.Lock()
//	snapshot.Status = Recycling
//	z.mux.Unlock()
//
//	return snapshot, nil
//}
//
//func (z *ZFSadm) destroySnapshot(s *Snapshot) {
//	options := "-dr"
//	execute := utilzfs.Zfs(z.zfs).DestroySnapshot(s.Name, options)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		klog.Errorf("failed to destroy snapshot '%s': %v", s.Name, err)
//		return
//	}
//	z.mux.Lock()
//	delete(z.System.Snapshots, s.Name)
//	delete(z.System.Pools[s.Pool].Devs, s.Name)
//	if fs, ok := z.System.FileSystems[s.Parent]; ok {
//		delete(fs.Snapshots, s.Name)
//	}
//	if vol, ok := z.System.Volumes[s.Parent]; ok {
//		delete(vol.Snapshots, s.Name)
//	}
//	z.mux.Unlock()
//}
//
//func (z *ZFSadm) CloneFileSystem(name, snap string, properties map[string]string) (*Filesystem, error) {
//
//	if _, err := z.GetFileSystem(name); err == nil {
//		return nil, fmt.Errorf("fileSystem '%s' is already exists", name)
//	}
//
//	if _, err := z.GetSnapshot(snap); err != nil {
//		return nil, err
//	}
//
//	pool := strings.SplitN(snap, "/", 2)[0]
//	name = path.Join(pool, name)
//	execute := utilzfs.Zfs(z.zfs).Clone(name, nil, snap)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create clone '%s' => '%s': %v", snap, name, err)
//	}
//
//	if properties == nil {
//		properties = map[string]string{}
//	}
//	properties[STATUS] = Ready
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create clone '%s' => '%s': %v", snap, name, err)
//	}
//
//	fileSystem := &Filesystem{Pool: pool, Name: name, Snapshots: map[string]string{}}
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//
//	z.mux.Lock()
//	setValue(data, fileSystem)
//	z.System.FileSystems[name] = fileSystem
//	z.System.Snapshots[snap].Clones[name] = name
//	z.System.Pools[pool].Devs[name] = name
//	z.mux.Unlock()
//	return fileSystem, nil
//}
//
//func (z *ZFSadm) CloneVolume(name, snap string, properties map[string]string) (*Volume, error) {
//	if _, err := z.GetVolume(name); err == nil {
//		return nil, fmt.Errorf("volume '%s' is already exists", name)
//	}
//
//	if _, err := z.GetSnapshot(snap); err != nil {
//		return nil, err
//	}
//
//	pool := strings.SplitN(snap, "/", 2)[0]
//	name = path.Join(pool, name)
//	execute := utilzfs.Zfs(z.zfs).Clone(name, nil, snap)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create clone '%s' => '%s': %v", snap, name, err)
//	}
//
//	if properties == nil {
//		properties = map[string]string{}
//	}
//	properties[STATUS] = Ready
//	execute = utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	if _, err := execute.Exec(); err != nil {
//		return nil, fmt.Errorf("failed to create clone '%s' => '%s': %v", snap, name, err)
//	}
//
//	volume := &Volume{Pool: pool, Name: name, Snapshots: map[string]string{}}
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//
//	z.mux.Lock()
//	setValue(data, volume)
//	z.System.Volumes[name] = volume
//	z.System.Snapshots[snap].Clones[name] = name
//	z.System.Pools[pool].Devs[name] = name
//	z.mux.Unlock()
//
//	return volume, nil
//}
//
//func (z *ZFSadm) SendFullSnapshot(srcSnapshot, pool, user, host string) (string, error) {
//	_, err := z.GetSnapshot(srcSnapshot)
//	if err != nil {
//		return "", err
//	}
//
//	name := srcSnapshot[strings.Index(srcSnapshot, "/")+1:]
//	dstSnapshot := path.Join(pool, name)
//	execute := utilzfs.Zfs(z.zfs).SendAndRecv(srcSnapshot, dstSnapshot, user, host)
//	klog.V(4).Infof(execute.Commit())
//	_, err = execute.Exec()
//
//	return dstSnapshot, err
//}
//
//func (z *ZFSadm) SendIncrementSnapshot(srcSnapshot, pool, user, host string) (string, error) {
//	_, err := z.GetSnapshot(srcSnapshot)
//	if err != nil {
//		return "", err
//	}
//
//	// srcSnapshot 的上一个快照，根据快照是否存在来判断是选择全量发送快照还是增量发送快照
//	lastSnapshot := ""
//
//	// 获取快照的源磁盘名称
//	source := strings.SplitN(srcSnapshot, "@", 2)[0]
//	// 通过源磁盘来查找 srcSnapshot 相关的其他快照
//	execute := utilzfs.Zfs(z.zfs).
//		List(source, "-H", "", []string{"name"}, "", "", "snap")
//	klog.V(4).Infof(execute.Commit())
//	data, err := execute.Exec()
//	if err != nil {
//		return "", err
//	}
//	// 保存源磁盘中关联的所有快照
//	snanpshots := strings.Split(string(data), "\n")
//	// index 用于存储 srcSnapshot 在 snapshots 列表中的位置
//	// index-1 就是上次个的快照
//	index := 0
//	for i, snapshot := range snanpshots {
//		if snapshot == srcSnapshot {
//			index = i
//			break
//		}
//	}
//	if index != 0 {
//		lastSnapshot = snanpshots[index-1]
//		klog.V(4).Infof("discovery last snapshot:%v", lastSnapshot)
//	}
//
//	if len(lastSnapshot) == 0 {
//		return "", fmt.Errorf("couldn't find last snapshot")
//	}
//
//	name := srcSnapshot[strings.Index(srcSnapshot, "/")+1:]
//	dstSnapshot := path.Join(pool, name)
//	execute = utilzfs.Zfs(z.zfs).IncrementSendAndRecv(srcSnapshot, lastSnapshot, dstSnapshot, user, host)
//	klog.V(4).Infof(execute.Commit())
//	_, err = execute.Exec()
//
//	return dstSnapshot, err
//}
//
//func (z *ZFSadm) SetSnapshot(name string) (*Snapshot, error) {
//	_, err := z.GetSnapshot(name)
//	if err == nil {
//		return nil, fmt.Errorf("snapshot '%s' is already exists", name)
//	}
//	parent := strings.SplitN(name, "@", 2)[0]
//	properties := map[string]string{RELATIONSHIP: parent}
//	execute := utilzfs.Zfs(z.zfs).Set(name, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, _ = execute.Exec()
//	properties = map[string]string{DISPLAY: Hidden}
//	execute = utilzfs.Zfs(z.zfs).Set(parent, properties)
//	klog.V(4).Infoln(execute.Commit())
//	_, _ = execute.Exec()
//	execute = utilzfs.Zfs(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
//	klog.V(4).Infoln(execute.Commit())
//	data, _ := execute.Exec()
//	pool := strings.SplitN(parent, "/", 2)[0]
//	snapshot := &Snapshot{Pool: pool, Parent: parent, Name: name, Clones: map[string]string{}}
//	z.mux.Lock()
//	setSnapshot(data, snapshot)
//	z.System.Snapshots[name] = snapshot
//	z.System.Pools[pool].Devs[name] = name
//	z.mux.Unlock()
//
//	return snapshot, nil
//}
//
//func (z *ZFSadm) Get(name, options, max string, out []string, t, s string, properties ...string) ([]byte, error) {
//	execute := utilzfs.Zfs(z.zfs).Get(name, options, max, out, t, s, properties...)
//	klog.V(4).Infoln(execute.Commit())
//	return execute.Exec()
//}
//
//func setSnapshot(data []byte, into *Snapshot) {
//	reader := bytes.NewReader(data)
//	rd := bufio.NewReader(reader)
//	for {
//		line, err := rd.ReadString('\n')
//		if err != nil && err != io.EOF {
//			// 读取中出现错误直接退出
//			break
//		}
//
//		line = strings.TrimSpace(line)
//		if len(line) == 0 {
//			continue
//		}
//
//		parts := strings.Split(line, "\t")
//		if parts[1] == "clones" {
//			clones := strings.Split(strings.TrimSpace(parts[2]), ",")
//			for _, clone := range clones {
//				if len(strings.TrimSpace(clone)) == 0 {
//					continue
//				}
//				into.Clones[clone] = clone
//			}
//		} else {
//			set(into, line, "\t")
//		}
//
//		if err == io.EOF {
//			return
//		}
//	}
//}
//
//// 从 <pool>/<name>@<snapshot> 中提取 <pool>/<name>
//func getName(line string) string {
//	i := strings.LastIndex(line, "@")
//	if i <= 0 {
//		i = len(line)
//	}
//	return line[:i]
//}
//
//// 从 <pool>/<name> 中提取 name
//func getPoolName(text string) string {
//	return strings.SplitN(text, "/", 2)[0]
//}
//
//func setValue(data []byte, into interface{}) {
//	reader := bytes.NewReader(data)
//	rd := bufio.NewReader(reader)
//	for {
//		line, err := rd.ReadString('\n')
//		if err != nil && err != io.EOF {
//			// 读取中出现错误直接退出
//			break
//		}
//
//		line = strings.TrimSpace(line)
//		if len(line) == 0 {
//			continue
//		}
//		set(into, line, "\t")
//		if err == io.EOF {
//			// 读取到尾部，返回
//			return
//		}
//	}
//}
//
//// 利用反射动态设置 target
//// @target: Target install
//// @data: 需要处理的字符串
//// @slim: 分隔符
//func set(target interface{}, data string, slim string) {
//	line := strings.Split(data, slim)
//	getType := reflect.TypeOf(target).Elem()
//	valueType := reflect.ValueOf(target).Elem()
//	for i := 0; i < getType.NumField(); i++ {
//		field := getType.Field(i)
//		if field.Tag.Get("zfs") == line[1] {
//			valueType.Field(i).SetString(line[2])
//		}
//	}
//}

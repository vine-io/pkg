package zfs

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	once sync.Once
	adm  = &ZFSadm{}
)

type ZFSadm struct {
	zpool string
	zfs   string

	//System *System
}

func Default() (*ZFSadm, error) {
	/**
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("zfs must work in linux")
	}
	**/

	adm.lazy()
	if len(adm.zfs) == 0 {
		return nil, fmt.Errorf("could't find command 'zfs', zfs is not installed")
	}
	if len(adm.zpool) == 0 {
		return nil, fmt.Errorf("could't find command 'zpool', zfs is not installed")
	}
	//go adm.GarbageCollecting()
	return adm, nil
}

func (z *ZFSadm) lazy() {
	once.Do(func() {
		if runtime.GOOS == "linux" {
			data, err := exec.Command("sh", "-c", `lsmod | grep zfs`).CombinedOutput()
			if err != nil || string(data) == "" {
				return
			}

			// 内核加载 zfs 模块
			_, _ = exec.Command("modprobe", "zfs").CombinedOutput()

			data, err = exec.Command("which", "zpool").CombinedOutput()
			if err != nil || string(data) == "" {
				return
			}
			adm.zpool = strings.TrimSuffix(string(data), "\n")
			data, err = exec.Command("which", "zfs").CombinedOutput()
			if err != nil || string(data) == "" {
				return
			}
			adm.zfs = strings.TrimSuffix(string(data), "\n")
		} else {
			adm.zfs = "zfs"
			adm.zpool = "zpool"
		}
	})
}

// GetPoolIO return zfs pool:
func (z *ZFSadm) GetPoolIO(ctx context.Context, name string) (string, string, string, string, error) {
	var (
		readIO  = "0"
		writeIO = "0"
		read    = "0"
		write   = "0"
	)

	//if _, err := z.getPool(ctx, name); err != nil {
	//	return "", "", "", "", err
	//}

	execute := ZPoolCtl(z.zpool).Iostat(name, "", "", "-Hp", `-n 5 2 | tail -1`)
	data, _ := execute.Exec()
	line := strings.Split(strings.TrimSuffix(string(data), "\n"), "\t")
	if len(line) > 6 {
		readIO = line[3]
		writeIO = line[4]
		read = line[5]
		write = line[6]
	}

	return readIO, writeIO, read, write, nil
}

func (z *ZFSadm) CreatePool(ctx context.Context, name, compression, raid string, quota float64, devices []string) (*Pool, error) {

	if pool, _ := z.getPool(ctx, name); pool != nil {
		return nil, fmt.Errorf("pool '%s' is already exists", name)
	}

	pool := &Pool{Name: name}
	execute := ZPoolCtl(z.zpool).Create(name, "-f", "", "", nil, raid, devices...)
	_, err := execute.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapPool(ctx, name, pool)

	poolSize, _ := strconv.ParseFloat(pool.Size_, 64)
	_quota := int64(poolSize * quota / 100)
	properties := map[string]string{
		"compression": compression,
		"quota":       strconv.FormatInt(_quota, 10),
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	_, err = execute.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}
	return pool, nil
}

func (z *ZFSadm) ExpensePool(ctx context.Context, name string, devices ...string) (*Pool, error) {

	pool, err := z.GetPool(ctx, name)
	if err != nil {
		return nil, err
	}

	execute := ZPoolCtl(z.zpool).Add(name, "", devices...)
	_, err = execute.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapPool(ctx, name, pool)

	return pool, nil
}

func (z *ZFSadm) SetPoolQuota(ctx context.Context, name string, quota float64) (*Pool, error) {

	pool, err := z.getPool(ctx, name)
	if err != nil {
		return nil, err
	}

	poolSize, _ := strconv.ParseFloat(pool.Size_, 64)
	_quota := int64(poolSize * quota / 100)
	properties := map[string]string{
		"quota": strconv.FormatInt(_quota, 10),
	}
	execute := ZFSCtl(z.zfs).Set(name, properties)
	_, err = execute.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapPool(ctx, name, pool)

	return pool, nil
}

func (z *ZFSadm) wrapPool(ctx context.Context, name string, out *Pool) {
	execute := ZPoolCtl(z.zpool).Get(name, "-Hp", nil, "all")
	data, _ := execute.Exec()
	SetValue(data, out)

	execute = ZFSCtl(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
	data, _ = execute.Exec()
	SetValue(data, out)
}

func (z *ZFSadm) GetPools(ctx context.Context) (map[string]*Pool, error) {
	pools := make(map[string]*Pool, 0)
	// zpool list -Hp -o name
	shell := ZPoolCtl(adm.zpool).
		List("", "-Hp", []string{"name"}, "")
	out, _ := shell.Exec()
	if len(out) == 0 {
		return pools, nil
	}

	names := strings.Split(string(out), "\n")
	for _, name := range names {
		pools[name] = &Pool{Name: name}
	}

	// 解析所有的pool
	for name, pool := range pools {
		z.wrapPool(ctx, name, pool)
	}
	return pools, nil
}

func (z *ZFSadm) getPool(ctx context.Context, name string) (*Pool, error) {
	shell := ZPoolCtl(adm.zpool).
		List(name, "-Hp", []string{"name"}, "")
	out, _ := shell.Exec()
	if len(out) == 0 {
		return nil, fmt.Errorf("pool '%s' is not exists", name)
	}
	return &Pool{Name: name}, nil
}

func (z *ZFSadm) GetPool(ctx context.Context, name string) (*Pool, error) {

	pool, err := z.getPool(ctx, name)
	if err != nil {
		return nil, err
	}

	z.wrapPool(ctx, name, pool)
	return pool, nil
}

func (z *ZFSadm) DeletePool(ctx context.Context, name string) (*Pool, error) {
	pool, err := z.getPool(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("pool '%s' is not exists: %v", name, err)
	}

	execute := ZPoolCtl(z.zpool).Destroy(name, false)
	_, err = execute.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	return pool, nil
}

func (z *ZFSadm) GetFileSystems(ctx context.Context) (map[string]*Volume, error) {
	filesystems := make(map[string]*Volume, 0)
	// 获取所有的 filesystem
	shell := ZFSCtl(adm.zfs).
		List("", "-Hp", "", []string{"name", "origin"}, "", "", "filesystem")
	out, _ := shell.Exec()
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		name := strings.TrimSpace(parts[0])
		if !strings.Contains(name, "/") {
			continue
		}
		var source string
		if v := strings.TrimSpace(parts[1]); v != "-" {
			source = v
		}
		filesystems[name] = &Volume{
			Name:   name,
			Source: source,
		}
	}
	// 解析所有的filesystem
	for name, filesystem := range filesystems {
		z.wrapFileSystem(ctx, name, filesystem)
	}
	return filesystems, nil
}

func (z *ZFSadm) GetFileSystem(ctx context.Context, name string) (*Volume, error) {
	fs, err := z.getFileSystem(ctx, name)
	if err != nil {
		return nil, err
	}
	z.wrapFileSystem(ctx, name, fs)
	return fs, nil
}

func (z *ZFSadm) getFileSystem(ctx context.Context, name string) (*Volume, error) {
	shell := ZFSCtl(adm.zfs).
		Get(name, "-Hp", "", []string{"name"}, "", "", "all")
	_, err := shell.Exec()
	if err != nil {
		return nil, err
	}
	return &Volume{Name: name}, nil
}

func (z *ZFSadm) wrapFileSystem(ctx context.Context, name string, out *Volume) {
	execute := ZFSCtl(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
	data, _ := execute.Exec()
	SetValue(data, out)
}

func (z *ZFSadm) CreateFileSystem(ctx context.Context, name string, properties map[string]string) (*Volume, error) {

	pool := strings.SplitN(name, "/", 2)[0]
	if _, err := z.getPool(ctx, pool); err != nil {
		return nil, err
	}
	if fs, _ := z.getFileSystem(ctx, name); fs != nil {
		return nil, fmt.Errorf("fileSystem '%s' is already exists", name)
	}

	fs := &Volume{Name: name}
	execute := ZFSCtl(z.zfs).CreateFileSystem(name, nil)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapFileSystem(ctx, name, fs)

	return fs, nil
}

func (z *ZFSadm) ShareFileSystem(ctx context.Context, name string, ips []string, mode string) error {
	// zfs set sharenfs='rw=@192.168.2.0/24,rw=@192.168.221.0/24,all_squash,insecure' tank/test
	if _, err := z.getFileSystem(ctx, name); err != nil {
		return err
	}

	properties := map[string]string{"sync": "disabled"}
	execute := ZFSCtl(z.zfs).Set(name, properties)
	_, _ = execute.Exec()

	args := ""
	for _, ip := range ips {
		args += fmt.Sprintf(`%s=@%s,`, mode, ip)
	}
	args += "no_root_squash,insecure"
	properties = map[string]string{"sharenfs": fmt.Sprintf(`'%s'`, args)}
	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	return nil
}

func (z *ZFSadm) UnShareFileSystem(ctx context.Context, name string) error {
	// zfs set sharenfs=off tank/test
	if _, err := z.getFileSystem(ctx, name); err != nil {
		return err
	}

	properties := map[string]string{"sharenfs": "off"}
	execute := ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return fmt.Errorf("%s: %v", execute.Commit(), err)
	}
	return nil
}

// DeleteFileSystem zfs 中删除 filesystem 时，如果 filesystem 存在 snapshot，强制删除时
// 会删除 snapshot 及其关联的 clone
func (z *ZFSadm) DeleteFileSystem(ctx context.Context, name string) error {

	options := "-rRf"
	execute := ZFSCtl(z.zfs).DestroyFileSystemOrVolume(name, options)
	if _, err := execute.Exec(); err != nil {
		return fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	return nil
}

func (z *ZFSadm) GetVolumes(ctx context.Context) (map[string]*Volume, error) {
	volumes := make(map[string]*Volume, 0)
	//// 获取所有的 volume
	shell := ZFSCtl(adm.zfs).
		List("", "-Hp", "", []string{"name", "origin"}, "", "", "volume")
	out, err := shell.Exec()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		name := strings.TrimSpace(parts[0])
		if !strings.Contains(name, "/") {
			continue
		}

		var source string
		if v := strings.TrimSpace(parts[1]); v != "-" {
			source = v
		}
		volumes[name] = &Volume{Name: name, Source: source}
	}

	// 解析所有的volume
	for name, volume := range volumes {
		z.wrapVolume(ctx, name, volume)
	}

	return volumes, nil
}

func (z *ZFSadm) GetVolume(ctx context.Context, name string) (*Volume, error) {
	volume, err := z.getVolume(ctx, name)
	if err != nil {
		return nil, err
	}
	z.wrapVolume(ctx, name, volume)
	return volume, nil
}

func (z *ZFSadm) getVolume(ctx context.Context, name string) (*Volume, error) {
	shell := ZFSCtl(adm.zfs).
		Get(name, "-Hp", "", []string{"name"}, "", "", "all")
	_, err := shell.Exec()
	if err != nil {
		return nil, fmt.Errorf("volume '%s' is not exists", name)
	}
	return &Volume{Name: name}, nil
}

func (z *ZFSadm) wrapVolume(ctx context.Context, name string, volume *Volume) {
	execute := ZFSCtl(z.zfs).Get(name, "-Hp", "", nil, "", "", "all")
	data, _ := execute.Exec()
	SetValue(data, volume)
}

func (z *ZFSadm) CreateVolume(ctx context.Context, name string, properties map[string]string, size int64) (*Volume, error) {

	pool := strings.SplitN(name, "/", 2)[0]
	if _, err := z.GetPool(ctx, pool); err != nil {
		return nil, err
	}
	if vol, _ := z.getVolume(ctx, name); vol != nil {
		return nil, fmt.Errorf("volume '%s' is alreay exists", name)
	}

	vol := &Volume{Name: name}
	execute := ZFSCtl(z.zfs).CreateVolume(name, 4096, nil, strconv.FormatInt(size, 10))
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapVolume(ctx, name, vol)
	return vol, nil
}

func (z *ZFSadm) DeleteVolume(ctx context.Context, name string) error {

	options := "-Rrf"
	execute := ZFSCtl(z.zfs).DestroyFileSystemOrVolume(name, options)
	if _, err := execute.Exec(); err != nil {
		return fmt.Errorf("%s: %v", execute.Commit(), err.Error())
	}

	return nil
}

func (z *ZFSadm) GetSnapshots(ctx context.Context) (map[string]*Snapshot, error) {
	snapshots := make(map[string]*Snapshot, 0)
	// 获取所有的 snapshot
	shell := ZFSCtl(adm.zfs).
		List("", "-Hp", "", []string{"name"}, "", "", "snapshot")
	out, err := shell.Exec()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		name := strings.TrimSpace(parts[0])
		if !strings.Contains(name, "/") {
			continue
		}
		parent := strings.Split(name, "@")[0]
		snapshots[name] = &Snapshot{Name: name, Parent: parent}
	}

	// 解析所有的snapshot
	for name, snapshot := range snapshots {
		z.wrapSnapshot(ctx, name, snapshot)
	}
	return snapshots, nil
}

func (z *ZFSadm) GetSnapshot(ctx context.Context, name string) (*Snapshot, error) {
	snapshot, err := z.getSnapshot(ctx, name)
	if err != nil {
		return nil, err
	}
	z.wrapSnapshot(ctx, name, snapshot)
	return snapshot, nil
}

func (z *ZFSadm) getSnapshot(ctx context.Context, name string) (*Snapshot, error) {
	shell := ZFSCtl(adm.zfs).
		Get(name, "-Hp", "", []string{"name"}, "", "", "all")
	_, err := shell.Exec()
	if err != nil {
		return nil, fmt.Errorf("volume '%s' is not exists", name)
	}
	parent := strings.Split(name, "@")[0]
	return &Snapshot{Name: name, Parent: parent}, nil
}

func (z *ZFSadm) wrapSnapshot(ctx context.Context, name string, out *Snapshot) {
	shell := ZFSCtl(adm.zfs).
		Get(name, "-Hp", "", nil, "", "", "all")
	data, _ := shell.Exec()
	setSnapshot(data, out)
}

func (z *ZFSadm) CreateSnapshot(ctx context.Context, parent, name string, properties map[string]string) (*Snapshot, error) {
	if !strings.Contains(name, "@") {
		return nil, fmt.Errorf("missing '@' in snapshot name")
	}

	if _, err := z.getFileSystem(ctx, parent); err != nil {
		if _, err := z.getVolume(ctx, parent); err != nil {
			return nil, fmt.Errorf("parent '%s' is not exists", parent)
		}
	}

	if ss, _ := z.getSnapshot(ctx, name); ss != nil {
		return nil, fmt.Errorf("snapshot '%s' is already exists", name)
	}

	snapshot := &Snapshot{Name: name, Parent: parent}
	execute := ZFSCtl(z.zfs).Snapshot(name, nil)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	z.wrapSnapshot(ctx, name, snapshot)
	return snapshot, nil
}

func (z *ZFSadm) DeleteSnapshot(ctx context.Context, name string) error {

	options := "-dr"
	execute := ZFSCtl(z.zfs).DestroySnapshot(name, options)
	if _, err := execute.Exec(); err != nil {
		return fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	return nil
}

func (z *ZFSadm) CloneFileSystem(ctx context.Context, name, snap string, properties map[string]string) (*Volume, error) {

	if v, _ := z.getFileSystem(ctx, name); v != nil {
		return nil, fmt.Errorf("fileSystem '%s' is already exists", name)
	}

	if _, err := z.getSnapshot(ctx, snap); err != nil {
		return nil, err
	}

	pool := strings.SplitN(snap, "/", 2)[0]
	name = path.Join(pool, name)
	execute := ZFSCtl(z.zfs).Clone(name, nil, snap)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	fileSystem := &Volume{Name: name, Source: snap}
	z.wrapFileSystem(ctx, name, fileSystem)
	return fileSystem, nil
}

func (z *ZFSadm) CloneVolume(ctx context.Context, name, snap string, properties map[string]string) (*Volume, error) {
	if v, _ := z.getVolume(ctx, name); v != nil {
		return nil, fmt.Errorf("volume '%s' is already exists", name)
	}

	if _, err := z.getSnapshot(ctx, snap); err != nil {
		return nil, err
	}

	pool := strings.SplitN(snap, "/", 2)[0]
	name = path.Join(pool, name)
	execute := ZFSCtl(z.zfs).Clone(name, nil, snap)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	execute = ZFSCtl(z.zfs).Set(name, properties)
	if _, err := execute.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %v", execute.Commit(), err)
	}

	volume := &Volume{Name: name, Source: snap}
	z.wrapVolume(ctx, name, volume)
	return volume, nil
}

func (z *ZFSadm) Get(ctx context.Context, name, options, max string, out []string, t, s string, properties ...string) ([]byte, error) {
	execute := ZFSCtl(z.zfs).Get(name, options, max, out, t, s, properties...)
	return execute.Exec()
}

func setSnapshot(data []byte, into *Snapshot) {
	reader := bytes.NewReader(data)
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			// 读取中出现错误直接退出
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		set(into, line, "\t")
		if err == io.EOF {
			return
		}
	}
}

// 从 <pool>/<name>@<snapshot> 中提取 <pool>/<name>
func getName(line string) string {
	i := strings.LastIndex(line, "@")
	if i <= 0 {
		i = len(line)
	}
	return line[:i]
}

// 从 <pool>/<name> 中提取 name
func getPoolName(text string) string {
	return strings.SplitN(text, "/", 2)[0]
}

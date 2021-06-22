// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package scst

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var once sync.Once

var scstcmd = &Scstcmd{
	mu: sync.RWMutex{},
}

type Scstcmd struct {
	mu sync.RWMutex

	scst string
}

func Default() (*Scstcmd, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("scstadmin must be in linux")
	}
	scstcmd.mu.RLock()
	scstcmd.lazy()
	scstcmd.mu.RUnlock()
	if scstcmd.scst == "" {
		return nil, fmt.Errorf("could't find scstadmin, scst is not installed")
	}

	return scstcmd, nil
}

func (s *Scstcmd) lazy() {
	once.Do(func() {
		data, err := exec.Command("which", "scstadmin").CombinedOutput()
		if err != nil || string(data) == "" {
			return
		}
		s.scst = strings.TrimSuffix(string(data), "\n")

		_ = ioutil.WriteFile("/sys/kernel/scst_tgt/targets/iscsi/enabled", []byte("1"), 0644)
	})
}

func (s *Scstcmd) CreateTarget(ctx context.Context, target, driver string) error {
	scst := NewCtl(s.scst).AddTarget(target).Driver(driver)
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) CreateDisk(ctx context.Context, name, block string) error {
	scst := NewCtl(s.scst).OpenDev(name).
		Handler("vdisk_blockio").
		Attr(map[string]string{"filename": block})
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) CreateGroup(ctx context.Context, group, target, driver string) error {
	scst := NewCtl(s.scst).AddGroup(group).Target(target).Driver(driver)
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) CreateLun(ctx context.Context, lun, target, driver, group, device string) error {

	scst := NewCtl(s.scst).AddLun(lun).Target(target).
		Driver(driver)
	if group != "" {
		scst = scst.Group(group)
	}

	scst = scst.Device(device)
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) AddInit(ctx context.Context, iqn, target, driver, group string) error {
	scst := NewCtl(s.scst).AddInit(iqn).Target(target).
		Driver(driver)
	if group != "" {
		scst = scst.Group(group)
	}
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) EnableTarget(ctx context.Context, target, driver string) error {
	scst1 := NewCtl(s.scst).EnableTarget(target).Driver(driver)
	if _, err := scst1.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst1.Commit(), err)
	}
	scst2 := NewCtl(s.scst).SetDrvAttr("iscsi").Attributes(map[string]string{"enabled": "1"})
	if _, err := scst2.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst2.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) Save(ctx context.Context) error {
	scst := NewCtl(s.scst).WriteConfig("/etc/scst.conf")
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) DeleteTarget(ctx context.Context, target, driver string) error {
	scst1 := NewCtl(s.scst).DisableTarget(target).Driver(driver)
	if _, err := scst1.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst1.Commit(), err)
	}
	scst2 := NewCtl(s.scst).RemoveTarget(target).Driver(driver)
	if _, err := scst2.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst2.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) DeleteInit(ctx context.Context, init, target, group, driver string) error {
	scst := NewCtl(s.scst).RemoveInit(init).Target(target)
	if group != "" {
		scst = scst.Group(group)
	}

	scst = scst.Driver(driver).Force()
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) DeleteLun(ctx context.Context, lun, target, group, device, driver string) error {
	scst := NewCtl(s.scst).RemoveLun(lun).Target(target)
	if group != "" {
		scst = scst.Group(group)
	}
	scst = scst.Device(device).Driver(driver)
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) DeleteGroup(ctx context.Context, group, target, driver string) error {
	scst := NewCtl(s.scst).RemoveGroup(group).Target(target).Driver(driver)
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) DeleteDisk(ctx context.Context, name string) error {
	scst := NewCtl(s.scst).CloseDev(name).Handler("vdisk_blockio")
	if _, err := scst.Execute(); err != nil {
		return fmt.Errorf("%s: %v", scst.Commit(), err)
	}
	return nil
}

func (s *Scstcmd) RawScst() *adm {
	return NewCtl(s.scst)
}

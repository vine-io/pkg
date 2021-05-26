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

package scst

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Manager struct {
	sync.RWMutex

	s *System
}

func NewManager() (*Manager, error) {
	s, err := FromKernel()
	if err != nil {
		return nil, err
	}

	return &Manager{s: s}, nil
}

func (m *Manager) GetHandlers() []*Handler {
	handlers := make([]*Handler, 0, len(m.s.Handlers))

	m.RLock()
	defer m.RUnlock()
	for _, h := range m.s.Handlers {
		handlers = append(handlers, h.DeepCopy())
	}
	return handlers
}

func (m *Manager) GetDevices() []*Device {
	devices := make([]*Device, 0)

	m.RLock()
	defer m.RUnlock()
	for _, h := range m.s.Handlers {
		for _, d := range h.Devices {
			devices = append(devices, d.DeepCopy())
		}
	}
	return devices
}

func (m *Manager) GetTargets() []*Target {
	targets := make([]*Target, 0)

	m.RLock()
	defer m.RUnlock()
	for _, d := range m.s.Drivers {
		for _, t := range d.Targets {
			targets = append(targets, t.DeepCopy())
		}
	}
	return targets
}

func (m *Manager) GetDrivers() []*Driver {
	drivers := make([]*Driver, 0)

	m.RLock()
	defer m.RUnlock()
	for _, d := range m.s.Drivers {
		drivers = append(drivers, d.DeepCopy())
	}
	return drivers
}

func (m *Manager) GetLuns() []*Lun {
	luns := make([]*Lun, 0)

	m.RLock()
	defer m.RUnlock()
	dr, ok := m.s.Drivers["copy_manager"]
	if !ok {
		return luns
	}
	tt, ok := dr.Targets["copy_manager_tgt"]
	if !ok {
		return luns
	}
	for _, lun := range tt.Luns {
		luns = append(luns, lun.DeepCopy())
	}
	return luns
}

func (m *Manager) GetGroups() []*Group {
	groups := make([]*Group, 0)

	m.RLock()
	defer m.RUnlock()

	for _, dr := range m.s.Drivers {
		for _, tt := range dr.Targets {
			for _, group := range tt.Groups {
				groups = append(groups, group.DeepCopy())
			}
		}
	}

	return groups
}

// OpenDev create device in HANDLER and return *Device, command and error.
func (m *Manager) OpenDev(handler, name, filename string) (*Device, string, error) {

	mgmt := filepath.Join(kernel, "handlers", handler, "mgmt")

	cmd := fmt.Sprintf("add_device %s filename=%s", name, filename)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	h, ok := m.s.Handlers[handler]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("handler '%s' not exists", handler)
	}
	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	device := &Device{
		Name:     name,
		Filename: filename,
	}

	size := readHeader(filepath.Join(kernel, handler, name, "size"))
	device.Size, _ = strconv.ParseInt(size, 10, 64)

	m.Lock()
	defer m.Unlock()
	h.Devices[name] = device

	return device, history, nil
}

// DelDev delete device in HANDLER
func (m *Manager) DelDev(handler, name string) (*Device, string, error) {

	mgmt := filepath.Join(kernel, "handlers", handler, "mgmt")
	cmd := fmt.Sprintf("add_device %s", name)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	h, ok := m.s.Handlers[handler]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("handler '%s' not exists", handler)
	}

	dev, ok := h.Devices[name]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("device '%s' not exists", name)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	m.Lock()
	defer m.Unlock()
	delete(h.Devices, name)

	return dev, history, nil
}

// CreateTarget create target in Driver
func (m *Manager) CreateTarget(driver, name string) (*Target, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, "mgmt")

	cmd := fmt.Sprintf("add_target %s", name)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}
	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	target := &Target{
		Name:    name,
		Enabled: 0,
		Groups:  map[string]*Group{},
		Luns:    []*Lun{},
	}

	m.Lock()
	defer m.Unlock()
	dr.Targets[name] = target

	return target, history, nil
}

// DelTarget delete target in Driver
func (m *Manager) DelTarget(driver, name string) (*Target, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, "mgmt")

	cmd := fmt.Sprintf("del_target %s", name)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	target, ok := dr.Targets[name]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", name)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	m.Lock()
	defer m.Unlock()
	delete(dr.Targets, name)

	return target, history, nil
}

// EnableTarget enabled target. If target id equal 0, updates it.
func (m *Manager) EnableTarget(driver, name string) (*Target, string, error) {
	enabled := filepath.Join(kernel, "targets", driver, name, "enabled")

	history := fmt.Sprintf(`echo "1" > %s`, enabled)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	target, ok := dr.Targets[name]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", name)
	}

	m.RUnlock()

	err := ioutil.WriteFile(enabled, []byte("1"), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	if target.Id == 0 {
		id := readHeader(filepath.Join(kernel, "targets", driver, name, "rel_tgt_id"))
		target.Id, _ = strconv.ParseInt(id, 10, 64)
	}
	target.Enabled = 1

	m.Lock()
	defer m.Unlock()
	dr.Targets[name] = target

	return target, history, nil
}

// DisableTarget disable target
func (m *Manager) DisableTarget(driver, name string) (*Target, string, error) {
	enabled := filepath.Join(kernel, "targets", driver, name, "enabled")

	history := fmt.Sprintf(`echo "0" > %s`, enabled)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	target, ok := dr.Targets[name]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", name)
	}

	m.RUnlock()

	err := ioutil.WriteFile(enabled, []byte("0"), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	target.Enabled = 0

	m.Lock()
	defer m.Unlock()
	dr.Targets[name] = target

	return target, history, nil
}

// CreateGroup create group
func (m *Manager) CreateGroup(driver, target, name string) (*Group, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, target, "ini_groups", "mgmt")
	cmd := fmt.Sprintf("create %s", name)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", target)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	group := &Group{
		Name:       name,
		Luns:       []*Lun{},
		Initiators: []string{},
	}

	tt.Groups[name] = group

	m.Lock()
	defer m.Unlock()
	dr.Targets[name] = tt

	return group, history, nil
}

// DelGroup delete group
func (m *Manager) DelGroup(driver, target, name string) (*Group, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, target, "ini_groups", "mgmt")
	cmd := fmt.Sprintf("del %s", name)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", target)
	}

	gg, ok := tt.Groups[name]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("group '%s' not exists", name)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	m.Lock()
	defer m.Unlock()
	delete(tt.Groups, name)

	return gg, history, nil
}

// CreateLun create logical unit in Target. If group is "", direct create lun in Target.
func (m *Manager) CreateLun(driver, target, group, device string, id int64) (*Lun, string, error) {
	var mgmt string
	if len(group) != 0 {
		mgmt = filepath.Join(kernel, "targets", driver, target, "ini_groups", group, "luns", "mgmt")
	} else {
		mgmt = filepath.Join(kernel, "targets", driver, target, "luns", "mgmt")
	}

	cmd := fmt.Sprintf("add %s %d", device, id)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	var gg *Group
	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", target)
	}

	exists := false
	for _, h := range m.s.Handlers {
		if _, ok := h.Devices[device]; ok {
			exists = true
			break
		}
	}
	if !exists {
		m.RUnlock()
		return nil, history, fmt.Errorf("device '%s' not exists", device)
	}

	if len(group) != 0 {
		var ok bool
		gg, ok = tt.Groups[group]
		if !ok {
			m.RUnlock()
			return nil, history, fmt.Errorf("group '%s' not exists", group)
		}
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	lun := &Lun{Id: id, Device: device}

	m.Lock()
	if gg != nil {
		gg.Luns = append(gg.Luns, lun)
	} else {
		tt.Luns = append(tt.Luns, lun)
	}
	m.Unlock()

	m.reloadTarget()

	return lun, history, nil
}

// DelLun delete logical unit.
func (m *Manager) DelLun(driver, target, group string, id int64) (*Lun, string, error) {
	var mgmt string
	if len(group) != 0 {
		mgmt = filepath.Join(kernel, "targets", driver, target, "ini_groups", group, "luns", "mgmt")
	} else {
		mgmt = filepath.Join(kernel, "targets", driver, target, "luns", "mgmt")
	}

	cmd := fmt.Sprintf("del %d", id)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	var gg *Group
	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return nil, history, fmt.Errorf("target '%s' not exists", target)
	}

	if len(group) != 0 {
		var ok bool
		gg, ok = tt.Groups[group]
		if !ok {
			m.RUnlock()
			return nil, history, fmt.Errorf("group '%s' not exists", group)
		}
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return nil, history, err
	}

	var curLun *Lun
	m.Lock()
	if gg != nil {
		luns := make([]*Lun, 0, len(gg.Luns)-1)
		for _, lun := range gg.Luns {
			if lun.Id == id {
				curLun = lun
				continue
			}
			luns = append(luns, lun)
		}
		gg.Luns = luns
	} else {
		luns := make([]*Lun, 0, len(tt.Luns)-1)
		for _, lun := range tt.Luns {
			if lun.Id == id {
				curLun = lun
				continue
			}
			luns = append(luns, lun)
		}
		tt.Luns = luns
	}
	m.Unlock()

	m.reloadTarget()

	return curLun, history, nil
}

func (m *Manager) reloadTarget() {
	driver := "copy_manager"
	target := "copy_manager_tgt"

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return
	}
	m.RUnlock()

	dirs := readDirs(kernel, "targets", driver, target, "luns")
	luns := make([]*Lun, 0, len(dirs))
	for _, id := range dirs {
		lunId, _ := strconv.ParseInt(id, 10, 64)
		device := readLink(filepath.Join(kernel, "targets", driver, target, "luns", id, "device"))
		if lIndex := strings.LastIndex(device, "/"); lIndex > 0 {
			device = device[lIndex+1:]
		}

		luns = append(luns, &Lun{Id: lunId, Device: device})
	}

	m.Lock()
	tt.Luns = luns
	m.Unlock()

	return
}

// AddInitiator add initiator to ini_group
func (m *Manager) AddInitiator(driver, target, group, initiator string) (string, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, target, "ini_groups", group, "initiators", "mgmt")
	cmd := fmt.Sprintf("add %s", initiator)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("target '%s' not exists", target)
	}

	gg, ok := tt.Groups[group]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("group '%s' not exists", group)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return "", history, err
	}

	m.Lock()
	defer m.Unlock()
	gg.Initiators = append(gg.Initiators, initiator)

	return initiator, history, nil
}

func (m *Manager) DelInitiator(driver, target, group, initiator string) (string, string, error) {
	mgmt := filepath.Join(kernel, "targets", driver, target, "ini_groups", group, "initiators", "mgmt")
	cmd := fmt.Sprintf("del %s", initiator)
	history := fmt.Sprintf(`echo "%s" > %s`, cmd, mgmt)

	m.RLock()
	dr, ok := m.s.Drivers[driver]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("driver '%s' not exists", driver)
	}

	tt, ok := dr.Targets[target]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("target '%s' not exists", target)
	}

	gg, ok := tt.Groups[group]
	if !ok {
		m.RUnlock()
		return "", history, fmt.Errorf("group '%s' not exists", group)
	}

	m.RUnlock()

	err := ioutil.WriteFile(mgmt, []byte(cmd), os.ModePerm)
	if err != nil {
		return "", history, err
	}

	m.Lock()
	defer m.Unlock()
	initiators := make([]string, 0, len(gg.Initiators)-1)
	for _, i := range gg.Initiators {
		if i != initiator {
			initiators = append(initiators, i)
		}
	}
	gg.Initiators = initiators

	return initiator, history, nil
}

// SaveToCfg save scst configuration to /etc/scst.conf
func (m *Manager) SaveToCfg() error {
	out, err := m.s.ToCfg()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(DefaultConf, out, os.ModePerm)
}

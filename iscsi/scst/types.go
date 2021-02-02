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
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"
)

type System struct {
	// scst version
	Version string `json:"version"`
	// the list of handler
	Handlers map[string]*Handler `json:"handlers"`
	// the list of drivers
	Drivers map[string]*Driver `json:"drivers"`
}

// Handler scst handler
// Example by scst config /etc/scst.conf:
// 	HANDLER vdisk_blockio {
//		DEVICE admin:ahtr {
//			filename /dev/sdc
//			size 10737418240
//		}
// 	}
type Handler struct {
	// the name of handler, example vdisk_blockio
	Name string `json:"name"`
	//
	Devices map[string]*Device `json:"devices"`
}

// Device scst device
// Example by scst configure /etc/scst.conf:
//	  DEVICE admin:ahtr {
//			filename /dev/sdc
//			size 10737418240
//	  }
type Device struct {
	// the name of device
	Name string `json:"name"`
	// the filename of device, the path of block file
	Filename string `json:"filename"`
	// the size of device (unit B)
	Size int64 `json:"size"`
}

// Driver scst
// 	TARGET_DRIVER iscsi {
//		enabled 1
//
//		TARGET iqn.2018-11.com.example:ahtr {
//			enabled 1
//			rel_tgt_id 23
//
//			GROUP ahtr {
//				LUN 0 ahtr
//
//				INITIATOR iqn.1991-05.com.microsoft:win-1bp99fqu2ri
//			}
//		}
//	}
type Driver struct {
	// the name of Driver
	Name string `json:"name"`
	//
	Enabled int8 `json:"enabled"`
	// Targets
	Targets map[string]*Target `json:"targets"`
}

// Target scst
// 	TARGET iqn.2018-11.com.example:bsic {
//		enabled 1
//		rel_tgt_id 28
//
//		GROUP bsic {
//			LUN 0 bsic
//
//			INITIATOR iqn.1991-05.com.microsoft:win-1bp99fqu2ri
//		}
//	}
type Target struct {
	Name string `json:"name"`

	Enabled int8 `json:"enabled"`

	Id int64 `json:"id"`

	Groups map[string]*Group `json:"groups"`

	Luns []*Lun `json:"luns"`
}

// Group scst resource group
// 	GROUP admin:bsic {
//		LUN 0 admin:bsic
//
//		INITIATOR iqn.1991-05.com.microsoft:win-1bp99fqu2ri
//	}
type Group struct {
	Name string `json:"name"`

	Luns []*Lun `json:"luns"`

	// iscst agent iqn
	Initiators []string `json:"initiators"`
}

// Lun scst logical unit
//
type Lun struct {
	// the id of lun
	Id int64 `json:"id"`
	// the name of device
	Device string `json:"name"`
}

const (
	_Version     = "version"
	_Handler     = "HANDLER"
	_Device      = "DEVICE"
	_CopyManager = "copy_manager"
	_CopyTgt     = "copy_manager_tgt"
	_Iscsi       = "iscsi"
	_IscsiTarget = "TARGET"
	_Group       = "GROUP"
)

// FromCfg get System by parsing /etc/scst.conf
func FromCfg() (*System, error) {
	return FromCfgFile(DefaultConf)
}

// FromCfgFile get System by parsing scst configuration
func FromCfgFile(cfg string) (*System, error) {
	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	system := System{
		Handlers: map[string]*Handler{},
		Drivers:  map[string]*Driver{},
	}

	var parent, kind string
	var n int
	rd := bufio.NewReader(bytes.NewBuffer(data))

	type resource struct {
		Kind string
		Name string
	}

	stack := list.New()
	/*
		stack = list.New()
		stack.PushBack(elem)
		stack.Remove(stack.Back())
	*/

	// *Ptr points the name of scst configuration file resource
	var handlerPtr, devicePtr, driverPtr, targetPtr, groupPtr string
	var EOF bool
	for {
		// update line number
		n++
		b, err := rd.ReadSlice('\n')
		if err == io.EOF {
			EOF = true
		}
		if err != nil && err != io.EOF {
			return nil, err
		}

		line := strings.TrimSpace(string(b))
		if len(line) == 0 {
			continue
		}

		check := func(line string, parts []string, cfg string, n int) error {
			if len(parts) != 3 {
				return fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
			}
			if !strings.HasSuffix(line, "{") {
				return fmt.Errorf("%w: missing '{' at %s:%d", ErrSyntax, cfg, n)
			}
			return nil
		}

		parts := strings.Split(line, " ")
		switch {
		case strings.HasPrefix(line, "#") && strings.Contains(line, "SCST"):
			kind = _Version
			system.Version = strings.TrimSuffix(parts[len(parts)-1], ".")

		case strings.HasPrefix(line, "HANDLER"):
			parent, kind = _Handler, _Handler

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			handlerPtr = parts[1]
			system.Handlers[handlerPtr] = &Handler{
				Name:    handlerPtr,
				Devices: map[string]*Device{},
			}
			stack.PushBack(&resource{Kind: _Handler, Name: handlerPtr})

		case parent == _Handler && strings.HasPrefix(line, "DEVICE"):
			kind = _Device

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			devicePtr = parts[1]
			system.Handlers[handlerPtr].Devices[devicePtr] = &Device{Name: devicePtr}
			stack.PushBack(&resource{Kind: _Device, Name: devicePtr})

		case strings.HasPrefix(line, "TARGET_DRIVER") && strings.Contains(line, "copy_manager"):
			parent, kind = _CopyManager, _CopyManager

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			driverPtr = parts[1]
			system.Drivers[driverPtr] = &Driver{
				Name:    driverPtr,
				Targets: map[string]*Target{},
			}
			stack.PushBack(&resource{Kind: _CopyManager, Name: driverPtr})

		case strings.HasPrefix(line, "TARGET_DRIVER") && strings.Contains(line, "iscsi"):
			parent, kind = _Iscsi, _Iscsi

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			driverPtr = parts[1]
			system.Drivers[driverPtr] = &Driver{
				Name:    driverPtr,
				Targets: map[string]*Target{},
			}
			stack.PushBack(&resource{Kind: _Iscsi, Name: driverPtr})

		case parent == _CopyManager && strings.HasPrefix(line, "TARGET"):
			kind = _CopyTgt

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			targetPtr = parts[1]
			system.Drivers[driverPtr].Targets[targetPtr] = &Target{
				Name:   targetPtr,
				Groups: map[string]*Group{},
				Luns:   make([]*Lun, 0),
			}
			stack.PushBack(&resource{Kind: _CopyTgt, Name: targetPtr})

		case parent == _Iscsi && strings.HasPrefix(line, "TARGET"):
			kind = _IscsiTarget

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			targetPtr = parts[1]
			system.Drivers[driverPtr].Targets[targetPtr] = &Target{
				Name:   targetPtr,
				Groups: map[string]*Group{},
				Luns:   make([]*Lun, 0),
			}
			stack.PushBack(&resource{Kind: _IscsiTarget, Name: targetPtr})

		case strings.HasPrefix(line, "GROUP"):
			kind = _Group

			if err := check(line, parts, cfg, n); err != nil {
				return nil, err
			}
			groupPtr = parts[1]
			system.Drivers[driverPtr].Targets[targetPtr].Groups[groupPtr] = &Group{
				Name:       groupPtr,
				Luns:       make([]*Lun, 0),
				Initiators: []string{},
			}
			stack.PushBack(&resource{Kind: _Group, Name: groupPtr})

		case strings.HasSuffix(line, "}"):
			if stack.Len() == 0 {
				return nil, fmt.Errorf("%w: don't match '}' at %s:%d", ErrNoScst, cfg, n)
			}
			stack.Remove(stack.Back())
		}

		switch kind {
		//case _Version:
		//case _Handler:
		case _Device:
			if strings.HasPrefix(line, "filename") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				system.Handlers[handlerPtr].Devices[devicePtr].Filename = parts[1]
			}
			if strings.HasPrefix(line, "size") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				system.Handlers[handlerPtr].Devices[devicePtr].Size, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		//case _CopyManager:
		case _CopyTgt:

			if strings.HasPrefix(line, "LUN") {
				if len(parts) != 3 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				target := system.Drivers[driverPtr].Targets[targetPtr]
				lun := &Lun{Device: parts[2]}
				lun.Id, _ = strconv.ParseInt(parts[1], 10, 64)
				target.Luns = append(target.Luns, lun)
			}

		case _Iscsi:

			if strings.HasPrefix(line, "enabled") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				enabled, _ := strconv.ParseInt(parts[1], 10, 64)
				system.Drivers[driverPtr].Enabled = int8(enabled)
			}

		case _IscsiTarget:

			if strings.HasPrefix(line, "enabled") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				enabled, _ := strconv.ParseInt(parts[1], 10, 64)
				system.Drivers[driverPtr].Targets[targetPtr].Enabled = int8(enabled)
			}
			if strings.HasPrefix(line, "rel_tgt_id") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				system.Drivers[driverPtr].Targets[targetPtr].Id, _ = strconv.ParseInt(parts[1], 10, 64)
			}
			if strings.HasPrefix(line, "LUN") {
				if len(parts) != 3 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				target := system.Drivers[driverPtr].Targets[targetPtr]
				lun := &Lun{Device: parts[2]}
				lun.Id, _ = strconv.ParseInt(parts[1], 10, 64)
				target.Luns = append(target.Luns, lun)
			}

		case _Group:

			if strings.HasPrefix(line, "LUN") {
				if len(parts) != 3 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				group := system.Drivers[driverPtr].Targets[targetPtr].Groups[groupPtr]
				lun := &Lun{Device: parts[2]}
				lun.Id, _ = strconv.ParseInt(parts[1], 10, 64)
				group.Luns = append(group.Luns, lun)
			}

			if strings.HasPrefix(line, "INITIATOR") {
				if len(parts) != 2 {
					return nil, fmt.Errorf("%w: bad format '%s' at %s:%d", ErrSyntax, line, cfg, n)
				}
				group := system.Drivers[driverPtr].Targets[targetPtr].Groups[groupPtr]
				group.Initiators = append(group.Initiators, parts[1])
			}
		}

		if EOF {
			break
		}
	}

	if stack.Len() != 0 {
		r := stack.Back().Value.(*resource)
		return nil, fmt.Errorf("%w: missing '}' for %s<%s>", ErrSyntax, r.Kind, r.Name)
	}

	// parse /etc/scst.conf
	return &system, nil
}

func (s *System) ToCfg() ([]byte, error) {
	tmpl, err := template.New("scst").Parse(ScstTmpl)
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte(""))
	err = tmpl.Execute(out, s)
	return out.Bytes(), err
}

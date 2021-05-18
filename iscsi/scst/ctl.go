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

// Package scst could not work in windows
package scst

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var DefaultConf = "/etc/scst.conf"

type adm struct {
	cmd string
}

func NewCtl(cmd string) *adm {
	return &adm{cmd: cmd}
}

func (a *adm) AddTarget(target string) *addTargetCmd {
	return &addTargetCmd{innerCmd{a.cmd, []string{"-add_target", target}}}
}

type innerCmd struct {
	cmd  string
	args []string
}

func (c *innerCmd) Commit() string {
	return fmt.Sprintf("%s %s", c.cmd, strings.Join(c.args, " "))
}

func (c *innerCmd) Execute() ([]byte, error) {
	out, err := exec.Command(c.cmd, c.args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", err, string(out))
	}
	return bytes.TrimSuffix(out, []byte("\n")), nil
}

type addTargetCmd struct {
	innerCmd
}

func (c *addTargetCmd) Driver(driver string) *addTargetCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (a *adm) OpenDev(dev string) *openDevCmd {
	return &openDevCmd{innerCmd{a.cmd, []string{"-open_dev", dev}}}
}

type openDevCmd struct {
	innerCmd
}

func (c *openDevCmd) Handler(handler string) *openDevCmd {
	c.args = append(c.args, "-handler", handler)
	return c
}

func (c *openDevCmd) Attr(attrs map[string]string) *openDevCmd {
	attributes := make([]string, 0)
	for k, v := range attrs {
		attributes = append(attributes, fmt.Sprintf("%s=%s", k, v))
	}
	c.args = append(c.args, "-attributes", strings.Join(attributes, ","))
	return c
}

func (a *adm) AddGroup(group string) *addGroupCmd {
	return &addGroupCmd{innerCmd{a.cmd, []string{"-add_group", group}}}
}

type addGroupCmd struct {
	innerCmd
}

func (c *addGroupCmd) Driver(driver string) *addGroupCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *addGroupCmd) Target(target string) *addGroupCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (a *adm) AddInit(iqn string) *addInitCmd {
	return &addInitCmd{innerCmd{a.cmd, []string{"-add_init", iqn}}}
}

type addInitCmd struct {
	innerCmd
}

func (c *addInitCmd) Driver(driver string) *addInitCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *addInitCmd) Target(target string) *addInitCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (c *addInitCmd) Group(group string) *addInitCmd {
	c.args = append(c.args, "-group", group)
	return c
}

func (a *adm) AddLun(lun string) *addLunCmd {
	return &addLunCmd{innerCmd{a.cmd, []string{"-add_lun", lun}}}
}

type addLunCmd struct {
	innerCmd
}

func (c *addLunCmd) Driver(driver string) *addLunCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *addLunCmd) Target(target string) *addLunCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (c *addLunCmd) Group(target string) *addLunCmd {
	c.args = append(c.args, "-group", target)
	return c
}

func (c *addLunCmd) Device(dev string) *addLunCmd {
	c.args = append(c.args, "-device", dev)
	return c
}

func (a *adm) EnableTarget(target string) *enableTargetCmd {
	return &enableTargetCmd{innerCmd{a.cmd, []string{"-enable_target", target}}}
}

type enableTargetCmd struct {
	innerCmd
}

func (c *enableTargetCmd) Driver(driver string) *enableTargetCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (a *adm) SetDrvAttr(driver string) *setDrvAttrCmd {
	return &setDrvAttrCmd{innerCmd{a.cmd, []string{"-set_drv_attr", driver}}}
}

type setDrvAttrCmd struct {
	innerCmd
}

func (c *setDrvAttrCmd) Attributes(attrs map[string]string) *setDrvAttrCmd {
	attributes := make([]string, 0)
	for k, v := range attrs {
		attributes = append(attributes, fmt.Sprintf("%s=%s", k, v))
	}
	c.args = append(c.args, "-attributes", strings.Join(attributes, ","))
	return c
}

func (c *setDrvAttrCmd) NoPrompt() *setDrvAttrCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) WriteConfig(cfg string) *writeCfgCmd {
	return &writeCfgCmd{innerCmd{a.cmd, []string{"-write_config", cfg}}}
}

type writeCfgCmd struct {
	innerCmd
}

func (a *adm) ReadCfg(cfg string) *readCfgCmd {
	return &readCfgCmd{innerCmd{a.cmd, []string{"-config", cfg}}}
}

type readCfgCmd struct {
	innerCmd
}

func (a *adm) ClearConfig() *writeCfgCmd {
	return &writeCfgCmd{innerCmd{a.cmd, []string{"-clear_config"}}}
}

type clearCfgCmd struct {
	innerCmd
}

func (c *clearCfgCmd) Force() *clearCfgCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *clearCfgCmd) NoPrompt() *clearCfgCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) DisableTarget(target string) *disableTargetCmd {
	return &disableTargetCmd{innerCmd{a.cmd, []string{"-target", target}}}
}

type disableTargetCmd struct {
	innerCmd
}

func (c *disableTargetCmd) Driver(driver string) *disableTargetCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *disableTargetCmd) Force() *disableTargetCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *disableTargetCmd) NoPrompt() *disableTargetCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) RemoveTarget(target string) *remTargetCmd {
	return &remTargetCmd{innerCmd{a.cmd, []string{"-rem_target", target}}}
}

type remTargetCmd struct {
	innerCmd
}

func (c *remTargetCmd) Driver(driver string) *remTargetCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *remTargetCmd) Force() *remTargetCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *remTargetCmd) NoPrompt() *remTargetCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) RemoveGroup(group string) *remGroupCmd {
	return &remGroupCmd{innerCmd{a.cmd, []string{"-rem_group", group}}}
}

type remGroupCmd struct {
	innerCmd
}

func (c *remGroupCmd) Driver(driver string) *remGroupCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *remGroupCmd) Target(target string) *remGroupCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (c *remGroupCmd) Force() *remGroupCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *remGroupCmd) NoPrompt() *remGroupCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) RemoveLun(lun string) *remLunCmd {
	return &remLunCmd{innerCmd{a.cmd, []string{"-rem_lun", lun}}}
}

type remLunCmd struct {
	innerCmd
}

func (c *remLunCmd) Driver(driver string) *remLunCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *remLunCmd) Target(target string) *remLunCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (c *remLunCmd) Group(group string) *remLunCmd {
	c.args = append(c.args, "-group", group)
	return c
}

func (c *remLunCmd) Device(device string) *remLunCmd {
	c.args = append(c.args, "-device", device)
	return c
}

func (c *remLunCmd) Force() *remLunCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *remLunCmd) NoPrompt() *remLunCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) RemoveInit(init string) *remInitCmd {
	return &remInitCmd{innerCmd{a.cmd, []string{"-rem_init", init}}}
}

type remInitCmd struct {
	innerCmd
}

func (c *remInitCmd) Driver(driver string) *remInitCmd {
	c.args = append(c.args, "-driver", driver)
	return c
}

func (c *remInitCmd) Target(target string) *remInitCmd {
	c.args = append(c.args, "-target", target)
	return c
}

func (c *remInitCmd) Group(group string) *remInitCmd {
	c.args = append(c.args, "-group", group)
	return c
}

func (c *remInitCmd) Force() *remInitCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *remInitCmd) NoPrompt() *remInitCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

func (a *adm) CloseDev(dev string) *closeDevCmd {
	return &closeDevCmd{innerCmd{a.cmd, []string{"-close_dev", dev}}}
}

type closeDevCmd struct {
	innerCmd
}

func (c *closeDevCmd) Handler(handler string) *closeDevCmd {
	c.args = append(c.args, "-handler", handler)
	return c
}

func (c *closeDevCmd) Force() *closeDevCmd {
	c.args = append(c.args, "-force")
	return c
}

func (c *closeDevCmd) NoPrompt() *closeDevCmd {
	c.args = append(c.args, "-noprompt")
	return c
}

package scst

import (
	"encoding/json"
	"testing"
)

func TestFromCfgFile(t *testing.T) {
	s, err := FromCfgFile("scst.conf")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	if len(s.Version) == 0 {
		t.Errorf("got version error")
	}

	v, _ := json.MarshalIndent(s, "", " ")
	t.Log(string(v))
}

func TestSystem_ToCfg(t *testing.T) {
	s, err := FromCfgFile("scst.conf")
	if err != nil {
		t.Error(err)
	}

	out, err := s.ToCfg()
	if err != nil {
		t.Error(err)
	}

	t.Log(string(out))
}

func ExampleNewCtl() {
	scst := "scstadmin"
	// CreateTarget
	NewCtl(scst).AddTarget("iqn.2018-11.com.vol").Driver("iscsi").Execute()

	// CreateDisk
	NewCtl(scst).OpenDev("vol").Handler("vdisk_blockio").Attr(map[string]string{"filename": "/dev/sdb"})

	// CreateGroup
	NewCtl(scst).AddGroup("vol_group").Target("iqn.2018-11.com.example.vol").Driver("iscsi")

	// CreateLun
	NewCtl(scst).AddLun("0").Target("iqn.2018-11.com.example.vol").Driver("iscsi").Group("vol_grup").Device("vol")

	// AddInit
	NewCtl(scst).
		AddLun("iqn.1991-05.com.microsoft:win-1bp99fqu2ri").
		Target("iqn.2018-11.com.example.vol").
		Driver("iscsi").
		Group("vol_group").
		Commit()

	// Enable target
	NewCtl(scst).
		EnableTarget("iqn.2018-11.com.example.vol").
		Driver("iscsi").
		Commit()

	// Save to scst configuration file
	DefaultConf = "/etc/scst.conf"
	NewCtl(scst).
		WriteConfig(DefaultConf).
		Commit()

	// Delete Target
	NewCtl(scst).
		RemoveTarget("iqn.2018-11.com.example.vol").
		Driver("iscsi").
		Commit()

	// Delete Init
	NewCtl(scst).
		RemoveInit("iqn.1991-05.com.microsoft:win-1bp99fqu2ri").
		Target("iqn.2018-11.com.example.vol").
		Group("vol_group").
		Driver("iscsi").
		Force().
		NoPrompt().
		Commit()

	// Delete Lun
	NewCtl(scst).
		RemoveLun("0").
		Target("iqn.2018-11.com.example.vol").
		Group("vol_group").
		Device("vol").
		Driver("iscsi").
		Commit()

	// DeleteGroup
	NewCtl(scst).
		RemoveGroup("vol_group").
		Target("iqn.2018-11.com.example.vol").
		Driver("iscsi").
		Commit()

	// DeleteDisk
	NewCtl(scst).
		CloseDev("vol").
		Handler("vdisk_blockio").
		Commit()
}

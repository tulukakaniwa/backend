package runner

import "os"

type flowParams struct {
	DesignName   string
	VerilogFiles string
	SdcFile      string
	DieArea      string
	CoreArea     string
}

func createConfigFile(filename string, params flowParams) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString("export DESIGN_NAME = " + params.DesignName + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export PLATFORM = nangate45\n\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export VERILOG_FILES = " + params.VerilogFiles + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export SDC_FILE = " + params.SdcFile + "\n\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export DIE_AREA = " + params.DieArea + "\n"); err != nil {
		return err
	}
	if _, err = f.WriteString("export CORE_AREA = " + params.CoreArea + "\n"); err != nil {
		return err
	}
	f.Sync()

	return nil
}

package runner

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type flowConf struct {
	DesignName  string   `yaml:"design_name"`
	DesignFiles []string `yaml:"verilog_files"`
	SdcFile     string   `yaml:"sdc_file"`
	DieArea     string   `yaml:"die_area"`
	CoreArea    string   `yaml:"core_area"`
}

func readFlowConf(filePath string) (*flowConf, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var c flowConf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

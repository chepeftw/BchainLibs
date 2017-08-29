package bchainlibs

import (
    "io/ioutil"
    "fmt"
    "gopkg.in/yaml.v2"
)

type Conf struct {
    TargetSync   float64  `yaml:"target"`
    Nodes        int    `yaml:"nodes"`
    MiningRetry  int    `yaml:"miningretry"`
    MiningWait   int    `yaml:"miningwait"`
    Timeout      int    `yaml:"timeout"`
    RootNode     int    `yaml:"rootnode"`
    Port         int    `yaml:"port"`
    LogPath      string    `yaml:"logpath"`
}

func (c *Conf) GetConf( filename string ) *Conf {

    yamlFile, err := ioutil.ReadFile(filename)
    if err != nil {
	fmt.Errorf("yamlFile.Get err #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
	fmt.Errorf("Unmarshal: %v", err)
    }

    return c
}
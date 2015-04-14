package main

import (
    "io/ioutil"
    "encoding/json"
    "os"
    "log"
)

var config logPullerConfig

type containerConfig struct {
    Image string `json:"image"`
    LogPaths []string `json:"log_paths"`

}

type logPullerConfig struct {
    LogDestination string `json:"log_destination"`
    Containers []containerConfig `json:"containers"`
}


func loadConfigFile() {
    file, e := ioutil.ReadFile("config.json")
    if e != nil {
        log.Fatal("Error ", e)
        os.Exit(1)
    }
    var pullerConfig logPullerConfig
    err := json.Unmarshal(file, &pullerConfig)
    if err != nil {
        log.Fatal(err)
    }
    config = pullerConfig
    setupLogDir(config.LogDestination)
}

func setupLogDir(directory string) {
    _, err := os.Stat(directory)
    if os.IsNotExist(err) {
        log.Printf("%v does not exist, creating",directory)
        err := os.MkdirAll(directory,0644)
        if err != nil {
            log.Fatal(err)
        }
    }
}
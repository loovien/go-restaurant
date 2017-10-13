package main

import (
	log "github.com/cihub/seelog"
	"path/filepath"
	"os"
)

// InitLog
func InitLog(filename string) {
	defer log.Flush()
	logConfigFile, err := getAbsPath(filename)
	if err != nil {
		log.Infof("filename 【%s】 not exists, uses default file config/log4g.xml ！", filename)
		logConfigFile = DEFAULT_LOG_CONFIG_FILE
	}
	logger, err := log.LoggerFromConfigAsFile(logConfigFile)
	if err != nil {
		log.Criticalf("【%s】 【%v】", filename, err)
	} else {
		log.UseLogger(logger)
	}
}

// ReloadLog config
func ReloadLog(filename string, reload bool) bool {
	defer log.Flush()
	logConfigFile, err := getAbsPath(filename)
	if err != nil && !reload {
		log.Infof("【%s】 not exists. Reload failed !", filename)
		return false
	}
	logger, err := log.LoggerFromConfigAsFile(logConfigFile)
	if err != nil {
		log.Criticalf("【%s】 【%v】", filename, err)
		return false
	}
	log.UseLogger(logger)
	return true
}

// getAbsPath get file absolute path
func getAbsPath(filename string) (string, error){
	if !filepath.IsAbs(filename) {
		filename, _ = filepath.Abs(filename)
	}
	 _, err := os.Stat(filename);
	return filename, err
}

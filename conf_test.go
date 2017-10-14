package main

import "testing"

func TestInitConf(t *testing.T) {
	InitConf("")
	conf, _ := GetConf()
	t.Log(conf)

}

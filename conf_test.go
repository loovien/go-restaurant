package main

import "testing"

func TestConf(t *testing.T)  {
	InitConf("conf/app.toml")
	conf, init := GetConf()
	t.Log(conf, init)
}

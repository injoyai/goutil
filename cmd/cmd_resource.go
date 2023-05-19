package main

import _ "embed"

//go:embed \resource\upx.exe
var upx []byte

//go:embed \resource\rsrc.exe
var rsrc []byte

//go:embed \resource\swag.exe
var swag []byte

//go:embed \resource\influxdb.zip
var influxdb []byte

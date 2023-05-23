package main

import _ "embed"

//============================install============================

//go:embed \in.exe
var in []byte

//go:embed \resource\upx.exe
var upx []byte

//go:embed \resource\rsrc.exe
var rsrc []byte

//go:embed \resource\swag.exe
var swag []byte

//============================demo============================

//go:embed \resource\build.sh
var build []byte

//go:embed \resource\Dockerfile
var dockerfile []byte

//go:embed \resource\service.service
var service []byte

//go:embed \resource\install_minio.sh
var installMinio []byte

//go:embed \resource\install_nodered.sh
var installNodeRed []byte

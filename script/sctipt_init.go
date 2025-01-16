package script

import (
	"context"
)

func WithObject(i Client) {
	i.Set("global", NewGlobal())
	i.Set("net", new(Net))
	i.Set("ios", new(Ios))
	i.Set("logs", NewLogs())
	i.Set("http", NewHTTP())
	i.Set("conv", NewConv())
	i.Set("os", NewOS())
	i.Set("cfg", NewCfg())
	i.Set("ctx", context.Background())
	i.Set("bytes", new(Bytes))
	i.Set("mux", new(Mux))
	i.Set("in", NewIn())
	i.Set("maps", new(Maps))
	i.Set("time", NewTime())
}

func WithFunc(i Client) {
	i.Set("go", funcGo)
	i.Set("print", funcPrint)
	i.Set("println", funcPrintln)
	i.Set("printf", funcPrintf)
	i.Set("sprintf", funcSprintf)
	i.Set("shell", funcShell)
	i.Set("speak", funcSpeak)
	i.Set("sleep", funcSleep)
	i.Set("len", funcLen)
	i.Set("int", funcToInt64)
	i.Set("intBytes", funcToInt64Bytes)
	i.Set("float", funcToFloat)
	i.Set("str", funcToString)
	i.Set("bool", funcToBool)
	i.Set("cut", funcCut)
	i.Set("bin", funcToBIN)
	i.Set("hex", funcToHex)
	i.Set("base64Encode", funcBase64Encode)
	i.Set("base64Decode", funcBase64Decode)
	i.Set("hexToBytes", funcHexToBytes)
	i.Set("hexDecodeString", funcHexToBytes)
	i.Set("hexToString", funcHexToString)
	i.Set("getJson", funcGetJson)
	i.Set("holdTime", funcHoldTime)
	i.Set("holdCount", funcHoldCount)
	i.Set("toInt", funcToInt)
	i.Set("toInt8", funcToInt8)
	i.Set("toInt16", funcToInt16)
	i.Set("toInt32", funcToInt32)
	i.Set("toInt64", funcToInt64)
	i.Set("toUint8", funcToUint8)
	i.Set("toUint16", funcToUint16)
	i.Set("toUint32", funcToUint32)
	i.Set("toUint64", funcToUint64)
	i.Set("toFloat", funcToFloat)
	i.Set("toFloat32", funcToFloat32)
	i.Set("toFloat64", funcToFloat64)
	i.Set("toString", funcToString)
	i.Set("toBool", funcToBool)
	i.Set("toBin", funcToBIN)
	i.Set("toHex", funcToHex)
	i.Set("getByte", funcGetByte)
	i.Set("sum", funcSum)
	i.Set("addInt", funcAddInt)
	i.Set("reverse", funcReverse)
	i.Set("rand", funcRand)
	i.Set("syncDate", funcSyncDate)
	i.Set("setDate", funcSetDate)
	i.Set("crc16", funcCrc16)
	i.Set("ping", funcPing)
}

package script

func WithBaseFunc(i Client) {
	setFunc := i.Set
	setFunc("sprintf", funcSprintf)
	setFunc("shell", funcShell)
	setFunc("speak", funcSpeak)
	setFunc("sleep", funcSleep)
	setFunc("len", funcLen)
	setFunc("int", funcToInt64)
	setFunc("float", funcToFloat)
	setFunc("str", funcToString)
	setFunc("bool", funcToBool)
	setFunc("bin", funcToBIN)
	setFunc("hex", funcToHex)
	setFunc("base64Encode", funcBase64Encode)
	setFunc("base64Decode", funcBase64Decode)
	setFunc("hexToBytes", funcHexToBytes)
	setFunc("hexToString", funcHexToString)
	setFunc("getJson", funcGetJson)
	setFunc("holdTime", funcHoldTime)
	setFunc("holdCount", funcHoldCount)
	setFunc("setCache", funcSetCache)
	setFunc("getCache", funcGetCache)
	setFunc("delCache", funcDelCache)
	setFunc("toInt", funcToInt)
	setFunc("toInt8", funcToInt8)
	setFunc("toInt16", funcToInt16)
	setFunc("toInt32", funcToInt32)
	setFunc("toInt64", funcToInt64)
	setFunc("toUint8", funcToUint8)
	setFunc("toUint16", funcToUint16)
	setFunc("toUint32", funcToUint32)
	setFunc("toUint64", funcToUint64)
	setFunc("toFloat", funcToFloat)
	setFunc("toFloat32", funcToFloat32)
	setFunc("toFloat64", funcToFloat64)
	setFunc("toString", funcToString)
	setFunc("toBool", funcToBool)
	setFunc("toBin", funcToBIN)
	setFunc("toHex", funcToHex)
	setFunc("getByte", funcGetByte)
	setFunc("http", funcHTTP)
	setFunc("udp", funcUDP)
	setFunc("sum", funcSum)
	setFunc("addInt", funcAddInt)
	setFunc("reverse", funcReverse)
	setFunc("rand", funcRand)
	setFunc("syncDate", funcSyncDate)
	setFunc("setDate", funcSetDate)
}

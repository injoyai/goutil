package router

import path2 "path"

//整理路径
func cleanPath(path string) string {
	if len(path) > 0 {
		path = path2.Clean(path)
		if len(path) > 0 && path[:1] != "/" {
			path = "/" + path
		}
	}
	return path
}

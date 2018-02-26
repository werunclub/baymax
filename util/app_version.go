package util

import (
	"strconv"
	"strings"
)

const APPVersionHeader = "X-App-Version"

// MAJOR.MINOR.PATCH
// 1.0.1
type APPVersion struct {
	Major int // 1
	Minor int // 2
	Patch int // 3
}

func NewAPPVersion(s string) APPVersion {
	versionList := strings.Split(s, ".")

	if len(versionList) >= 2 {
		major, _ := strconv.Atoi(versionList[0])
		minor, _ := strconv.Atoi(versionList[1])
		patch := 0

		if len(versionList) >= 3 {
			patch, _ = strconv.Atoi(versionList[2])
		}

		return APPVersion{Major: major, Minor: minor, Patch: patch}
	}

	if len(versionList) == 2 {
		major, _ := strconv.Atoi(versionList[0])
		minor, _ := strconv.Atoi(versionList[1])
		return APPVersion{Major: major, Minor: minor, Patch: 0}
	}

	return APPVersion{}
}

// 是否比 v 更新的版本
func (version *APPVersion) Newer(v APPVersion) bool {
	if version.Major < v.Major {
		return false
	}

	if version.Major > v.Major {
		return true
	}

	if version.Major == v.Major && version.Minor < v.Minor {
		return false
	}

	if version.Major == v.Major && version.Minor > v.Minor {
		return true
	}

	if version.Major == v.Major && version.Minor == v.Minor && version.Patch < v.Patch {
		return false
	}

	if version.Major == v.Major && version.Minor == v.Minor && version.Patch > v.Patch {
		return true
	}

	return false
}

// version 比指定版本 (v) 更旧的一个版本
func (version *APPVersion) Older(v APPVersion) bool {

	if version.Major < v.Major {
		return true
	} else if version.Major > v.Major {
		return false
	} else {
		// Major Equal
		if version.Minor < v.Minor {
			return true
		} else if version.Minor > v.Minor {
			return false
		} else {
			// Minor Equal
			if version.Patch < v.Patch {
				return true
			} else if version.Patch > v.Patch {
				return false
			} else {
				// Patch Equal
				return false
			}
		}
	}

}

// 版本号相同
func (version *APPVersion) Equal(v APPVersion) bool {
	return version.Major == v.Major && version.Minor == v.Minor && version.Patch == v.Patch
}

// version 版本号 >= 指定版本
func (version *APPVersion) NewOrEqual(v APPVersion) bool {
	return version.Equal(v) || version.Newer(v)
}

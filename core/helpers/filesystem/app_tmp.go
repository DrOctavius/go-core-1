package filesystem

import (
	"github.com/kyaxcorp/go-core/core/config"
	"github.com/kyaxcorp/go-core/core/helpers/err/define"
	fsPath "github.com/kyaxcorp/go-core/core/helpers/filesystem/path"
	"strings"
)

func GetAppTmpPath(paths ...string) (string, error) {
	var _err error = nil
	itemPath := config.GetConfig().Application.TempPath
	itemPath, _err = fsPath.GenRealPath(itemPath, true)

	if _err != nil {
		//log.Println(err)
		return "", _err
	}

	if len(paths) > 0 {
		itemPath = itemPath + strings.Join(paths, DirSeparator())
	}

	if !Exists(itemPath) {
		if !MkDir(itemPath) {
			return "", define.Err(0, "failed to create path -> ", itemPath)
		}
	}

	return itemPath, nil
}

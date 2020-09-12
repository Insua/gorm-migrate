package gorm_migrate

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

func Create(baseDir, packageName, fileName string) error {
	if !gregex.IsMatchString("^[a-zA-Z][0-9a-zA-Z_]{0,}[0-9a-zA-Z_]{0,}$", fileName) {
		return errors.New("filename should match go variable standard")
	}
	if gstr.HasSuffix(fileName, "_test") {
		return errors.New("filename end with _test will be used as test file")
	}
	timeStamp := fmt.Sprintf("%s%v", time.Now().Format("20060102150405"), time.Now().Nanosecond())
	errDir := gfile.Mkdir(baseDir)
	if errDir != nil {
		return errDir
	}

	f, _ := os.Create(path.Join(baseDir, timeStamp+"_"+fileName+".go"))

	defer func() {
		_ = f.Close()
	}()

	tpl := template.Must(template.New("file").Parse(fileContent))

	type Tpl struct {
		Package    string
		StructName string
		FileName   string
	}

	t := Tpl{
		Package:    packageName,
		StructName: gstr.CamelCase("Migration" + timeStamp + fileName),
		FileName:   timeStamp + "_" + fileName,
	}

	errExe := tpl.Execute(f, t)

	return errExe

}

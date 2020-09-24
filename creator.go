package migrate

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"github.com/gogf/gf/util/gconv"

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
	timeStamp := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), nanoString())

	if err := makeBasedir(baseDir); err != nil {
		return err
	}

	if err := createMigrationsVariableFile(baseDir, packageName); err != nil {
		return err
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

func nanoString() string {
	nano := gconv.String(time.Now().Nanosecond())
	if len(nano) > 4 {
		return nano[0:4]
	} else {
		str := nano[0:]
		for i := 0; i < 4-len(nano); i++ {
			str += "0"
		}
		return str
	}
}

func makeBasedir(baseDir string) error {
	if exist := gfile.Exists(baseDir); !exist {
		errDir := gfile.Mkdir(baseDir)
		return errDir
	}
	return nil
}

func createMigrationsVariableFile(baseDir, packageName string) error {
	exist := gfile.Exists(path.Join(baseDir, "migration.go"))
	if exist {
		return nil
	}
	f, _ := os.Create(path.Join(baseDir, "migration.go"))
	defer func() {
		_ = f.Close()
	}()

	tpl := template.Must(template.New("migration").Parse(migrationContent))

	type Tpl struct {
		Package string
	}

	t := Tpl{
		Package: packageName,
	}

	errExe := tpl.Execute(f, t)

	return errExe
}

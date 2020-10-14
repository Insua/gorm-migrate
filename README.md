# gorm migrate tool

### how to use

gf by example  

create a cmd.go file  
```go
package main

import (
	"gf-cms/app/command/migrate"

	"github.com/gookit/color"

	"github.com/gogf/gf/os/gcmd"
)

func main() {
	command := gcmd.GetArg(1)

	switch command {
	case "migrate":
		sub := gcmd.GetArg(2)
		switch sub {
		case "new":
			fileName := gcmd.GetArg(3)
			if len(fileName) == 0 {
				color.Error.Prompt("Doesn't have fileName")
			} else {
				migrate.New(fileName)
			}
		case "up":
			migrate.Up()
		case "down":
			migrate.Down()
		default:
			color.Warn.Prompt("Wrong Process method")
		}

	default:
		color.Warn.Prompt("Not correct command")
	}
}
```

in app/command/migrate/new.go  
```go
package migrate

import (
	migrate "github.com/Insua/gorm-migrate"
	"github.com/gookit/color"
)

func New(fileName string) {
	err := migrate.Create("database/migrations", "migrations", fileName)
	if err != nil {
		color.Red.Println(err)
	} else {
		color.Green.Println(fileName + " has be created")
	}
}
```

in app/command/migrate/migrate.go  
```go
package migrate

import (
	_ "gf-cms/boot"
	ms "gf-cms/database/migrations"
	"gf-cms/global"

	migrate "github.com/Insua/gorm-migrate"
	"github.com/gookit/color"
)

func Up() {
	migrations := ms.Migrations
	err := migrate.Up(global.DB, migrations)
	if err != nil {
		color.Red.Println(err)
	}
}

func Down() {
	migrations := ms.Migrations
	err := migrate.Down(global.DB, migrations)
	if err != nil {
		color.Red.Println(err)
	}
}
```

####  command
now you can use ```go run cmd.go migrate new file new``` to create a migrate file, and ```go run cmd.go migrate up``` to migrate ```go rum cmd.go migrate down ``` to rollback  

for project example, a [go frame](https://github.com/gogf/gf) cms called [gf-cms](https://github.com/Insua/gf-cms) which will specify the detail.


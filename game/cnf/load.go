package cnf

import (
	"encoding/csv"
	"os"

	"github.com/scp1513/san/game/g"
	"github.com/scp1513/san/game/opt"
)

// LoadAll 加载所有配置
func LoadAll() error {
	g.L.Println("Load all config.")
	if err := load(loadConstTable, "const"); err != nil {
		return err
	}
	return nil
}

//
func load(loader func(*csv.Reader) error, fn string) error {
	g.L.Printf("Load csv %s.", fn)
	file, err := os.Open(opt.Val().ResDir + "/" + fn + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	return loader(r)
}

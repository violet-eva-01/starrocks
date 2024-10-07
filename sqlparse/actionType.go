package sqlparse

import (
	"github.com/violet-eva-01/starrocks/util"
	"regexp"
	"strings"
)

const (
	extractTime = 96
	with        = 97
	from        = 98
	join        = 99
	Select      = 100
	Alter       = 101
	Insert      = 102
	Create      = 103
	Drop        = 104
	Update      = 105
	Delete      = 106
	Truncate    = 107
)

func (p *Parse) assignment(code int, tbl []string) {

	tbl = util.RemoveRepeatElementAndToLower(tbl)

	switch code {
	case extractTime:
		p.extractTime = tbl
	case with:
		p.withTablaName = tbl
	case from:
		p.fromTableName = tbl
	case join:
		p.joinTableName = tbl
	case Select:
		otherArr := append(p.extractTime, p.selectExcludeTables...)
		otherArr = append(otherArr, p.DeleteTableName...)
		tbl = util.RemoveCoincideElement(append(p.fromTableName, p.joinTableName...), append(p.withTablaName, otherArr...))
		p.SelectTableName = tbl
	case Insert:
		p.InsertTableName = tbl
	case Create:
		p.CreatTableName = tbl
	case Drop:
		tbl = util.RemoveIfExistsElement(tbl, p.dropExcludeTables)
		p.DropTableName = tbl
	case Alter:
		p.AlterTableName = tbl
	case Delete:
		p.DeleteTableName = tbl
	case Update:
		p.UpdateTableName = tbl
	case Truncate:
		p.TruncateTableName = tbl
	default:
		return
	}

	p.assign(code, tbl)
}

func (p *Parse) assign(code int, tbl []string) {
	switch code {
	case Create, Drop, Insert, Select, Alter, Delete, Update, Truncate:
		i := 0
		for _, table := range tbl {
			var t Table
			table = strings.ReplaceAll(table, " ", "")
			compile, _ := regexp.Compile("(`[^`]+`|[a-z0-9_]+)")
			allString := compile.FindAllString(table, -1)
			length := len(allString)
			switch length {
			case 1:
				if p.DbName == "" {
					p.ErrorTables = append(p.ErrorTables, table)
					continue
				}
				t.TableName = strings.ReplaceAll(allString[0], "`", "")
				t.DbName = strings.ReplaceAll(p.DbName, "`", "")
				t.Catalog = strings.ReplaceAll(p.Catalog, "`", "")
			case 2:
				t.TableName = strings.ReplaceAll(allString[1], "`", "")
				t.DbName = strings.ReplaceAll(allString[0], "`", "")
				t.Catalog = strings.ReplaceAll(p.Catalog, "`", "")
			case 3:
				t.TableName = strings.ReplaceAll(allString[2], "`", "")
				t.DbName = strings.ReplaceAll(allString[1], "`", "")
				t.Catalog = strings.ReplaceAll(allString[0], "`", "")
			default:
				p.ErrorTables = append(p.ErrorTables, table)
				continue
			}
			t.Action = code
			t.Index = i
			p.ParseTables = append(p.ParseTables, t)
			i += 1
		}
	default:
		return
	}
}

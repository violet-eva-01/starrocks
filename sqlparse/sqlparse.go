package sqlparse

import (
	"fmt"
	"github.com/violet-eva-01/starrocks/util"
	"regexp"
	"sort"
	"strings"
)

type Parse struct {
	Query             string
	Catalog           string
	DbName            string
	ParseTables       []Table
	SelectTableName   []string
	AlterTableName    []string
	InsertTableName   []string
	CreatTableName    []string
	DropTableName     []string
	DeleteTableName   []string
	UpdateTableName   []string
	TruncateTableName []string
	ErrorTables       []string
	withTablaName     []string
	fromTableName     []string
	joinTableName     []string
	extractTime       []string
	excludeTables     []string
	excludeSign       []string
}

type Table struct {
	Catalog   string
	DbName    string
	Action    int
	TableName string
	Index     int
}

type sqlParseRegexp struct {
	Reg *regexp.Regexp
	New string
}

func newRegexp(reg string, new string) *sqlParseRegexp {
	compile := regexp.MustCompile(reg)
	return &sqlParseRegexp{
		Reg: compile,
		New: new,
	}
}

func findAllStrings(str string, regArr ...*regexp.Regexp) (result []string) {
	for _, reg := range regArr {
		findAllString := reg.FindAllString(str, -1)
		for _, f := range findAllString {
			if len(f) > 0 {
				result = append(result, f)
			}
		}
	}
	return
}

func regexpReplaceAllStrings(strArr []string, regArr ...*sqlParseRegexp) (result []string) {

	for _, str := range strArr {
		for _, reg := range regArr {
			str = reg.Reg.ReplaceAllString(str, reg.New)
		}
		if len(str) > 0 {
			result = append(result, str)
		}
	}
	return
}

func NewParse(query string, catalog string, dbName string, defaultCatalog string) *Parse {

	if len(strings.ReplaceAll(catalog, " ", "")) < 1 {
		if len(defaultCatalog) > 0 {
			catalog = defaultCatalog
		} else {
			catalog = "default_catalog"
		}
	}

	return &Parse{
		Query:   query,
		Catalog: catalog,
		DbName:  dbName,
	}
}

func (p *Parse) QueryClearAnnotation(isClean bool) {

	var (
		tmpStrArr   []string
		finalStrArr []string
	)

	replaceRegexp1 := regexp.MustCompile("(\\\\n|/\\*([^*]|\\*[^/])*\\*/)")
	tmpQuery := replaceRegexp1.ReplaceAllString(p.Query, "\n")

	replaceRegexp2 := newRegexp(`'((?:\\.|[^\\'])*)'`, " ")
	replaceRegexp3 := newRegexp(`"((?:\\.|[^\\"])*)"`, " ")
	replaceRegexp4 := newRegexp("--.*$", " ")

	if isClean {
		tmpStrArr = regexpReplaceAllStrings(strings.Split(tmpQuery, "\n"), replaceRegexp2, replaceRegexp3, replaceRegexp4)
	} else {
		tmpStrArr = regexpReplaceAllStrings(strings.Split(tmpQuery, "\n"), replaceRegexp4)
	}

	for _, str := range tmpStrArr {
		if len(strings.TrimSpace(str)) > 0 {
			finalStrArr = append(finalStrArr, str)
		}
	}

	if isClean {
		finalStrArr = regexpReplaceAllStrings([]string{strings.Join(finalStrArr, "\n")}, replaceRegexp2, replaceRegexp3)
	}

	p.Query = strings.Join(finalStrArr, "\n")
}

func (p *Parse) GetCatalogDB() {
	queryArr := strings.Split(p.Query, ";")
	if len(queryArr) < 2 {
		return
	}
	for _, query := range queryArr {
		p.getSet(query)
		p.getUse(query)
	}
}

func (p *Parse) getSet(str string) {

	parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)set\\s+catalog(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)\\s*")
	result := findAllStrings(str, parseFindRegexp)
	if len(result) <= 0 {
		return
	}
	parseReplaceRegexp1 := newRegexp("(?i)((^|\\s+|\\\\n)set\\s+catalog\\s+|\\s*)", "")
	parseReplaceRegexp2 := newRegexp("(?i)(^|\\s+|\\\\n)set\\s+catalog`", "`")
	tmpStrArr := regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2)
	if len(tmpStrArr) <= 0 {
		return
	}
	p.Catalog = strings.ToLower(strings.ReplaceAll(tmpStrArr[0], "`", ""))

}

func (p *Parse) getUse(str string) {
	parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)use(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))?\\s*")
	result := findAllStrings(str, parseFindRegexp)
	if len(result) <= 0 {
		return
	}
	parseReplaceRegexp1 := newRegexp("(?i)((^|\\s+|\\\\n)use\\s+|\\s*)", "")
	parseReplaceRegexp2 := newRegexp("(?i)(^|\\s+|\\\\n)use`", "`")
	tmpStrArr := regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2)
	if len(tmpStrArr) <= 0 {
		return
	}
	catalogDB := strings.ReplaceAll(tmpStrArr[0], "`", "")
	strArr := strings.Split(strings.ToLower(catalogDB), ".")
	switch len(strArr) {
	case 1:
		p.DbName = strArr[0]
	case 2:
		p.Catalog = strArr[0]
		p.DbName = strArr[1]
	default:
		return
	}
}

// getTableNames
// @Description
// @param action
// @return error
func (p *Parse) getTableNames(action int, isExists bool) {

	var (
		tableNames []string
	)

	switch action {
	case extractTime:
		parseFindRegexp := regexp.MustCompile("(?i)extract\\s*\\([^)]+from(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)extract\\s*\\([^)]+from\\s+", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case from:
		parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)from(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp1 := newRegexp("(?i)((^|\\s+|\\\\n)from\\s+|\\s*)", "")
		parseReplaceRegexp2 := newRegexp("(?i)(^|\\s+|\\\\n)from`", "`")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2)
	case with:
		parseFindRegexp := regexp.MustCompile("(?i)(with(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\([^)]+\\))?\\s+as\\s*\\(|,\\s*([a-z0-9_\\p{L}]+|`[^`]+`)(\\s*\\([^)]+\\))?\\s+as\\s*\\()")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)(with\\s+|,\\s*|(\\s*\\([^)]+\\))?\\s+as\\s*\\()", "")
		parseReplaceRegexp2 := newRegexp("(?i)with`", "`")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp, parseReplaceRegexp2)
	case Insert:
		parseFindRegexp := regexp.MustCompile("(?i)insert\\s+(into|overwrite)(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)insert\\s+(into|overwrite)\\s+", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Drop:
		parseFindRegexp := regexp.MustCompile("(?i)drop\\s+(temporary\\s+)?(table|view|materialized\\s+view)+(\\s+if\\s+exists)?(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)drop\\s+(temporary\\s+)?(table|view|materialized\\s+view)+(\\s+if\\s+exists)?\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Create:
		parseFindRegexp := regexp.MustCompile("(?i)create\\s+(table|view|materialized\\s+view)+(\\s+if\\s+not\\s+exists)?(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)create\\s+(table|view|materialized\\s+view)+(\\s+if\\s+not\\s+exists)?\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case join:
		tableName := "(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))*(\\s*as)?(\\s*[a-z0-9]+)?"
		parseFindRegexp1 := regexp.MustCompile(fmt.Sprintf("(?i)from%s(\\s*,\\s*%s)*", tableName, tableName))
		parseReplaceRegexp1 := newRegexp(fmt.Sprintf("(?i)from%s", tableName), "")
		parseReplaceRegexp2 := newRegexp("\\s*,\\s*", ",")
		parseReplaceRegexp3 := newRegexp("\\s*\\.\\s*", ".")
		parseFindRegexp2 := regexp.MustCompile("(?i)([a-z0-9_\\p{L}]+|`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp1)
		tmpTables := regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2, parseReplaceRegexp3)
		for _, tmpTable := range tmpTables {
			for _, tmpTBL := range strings.Split(tmpTable, ",") {
				if strings.ReplaceAll(tmpTBL, " ", "") == "" {
					continue
				}
				tmpTableNames := findAllStrings(tmpTBL, parseFindRegexp2)
				if len(tmpTableNames) > 0 {
					tableNames = append(tableNames, tmpTableNames[0])
				}
			}
		}
		parseFindRegexp3 := regexp.MustCompile("(?i)(^|\\s+|\\\\n)join(\\s+[a-z0-9_\\p{L}]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_\\p{L}]+|`[^`]+`))*")
		result1 := findAllStrings(p.Query, parseFindRegexp3)
		parseReplaceRegexp5 := newRegexp("(?i)(^|\\s+|\\\\n)join\\s+", "")
		parseReplaceRegexp6 := newRegexp("(?i)(^|\\s+|\\\\n)join`", "`")
		tmpTableNames := regexpReplaceAllStrings(result1, parseReplaceRegexp5, parseReplaceRegexp6, parseReplaceRegexp3)
		tableNames = append(tableNames, tmpTableNames...)
	case Select:
	case Alter:
		parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)alter\\s+(table|view|materialized\\s+view)+(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)(^|\\s+|\\\\n)alter\\s+(table|view|materialized\\s+view)+\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Delete:
		parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)delete\\s+from(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)(^|\\s+|\\\\n)delete\\s+from\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Update:
		parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)update(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp1 := newRegexp("(?i)((^|\\s+|\\\\n)update\\s+|\\s*)", "")
		parseReplaceRegexp2 := newRegexp("(?i)(^|\\s+|\\\\n)update`", "`")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2)
	case Truncate:
		parseFindRegexp := regexp.MustCompile("(?i)(^|\\s+|\\\\n)truncate\\s+table(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)(^|\\s+|\\\\n)truncate\\s+table\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	default:
		return
	}

	p.assignment(action, tableNames, isExists)

	return
}

func (p *Parse) AddExcludeTables(excludeTables ...string) {
	p.excludeTables = util.RemoveRepeatElementAndToLower(append(p.excludeTables, addSpecialCharacters(excludeTables)...))
	sort.Strings(p.excludeTables)
}

func (p *Parse) InitExcludeTables(excludeTables ...string) {
	p.excludeTables = []string{"`dual`", "`unnest`", "`files`", "`generate_series`"}
	if len(excludeTables) > 0 {
		p.AddExcludeTables(excludeTables...)
	} else {
		sort.Strings(p.excludeTables)
	}
}

func (p *Parse) AddExcludeSign(excludeSign ...string) {
	p.excludeSign = util.RemoveRepeatElementAndToLower(append(p.excludeTables, addSpecialCharacters(excludeSign)...))
	sort.Strings(p.excludeSign)
}

func (p *Parse) InitExcludeSign(excludeSign ...string) {
	p.excludeSign = []string{"#tableau_"}
	if len(excludeSign) > 0 {
		p.AddExcludeSign(excludeSign...)
	} else {
		sort.Strings(p.excludeSign)
	}
}

func (p *Parse) GetSelectTables(isExists bool) {
	p.getTableNames(extractTime, isExists)
	p.getTableNames(with, isExists)
	p.getTableNames(Delete, isExists)
	p.getTableNames(from, isExists)
	p.getTableNames(join, isExists)
	p.getTableNames(Select, isExists)
}

func (p *Parse) GetCreateTables(isExists bool) {
	p.getTableNames(Create, isExists)
}

func (p *Parse) GetDropTables(isExists bool) {
	p.getTableNames(Drop, isExists)
}

func (p *Parse) GetInsertTables(isExists bool) {
	p.getTableNames(Insert, isExists)
}

func (p *Parse) GetUpdateTables(isExists bool) {
	p.getTableNames(Update, isExists)
}

func (p *Parse) GetDeleteTables(isExists bool) {
	p.getTableNames(Delete, isExists)
}

func (p *Parse) GetTruncateTables(isExists bool) {
	p.getTableNames(Truncate, isExists)
}

func (p *Parse) GetAlterTables(isExists bool) {
	p.getTableNames(Alter, isExists)
}

func (p *Parse) InitAllUseTable() {
	p.QueryClearAnnotation(true)
	p.GetCatalogDB()
	p.GetSelectTables(false)
	p.GetAlterTables(false)
	p.GetCreateTables(false)
	p.GetDropTables(true)
	p.GetInsertTables(false)
	p.GetUpdateTables(false)
	p.GetTruncateTables(false)
}

func (p *Parse) DebugGetSelectTables() {
	fmt.Println("clean 后的query")
	fmt.Println(p.Query)
	fmt.Println("查询表名")
	for _, i := range p.fromTableName {
		fmt.Println("fromTableName : ", i)
	}
	for _, i := range p.joinTableName {
		fmt.Println("joinTableName : ", i)
	}
	fmt.Println("除外表名")
	for _, i := range p.extractTime {
		fmt.Println("extractTime : ", i)
	}
	for _, i := range p.withTablaName {
		fmt.Println("withTablaName : ", i)
	}
	for _, i := range p.DeleteTableName {
		fmt.Println("DeleteTableName : ", i)
	}
	fmt.Println("除外常量表名")
	for _, i := range p.excludeTables {
		fmt.Println("excludeTables : ", i)
	}
	fmt.Println("除外常量标志")
	for _, i := range p.excludeSign {
		fmt.Println("excludeSign : ", i)
	}
	fmt.Println("最终查询表名")
	for _, i := range p.SelectTableName {
		fmt.Println("SelectTableName : ", i)
	}
}

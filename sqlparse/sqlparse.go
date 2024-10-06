package sqlparse

import (
	"fmt"
	"github.com/violet-eva-01/starrocks/util"
	"regexp"
	"strings"
)

type Parse struct {
	Query           string
	Catalog         string
	DbName          string
	ParseTables     []Table
	SelectTableName []string
	InsertTableName []string
	CreatTableName  []string
	DropTableName   []string
	WithTablaName   []string
	FromTableName   []string
	JoinTableName   []string
	ExtractTime     []string
	ErrorTables     []string
	DirtyData       []string
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

func (p *Parse) StmtClearAnnotation() {

	var (
		finalStrArr []string
	)

	//replaceRegexp, _ := regexp.Compile(`(--.*$|/\*([^*]|\*[^/])*\*/|\n)`)
	replaceRegexp1, _ := regexp.Compile("(\\\\n|/\\*([^*]|\\*[^/])*\\*/)")

	tmpQuery := replaceRegexp1.ReplaceAllString(p.Query, "")

	replaceRegexp2, _ := regexp.Compile("--.*$")

	for _, tmpStr := range strings.Split(tmpQuery, "\n") {
		tmpRRStr := replaceRegexp2.ReplaceAllString(tmpStr, "")
		if len(tmpRRStr) > 0 {
			finalStrArr = append(finalStrArr, tmpRRStr)
		}
	}

	p.Query = strings.Join(finalStrArr, "\n")

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

func newRegexp(reg string, new string) *sqlParseRegexp {
	compile, _ := regexp.Compile(reg)
	return &sqlParseRegexp{
		Reg: compile,
		New: new,
	}
}

// getTableNames
// @Description
// @param action
// @return error
func (p *Parse) getTableNames(action int) {

	var (
		tableNames []string
	)

	switch action {
	case ExtractTime:
		parseFindRegexp, _ := regexp.Compile("(?i)extract\\s*\\([^)]+from(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)extract\\s*\\([^)]+from\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case From:
		parseFindRegexp, _ := regexp.Compile("(?i)from(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp1 := newRegexp("(?i)(from\\s+|\\s*)", "")
		parseReplaceRegexp2 := newRegexp("(?i)from`", "`")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp1, parseReplaceRegexp2)
	case With:
		parseFindRegexp, _ := regexp.Compile("(?i)(with(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\([^)]+\\))?\\s+as\\s*\\(|,\\s*([a-z0-9_]+|`[^`]+`)(\\s*\\([^)]+\\))?\\s+as\\s*\\()")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)(with\\s+|,\\s*|(\\s*\\([^)]+\\))?\\s+as\\s*\\()", "")
		parseReplaceRegexp2 := newRegexp("(?i)with`", "`")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp, parseReplaceRegexp2)
	case Insert:
		parseFindRegexp, _ := regexp.Compile("(?i)insert\\s+(into|overwrite)(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)insert\\s+(into|overwrite)\\s+", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Drop:
		parseFindRegexp, _ := regexp.Compile("(?i)drop\\s+(table|view|materialized\\s+view)+(\\s+if\\s+exists)?(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)drop\\s+(table|view|materialized\\s+view)+(\\s+if\\s+exists)?\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Create:
		parseFindRegexp, _ := regexp.Compile("(?i)create\\s+(table|view|materialized\\s+view)+(\\s+if\\s+not\\s+exists)?(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result := findAllStrings(p.Query, parseFindRegexp)
		parseReplaceRegexp := newRegexp("(?i)create\\s+(table|view|materialized\\s+view)+(\\s+if\\s+not\\s+exists)?\\s*", "")
		tableNames = regexpReplaceAllStrings(result, parseReplaceRegexp)
	case Join:
		tableName := "(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*(\\s*as)?(\\s*[a-z0-9]+)?"
		parseFindRegexp1, _ := regexp.Compile(fmt.Sprintf("(?i)from%s(\\s*,\\s*%s)*", tableName, tableName))
		parseReplaceRegexp1 := newRegexp(fmt.Sprintf("(?i)from%s", tableName), "")
		parseReplaceRegexp2 := newRegexp("\\s*,\\s*", ",")
		parseReplaceRegexp3 := newRegexp("\\s*\\.\\s*", ".")
		parseFindRegexp2, _ := regexp.Compile("(?i)([a-z0-9_]+|`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
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
		parseFindRegexp3, _ := regexp.Compile("(?i)join(\\s+[a-z0-9_]+|\\s*`[^`]+`)(\\s*\\.\\s*([a-z0-9_]+|`[^`]+`))*")
		result1 := findAllStrings(p.Query, parseFindRegexp3)
		parseReplaceRegexp5 := newRegexp("(?i)join\\s+", "")
		parseReplaceRegexp6 := newRegexp("(?i)join`", "`")
		tmpTableNames := regexpReplaceAllStrings(result1, parseReplaceRegexp5, parseReplaceRegexp6, parseReplaceRegexp3)
		tableNames = append(tableNames, tmpTableNames...)

	case Select:
	default:
		return
	}

	p.assignment(action, tableNames)

	return
}

func (p *Parse) AddDirtyData(strArr []string) {
	p.DirtyData = util.RemoveRepeatElementAndToLower(append(p.DirtyData, strArr...))
}

func (p *Parse) GetSelectFromTables() {
	p.getTableNames(ExtractTime)
	p.getTableNames(With)
	p.getTableNames(From)
	p.getTableNames(Join)
	p.getTableNames(Select)
}

func (p *Parse) GetCreateTables() {
	p.getTableNames(Create)
}

func (p *Parse) GetDropTables() {
	p.getTableNames(Drop)
}

func (p *Parse) GetInsertTables() {
	p.getTableNames(Insert)
}

func (p *Parse) initDirtyData() {
	p.DirtyData = []string{"dual", "unnest"}
}

func (p *Parse) InitAllUseTable(isInitDirtyData bool) {
	if isInitDirtyData {
		p.initDirtyData()
	}
	p.StmtClearAnnotation()
	p.GetSelectFromTables()
	p.GetCreateTables()
	p.GetDropTables()
	p.GetInsertTables()
}

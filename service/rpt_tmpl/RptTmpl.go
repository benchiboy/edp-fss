package rpttmpl

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	
	
)

const (
	SQL_NEWDB	= "NewDB  ===>"
	SQL_INSERT  = "Insert ===>"
	SQL_UPDATE  = "Update ===>"
	SQL_SELECT  = "Select ===>"
	SQL_DELETE  = "Delete ===>"
	SQL_ELAPSED = "Elapsed===>"
	SQL_ERROR   = "Error  ===>"
	SQL_TITLE   = "===================================="
	DEBUG       = 1
	INFO        = 2
)

type Search struct {
	
	Id	int64	`json:"id"`
	TmplNo	string	`json:"tmpl_no"`
	TmplName	string	`json:"tmpl_name"`
	TmplType	string	`json:"tmpl_type"`
	Content	string	`json:"content"`
	Status	string	`json:"status"`
	StartDate	int64	`json:"start_date"`
	ExpireDate	int64	`json:"expire_date"`
	UpdateDate	int64	`json:"update_date"`
	InsertDate	int64	`json:"insert_date"`
	Version	int64	`json:"version"`
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
	ExtraWhere   string `json:"extra_where"`
	SortFld  string `json:"sort_fld"`
}

type RptTmplList struct {
	DB      *sql.DB
	Level   int
	Total   int      `json:"total"`
	RptTmpls []RptTmpl `json:"RptTmpl"`
}

type RptTmpl struct {
	
	Id	int64	`json:"id"`
	TmplNo	string	`json:"tmpl_no"`
	TmplName	string	`json:"tmpl_name"`
	TmplType	string	`json:"tmpl_type"`
	Content	string	`json:"content"`
	Status	string	`json:"status"`
	StartDate	int64	`json:"start_date"`
	ExpireDate	int64	`json:"expire_date"`
	UpdateDate	int64	`json:"update_date"`
	InsertDate	int64	`json:"insert_date"`
	Version	int64	`json:"version"`
}


type Form struct {
	Form   RptTmpl `json:"RptTmpl"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *RptTmplList {
	if db==nil{
		log.Println(SQL_SELECT,"Database is nil")
		return nil
	}
	return &RptTmplList{DB: db, Total: 0, RptTmpls: make([]RptTmpl, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *RptTmplList {
	var err error
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println(SQL_SELECT,"Open database error:", err)
		return nil
	}
	if err = db.Ping(); err != nil {
		log.Println(SQL_SELECT,"Ping database error:", err)
		return nil
	}
	return &RptTmplList{DB: db, Total: 0, RptTmpls: make([]RptTmpl, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *RptTmplList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.TmplNo != "" {
		where += " and tmpl_no='" + s.TmplNo + "'"
	}	
	
			
	if s.TmplName != "" {
		where += " and tmpl_name='" + s.TmplName + "'"
	}	
	
			
	if s.TmplType != "" {
		where += " and tmpl_type='" + s.TmplType + "'"
	}	
	
			
	if s.Content != "" {
		where += " and content='" + s.Content + "'"
	}	
	
			
	if s.Status != "" {
		where += " and status='" + s.Status + "'"
	}	
	
	
	if s.StartDate != 0 {
		where += " and start_date=" + fmt.Sprintf("%d", s.StartDate)
	}			
	
	
	if s.ExpireDate != 0 {
		where += " and expire_date=" + fmt.Sprintf("%d", s.ExpireDate)
	}			
	
	
	if s.UpdateDate != 0 {
		where += " and update_date=" + fmt.Sprintf("%d", s.UpdateDate)
	}			
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select count(1) as total from rpt_tmpl   where 1=1 %s", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return 0, err
	}
	defer rows.Close()
	var total int
	for rows.Next() {
		rows.Scan(&total)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return total, nil
}

/*
	说明：根据主键查询符合条件的条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r RptTmplList) Get(s Search) (*RptTmpl, error) {
	var where string
	l := time.Now()
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.TmplNo != "" {
		where += " and tmpl_no='" + s.TmplNo + "'"
	}	
	
			
	if s.TmplName != "" {
		where += " and tmpl_name='" + s.TmplName + "'"
	}	
	
			
	if s.TmplType != "" {
		where += " and tmpl_type='" + s.TmplType + "'"
	}	
	
			
	if s.Content != "" {
		where += " and content='" + s.Content + "'"
	}	
	
			
	if s.Status != "" {
		where += " and status='" + s.Status + "'"
	}	
	
	
	if s.StartDate != 0 {
		where += " and start_date=" + fmt.Sprintf("%d", s.StartDate)
	}			
	
	
	if s.ExpireDate != 0 {
		where += " and expire_date=" + fmt.Sprintf("%d", s.ExpireDate)
	}			
	
	
	if s.UpdateDate != 0 {
		where += " and update_date=" + fmt.Sprintf("%d", s.UpdateDate)
	}			
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}
	
	qrySql := fmt.Sprintf("Select id,tmpl_no,tmpl_name,tmpl_type,content,status,start_date,expire_date,update_date,insert_date,version from rpt_tmpl where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p  RptTmpl
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err:=rows.Scan(&p.Id,&p.TmplNo,&p.TmplName,&p.TmplType,&p.Content,&p.Status,&p.StartDate,&p.ExpireDate,&p.UpdateDate,&p.InsertDate,&p.Version)
		if err != nil {
			log.Println(SQL_ERROR, err.Error())
			return nil, err
		}
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return &p, nil
}

/*
	说明：根据条件查询复核条件对象列表，支持分页查询
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象列表, 参数2：如果错误返回错误对象
*/

func (r *RptTmplList) GetList(s Search) ([]RptTmpl, error) {
	var where string
	l := time.Now()
	
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.TmplNo != "" {
		where += " and tmpl_no='" + s.TmplNo + "'"
	}	
	
			
	if s.TmplName != "" {
		where += " and tmpl_name='" + s.TmplName + "'"
	}	
	
			
	if s.TmplType != "" {
		where += " and tmpl_type='" + s.TmplType + "'"
	}	
	
			
	if s.Content != "" {
		where += " and content='" + s.Content + "'"
	}	
	
			
	if s.Status != "" {
		where += " and status='" + s.Status + "'"
	}	
	
	
	if s.StartDate != 0 {
		where += " and start_date=" + fmt.Sprintf("%d", s.StartDate)
	}			
	
	
	if s.ExpireDate != 0 {
		where += " and expire_date=" + fmt.Sprintf("%d", s.ExpireDate)
	}			
	
	
	if s.UpdateDate != 0 {
		where += " and update_date=" + fmt.Sprintf("%d", s.UpdateDate)
	}			
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	
	
	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	var qrySql string
	if s.PageSize==0 &&s.PageNo==0{
		qrySql = fmt.Sprintf("Select id,tmpl_no,tmpl_name,tmpl_type,content,status,start_date,expire_date,update_date,insert_date,version from rpt_tmpl where 1=1 %s", where)
	}else{
		qrySql = fmt.Sprintf("Select id,tmpl_no,tmpl_name,tmpl_type,content,status,start_date,expire_date,update_date,insert_date,version from rpt_tmpl where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
	}
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p RptTmpl
	for rows.Next() {
		rows.Scan(&p.Id,&p.TmplNo,&p.TmplName,&p.TmplType,&p.Content,&p.Status,&p.StartDate,&p.ExpireDate,&p.UpdateDate,&p.InsertDate,&p.Version)
		r.RptTmpls = append(r.RptTmpls, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.RptTmpls, nil
}


/*
	说明：根据主键查询符合条件的记录，并保持成MAP
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r *RptTmplList) GetExt(s Search) (map[string]string, error) {
	var where string
	l := time.Now()

	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.TmplNo != "" {
		where += " and tmpl_no='" + s.TmplNo + "'"
	}	
	
			
	if s.TmplName != "" {
		where += " and tmpl_name='" + s.TmplName + "'"
	}	
	
			
	if s.TmplType != "" {
		where += " and tmpl_type='" + s.TmplType + "'"
	}	
	
			
	if s.Content != "" {
		where += " and content='" + s.Content + "'"
	}	
	
			
	if s.Status != "" {
		where += " and status='" + s.Status + "'"
	}	
	
	
	if s.StartDate != 0 {
		where += " and start_date=" + fmt.Sprintf("%d", s.StartDate)
	}			
	
	
	if s.ExpireDate != 0 {
		where += " and expire_date=" + fmt.Sprintf("%d", s.ExpireDate)
	}			
	
	
	if s.UpdateDate != 0 {
		where += " and update_date=" + fmt.Sprintf("%d", s.UpdateDate)
	}			
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	qrySql := fmt.Sprintf("Select id,tmpl_no,tmpl_name,tmpl_type,content,status,start_date,expire_date,update_date,insert_date,version from rpt_tmpl where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()


	Columns, _ := rows.Columns()

	values := make([]sql.RawBytes, len(Columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err = rows.Scan(scanArgs...)
	}

	fldValMap := make(map[string]string)
	for k, v := range Columns {
		fldValMap[v] = string(values[k])
	}

	log.Println(SQL_ELAPSED, "==========>>>>>>>>>>>", fldValMap)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return fldValMap, nil

}

/*
	说明：插入对象到数据表中，这个方法要求对象的各个属性必须赋值
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r RptTmplList) Insert(p RptTmpl) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  rpt_tmpl(tmpl_no,tmpl_name,tmpl_type,content,status,start_date,expire_date,update_date,insert_date,version)  values(?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.TmplNo,p.TmplName,p.TmplType,p.Content,p.Status,p.StartDate,p.ExpireDate,p.UpdateDate,p.InsertDate,p.Version)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}


/*
	说明：插入对象到数据表中，这个方法会判读对象的各个属性，如果属性不为空，才加入插入列中；
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/


func (r RptTmplList) InsertEntity(p RptTmpl, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)
	
		
	if p.TmplNo != "" {
		colNames += "tmpl_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.TmplNo)
	}			
		
	if p.TmplName != "" {
		colNames += "tmpl_name,"
		colTags += "?,"
		valSlice = append(valSlice, p.TmplName)
	}			
		
	if p.TmplType != "" {
		colNames += "tmpl_type,"
		colTags += "?,"
		valSlice = append(valSlice, p.TmplType)
	}			
		
	if p.Content != "" {
		colNames += "content,"
		colTags += "?,"
		valSlice = append(valSlice, p.Content)
	}			
		
	if p.Status != "" {
		colNames += "status,"
		colTags += "?,"
		valSlice = append(valSlice, p.Status)
	}			
	
	if p.StartDate != 0 {
		colNames += "start_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.StartDate)
	}				
	
	if p.ExpireDate != 0 {
		colNames += "expire_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.ExpireDate)
	}				
	
	if p.UpdateDate != 0 {
		colNames += "update_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.UpdateDate)
	}				
	
	if p.InsertDate != 0 {
		colNames += "insert_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.InsertDate)
	}				
	
	if p.Version != 0 {
		colNames += "version,"
		colTags += "?,"
		valSlice = append(valSlice, p.Version)
	}				
	
	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")
	exeSql := fmt.Sprintf("Insert into  rpt_tmpl(%s)  values(%s)", colNames, colTags)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "Insert data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：插入一个MAP到数据表中；
	入参：m:插入的Map
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r RptTmplList) InsertMap(m map[string]interface{},tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + ","
		colTags += "?,"
		valSlice = append(valSlice, v)
	}
	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")

	exeSql := fmt.Sprintf("Insert into  rpt_tmpl(%s)  values(%s)", colNames, colTags)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "insert data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}



/*
	说明：插入对象到数据表中，这个方法会判读对象的各个属性，如果属性不为空，才加入插入列中；
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/


func (r RptTmplList) UpdataEntity(keyNo string,p RptTmpl,tr *sql.Tx) error {
	l := time.Now()
	var colNames string
	valSlice := make([]interface{}, 0)
	
	
	if p.Id != 0 {
		colNames += "id=?,"
		valSlice = append(valSlice, p.Id)
	}				
		
	if p.TmplNo != "" {
		colNames += "tmpl_no=?,"
		
		valSlice = append(valSlice, p.TmplNo)
	}			
		
	if p.TmplName != "" {
		colNames += "tmpl_name=?,"
		
		valSlice = append(valSlice, p.TmplName)
	}			
		
	if p.TmplType != "" {
		colNames += "tmpl_type=?,"
		
		valSlice = append(valSlice, p.TmplType)
	}			
		
	if p.Content != "" {
		colNames += "content=?,"
		
		valSlice = append(valSlice, p.Content)
	}			
		
	if p.Status != "" {
		colNames += "status=?,"
		
		valSlice = append(valSlice, p.Status)
	}			
	
	if p.StartDate != 0 {
		colNames += "start_date=?,"
		valSlice = append(valSlice, p.StartDate)
	}				
	
	if p.ExpireDate != 0 {
		colNames += "expire_date=?,"
		valSlice = append(valSlice, p.ExpireDate)
	}				
	
	if p.UpdateDate != 0 {
		colNames += "update_date=?,"
		valSlice = append(valSlice, p.UpdateDate)
	}				
	
	if p.InsertDate != 0 {
		colNames += "insert_date=?,"
		valSlice = append(valSlice, p.InsertDate)
	}				
	
	if p.Version != 0 {
		colNames += "version=?,"
		valSlice = append(valSlice, p.Version)
	}				
	
	colNames = strings.TrimRight(colNames, ",")
	valSlice = append(valSlice, keyNo)

	exeSql := fmt.Sprintf("update  rpt_tmpl  set %s  where id=? ", colNames)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "Update data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：根据更新主键及更新Map值更新数据表；
	入参：keyNo:更新数据的关键条件，m:更新数据列的Map
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r RptTmplList) UpdateMap(keyNo string, m map[string]interface{},tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update rpt_tmpl set %s where id=?", colNames)
	if r.Level == DEBUG {
		log.Println(SQL_UPDATE, updateSql)
	}
	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(updateSql)
	} else {
		stmt, err = tr.Prepare(updateSql)
	}
	
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_UPDATE, "Update data error: %v\n", err)
		return err
	}
	defer stmt.Close()

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_UPDATE, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_UPDATE, "RowsAffected:", RowsAffected)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}


/*
	说明：根据主键删除一条数据；
	入参：keyNo:要删除的主键值
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r RptTmplList) Delete(keyNo string,tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  rpt_tmpl  where id=?")
	if r.Level == DEBUG {
		log.Println(SQL_UPDATE, delSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(delSql)
	} else {
		stmt, err = tr.Prepare(delSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	ret, err := stmt.Exec(keyNo)
	if err != nil {
		log.Println(SQL_DELETE, "Delete error: %v\n", err)
		return err
	}
	defer stmt.Close()

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_DELETE, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_DELETE, "RowsAffected:", RowsAffected)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}


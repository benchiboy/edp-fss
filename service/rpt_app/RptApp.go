package rptapp

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
	AppNo	string	`json:"app_no"`
	AppName	string	`json:"app_name"`
	AppKey	string	`json:"app_key"`
	AppIp	string	`json:"app_ip"`
	BackUrl	string	`json:"back_url"`
	Status	string	`json:"status"`
	StartDate	int64	`json:"start_date"`
	ExpireDate	int64	`json:"expire_date"`
	InsertDate	int64	`json:"insert_date"`
	UdpateDate	int64	`json:"udpate_date"`
	Version	int64	`json:"version"`
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
	ExtraWhere   string `json:"extra_where"`
	SortFld  string `json:"sort_fld"`
}

type RptAppList struct {
	DB      *sql.DB
	Level   int
	Total   int      `json:"total"`
	RptApps []RptApp `json:"RptApp"`
}

type RptApp struct {
	
	Id	int64	`json:"id"`
	AppNo	string	`json:"app_no"`
	AppName	string	`json:"app_name"`
	AppKey	string	`json:"app_key"`
	AppIp	string	`json:"app_ip"`
	BackUrl	string	`json:"back_url"`
	Status	string	`json:"status"`
	StartDate	int64	`json:"start_date"`
	ExpireDate	int64	`json:"expire_date"`
	InsertDate	int64	`json:"insert_date"`
	UdpateDate	int64	`json:"udpate_date"`
	Version	int64	`json:"version"`
}


type Form struct {
	Form   RptApp `json:"RptApp"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *RptAppList {
	if db==nil{
		log.Println(SQL_SELECT,"Database is nil")
		return nil
	}
	return &RptAppList{DB: db, Total: 0, RptApps: make([]RptApp, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *RptAppList {
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
	return &RptAppList{DB: db, Total: 0, RptApps: make([]RptApp, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *RptAppList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.AppNo != "" {
		where += " and app_no='" + s.AppNo + "'"
	}	
	
			
	if s.AppName != "" {
		where += " and app_name='" + s.AppName + "'"
	}	
	
			
	if s.AppKey != "" {
		where += " and app_key='" + s.AppKey + "'"
	}	
	
			
	if s.AppIp != "" {
		where += " and app_ip='" + s.AppIp + "'"
	}	
	
			
	if s.BackUrl != "" {
		where += " and back_url='" + s.BackUrl + "'"
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
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.UdpateDate != 0 {
		where += " and udpate_date=" + fmt.Sprintf("%d", s.UdpateDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select count(1) as total from rpt_app   where 1=1 %s", where)
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

func (r RptAppList) Get(s Search) (*RptApp, error) {
	var where string
	l := time.Now()
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.AppNo != "" {
		where += " and app_no='" + s.AppNo + "'"
	}	
	
			
	if s.AppName != "" {
		where += " and app_name='" + s.AppName + "'"
	}	
	
			
	if s.AppKey != "" {
		where += " and app_key='" + s.AppKey + "'"
	}	
	
			
	if s.AppIp != "" {
		where += " and app_ip='" + s.AppIp + "'"
	}	
	
			
	if s.BackUrl != "" {
		where += " and back_url='" + s.BackUrl + "'"
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
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.UdpateDate != 0 {
		where += " and udpate_date=" + fmt.Sprintf("%d", s.UdpateDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}
	
	qrySql := fmt.Sprintf("Select id,app_no,app_name,app_key,app_ip,back_url,status,start_date,expire_date,insert_date,udpate_date,version from rpt_app where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p  RptApp
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err:=rows.Scan(&p.Id,&p.AppNo,&p.AppName,&p.AppKey,&p.AppIp,&p.BackUrl,&p.Status,&p.StartDate,&p.ExpireDate,&p.InsertDate,&p.UdpateDate,&p.Version)
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

func (r *RptAppList) GetList(s Search) ([]RptApp, error) {
	var where string
	l := time.Now()
	
	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.AppNo != "" {
		where += " and app_no='" + s.AppNo + "'"
	}	
	
			
	if s.AppName != "" {
		where += " and app_name='" + s.AppName + "'"
	}	
	
			
	if s.AppKey != "" {
		where += " and app_key='" + s.AppKey + "'"
	}	
	
			
	if s.AppIp != "" {
		where += " and app_ip='" + s.AppIp + "'"
	}	
	
			
	if s.BackUrl != "" {
		where += " and back_url='" + s.BackUrl + "'"
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
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.UdpateDate != 0 {
		where += " and udpate_date=" + fmt.Sprintf("%d", s.UdpateDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	
	
	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	var qrySql string
	if s.PageSize==0 &&s.PageNo==0{
		qrySql = fmt.Sprintf("Select id,app_no,app_name,app_key,app_ip,back_url,status,start_date,expire_date,insert_date,udpate_date,version from rpt_app where 1=1 %s", where)
	}else{
		qrySql = fmt.Sprintf("Select id,app_no,app_name,app_key,app_ip,back_url,status,start_date,expire_date,insert_date,udpate_date,version from rpt_app where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
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

	var p RptApp
	for rows.Next() {
		rows.Scan(&p.Id,&p.AppNo,&p.AppName,&p.AppKey,&p.AppIp,&p.BackUrl,&p.Status,&p.StartDate,&p.ExpireDate,&p.InsertDate,&p.UdpateDate,&p.Version)
		r.RptApps = append(r.RptApps, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.RptApps, nil
}


/*
	说明：根据主键查询符合条件的记录，并保持成MAP
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r *RptAppList) GetExt(s Search) (map[string]string, error) {
	var where string
	l := time.Now()

	
	
	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}			
	
			
	if s.AppNo != "" {
		where += " and app_no='" + s.AppNo + "'"
	}	
	
			
	if s.AppName != "" {
		where += " and app_name='" + s.AppName + "'"
	}	
	
			
	if s.AppKey != "" {
		where += " and app_key='" + s.AppKey + "'"
	}	
	
			
	if s.AppIp != "" {
		where += " and app_ip='" + s.AppIp + "'"
	}	
	
			
	if s.BackUrl != "" {
		where += " and back_url='" + s.BackUrl + "'"
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
	
	
	if s.InsertDate != 0 {
		where += " and insert_date=" + fmt.Sprintf("%d", s.InsertDate)
	}			
	
	
	if s.UdpateDate != 0 {
		where += " and udpate_date=" + fmt.Sprintf("%d", s.UdpateDate)
	}			
	
	
	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}			
	

	qrySql := fmt.Sprintf("Select id,app_no,app_name,app_key,app_ip,back_url,status,start_date,expire_date,insert_date,udpate_date,version from rpt_app where 1=1 %s ", where)
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

func (r RptAppList) Insert(p RptApp) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  rpt_app(app_no,app_name,app_key,app_ip,back_url,status,start_date,expire_date,insert_date,udpate_date,version)  values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.AppNo,p.AppName,p.AppKey,p.AppIp,p.BackUrl,p.Status,p.StartDate,p.ExpireDate,p.InsertDate,p.UdpateDate,p.Version)
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


func (r RptAppList) InsertEntity(p RptApp, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)
	
		
	if p.AppNo != "" {
		colNames += "app_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.AppNo)
	}			
		
	if p.AppName != "" {
		colNames += "app_name,"
		colTags += "?,"
		valSlice = append(valSlice, p.AppName)
	}			
		
	if p.AppKey != "" {
		colNames += "app_key,"
		colTags += "?,"
		valSlice = append(valSlice, p.AppKey)
	}			
		
	if p.AppIp != "" {
		colNames += "app_ip,"
		colTags += "?,"
		valSlice = append(valSlice, p.AppIp)
	}			
		
	if p.BackUrl != "" {
		colNames += "back_url,"
		colTags += "?,"
		valSlice = append(valSlice, p.BackUrl)
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
	
	if p.InsertDate != 0 {
		colNames += "insert_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.InsertDate)
	}				
	
	if p.UdpateDate != 0 {
		colNames += "udpate_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.UdpateDate)
	}				
	
	if p.Version != 0 {
		colNames += "version,"
		colTags += "?,"
		valSlice = append(valSlice, p.Version)
	}				
	
	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")
	exeSql := fmt.Sprintf("Insert into  rpt_app(%s)  values(%s)", colNames, colTags)
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

func (r RptAppList) InsertMap(m map[string]interface{},tr *sql.Tx) error {
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

	exeSql := fmt.Sprintf("Insert into  rpt_app(%s)  values(%s)", colNames, colTags)
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


func (r RptAppList) UpdataEntity(keyNo string,p RptApp,tr *sql.Tx) error {
	l := time.Now()
	var colNames string
	valSlice := make([]interface{}, 0)
	
	
	if p.Id != 0 {
		colNames += "id=?,"
		valSlice = append(valSlice, p.Id)
	}				
		
	if p.AppNo != "" {
		colNames += "app_no=?,"
		
		valSlice = append(valSlice, p.AppNo)
	}			
		
	if p.AppName != "" {
		colNames += "app_name=?,"
		
		valSlice = append(valSlice, p.AppName)
	}			
		
	if p.AppKey != "" {
		colNames += "app_key=?,"
		
		valSlice = append(valSlice, p.AppKey)
	}			
		
	if p.AppIp != "" {
		colNames += "app_ip=?,"
		
		valSlice = append(valSlice, p.AppIp)
	}			
		
	if p.BackUrl != "" {
		colNames += "back_url=?,"
		
		valSlice = append(valSlice, p.BackUrl)
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
	
	if p.InsertDate != 0 {
		colNames += "insert_date=?,"
		valSlice = append(valSlice, p.InsertDate)
	}				
	
	if p.UdpateDate != 0 {
		colNames += "udpate_date=?,"
		valSlice = append(valSlice, p.UdpateDate)
	}				
	
	if p.Version != 0 {
		colNames += "version=?,"
		valSlice = append(valSlice, p.Version)
	}				
	
	colNames = strings.TrimRight(colNames, ",")
	valSlice = append(valSlice, keyNo)

	exeSql := fmt.Sprintf("update  rpt_app  set %s  where id=? ", colNames)
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

func (r RptAppList) UpdateMap(keyNo string, m map[string]interface{},tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update rpt_app set %s where id=?", colNames)
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

func (r RptAppList) Delete(keyNo string,tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  rpt_app  where id=?")
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


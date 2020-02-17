package rpt

import (
	"bytes"
	"edp-fss/common"
	"edp-fss/service/dbcomm"
	"edp-fss/service/rpt_app"
	"edp-fss/service/rpt_order"
	"edp-fss/service/rpt_tmpl"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
)

/*
	Pdf模板-章节定义
*/

type Chapter_Define struct {
	Title       string  `json:"title"`
	TitleSize   float64 `json:"title_size"`
	Content     string  `json:"content"`
	ContentSize float64 `json:"content_size"`

	Table          string    `json:"table"`
	TableCols      []float64 `json:"table_cols"`
	TableTitleSize float64   `json:"table_title_size"`
	TableDataSize  float64   `json:"table_data_size"`

	Chart          string  `json:"chart"`
	ChartTitleSize float64 `json:"chart_title_size"`
	ChartDataSize  float64 `json:"chart_data_size"`
}

/*
	Pdf模板-签章结构定义
*/
type PdfTmpl_Sign struct {
	PageNo   string `json:"page_no"`
	SignX    string `json:"sign_x"`
	SignY    string `json:"sign_y"`
	CorpNo   string `json:"corp_no"`
	CorpName string `json:"corp_name"`
}

/*
	Pdf模板定义
*/
type PdfTmpl_Define struct {
	Title       string           `json:"title"`
	TitleSize   float64          `json:"title_size"`
	Preface     string           `json:"preface"`
	PrefaceSize float64          `json:"preface_size"`
	Logo        string           `json:"logo"`
	Chapters    []Chapter_Define `json:"chapters"`
	Signs       []PdfTmpl_Sign   `json:"sings"`
	Footer      string           `json:"footer"`
	FooterSize  float64          `json:"footer_size"`
}

type CellData struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
}

type SeriesNode struct {
	Legend string   `json:"legend"`
	Data   []string `json:"data"`
}

type ChartData struct {
	Legend []string     `json:"legend"`
	Xaxis  []string     `json:"xaxis"`
	Series []SeriesNode `json:"series"`
}

/*
  创建报告
*/
type CrtRptReq struct {
	AppNo    string                  `json:"app_no"`
	AppReqNo string                  `json:"app_req_no"`
	TmplNo   string                  `json:"tmpl_no"`
	BaseMap  map[string]string       `json:"base_info"`
	TableMap map[string][][]CellData `json:"table_info"`
	ChartMap map[string]ChartData    `json:"chart_info"`
}

type CrtRptResp struct {
	AppReqNo string `json:"app_req_no"`
	RptNo    string `json:"rpt_no"`
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
}

/*
  查看协议
*/
type QryAgrtReq struct {
	AppNo    string            `json:"app_no"`
	AppReqNo string            `json:"app_req_no"`
	TmplNo   string            `json:"tmpl_no"`
	BaseMap  map[string]string `json:"base_info"`
}

type QryAgrtResp struct {
}

/*
  查询报告状态
*/

type QryRptStaReq struct {
	AppNo string `json:"app_no"`
	RptNo string `json:"rpt_no"`
}

type QryRptStaResp struct {
	RptNo   string `json:"rpt_no"`
	Status  string `json:"status"`
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
  查询报告
*/
type QryRptReq struct {
	AppNo string `json:"app_no"`
	RptNo string `json:"rpt_no"`
}

/*
  返回流
*/
type QryRptResp struct {
}

/*
  ##### ##### ##### ##### ##### ##### ##### ##### ##### ##### ##### ##### #####
*/
var (
	preFix        = "空格"
	FontName      = "NotoSansSC-Regular"
	leftMargin    float64
	rightMargin   float64
	topMargin     float64
	bottomMargin  float64
	pageWidth     float64
	pageHeight    float64
	pageWidthAll  float64
	pageHeightAll float64
)

func initPdf(pdf *gofpdf.Fpdf) {
	common.PrintHead("initPdf")
	leftMargin, topMargin, rightMargin, bottomMargin = pdf.GetMargins()
	pageWidthAll, pageHeightAll = pdf.GetPageSize()
	pageWidth = pageWidthAll - leftMargin - rightMargin
	pageHeight = pageHeightAll - topMargin - bottomMargin
	log.Println("pageWidthAll,pageHeightAll==>", pageWidthAll, pageHeightAll)
	log.Println("pageWidth,pageHeight==>", pageWidth, pageHeight)
	log.Println("ltrb==>", leftMargin, topMargin, rightMargin, bottomMargin)
	common.PrintTail("initPdf")
}

func CrtPdfByTmpl(pdfTmpl PdfTmpl_Define, baseMap map[string]string, tableMap map[string][][]CellData, chartMap map[string]ChartData) error {
	common.PrintHead("CrtPdfByTmpl")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font(FontName, "", "./NotoSansSC-Regular.ttf")
	initPdf(pdf)
	pdf.AddPage()

	if err, titleMsg := parserTmpl("FIX001", pdfTmpl.Title, baseMap); err != nil {
	} else {
		drawTitle(pdf, pdfTmpl.TitleSize, titleMsg)
	}

	if err, titleMsg := parserTmpl("FIX001", pdfTmpl.Preface, baseMap); err != nil {
	} else {
		drawText(pdf, pdfTmpl.PrefaceSize, titleMsg)
	}

	for _, v := range pdfTmpl.Chapters {

		drawChapter(pdf, v.TitleSize, v.Title)

		drawText(pdf, v.ContentSize, v.Content)

		if v.Table != common.EMPTY_STRING {
			drawDTable(pdf, v.TableTitleSize, v.TableDataSize, tableMap[v.Table])
		}

		if v.Chart != common.EMPTY_STRING {

		}
	}

	drawFooter(pdf, pdfTmpl.FooterSize, pdfTmpl.Footer)

	err := pdf.OutputFileAndClose(common.DEFAULT_PATH + fmt.Sprintf("pdf_%d.pdf", 123))
	log.Println(err)
	return nil
}

/*
	在PDF画二维表格
*/
func drawDTable(pdf *gofpdf.Fpdf, titleSize float64, dataSize float64, tableData [][]CellData) {
	w := []float64{30, 30, 40, 40, 50}
	fill := false
	for i, v := range tableData {
		if i == 0 {
			pdf.SetFillColor(255, 0, 0)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetDrawColor(128, 0, 0)
			pdf.SetLineWidth(.2)
			pdf.SetFont(FontName, "", titleSize)
			for j, vv := range v {
				pdf.CellFormat(w[j], 12, vv.Desc, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFillColor(224, 235, 255)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont(FontName, "", dataSize)
			for j, vv := range v {
				pdf.CellFormat(w[j], 10, vv.Value, "1", 0, "L", true, 0, "")
			}
		} else {
			pdf.Ln(-1)
			pdf.SetFillColor(224, 235, 255)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont(FontName, "", dataSize)
			for j, vv := range v {
				pdf.CellFormat(w[j], 10, vv.Value, "1", 0, "L", fill, 0, "")
			}
			fill = !fill
		}
	}
	pdf.Ln(-1)
}

/*
  根据基础MAP，画静态TABLE
*/
func drawSTable(pdf *gofpdf.Fpdf, lableName []string, baseInfo map[string]string) {
	w := []float64{20, 30, 20, 30, 20, 30, 20, 20}
	fill := true
	index := 0
	var totalLen float64
	pdf.Ln(2)
	for i, v := range lableName {
		if i%4 == 0 {
			pdf.Ln(-1)
			totalLen = 0
		}
		index = i % 4
		//Label
		totalLen += w[index*2]
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont(FontName, "", 9)
		pdf.CellFormat(w[index*2], 6, lableName[i], "1", 0, "L", fill, 0, "")

		if i == len(lableName)-1 {
			if totalLen != 190.00 {
				pdf.CellFormat(190.00-totalLen, 6, baseInfo[v], "1", 0, "L", !fill, 0, "")
			} else {
				pdf.CellFormat(w[index*2+1], 6, baseInfo[v], "1", 0, "L", !fill, 0, "")
			}
		} else {
			totalLen += w[index*2+1]
			pdf.CellFormat(w[index*2+1], 6, baseInfo[v], "1", 0, "L", !fill, 0, "")
		}
	}
	pdf.Ln(-1)
}

func drawLogo(pdf *gofpdf.Fpdf) {
	pdf.Image("./bd.png", 180, -10, 10, 10, true, "", 0, "http://www.fpdf.org")
}

var nameItems = []string{"衬衫", "牛仔裤", "运动裤", "袜子", "冲锋衣", "羊毛衫"}
var seed = rand.NewSource(time.Now().UnixNano())

func randInt() []int {
	cnt := len(nameItems)
	r := make([]int, 0)
	for i := 0; i < cnt; i++ {
		r = append(r, int(seed.Int63())%50)
	}
	return r
}

func getZWFont() *truetype.Font {

	fontFile := "../../NotoSansSC-Regular.ttf"
	//fontFile := "/Library/Fonts/AppleMyungjo.ttf"

	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}

func drawChart(pdf *gofpdf.Fpdf) {
	//	f := getZWFont() // 用自己的字体
	graph := chart.BarChart{
		Title: "Test Bar Chart",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 10,
			},
		},
		Height:   300,
		BarWidth: 300,
		Bars: []chart.Value{
			{Value: 5.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "??"},
			{Value: 10, Label: "!!"},
		},
	}

	ff, _ := os.Create("output.png")
	defer ff.Close()
	graph.Render(chart.PNG, ff)
	pdf.Ln(-1)
	pdf.Image("./output.png", 10, pdf.GetY(), 170, 40, true, "", 0, "http://www.fpdf.org")

}

func drawChapter(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.Ln(2)
	pdf.SetFont(FontName, "", fontSize)
	_, lineHt := pdf.GetFontSize()
	pdf.SetFillColor(200, 220, 255)
	pdf.CellFormat(pageWidth, lineHt+2, text,
		"", 1, "L", true, 0, "")
	pdf.Ln(2)
}

func drawTheme(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.SetFont(FontName, "U", fontSize)
	pdf.Ln(2)
	_, lineHt := pdf.GetFontSize()
	pdf.SetUnderlineThickness(1)
	pdf.SetTextColor(255, 0, 0)
	pdf.WriteAligned(pageWidth, lineHt, text, "L")
	pdf.Ln(lineHt + 2)
}

func drawHeader(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.SetTopMargin(20)
	pdf.SetHeaderFuncMode(func() {
		//	pdf.Image("./bd.png", 10, 2, 20, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("NotoSansSC-Regular", "U", fontSize)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(20, 10, text, "", 0, "C", false, 0, "")
	}, true)
}

func drawFooter(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont(FontName, "", fontSize)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 10, text+fmt.Sprintf("%d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
}

func drawTitle(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.SetFont(FontName, "", fontSize)
	_, lineHt := pdf.GetFontSize()
	pdf.WriteAligned(pageWidth, lineHt, text, "C")
	pdf.Ln(lineHt + 2)
}

/*
	显示文本信息
*/
func drawText(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	pdf.SetFont(FontName, "", fontSize)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(0, 0, 0)
	lines := pdf.SplitText(preFix+text, pageWidth)
	lineHt += 2
	for r, i := range lines {
		if r == 0 {
			firstLine := strings.TrimLeft(i, preFix)
			pdf.SetX(leftMargin + pdf.GetStringWidth(preFix))
			pdf.Cell(pageWidth-pdf.GetStringWidth(preFix), lineHt, firstLine)
		} else {
			pdf.Cell(pageWidth, lineHt, i)
		}
		pdf.Ln(lineHt)
	}
}

/*
	说明：判断
	出参：参数1：返回符合条件的对象列表
*/

func qryTmpl(tmplNo string) (string, string) {
	common.PrintHead("qryTmpl")
	r := rpttmpl.New(dbcomm.GetDB(), rpttmpl.DEBUG)
	var search rpttmpl.Search
	search.TmplNo = tmplNo
	search.Status = common.STATUS_ENABLED
	if u, err := r.Get(search); err != nil {
		return common.ERR_CODE_NOTFIND, common.EMPTY_STRING
	} else {
		return common.ERR_CODE_SUCCESS, u.Content
	}
}

/*
	说明：判断
	出参：参数1：返回符合条件的对象列表
*/

func parserTmpl(tmplNo string, tmplMsg string, baseMap map[string]string) (error, string) {
	common.PrintHead("parserTmpl")
	t, err := template.New(tmplNo).Parse(tmplMsg)
	if err != nil {
		return err, common.EMPTY_STRING
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, baseMap); err != nil {
		return err, common.EMPTY_STRING
	}
	common.PrintTail("parserTmpl")
	return nil, buf.String()
}

/*
	说明：判断
	出参：参数1：返回符合条件的对象列表
*/
func qryOrder(rptNo string) (string, string) {
	common.PrintHead("QryAgrt")
	r := rptorder.New(dbcomm.GetDB(), rptorder.DEBUG)
	var search rptorder.Search
	search.AppReqNo = rptNo
	if u, err := r.Get(search); err != nil {
		return common.ERR_CODE_NOTFIND, common.EMPTY_STRING
	} else {
		return common.ERR_CODE_SUCCESS, u.RptNo
	}
}

/*
	说明：判断
	出参：参数1：返回符合条件的对象列表
*/

func qryApp(appNo string) (string, string) {
	common.PrintHead("QryAgrt")
	r := rptapp.New(dbcomm.GetDB(), rptapp.DEBUG)
	var search rptapp.Search
	search.AppNo = appNo
	if u, err := r.Get(search); err != nil {
		return common.ERR_CODE_NOTFIND, common.EMPTY_STRING
	} else {
		return common.ERR_CODE_SUCCESS, u.Status
	}
}

/*
	说明：发送短信
	出参：参数1：返回符合条件的对象列表
*/

func QryAgrt(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("QryAgrt")
	var qryReq QryAgrtReq
	err := json.NewDecoder(req.Body).Decode(&qryReq)
	if err != nil {
		log.Println("NewDecoder Json Error:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	defer req.Body.Close()

	var tmplMsg, errCode string
	if errCode, tmplMsg = qryTmpl(qryReq.TmplNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query Template  Error:", errCode)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err, tmplMsg = parserTmpl(qryReq.TmplNo, tmplMsg, qryReq.BaseMap)
	if err != nil {
		log.Println("Parser Template  Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(tmplMsg))
	common.PrintTail("QryAgrt")
}

/*
	说明：发送短信
	出参：参数1：返回符合条件的对象列表
*/

func CrtRptFile(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CrtRptFile")
	var crtReq CrtRptReq
	var crtResp CrtRptResp
	err := json.NewDecoder(req.Body).Decode(&crtReq)
	if err != nil {
		crtResp.ErrCode = common.ERR_CODE_JSONERR
		crtResp.ErrCode = common.ERROR_MAP[common.ERR_CODE_JSONERR]
		common.Write_Response(crtResp, w, req)
		return
	}
	defer req.Body.Close()
	var tmplMsg, errCode string
	if errCode, tmplMsg = qryTmpl(crtReq.TmplNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query Template  Error:", errCode)
		crtResp.ErrCode = common.ERR_CODE_NOTFIND
		crtResp.ErrCode = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(crtResp, w, req)
		return
	}
	var pdfTmpl PdfTmpl_Define
	if err = json.Unmarshal([]byte(tmplMsg), &pdfTmpl); err != nil {
		log.Println("Parser Template  Error:", err)
		crtResp.ErrCode = common.ERR_CODE_JSONERR
		crtResp.ErrMsg = "模板" + common.ERROR_MAP[common.ERR_CODE_JSONERR]
		common.Write_Response(crtResp, w, req)
		return
	}
	if err = CrtPdfByTmpl(pdfTmpl, crtReq.BaseMap, crtReq.TableMap, crtReq.ChartMap); err != nil {
		log.Println("Parser Template  Error:", err)
		crtResp.ErrCode = common.ERR_CODE_PDFERR
		crtResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PDFERR]
		common.Write_Response(crtResp, w, req)
		return
	}
	crtResp.ErrCode = common.ERR_CODE_SUCCESS
	crtResp.ErrCode = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(crtResp, w, req)

	common.PrintTail("CrtRptFile")
}

/*
	说明：发送短信
	出参：参数1：返回符合条件的对象列表
*/

func QryRptFile(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("QryRptFile")
	var qryReq QryRptReq
	err := json.NewDecoder(req.Body).Decode(&qryReq)
	if err != nil {
		log.Println("NewDecoder Json Error:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	defer req.Body.Close()

	r := rptorder.New(dbcomm.GetDB(), rptorder.DEBUG)
	var search rptorder.Search
	search.RptNo = qryReq.RptNo
	search.AppNo = qryReq.AppNo

	if u, err := r.Get(search); err != nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		if buf, err := ioutil.ReadFile(common.DEFAULT_PATH + u.FileUrl); err != nil {
			log.Println("Open File Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Write(buf)
		}
	}
	common.PrintTail("QryRptFile")
}

/*
	说明：查询生成数据报告的状态
	入参：APP_NO,RPT_NO
	返回：ERR_CODE,ERR_MSG,STATUS
*/

func QryRptStatus(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("QryRptStatus")
	var qryReq QryRptStaReq
	var qryResp QryRptStaResp
	err := json.NewDecoder(req.Body).Decode(&qryReq)
	if err != nil {
		qryResp.ErrCode = common.ERR_CODE_JSONERR
		qryResp.ErrCode = common.ERROR_MAP[common.ERR_CODE_JSONERR]
		common.Write_Response(qryResp, w, req)
		return
	}
	defer req.Body.Close()

	r := rptorder.New(dbcomm.GetDB(), rptorder.DEBUG)
	var search rptorder.Search
	search.AppNo = qryReq.AppNo
	search.RptNo = qryReq.RptNo
	if u, err := r.Get(search); err != nil {
		qryResp.Status = common.EMPTY_STRING
		qryResp.ErrCode = common.ERR_CODE_NOTFIND
		qryResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
	} else {
		qryResp.Status = u.Status
		qryResp.ErrCode = u.ErrCode
		qryResp.ErrCode = u.ErrMsg
	}
	qryResp.RptNo = qryReq.RptNo
	common.Write_Response(qryResp, w, req)
	common.PrintTail("QryRptStatus")
}

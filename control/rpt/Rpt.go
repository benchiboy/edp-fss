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
	"github.com/wcharczuk/go-chart/drawing"
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
	AppNo    string                  `json:"systemNo"`
	AppReqNo string                  `json:"requestRefNo"`
	TmplNo   string                  `json:"templateNo"`
	BaseMap  map[string]string       `json:"baseInfo"`
	TableMap map[string][][]CellData `json:"tableInfo"`
	ChartMap map[string]ChartData    `json:"chartInfo"`
}

type CrtRptResp struct {
	AppReqNo string `json:"requestRefNo"`
	RptNo    string `json:"reportId"`
	ErrCode  string `json:"errorCode"`
	ErrMsg   string `json:"errorMessage"`
}

/*
  查看协议
*/
type QryAgrtReq struct {
	AppNo    string            `json:"systemNo"`
	AppReqNo string            `json:"requestRefNo"`
	TmplNo   string            `json:"templateNo"`
	BaseMap  map[string]string `json:"baseInfo"`
}

type QryAgrtResp struct {
}

/*
  查询报告状态
*/

type QryRptStaReq struct {
	AppNo string `json:"systemNo"`
	RptNo string `json:"reportId"`
}

type QryRptStaResp struct {
	RptNo   string `json:"reportId"`
	Status  string `json:"status"`
	ErrCode string `json:"errorCode"`
	ErrMsg  string `json:"errorMessage"`
}

/*
  查询报告
*/
type QryRptReq struct {
	AppNo string `json:"systemNo"`
	RptNo string `json:"reportId"`
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

func CrtPdfByTmpl(pdfTmpl PdfTmpl_Define, baseMap map[string]string, tableMap map[string][][]CellData, chartMap map[string]ChartData) (error, string) {
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
			log.Println("Draw a table...", v.TableCols, tableMap[v.Table])
			drawDTable(pdf, v.TableCols, v.TableTitleSize, v.TableDataSize, tableMap[v.Table])
		}

		if v.Chart != common.EMPTY_STRING {
			DrawPieChart(pdf)
		}
	}

	drawFooter(pdf, pdfTmpl.FooterSize, pdfTmpl.Footer)
	pdfName := common.DEFAULT_PATH + fmt.Sprintf("pdf_%d.pdf", time.Now().UnixNano())
	err := pdf.OutputFileAndClose(pdfName)
	log.Println(err)
	return nil, pdfName
}

/*
	在PDF画二维表格
*/
func drawDTable(pdf *gofpdf.Fpdf, colLens []float64, titleSize float64, dataSize float64, tableData [][]CellData) error {
	if titleSize == 0 {
		titleSize = common.DEFAULT_FONT_SIZE
	}
	if dataSize == 0 {
		titleSize = common.DEFAULT_FONT_SIZE
	}
	w := colLens
	var totalLen float64
	for _, v := range w {
		totalLen += v
	}
	if totalLen == 0 {
		log.Println("表格合计长度不能为", totalLen)
		return fmt.Errorf("表格合计长度不能为%f", totalLen)
	}
	if totalLen > pageWidth+1 {
		log.Println("表格合计长度不能为", totalLen)
		return fmt.Errorf("表格合计长度%f长度大于pageWidth", totalLen)
	}
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
	return nil
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

func drawLogo(pdf *gofpdf.Fpdf, imgFile string) {
	pdf.Image(imgFile, 180, -10, 10, 10, true, "", 0, "http://www.crfchina.com")
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

func getCustFont() *truetype.Font {
	fontFile := "../../NotoSansSC-Regular.ttf"
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

func releases() []chart.GridLine {
	return []chart.GridLine{
		{Value: 1.00},
		{Value: 2.00},
	}
}

func DrawSeriesChart() {
	custFont := getCustFont()

	graph := chart.Chart{
		//	Title:      "信用分布情况",
		Font:       custFont,
		TitleStyle: chart.StyleShow(),

		Background: chart.Style{
			Padding: chart.Box{
				Top: 10,
				//Left: 10,
			},
		},

		XAxis: chart.XAxis{
			Name:           "hello",
			NameStyle:      chart.StyleShow(),
			ValueFormatter: chart.TimeValueFormatter,
			Style:          chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 0.0,
			},
			Ticks: []chart.Tick{
				{Value: 1.0, Label: "2001年"},
				{Value: 2.0, Label: "2002年"},
				{Value: 3.0, Label: "2003年"},
				{Value: 4.0, Label: "2004年"},
				{Value: 5.0, Label: "2005年"},
			},
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			GridLines: releases(),
		},
		Width:  412,
		Height: 390,
		YAxis: chart.YAxis{
			//Name:      "yvalue",
			//Style: {chart.StyleShow(),
			Style: chart.Style{
				Show: true,
				//FontSize: 16.0,
				//StrokeColor: drawing.ColorRed, // will supercede defaults
				//FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
			},
			//NameStyle: chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 30.0,
				FontSize:            9.0,
				//FontColor:       drawing.ColorRed,
				//TextLineSpacing: 5,
			},
			Ticks: []chart.Tick{
				{Value: 1.0, Label: "1.00"},
				{Value: 2.0, Label: "2.00"},
				{Value: 3.0, Label: "3.00"},
				{Value: 4.0, Label: "4.00"},
				{Value: 5.0, Label: "5.00"},
				{Value: 6.0, Label: "6.00"},
				{Value: 7.0, Label: "7.00"},
				{Value: 8.0, Label: "8.00"},
				{Value: 9.0, Label: "9.00"},
				{Value: 10.0, Label: "10.00"},
				{Value: 11.0, Label: "11.00"},
				{Value: 12.0, Label: "12.00"},
			},
			GridMajorStyle: chart.Style{
				StrokeColor:     drawing.ColorRed,
				StrokeWidth:     0.4,
				Show:            true,
				FillColor:       drawing.ColorFromHex("efefef"),
				StrokeDashArray: []float64{2.0, 7.0},
			},
			GridLines: []chart.GridLine{{Value: 6}},
		},
		Canvas: chart.Style{
			//FillColor: drawing.ColorFromHex("efe1ef"),
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					FontSize:    16.0,
					StrokeWidth: 1.0,
					//	StrokeColor: drawing.ColorRed,               // will supercede defaults
					//	FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
				},
				Name:    "工业指数",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{6.0, 7.0, 10.0, 9.0, 6.0},
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:            true,
					StrokeDashArray: []float64{5.0, 4.0},
					DotWidth:        1.0,
					DotColor:        drawing.ColorRed,
					//StrokeColor:     drawing.ColorBlue,               // will supercede defaults
					//FillColor:       drawing.ColorRed.WithAlpha(164), // will supercede defaults
				},
				Name:    "工业指数1",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{2.0, 5.0, 4.0, 10.0, 3.0},
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create("output11.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}

func DrawChart1() {
	//custFont := getCustFont()

	graph := chart.Chart{
		//		Background: chart.Style{
		//			FillColor: drawing.ColorBlue,
		//		},
		//		Canvas: chart.Style{
		//			FillColor: drawing.ColorFromHex("efefef"),
		//		},

		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},

		YAxis: chart.YAxis{
			Name:      "hello",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			GridMajorStyle: chart.Style{
				StrokeColor:     drawing.ColorRed,
				StrokeWidth:     0.4,
				Show:            true,
				FillColor:       drawing.ColorFromHex("efefef"),
				StrokeDashArray: []float64{2.0, 7.0},
			},
			GridLines: []chart.GridLine{{Value: 2.0}, {Value: 3}, {Value: 4}, {Value: 5}},
		},
		XAxis: chart.XAxis{
			Name:      "hello",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			//			GridMajorStyle: chart.Style{
			//				StrokeColor: chart.ColorAlternateGray,
			//				StrokeWidth: 1.0,
			//				Show:        true,
			//			},
			//			GridLines: []chart.GridLine{{Value: 4.0}},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{2.0, 4.0, 1.0, 8.0, 7.0},
			},
		},
	}

	//	graph := chart.Chart{
	//		Font: custFont,
	//		XAxis: chart.XAxis{
	//			Name:  "The YAxis",
	//			Style: chart.StyleShow(),
	//		},
	//		Series: []chart.Series{
	//			chart.ContinuousSeries{
	//				Style: chart.Style{
	//					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
	//					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
	//				},
	//				XValues: []float64{2001, 2002, 2003, 2004, 2005},
	//				YValues: []float64{1.0, 8.0, 3.0, 4.0, 10.0},
	//			},
	//		},
	//	}

	//	graph.Elements = []chart.Renderable{
	//		chart.LegendThin(&graph),
	//	}

	f, _ := os.Create("output11.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}

func DrawBarChart(pdf *gofpdf.Fpdf) {
	//	custFont := getCustFont()
	profitStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("13c158"),
		StrokeColor: drawing.ColorFromHex("13c158"),
		StrokeWidth: 0,
	}

	lossStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("c11313"),
		StrokeColor: drawing.ColorFromHex("c11313"),
		StrokeWidth: 0,
	}

	sbc := chart.BarChart{
		Title: "Bar Chart Using BaseValue",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		YAxis: chart.YAxis{
			Ticks: []chart.Tick{
				{Value: -4.0, Label: "-4"},
				{Value: -2.0, Label: "-2"},
				{Value: 0, Label: "0"},
				{Value: 2.0, Label: "2"},
				{Value: 4.0, Label: "4"},
				{Value: 6.0, Label: "6"},
				{Value: 8.0, Label: "8"},
				{Value: 10.0, Label: "10"},
				{Value: 12.0, Label: "12"},
			},
		},
		Bars: []chart.Value{
			{Value: 10.0, Style: profitStyle, Label: "Profit"},
			{Value: 12.0, Style: profitStyle, Label: "More Profit"},
			{Value: 8.0, Style: profitStyle, Label: "Still Profit"},
			{Value: -4.0, Style: lossStyle, Label: "Loss!"},
			{Value: 3.0, Style: profitStyle, Label: "Phew Ok"},
			{Value: -2.0, Style: lossStyle, Label: "Oh No!"},
		},
	}
	f, _ := os.Create("output11.png")
	defer f.Close()
	sbc.Render(chart.PNG, f)

	pdf.Image("./output11.png", 40, 10, 80, 100, true, "", 0, "http://www.crfchina.com")

}
func DrawPieChart(pdf *gofpdf.Fpdf) {
	custFont := getCustFont()
	pie := chart.PieChart{
		Title:      "信用分布情况",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 35,
			},
		},
		Font:   custFont,
		Width:  512,
		Height: 590,
		Values: []chart.Value{
			{Value: 40, Label: "信用低"},
			{Value: 30, Label: "Two"},
			{Value: 30, Label: "One"},
			{Value: 30, Label: "One"},
		},
	}

	f, _ := os.Create("output11.png")
	defer f.Close()
	pie.Render(chart.PNG, f)

	pdf.Image("./output11.png", 40, 10, 80, 100, true, "", 0, "http://www.crfchina.com")

}

func drawChapter(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	if fontSize == 0 {
		fontSize = common.DEFAULT_FONT_SIZE
	}
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
	if fontSize == 0 {
		fontSize = common.DEFAULT_FONT_SIZE
	}
	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont(FontName, "", fontSize)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(0, 10, text+fmt.Sprintf("%d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
}

func drawTitle(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	if fontSize == 0 {
		fontSize = common.DEFAULT_FONT_SIZE
	}
	pdf.SetFont(FontName, "", fontSize)
	_, lineHt := pdf.GetFontSize()
	pdf.WriteAligned(pageWidth, lineHt, text, "C")
	pdf.Ln(lineHt + 2)
}

/*
	显示文本信息
*/
func drawText(pdf *gofpdf.Fpdf, fontSize float64, text string) {
	if fontSize == 0 {
		fontSize = common.DEFAULT_FONT_SIZE
	}
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
	search.Status = common.STATUS_ENABLED
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

	if req.Method == "GET" {
		appNo, ok := req.URL.Query()["systemNo"]
		if !ok || len(appNo) < 1 || appNo[0] == "" {
			log.Println("URL.Query Json Error:")
			w.WriteHeader(http.StatusForbidden)
			return

		}
		qryReq.AppNo = appNo[0]

		tmplNo, ok := req.URL.Query()["templateNo"]
		if !ok || len(appNo) < 1 || appNo[0] == "" {
			log.Println("URL.Query Json Error:")
			w.WriteHeader(http.StatusForbidden)
			return

		}
		qryReq.TmplNo = tmplNo[0]
	} else {
		err := json.NewDecoder(req.Body).Decode(&qryReq)
		if err != nil {
			log.Println("NewDecoder Json Error:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		defer req.Body.Close()
	}

	var errCode string
	if errCode, _ = qryApp(qryReq.AppNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query AppNo  Error:", errCode)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var tmplMsg string
	if errCode, tmplMsg = qryTmpl(qryReq.TmplNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query Template  Error:", errCode)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err, tmplMsg := parserTmpl(qryReq.TmplNo, tmplMsg, qryReq.BaseMap)
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

	var errCode string
	if errCode, _ = qryApp(crtReq.AppNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query AppNo  Error:", errCode)
		crtResp.ErrCode = common.ERR_CODE_NOTFIND
		crtResp.ErrMsg = crtReq.AppNo + common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(crtResp, w, req)
		return
	}

	var tmplMsg string
	if errCode, tmplMsg = qryTmpl(crtReq.TmplNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query Template  Error:", errCode)
		crtResp.ErrCode = common.ERR_CODE_NOTFIND
		crtResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
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

	r := rptorder.New(dbcomm.GetDB(), rptorder.DEBUG)
	var search rptorder.Search
	search.AppNo = crtReq.AppNo
	search.RptNo = crtReq.AppReqNo
	var e rptorder.RptOrder
	if u, err := r.Get(search); err != nil {
		e.AppNo = crtReq.AppNo
		e.AppReqNo = crtReq.AppReqNo
		e.RptNo = fmt.Sprintf("A%d", time.Now().UnixNano())
		e.TmplNo = crtReq.TmplNo
		buf, _ := json.Marshal(crtReq.BaseMap)
		e.BaseInfo = string(buf)
		buf, _ = json.Marshal(crtReq.TableMap)
		e.TableInfo = string(buf)
		buf, _ = json.Marshal(crtReq.ChartMap)
		e.ChartInfo = string(buf)
		e.Status = common.STATUS_INIT
		e.InsertDate = time.Now().Unix()
		e.Version = 1
		r.InsertEntity(e, nil)

	} else {
		crtResp.AppReqNo = u.AppReqNo
		crtResp.RptNo = u.RptNo
		crtResp.ErrCode = u.ErrCode
		crtResp.ErrCode = u.ErrMsg
		common.Write_Response(crtResp, w, req)
		return
	}
	var pdfName string
	if err, pdfName = CrtPdfByTmpl(pdfTmpl, crtReq.BaseMap, crtReq.TableMap, crtReq.ChartMap); err != nil {
		log.Println("Parser Template  Error:", err)
		crtResp.ErrCode = common.ERR_CODE_PDFERR
		crtResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PDFERR]
		common.Write_Response(crtResp, w, req)
		return
	}
	tmpMap := map[string]interface{}{"err_code": common.ERR_CODE_SUCCESS,
		"err_msg":  common.ERROR_MAP[common.ERR_CODE_SUCCESS],
		"file_url": pdfName}

	r.UpdateMap(e.RptNo, tmpMap, nil)

	crtResp.AppReqNo = e.AppReqNo
	crtResp.RptNo = e.RptNo
	crtResp.ErrCode = common.ERR_CODE_SUCCESS
	crtResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
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

	if req.Method == "GET" {
		appNo, ok := req.URL.Query()["systemNo"]
		if !ok || len(appNo) < 1 || appNo[0] == "" {
			log.Println("URL.Query Json Error:")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		qryReq.AppNo = appNo[0]
		rptNo, ok := req.URL.Query()["reportId"]
		if !ok || len(appNo) < 1 || appNo[0] == "" {
			log.Println("URL.Query Json Error:")
			w.WriteHeader(http.StatusForbidden)
			return

		}
		qryReq.RptNo = rptNo[0]
	} else {
		err := json.NewDecoder(req.Body).Decode(&qryReq)
		if err != nil {
			log.Println("NewDecoder Json Error:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		defer req.Body.Close()
	}

	var errCode string
	if errCode, _ = qryApp(qryReq.AppNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query AppNo  Error:", errCode)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r := rptorder.New(dbcomm.GetDB(), rptorder.DEBUG)
	var search rptorder.Search
	search.RptNo = qryReq.RptNo
	search.AppNo = qryReq.AppNo
	if u, err := r.Get(search); err != nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		log.Println(u.FileUrl)
		if buf, err := ioutil.ReadFile(u.FileUrl); err != nil {
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
		qryResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR]
		common.Write_Response(qryResp, w, req)
		return
	}
	defer req.Body.Close()

	var errCode string
	if errCode, _ = qryApp(qryReq.AppNo); errCode != common.ERR_CODE_SUCCESS {
		log.Println("Query AppNo  Error:", errCode)
		qryResp.ErrCode = common.ERR_CODE_NOTFIND
		qryResp.ErrMsg = qryReq.AppNo + common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(qryResp, w, req)
		return
	}

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
		qryResp.ErrMsg = u.ErrMsg
	}
	qryResp.RptNo = qryReq.RptNo
	common.Write_Response(qryResp, w, req)
	common.PrintTail("QryRptStatus")
}

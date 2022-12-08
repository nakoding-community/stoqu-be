package usecase

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"golang.org/x/exp/maps"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

type Report interface {
	FindOrder(ctx context.Context, filterParam abstraction.Filter) (result dto.OrderReportResponse, pagination abstraction.PaginationInfo, err error)
	FindOrderExcel(ctx context.Context, filterParam abstraction.Filter) (f *os.File, err error)
	FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error)
	FindOrderProductExcel(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (f *os.File, err error)
}

type report struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewReport(cfg *config.Configuration, f repository.Factory) Report {
	return &report{cfg, f}
}

func (u *report) FindOrder(ctx context.Context, filterParam abstraction.Filter) (result dto.OrderReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result.Orders = append(result.Orders, order)
	}

	totalOrder, err := u.Repo.OrderTrx.Count(ctx)
	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	totalIncome, err := u.Repo.OrderTrx.CountIncome(ctx)
	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	result.TotalIncome = totalIncome
	result.TotalOrder = totalOrder

	return result, pagination, nil
}

func (u *report) orderExcelData(ctx context.Context, filterParam abstraction.Filter) ([]entity.OrderView, error) {
	var search *abstraction.Search

	// !TODO, able to exclude pagination or using concurrent instead
	defaultLimit := 50
	limit := func() int {
		if filterParam.Limit == nil {
			return defaultLimit
		}
		if *filterParam.Limit > defaultLimit {
			return *filterParam.Limit
		}
		return defaultLimit
	}()

	filterParam.Limit = &limit

	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return nil, err
	}

	if info.Count > int64(limit) {
		limit = int(info.Count)
		filterParam.Limit = &limit
		orders, err = u.orderExcelData(ctx, filterParam)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (u *report) FindOrderExcel(ctx context.Context, filterParam abstraction.Filter) (f *os.File, err error) {
	orders, err := u.orderExcelData(ctx, filterParam)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	// !TODO, should implement flexible index excel reading, now redundant with excel cell value loop
	type ExcelEntity struct {
		CreatedAt     string
		CustomerName  string
		CustomerPhone string
		PriceFinal    float64
		Status        string
		Notes         string
	}
	var excelData []ExcelEntity
	for _, v := range orders {
		excelData = append(excelData, ExcelEntity{
			CreatedAt:     v.CreatedAt.Format(time.RFC822Z),
			CustomerName:  v.CustomerName,
			CustomerPhone: v.CustomerPhone,
			PriceFinal:    v.FinalPrice,
			Status:        v.Status,
			Notes:         v.Notes,
		})
	}

	headers := map[string]string{
		"A1": "No",
		"B1": "Tanggal",
		"C1": "Customer",
		"D1": "No. Handphone",
		"E1": "Harga",
		"F1": "Status",
		"G1": "Keterangan",
	}

	return u.GenerateExcelReport(
		ctx,
		dto.GenerateExcelReportInput{
			SheetName: "Laporan Pemesanan",
			FileName:  "report-summaries-1.xlsx",
			Headers:   headers,
			Data:      excelData,
		},
	)
}

func (u *report) FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search

	var (
		orders []entity.OrderViewProduct
		info   *abstraction.PaginationInfo
		count  int64
	)
	switch query.Group {
	case constant.GROUP_BY_VARIANT:
		logrus.Info("FindGroupByVariant")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
	case constant.GROUP_BY_PACKET:
		logrus.Info("FindGroupByPacket")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
	default: //constant.GROUP_BY_BRAND
		logrus.Info("FindGroupByBrand")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByBrand(ctx, filterParam, search)
	}

	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result.Orders = append(result.Orders, order)
	}

	result.Total = count

	return result, pagination, nil
}

func (u *report) orderProductExcelData(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (orders []entity.OrderViewProduct, err error) {
	var search *abstraction.Search

	// !TODO, able to exclude pagination or using concurrent instead
	defaultLimit := 50
	limit := func() int {
		if filterParam.Limit == nil {
			return defaultLimit
		}
		if *filterParam.Limit > defaultLimit {
			return *filterParam.Limit
		}
		return defaultLimit
	}()

	filterParam.Limit = &limit

	var info *abstraction.PaginationInfo

	switch query.Group {
	case constant.GROUP_BY_VARIANT:
		logrus.Info("orderProductExcelDataByVariant")
		orders, _, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
	case constant.GROUP_BY_PACKET:
		logrus.Info("orderProductExcelDataByPacket")
		orders, _, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
	default:
		logrus.Info("orderProductExcelDataByBrand")
		orders, _, info, err = u.Repo.OrderTrx.FindGroupByBrand(ctx, filterParam, search)
	}

	if err != nil {
		return nil, err
	}

	if info.Count > int64(limit) {
		limit = int(info.Count)
		filterParam.Pagination.Limit = &limit
		orders, err = u.orderProductExcelData(ctx, filterParam, query)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (u *report) FindOrderProductExcel(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (f *os.File, err error) {
	orders, err := u.orderProductExcelData(ctx, filterParam, query)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	type ExcelEntityVariant struct {
		BrandID     string
		BrandName   string
		PacketID    string
		PacketName  string
		VariantID   string
		VariantName string
		Qty         float64
	}

	type ExcelEntityPacket struct {
		PacketID   string
		PacketName string
		Qty        float64
	}

	type ExcelEntityBrand struct {
		BrandID    string
		BrandName  string
		PacketID   string
		PacketName string
		Qty        float64
	}

	switch query.Group {
	case constant.GROUP_BY_VARIANT:
		var excelData []ExcelEntityVariant
		headers := map[string]string{
			"A1": "No",
			"B1": "Brand ID",
			"C1": "Brand Name",
			"D1": "Packet ID",
			"E1": "Packet Name",
			"F1": "Variant ID",
			"G1": "Variant Name",
			"H1": "Quantity",
		}

		for _, v := range orders {
			excelData = append(excelData, ExcelEntityVariant{
				BrandID:     v.BrandID,
				BrandName:   v.BrandName,
				PacketID:    v.PacketID,
				PacketName:  v.PacketName,
				VariantID:   v.VariantID,
				VariantName: v.VariantName,
				Qty:         v.Count,
			})
		}

		return u.GenerateExcelReport(
			ctx,
			dto.GenerateExcelReportInput{
				SheetName: "Laporan Produk",
				FileName:  "report-summaries-1.xlsx",
				Headers:   headers,
				Data:      excelData,
			},
		)

	case constant.GROUP_BY_PACKET:
		var excelData []ExcelEntityPacket
		headers := map[string]string{
			"A1": "No",
			"B1": "Packet ID",
			"C1": "Packet Name",
			"D1": "Quantity",
		}

		for _, v := range orders {
			excelData = append(excelData, ExcelEntityPacket{
				PacketID:   v.PacketID,
				PacketName: v.PacketName,
				Qty:        v.Count,
			})
		}

		return u.GenerateExcelReport(
			ctx,
			dto.GenerateExcelReportInput{
				SheetName: "Laporan Produk",
				FileName:  "report-summaries-1.xlsx",
				Headers:   headers,
				Data:      excelData,
			},
		)
	default:
		var excelData []ExcelEntityBrand
		headers := map[string]string{
			"A1": "No",
			"B1": "Brand ID",
			"C1": "Brand Name",
			"D1": "Packet ID",
			"E1": "Packet Name",
			"F1": "Quantity",
		}

		for _, v := range orders {
			excelData = append(excelData, ExcelEntityBrand{
				BrandID:    v.BrandID,
				BrandName:  v.BrandName,
				PacketID:   v.PacketID,
				PacketName: v.PacketName,
				Qty:        v.Count,
			})
		}

		return u.GenerateExcelReport(
			ctx,
			dto.GenerateExcelReportInput{
				SheetName: "Laporan Produk",
				FileName:  "report-summaries-1.xlsx",
				Headers:   headers,
				Data:      excelData,
			},
		)
	}
}

func (u *report) GenerateExcelReport(ctx context.Context, excel dto.GenerateExcelReportInput) (f *os.File, err error) {
	xlsx := excelize.NewFile()
	sheet := excel.SheetName
	xlsx.NewSheet(sheet)
	xlsx.DeleteSheet("Sheet1")

	styleHeaderAndBorder, err := xlsx.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"bdc3c7"},
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  13,
			Color: "000000",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	borderStyle, err := xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	// set headers
	headerIdx := maps.Keys(excel.Headers)
	for k, v := range excel.Headers {
		xlsx.SetCellValue(sheet, k, v)
		xlsx.SetCellStyle(sheet, k, k, styleHeaderAndBorder)
	}

	// set content
	// !FIXME, the excel.Data value might not be ordered correctly, so the content messed up
	// GC (standard go compiler) will reorder the struct for memory efficiency
	// https://stackoverflow.com/questions/9486364/why-cant-c-compilers-rearrange-struct-members-to-eliminate-alignment-padding
	// possible solution: set cell value base on map key
	slc := reflect.ValueOf(excel.Data)
	for i := 0; i < slc.Len(); i++ {
		for j := 0; j < slc.Index(i).NumField()+1; j++ {
			key := strings.Split(headerIdx[j], "1")[0]
			if j == 0 {
				xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", key, i+2), i+1) //numbering Ex: No. 1
			} else {
				xlsx.SetCellValue(sheet, fmt.Sprintf("%s%x", key, i+2), slc.Index(i).Field(j-1).Interface())
			}
			xlsx.SetCellStyle(sheet, fmt.Sprintf("%s%d", key, i+2), fmt.Sprintf("%s%d", key, i+2), borderStyle)
		}
	}

	f, err = os.Create(excel.FileName)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	err = xlsx.Write(f)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	return f, nil
}

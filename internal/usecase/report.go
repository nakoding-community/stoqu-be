package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
	"gitlab.com/stoqu/stoqu-be/internal/config"

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

	var excelData []map[string]interface{}
	for i, v := range orders {
		data := v.ToMap()
		data["excel_no"] = i + 1
		excelData = append(excelData, data)
	}

	dataMapping := dto.DataMapping{
		Headers: map[string]string{
			"A1": "No",
			"B1": "Tanggal",
			"C1": "Customer",
			"D1": "No. Handphone",
			"E1": "Harga",
			"F1": "Status",
			"G1": "Keterangan",
		},
		Body: map[string]string{
			"A": "excel_no",
			"B": "created_at",
			"C": "customer_name",
			"D": "customer_phone",
			"E": "final_price",
			"F": "status",
			"G": "notes",
		},
	}

	return u.GenerateExcelReport(
		ctx,
		dto.GenerateExcelReportInput{
			SheetName:   "Laporan Pemesanan",
			FileName:    "report-pemesanan-1.xlsx",
			Data:        excelData,
			DataMapping: dataMapping,
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
	switch query.Category {
	case constant.GROUP_BY_VARIANT:
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
	case constant.GROUP_BY_PACKET:
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
	default: //constant.GROUP_BY_BRAND
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

	switch query.Category {
	case constant.GROUP_BY_VARIANT:
		orders, _, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
	case constant.GROUP_BY_PACKET:
		orders, _, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
	default: //constant.GROUP_BY_BRAND
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

	var excelData []map[string]interface{}
	for i, v := range orders {
		data := v.ToMap()
		data["excel_no"] = i + 1
		excelData = append(excelData, data)
	}

	dataMappers := map[string]dto.DataMapping{
		constant.GROUP_BY_VARIANT: {
			Headers: map[string]string{
				"A1": "No",
				"B1": "Brand ID",
				"C1": "Brand Name",
				"D1": "Packet ID",
				"E1": "Packet Name",
				"F1": "Variant ID",
				"G1": "Variant Name",
				"H1": "Quantity",
			},
			Body: map[string]string{
				"A": "excel_no",
				"B": "brand_id",
				"C": "brand_name",
				"D": "packet_id",
				"E": "packet_name",
				"F": "variant_id",
				"G": "variant_name",
				"H": "count",
			},
		},
		constant.GROUP_BY_PACKET: {
			Headers: map[string]string{
				"A1": "No",
				"B1": "Packet ID",
				"C1": "Packet Name",
				"D1": "Quantity",
			},
			Body: map[string]string{
				"A": "excel_no",
				"B": "packet_id",
				"C": "packet_name",
				"D": "count",
			},
		},
		constant.GROUP_BY_BRAND: {
			Headers: map[string]string{
				"A1": "No",
				"B1": "Brand ID",
				"C1": "Brand Name",
				"D1": "Packet ID",
				"E1": "Packet Name",
				"F1": "Quantity",
			},
			Body: map[string]string{
				"A": "excel_no",
				"B": "brand_id",
				"C": "brand_name",
				"D": "packet_id",
				"E": "packet_name",
				"F": "count",
			},
		},
	}

	return u.GenerateExcelReport(
		ctx,
		dto.GenerateExcelReportInput{
			SheetName:   "Laporan Produk",
			FileName:    "laporan-produk-1.xlsx",
			Data:        excelData,
			DataMapping: dataMappers[query.Category],
		},
	)
}

func (u *report) GenerateExcelReport(ctx context.Context, excel dto.GenerateExcelReportInput) (f *os.File, err error) {
	xlsx := excelize.NewFile()
	sheet := excel.SheetName
	xlsx.NewSheet(sheet)
	xlsx.DeleteSheet("Sheet1")

	headerStyle, err := xlsx.NewStyle(constant.EXCEL_HEADER_STYLE)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	bodyStyle, err := xlsx.NewStyle(constant.EXCEL_BODY_STYLE)
	if err != nil {
		return nil, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	// set headers
	for ax, v := range excel.DataMapping.Headers {
		xlsx.SetCellValue(sheet, ax, v)
		xlsx.SetCellStyle(sheet, ax, ax, headerStyle)
	}

	// set content
	for i, d := range excel.Data {
		r := i + 2 //start from second row
		for k, v := range excel.DataMapping.Body {
			ax := fmt.Sprintf("%s%d", k, r)
			xlsx.SetCellValue(sheet, ax, d[v])
			xlsx.SetCellStyle(sheet, ax, ax, bodyStyle)
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

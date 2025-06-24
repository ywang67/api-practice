package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CableModemsByMac 是对应 GraphQL ByMac 的 RESTful 版本
func CableModemsByMac(c *gin.Context) {
	dbVal, ok := c.Get("dbRead")
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "database connection not available",
		})
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid database connection type",
		})
		return
	}
	// 获取 mac 参数（支持逗号分隔多个 mac）
	macParam := c.Query("mac")
	if macParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mac is required"})
		return
	}

	macs := strings.Split(macParam, ",")
	for i := range macs {
		macs[i] = strings.TrimSpace(macs[i])
	}

	// 调用 resolver
	modems, err := ByMacRds(c.Request.Context(), db, macs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modems)
}

// ByMacRds 是直接移植的 GraphQL resolver 核心逻辑
func ByMacRds(ctx context.Context, db *sql.DB, macAddresses []string) ([]*CableModem, error) {
	if len(macAddresses) == 0 {
		return nil, sql.ErrNoRows
	}

	// 分页批量 IN 查询
	return inRds(ctx, db, "mac", macAddresses, false)
}

func CableModemsByCmts(c *gin.Context) {
	// 获取 mac 参数
	macStr := c.Query("mac")
	if macStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mac is required"})
		return
	}
	macs := strings.Split(macStr, ",") // 支持多个 mac 用逗号分隔

	// TODO: 数据库查询 macs 相关数据

	c.JSON(http.StatusOK, gin.H{
		"macs": macs,
		// "data": 查询到的结果
	})
}

func CableModemsByPoller(c *gin.Context) {
	// 获取 mac 参数
	macStr := c.Query("mac")
	if macStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mac is required"})
		return
	}
	macs := strings.Split(macStr, ",") // 支持多个 mac 用逗号分隔

	// TODO: 数据库查询 macs 相关数据

	c.JSON(http.StatusOK, gin.H{
		"macs": macs,
		// "data": 查询到的结果
	})
}

func inRds(ctx context.Context, db *sql.DB, field string, values []string, single bool) ([]*CableModem, error) {
	if len(values) == 0 {
		return nil, errors.New("values slice is empty")
	}

	result := []*CableModem{}

	tempValue := values[0]
	checkPpodQuery := ""
	if field == "ppod" {
		tempValue = strings.ToUpper(values[0])
	}

	if !(strings.HasPrefix(tempValue, "acr") || strings.HasPrefix(tempValue, "cbr") || strings.HasPrefix(tempValue, "smi")) {
		checkPpodQuery = "AND ppod IS NOT NULL"
	}

	if single && len(values) > 0 {
		query := fmt.Sprintf(`
			SELECT * 
			FROM cablemodems
			WHERE %s = '%s'
			AND not_found_date IS NULL
			AND (
				(ipv4 IS NOT NULL AND ipv4 != '0.0.0.0' AND ipv4 != '') 
				OR 
				(ipv6 IS NOT NULL AND ipv6 != '0000:0000:0000:0000:0000:0000:0000:0000' AND ipv6 != '')
			)
			%s
			LIMIT 1;
		`, field, tempValue, checkPpodQuery)

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var cablemodems []*CableModem
		for rows.Next() {
			cablemodem := &CableModem{}
			err := rows.Scan(
				&cablemodem.Mac,
				&cablemodem.CpeMac,
				&cablemodem.MacDomain,
				&cablemodem.CableModemIndex,
				&cablemodem.ConfigFile,
				&cablemodem.Model,
				&cablemodem.FiberNode,
				&cablemodem.Ipv4,
				&cablemodem.Ipv6,
				&cablemodem.CpeIpv4,
				&cablemodem.Transponder,
				&cablemodem.DocsisVersion,
				&cablemodem.Ppod,
				&cablemodem.Fqdn,
				&cablemodem.State,
				&cablemodem.NotFoundDate,
				&cablemodem.RegState,
				&cablemodem.FnName,
				&cablemodem.NumberOfGenerators,
				&cablemodem.RpdName,
				&cablemodem.UpdatedAt,
				&cablemodem.Bootr,
				&cablemodem.Vendor,
				&cablemodem.SwRev,
				&cablemodem.OltName,
				&cablemodem.PonName,
				&cablemodem.UpdatedAtTs,
				&cablemodem.IsCpe,
				&cablemodem.CmtsType,
				&cablemodem.DeviceType,
			)
			if err != nil {
				return nil, err
			}
			cablemodems = append(cablemodems, cablemodem)
		}

		err = rows.Err()
		if err != nil {
			return nil, err
		}

		return cablemodems, nil
	}

	getRecordsIn := func(ctxArg context.Context, db *sql.DB, field string, values []string, page, pageSize int) ([]*CableModem, error) {
		offset := page * pageSize
		placeholderArgs := make([]string, len(values))
		queryArgs := make([]interface{}, len(values)+2)
		for i, value := range values {
			placeholderArgs[i] = fmt.Sprintf("$%d", i+1)
			queryArgs[i] = value
		}
		queryArgs[len(values)] = pageSize
		queryArgs[len(values)+1] = offset

		query := fmt.Sprintf(`
		SELECT *
		FROM cablemodems
		WHERE %s IN (%s)
		ORDER BY fqdn
		LIMIT $%d OFFSET $%d;
	`, field, strings.Join(placeholderArgs, ", "), len(values)+1, len(values)+2)

		rows, err := db.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var cablemodems []*CableModem
		for rows.Next() {
			cablemodem := &CableModem{}
			err := rows.Scan(
				&cablemodem.Mac,
				&cablemodem.CpeMac,
				&cablemodem.MacDomain,
				&cablemodem.CableModemIndex,
				&cablemodem.ConfigFile,
				&cablemodem.Model,
				&cablemodem.FiberNode,
				&cablemodem.Ipv4,
				&cablemodem.Ipv6,
				&cablemodem.CpeIpv4,
				&cablemodem.Transponder,
				&cablemodem.DocsisVersion,
				&cablemodem.Ppod,
				&cablemodem.Fqdn,
				&cablemodem.State,
				&cablemodem.NotFoundDate,
				&cablemodem.RegState,
				&cablemodem.FnName,
				&cablemodem.NumberOfGenerators,
				&cablemodem.RpdName,
				&cablemodem.UpdatedAt,
				&cablemodem.Bootr,
				&cablemodem.Vendor,
				&cablemodem.SwRev,
				&cablemodem.OltName,
				&cablemodem.PonName,
				&cablemodem.UpdatedAtTs,
				&cablemodem.IsCpe,
				&cablemodem.CmtsType,
				&cablemodem.DeviceType,
			)
			if err != nil {
				return nil, err
			}
			cablemodems = append(cablemodems, cablemodem)
		}

		err = rows.Err()
		if err != nil {
			return nil, err
		}

		return cablemodems, nil
	}

	page := 0
	pageSize := 100_000
	for {
		records, err := getRecordsIn(ctx, db, field, values, page, pageSize)
		if err != nil {
			return nil, err
		}

		if len(records) == 0 {
			break
		}

		result = append(result, records...)

		page++
	}

	return result, nil
}

package methods

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"api-project/grpc-api/gen/cablemodems"
	"api-project/grpc-api/helpers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CableModemMethod 实现 CableModemServiceServer 接口
type CableModemMethod struct {
	cablemodems.UnimplementedCableModemServiceServer
	Db *sql.DB
}

// 示例实现：GetCableModem
func (h *CableModemMethod) ByMac(ctx context.Context, req *cablemodems.ByMacRequest) (*cablemodems.ByMacResponse, error) {
	if h.Db == nil {
		return nil, status.Error(codes.Internal, "database connection is nil")
	}
	if len(req.MacAddress) == 0 {
		return nil, status.Error(codes.InvalidArgument, "macAddress list is empty")
	}

	// 构建查询语句
	placeholders := make([]string, len(req.MacAddress))
	args := make([]interface{}, len(req.MacAddress))
	for i, mac := range req.MacAddress {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = mac
	}

	query := fmt.Sprintf(`
		SELECT mac, cpe_mac, mac_domain, cable_modem_index, config_file, model,
		       fiber_node, ipv4, ipv6, cpe_ipv4, transponder, docsis_version,
		       ppod, fqdn, state, not_found_date, reg_state, fn_name,
		       number_of_generators, rpd_name, updated_at, bootr, vendor,
		       sw_rev, olt_name, pon_name, updated_at_ts, is_cpe, cmts_type, device_type
		FROM cablemodems
		WHERE mac IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := h.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query error: %v", err)
	}
	defer rows.Close()

	var modems []*cablemodems.CableModem

	for rows.Next() {
		var (
			m         cablemodems.CableModem
			stateStr  string
			docsisStr string
		)

		err := rows.Scan(
			&m.Mac,
			&m.CpeMac,
			&m.MacDomain,
			&m.CableModemIndex,
			&m.ConfigFile,
			&m.Model,
			&m.FiberNode,
			&m.Ipv4,
			&m.Ipv6,
			&m.CpeIpv4,
			&m.Transponder,
			&docsisStr,
			&m.Ppod,
			&m.Fqdn,
			&stateStr,
			&m.NotFoundDate,
			&m.RegState,
			&m.FnName,
			&m.NumberOfGenerators,
			&m.RpdName,
			&m.UpdatedAt,
			&m.Bootr,
			&m.Vendor,
			&m.SwRev,
			&m.OltName,
			&m.PonName,
			&m.UpdatedAtTs,
			&m.IsCpe,
			&m.CmtsType,
			&m.DeviceType,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "row scan error: %v", err)
		}

		stateEnum, err := helpers.ParseStateFromString(stateStr)
		if err != nil {
			stateEnum = cablemodems.State_UNKNOWN
		}

		docsisEnum, err := helpers.ParseDocsisVersionFromString(docsisStr)
		if err != nil {
			docsisEnum = cablemodems.DocsisVersion_DOCSIS_UNKNOWN
		}

		m.State = &stateEnum
		m.DocsisVersion = &docsisEnum

		modems = append(modems, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "rows iteration error: %v", err)
	}

	return &cablemodems.ByMacResponse{
		Modems: modems,
		Error:  nil,
	}, nil
}

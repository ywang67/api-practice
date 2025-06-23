package cablemodems

import (
	"api-project/graphql-api/gql/graph/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// CableModems does nothing and just exists to generate the nice namespacing for gqlgen.
// i.e, query{transponders{byBucket(200)}}{FiberNode}
type CableModems struct{}

func ByMacRds(ctx context.Context, db *sql.DB, macAddresses []string) (modems []*model.CableModem, err error) {
	if db == nil {
		return nil, errors.New("nil db connection")
	}

	modems, err = inRds(ctx, db, "mac", macAddresses, false)

	return modems, err
}

func inRds(ctx context.Context, db *sql.DB, field string, values []string, single bool) ([]*model.CableModem, error) {
	if len(values) == 0 {
		return nil, errors.New("values slice is empty")
	}

	result := []*model.CableModem{}

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

		var cablemodems []*model.CableModem
		for rows.Next() {
			cablemodem := &model.CableModem{}
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

	getRecordsIn := func(ctxArg context.Context, db *sql.DB, field string, values []string, page, pageSize int) ([]*model.CableModem, error) {
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

		var cablemodems []*model.CableModem
		for rows.Next() {
			cablemodem := &model.CableModem{}
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

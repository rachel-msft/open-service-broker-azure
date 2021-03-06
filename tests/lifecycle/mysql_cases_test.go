// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL SQL driver
)

var mysqlTestCases = []serviceLifecycleTestCase{
	{
		group:     "mysql",
		name:      "all-in-one",
		serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
		planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
		location:  "southcentralus",
		provisioningParameters: service.CombinedProvisioningParameters{
			"sslEnforcement": "disabled",
			"firewallRules": []map[string]string{
				{
					"name":           "AllowSome",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "35.0.0.0",
				},
				{
					"name":           "AllowMore",
					"startIPAddress": "35.0.0.1",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		testCredentials: testMySQLCreds,
	},
	{
		group:     "mysql",
		name:      "dbms-only",
		serviceID: "30e7b836-199d-4335-b83d-adc7d23a95c2",
		planID:    "3f65ebf9-ac1d-4e77-b9bf-918889a4482b",
		location:  "eastus",
		provisioningParameters: service.CombinedProvisioningParameters{
			"firewallRules": []map[string]string{
				{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "mysql",
				name:            "database-only",
				serviceID:       "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
				planID:          "ec77bd04-2107-408e-8fde-8100c1ce1f46",
				location:        "", // This is actually irrelevant for this test
				testCredentials: testMySQLCreds,
			},
		},
	},
}

func testMySQLCreds(credentials map[string]interface{}) error {

	var connectionStrTemplate string
	if credentials["sslRequired"].(bool) {
		connectionStrTemplate =
			"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=true"
	} else {
		connectionStrTemplate =
			"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true"
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		connectionStrTemplate,
		credentials["username"].(string),
		credentials["password"].(string),
		credentials["host"].(string),
		credentials["database"].(string),
	))
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}
	defer db.Close() // nolint: errcheck
	rows, err := db.Query("SELECT * from INFORMATION_SCHEMA.TABLES")
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error could not select from INFORMATION_SCHEMA.TABLES'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			`error iterating rows`,
		)
	}
	return nil
}

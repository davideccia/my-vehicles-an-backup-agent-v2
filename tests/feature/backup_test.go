package feature

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"my_vehicles_an_backup_agent/tests"
)

// BackupTestSuite exercises the /backup/{table}, /dashboard/stats and /ping
// contract against the Android app's expectations: mirror-sync semantics
// (upsert + delete-not-in), the API key gate, the schema version gate, and
// the strict envelope validation (unknown fields, duplicate ids).
type BackupTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestBackupTestSuite(t *testing.T) {
	suite.Run(t, new(BackupTestSuite))
}

func (s *BackupTestSuite) SetupTest() {
	s.RefreshDatabase()
}

// testAPIKey matches .env's API_KEY, which is loaded once at boot (see
// tests/test_case.go's init()).
const testAPIKey = "changeme"

func vehicleRow(id, plate string) map[string]any {
	return map[string]any{
		"id":           id,
		"plate_number": plate,
		"brand":        "Fiat",
		"model":        "Panda",
		"year":         2020,
		"fuel_type":    "PETROL",
		"fuel_unit":    "LITERS",
		"color":        "#4CAF50",
	}
}

func envelope(version string, rows ...map[string]any) *bytes.Reader {
	if rows == nil {
		rows = []map[string]any{}
	}
	buf, err := json.Marshal(map[string]any{"version": version, "data": rows})
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(buf)
}

func (s *BackupTestSuite) TestPostThenGetRoundTrip() {
	postResp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA"), vehicleRow("v2", "BB222BB")))
	s.Require().NoError(err)
	postResp.AssertCreated()

	postJSON, err := postResp.Json()
	s.Require().NoError(err)
	s.EqualValues(2, postJSON["upserted_count"])
	s.EqualValues(0, postJSON["deleted_count"])

	getResp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Get("/backup/vehicles?version=6")
	s.Require().NoError(err)
	getResp.AssertOk()

	getJSON, err := getResp.Json()
	s.Require().NoError(err)
	s.Equal("6", getJSON["version"])
	data, ok := getJSON["data"].([]any)
	s.Require().True(ok)
	s.Len(data, 2)
}

func (s *BackupTestSuite) TestSecondPostDropsMissingRecords() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA"), vehicleRow("v2", "BB222BB"), vehicleRow("v3", "CC333CC")))
	s.Require().NoError(err)
	resp.AssertCreated()

	// v3 is dropped, v4 is new: v1/v2 survive, v3 is deleted, v4 is inserted.
	resp, err = s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA"), vehicleRow("v2", "BB222BB"), vehicleRow("v4", "DD444DD")))
	s.Require().NoError(err)
	resp.AssertCreated()

	body, err := resp.Json()
	s.Require().NoError(err)
	s.EqualValues(3, body["upserted_count"])
	s.EqualValues(1, body["deleted_count"])

	getResp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Get("/backup/vehicles?version=6")
	s.Require().NoError(err)
	getJSON, err := getResp.Json()
	s.Require().NoError(err)
	data, _ := getJSON["data"].([]any)
	s.Len(data, 3)
}

func (s *BackupTestSuite) TestEmptyPostWipesTable() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA")))
	s.Require().NoError(err)
	resp.AssertCreated()

	resp, err = s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Post("/backup/vehicles", envelope("6"))
	s.Require().NoError(err)
	resp.AssertCreated()

	body, err := resp.Json()
	s.Require().NoError(err)
	s.EqualValues(0, body["upserted_count"])
	s.EqualValues(1, body["deleted_count"])
}

func (s *BackupTestSuite) TestDuplicateIDIsRejected() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA"), vehicleRow("v1", "BB222BB")))
	s.Require().NoError(err)
	resp.AssertBadRequest()
}

func (s *BackupTestSuite) TestUnknownFieldIsRejected() {
	row := vehicleRow("v1", "AA111AA")
	row["plateNumber"] = "camelCaseGhost" // unknown key: the correct one is plate_number
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Post("/backup/vehicles", envelope("6", row))
	s.Require().NoError(err)
	resp.AssertBadRequest()
}

func (s *BackupTestSuite) TestMissingRequiredFieldIsRejected() {
	row := vehicleRow("v1", "AA111AA")
	delete(row, "color")
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Post("/backup/vehicles", envelope("6", row))
	s.Require().NoError(err)
	resp.AssertBadRequest()
}

func (s *BackupTestSuite) TestInvalidEnumIsRejected() {
	row := vehicleRow("v1", "AA111AA")
	row["fuel_type"] = "ELECTRIC"
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Post("/backup/vehicles", envelope("6", row))
	s.Require().NoError(err)
	resp.AssertBadRequest()
}

func (s *BackupTestSuite) TestVersionMismatchIsRejectedAndDoesNotWrite() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("5", vehicleRow("v1", "AA111AA")))
	s.Require().NoError(err)
	resp.AssertConflict()

	getResp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Get("/backup/vehicles?version=6")
	s.Require().NoError(err)
	getJSON, err := getResp.Json()
	s.Require().NoError(err)
	data, _ := getJSON["data"].([]any)
	s.Empty(data)
}

func (s *BackupTestSuite) TestMissingAPIKeyIsRejected() {
	resp, err := s.Http(s.T()).Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA")))
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *BackupTestSuite) TestWrongAPIKeyIsRejected() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", "wrong-key").
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA")))
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *BackupTestSuite) TestHealthcheckIsPublicAndUnversioned() {
	resp, err := s.Http(s.T()).Get("/healthcheck")
	s.Require().NoError(err)
	resp.AssertOk()

	body, err := resp.Json()
	s.Require().NoError(err)
	s.Equal("ok", body["status"])
}

func (s *BackupTestSuite) TestPingIsPublicAndVersionGated() {
	resp, err := s.Http(s.T()).Get("/ping?version=6")
	s.Require().NoError(err)
	resp.AssertOk()

	body, err := resp.Json()
	s.Require().NoError(err)
	s.Equal("ok", body["status"])
	s.Equal("6", body["schema_version"])

	resp, err = s.Http(s.T()).Get("/ping?version=5")
	s.Require().NoError(err)
	resp.AssertConflict()
}

func (s *BackupTestSuite) TestDashboardStats() {
	resp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).
		Post("/backup/vehicles", envelope("6", vehicleRow("v1", "AA111AA"), vehicleRow("v2", "BB222BB")))
	s.Require().NoError(err)
	resp.AssertCreated()

	statsResp, err := s.Http(s.T()).WithHeader("x-api-key", testAPIKey).Get("/dashboard/stats?version=6")
	s.Require().NoError(err)
	statsResp.AssertOk()

	statsJSON, err := statsResp.Json()
	s.Require().NoError(err)
	data, ok := statsJSON["data"].(map[string]any)
	s.Require().True(ok)
	s.EqualValues(2, data["vehicles"])
	s.EqualValues(0, data["refuels"])
}

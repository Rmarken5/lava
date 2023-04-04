package data_access

//
//import (
//	"database/sql"
//	"fmt"
//	"log"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//)
//
//func TestInspectorImpl_InspectTable(t *testing.T) {
//
//	testCases := map[string]struct {
//		name           string
//		expectedResult Table
//		expectedError  error
//	}{
//		"ValidTable": {
//			name: "my_table",
//			expectedResult: Table{
//				Columns: []Column{
//					{Name: "id", ColumnType: &sql.ColumnType{}},
//					{Name: "name", ColumnType: &sql.ColumnType{DatabaseTypeName: "varchar"}},
//					{Name: "age", ColumnType: &sql.ColumnType{DatabaseTypeName: "integer"}},
//				},
//			},
//			expectedError: nil,
//		},
//		"InvalidTable": {
//			name:           "non_existent_table",
//			expectedResult: Table{},
//			expectedError:  fmt.Errorf("sql: no rows in result set"),
//		},
//	}
//
//	for name, testCase := range testCases {
//		t.Run(name, func(t *testing.T) {
//			db, mock, err := sqlmock.New()
//			if err != nil {
//				t.Fatalf("Failed to create mock database: %s", err)
//			}
//			defer db.Close()
//
//			expectedQuery := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '%s'", testCase.name)
//			rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "Alice", 30)
//			if testCase.expectedError != nil {
//				rows = sqlmock.NewRows([]string{})
//			}
//			mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
//
//			impl := &InspectorImpl{db: *db, Logger: *log.Default()}
//			result, err := impl.InspectTable(testCase.name)
//			if err != testCase.expectedError {
//				t.Errorf("Expected error %v, got %v", testCase.expectedError, err)
//			}
//			if fmt.Sprintf("%+v", result) != fmt.Sprintf("%+v", testCase.expectedResult) {
//				t.Errorf("Expected %+v, got %+v", testCase.expectedResult, result)
//			}
//			mock.ExpectationsWereMet()
//		})
//	}
//}

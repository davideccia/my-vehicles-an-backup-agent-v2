// Package services holds the backup sync logic shared by every /backup/{table}
// endpoint.
package services

import (
	ormcontract "github.com/goravel/framework/contracts/database/orm"

	"my_vehicles_an_backup_agent/app/facades"
)

// chunkSize caps how many rows go into a single INSERT statement, to stay
// under each dialect's bind-variable ceiling (SQLite: 999, MySQL: 65535,
// Postgres: 65535) regardless of how many columns a table has.
const chunkSize = 200

// SyncTable replaces the entire contents of a table with records, inside one
// transaction. This mirrors the NestJS reference server's semantics — the
// posted array becomes the table's full content — but implemented as
// delete-all + insert-all rather than upsert + delete-not-in, because:
//
//   - Goravel v1.18 has no bulk upsert (only a row-at-a-time UpdateOrCreate),
//     so a literal port would mean one query per row instead of one query
//     per chunk.
//   - None of these tables have server-managed columns (no created_at/
//     updated_at) that a partial upsert would need to preserve across
//     unchanged rows, so the two approaches leave the table in the same
//     final state.
//
// model must be a pointer to the model being synced (e.g. &models.Vehicle{});
// records is the full replacement content; ids are the records' ids, in the
// same order, used only to compute deletedCount before the delete happens.
func SyncTable[T any](model *T, records []T, ids []string) (upsertedCount, deletedCount int, err error) {
	err = facades.Orm().Transaction(func(tx ormcontract.Query) error {
		var existingIDs []string
		if pluckErr := tx.Model(model).Pluck("id", &existingIDs); pluckErr != nil {
			return pluckErr
		}

		keep := make(map[string]struct{}, len(ids))
		for _, id := range ids {
			keep[id] = struct{}{}
		}
		for _, id := range existingIDs {
			if _, ok := keep[id]; !ok {
				deletedCount++
			}
		}

		// GORM refuses an unconditional Delete (global-update guard), even
		// though the contract's Delete doc says "if the conditions are empty
		// will delete all records" — so make the condition explicit.
		if _, delErr := tx.Model(model).Where("1 = 1").Delete(); delErr != nil {
			return delErr
		}

		for start := 0; start < len(records); start += chunkSize {
			end := start + chunkSize
			if end > len(records) {
				end = len(records)
			}

			chunk := records[start:end]
			if createErr := tx.Create(&chunk); createErr != nil {
				return createErr
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return len(records), deletedCount, nil
}

// FetchTable returns the entire contents of a table, unordered and
// unpaginated — mirroring the NestJS reference server's `repo.find()`.
func FetchTable[T any](model *T) ([]T, error) {
	var records []T
	if err := facades.Orm().Query().Model(model).Get(&records); err != nil {
		return nil, err
	}
	return records, nil
}

// CountTable returns the row count of a table, for the dashboard stats
// endpoint's five `COUNT(*)` queries.
func CountTable[T any](model *T) (int64, error) {
	return facades.Orm().Query().Model(model).Count()
}

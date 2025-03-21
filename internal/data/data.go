package data

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type ContextKey string

const ContextKeyStore ContextKey = "ContextKeyStore"

type Store struct {
	User         UserStore
	Organization OrganizationStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		User:         NewUserStore(db),
		Organization: NewOrganizationStore(db),
	}
}

func WithStore(ctx context.Context, store *Store) context.Context {
	return context.WithValue(ctx, ContextKeyStore, store)
}

func GetStore(ctx context.Context) (*Store, bool) {
	store, ok := ctx.Value(ContextKeyStore).(*Store)
	return store, ok
}

// BuildInsertQuery builds an INSERT query for the given table and data map.
// It returns a query string with numbered placeholders and a slice of argument values.
// Note the use of "RETURNING *", which allows you to return the full inserted row.
//
// Example usage:
//
//	data := map[string]any{
//	     "name": "Alice",
//	     "age":  30,
//	}
//	query, args := BuildInsertQuery("users", data)
//	// query => "INSERT INTO users (age, name) VALUES ($1, $2) RETURNING *"
//	// args  => []any{30, "Alice"}
func BuildInsertQuery(tableName string, data map[string]any) (string, []any, error) {
	if !isValidTable(tableName) {
		return "", nil, fmt.Errorf("invalid table name: %s", tableName)
	}
	if len(data) == 0 {
		return "", nil, fmt.Errorf("no data provided for insert query")
	}

	// Use sorted keys to ensure deterministic output.
	keys := sortedKeys(data)

	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]any, 0, len(data))
	for i, col := range keys {
		columns = append(columns, col)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		values = append(values, data[col])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values, nil
}

// BuildUpdateQuery builds an UPDATE query for a given table, data map and conditions map.
// Both data and conditions maps are sorted alphabetically (to guarantee consistent ordering)
// and then converted to placeholder queries. A non-empty conditions map is required to prevent accidental updates.
// The returned query uses "RETURNING *" to retrieve the updated row.
//
// Example usage:
//
//	updateData := map[string]any{
//	     "name": "Bob",
//	}
//	conditions := map[string]any{
//	     "id": 1,
//	}
//	query, args := BuildUpdateQuery("users", updateData, conditions)
//	// query => "UPDATE users SET name = $1 WHERE id = $2 RETURNING *"
//	// args  => []any{"Bob", 1}
func BuildUpdateQuery(tableName string, updateData, conditions map[string]any) (string, []any, error) {
	if !isValidTable(tableName) {
		return "", nil, fmt.Errorf("invalid table name: %s", tableName)
	}
	if len(updateData) == 0 {
		return "", nil, fmt.Errorf("update data cannot be empty")
	}
	if len(conditions) == 0 {
		return "", nil, fmt.Errorf("conditions cannot be empty for update query")
	}

	dataKeys := sortedKeys(updateData)
	conditionKeys := sortedKeys(conditions)

	setClauses := make([]string, 0, len(updateData))
	values := make([]any, 0, len(updateData)+len(conditions))
	// Build SET clause using sorted data keys.
	for i, col := range dataKeys {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i+1))
		values = append(values, updateData[col])
	}

	// Build WHERE clause using sorted condition keys.
	whereClauses := make([]string, 0, len(conditions))
	for j, col := range conditionKeys {
		// Note: j+1+len(updateData) is used to start numbering after data placeholders.
		placeholderIdx := len(updateData) + j + 1
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", col, placeholderIdx))
		values = append(values, conditions[col])
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s RETURNING *",
		tableName,
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "),
	)

	return query, values, nil
}

// BuildSelectQuery builds a generic SELECT query for the given table.
// If a non-empty conditions map is provided, it will be used to generate a WHERE clause.
// This query returns all columns, making it a good match for a GET endpoint.
//
// Example usage:
//
//	conditions := map[string]any{
//	     "id": 1,
//	}
//	query, args := BuildSelectQuery("users", conditions)
//	// query => "SELECT * FROM users WHERE id = $1"
//	// args  => []any{1}
func BuildSelectQuery(tableName string, conditions map[string]any) (string, []any, error) {
	if !isValidTable(tableName) {
		return "", nil, fmt.Errorf("invalid table name: %s", tableName)
	}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	values := []any{}

	if len(conditions) > 0 {
		condKeys := sortedKeys(conditions)
		whereClauses := make([]string, 0, len(conditions))

		for i, col := range condKeys {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", col, i+1))
			values = append(values, conditions[col])
		}

		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(whereClauses, " AND "))
	}

	return query, values, nil
}

var validTables = map[string]struct{}{
	"organizations":      {},
	"uld_inventories":    {},
	"delivery_manifests": {},
	"manifest_items":     {},
	"warehouses":         {},
	"airlines":           {},
	"carriers":           {},
	"users":              {},
	"user_associations":  {},
}

func isValidTable(tableName string) bool {
	_, ok := validTables[tableName]
	return ok
}

func sortedKeys(data map[string]any) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GenerateRandomString(n int) string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

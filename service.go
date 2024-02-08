package qparser

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// parseQuery parses the given query string and returns a Field object representing the parsed query.
// The query string should be in the format "operator:value".
// The name parameter specifies the name of the field being queried.
// If the query string is not in the correct format, an error is returned.
func parseQuery(name, query string) (*Field, error) {
	args := strings.Split(query, ":")
	if len(args) < 2 {
		return nil, fmt.Errorf("bad query, use operator:value")
	}

	if len(strings.Split(args[0], " ")) > 1 {
		return nil, fmt.Errorf("bad query, use operator:value")
	}

	operator, err := convertOperator(args[0])
	if err != nil {
		return nil, err

	}

	return &Field{
		Name:     name,
		Operator: operator,
		Value:    strings.Join(args[1:], " "),
	}, nil
}

// ParseStruct parses the given data and returns an Options struct and an error.
// It iterates over the fields of the data structure and populates the Options struct accordingly.
// The "query" tag is used to specify the behavior for each field.
// The "limit" tag is used to set the limit value for the Options struct.
// The "offset" tag is used to set the offset value for the Options struct.
// For other fields, the parseQuery function is used to parse the field value and add it to the Options struct.
// If any parsing or validation error occurs, an error is returned.
func ParseStruct(data interface{}) (*Options, error) {
	filterValue := reflect.ValueOf(data)
	filterType := filterValue.Type()

	opt := &Options{
		limit:  0,
		offset: 0,
		fields: make([]*Field, 0),
	}

	for i := 0; i < filterType.NumField(); i++ {
		field := filterType.Field(i)
		value := filterValue.Field(i)

		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		tag := field.Tag.Get("query")
		fieldValue := reflect.Indirect(value).Interface()

		switch tag {
		case "limit":
			l, ok := fieldValue.(int)

			if !ok {
				return nil, fmt.Errorf("failed to parse limit")
			}

			if l < 0 {
				return nil, fmt.Errorf("limit must be greater than 0")
			}

			opt.limit = l

			continue
		case "offset":
			o, ok := fieldValue.(int)

			if !ok {
				return nil, fmt.Errorf("failed to parse offset")
			}

			if o < 0 {
				return nil, fmt.Errorf("offset must be greater than 0")
			}

			opt.offset = o

			continue
		}

		switch field.Type {
		case reflect.TypeOf((*bool)(nil)):
			{
				if err := opt.AddField(tag, fmt.Sprint(fieldValue), operatorEqual); err != nil {
					return nil, err
				}
			}
		default:
			{
				fieldValueStr := fmt.Sprint(fieldValue)

				if len(fieldValueStr) == 0 {
					continue
				}

				field, err := parseQuery(tag, fieldValueStr)
				if err != nil {
					return nil, err
				}

				if err := opt.AddField(field.Name, field.Value, field.Operator); err != nil {
					return nil, err
				}
			}
		}
	}

	return opt, nil
}

// validateOperator validates the given operator string.
// It checks if the operator is one of the supported SQL operators.
// If the operator is not supported, it returns an error.
func validateOperator(operator string) error {
	switch operator {
	case sqlOperatorEqual:
	case sqlOperatorNotEqual:
	case sqlOperatorGreaterThan:
	case sqlOperatorGreaterThanEqual:
	case sqlOperatorLowerThan:
	case sqlOperatorLowerThanEqual:
	case sqlOperatorLike:
	case sqlOperatorRange:
	default:
		return fmt.Errorf("bad operator")
	}
	return nil
}

// convertOperator converts a given operator string to its corresponding SQL operator.
// It returns the SQL operator as a string and an error if the operator is not recognized.
func convertOperator(operator string) (string, error) {
	switch operator {
	case operatorEqual:
		return sqlOperatorEqual, nil
	case operatorNotEqual:
		return sqlOperatorNotEqual, nil
	case operatorGreaterThan:
		return sqlOperatorGreaterThan, nil
	case operatorGreaterThanEqual:
		return sqlOperatorGreaterThanEqual, nil
	case operatorLowerThan:
		return sqlOperatorLowerThan, nil
	case operatorLowerThanEqual:
		return sqlOperatorLowerThanEqual, nil
	case operatorLike:
		return sqlOperatorLike, nil
	case operatorRange:
		return sqlOperatorRange, nil
	default:
		return "", fmt.Errorf("bad operator")
	}
}

// AddField adds a field to the Options struct.
// It takes the name, value, and operator of the field as parameters.
// The operator is validated, and if it is invalid, an error is returned.
// If the operator is "like" and the value does not contain "%", the value is modified to include "%" at the beginning and end.
// If the operator is "range", the value is split into two parts using " to " as the delimiter.
// If the value does not contain exactly two parts, an error is returned.
// The field is then appended to the fields slice in the Options struct.
// Returns nil if successful, otherwise returns an error.
func (o *Options) AddField(name, value, operator string) error {
	if err := validateOperator(operator); err != nil {
		return err
	}

	if operator == sqlOperatorLike && !strings.ContainsAny(value, "%") {
		value = fmt.Sprintf("%%%s%%", value)
	}

	if operator == sqlOperatorRange {
		args := strings.Split(value, " to ")
		if len(args) != 2 {
			return fmt.Errorf("invalid usage of operator rng. rng:value1:to:value2")
		}

		value = fmt.Sprintf("%s %s", args[0], args[1])
	}

	o.fields = append(o.fields, &Field{
		Name:     name,
		Value:    value,
		Operator: operator,
	})

	return nil
}

// Apply applies the options to the given GORM transaction.
// It iterates through each option and applies the corresponding condition to the transaction.
// If the option's operator is "range", it splits the option value by space and applies a range condition.
// Otherwise, it applies a regular condition using the option's name, operator, and value.
// It also sets the offset and limit of the transaction based on the options.
// Finally, it returns the modified transaction.
func (o *Options) Apply(tx *gorm.DB) *gorm.DB {
	for _, option := range o.fields {
		if option.Operator == sqlOperatorRange {
			args := strings.Split(option.Value, " ")

			tx = tx.Where(fmt.Sprintf("%s %s ? AND ?", option.Name, option.Operator), args[0], args[1])
			continue
		}

		tx = tx.Where(fmt.Sprintf("%s %s ?", option.Name, option.Operator), option.Value)
	}

	tx = tx.Offset(o.offset)

	if o.limit > 0 {
		tx = tx.Limit(o.limit)
	}

	tx = tx.Offset(o.offset)

	return tx
}

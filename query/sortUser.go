package query

import (
	"errors"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/ONSdigital/dp-identity-api/models"
)

// LessFunc used by MultiSorter OrderedBy  used to hold the seq of sort
type LessFunc func(p1, p2 *models.UserParams) bool

// MultiSorter structure to hold input array and the rwquest query parameters converted see GetLessFunc
type MultiSorter struct {
	changes []models.UserParams
	less    []LessFunc
}

// Sort some description of the function
func (ms *MultiSorter) Sort(changes []models.UserParams) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy some description of the function
func OrderedBy(less ...LessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}

// Len function to produce length of changes required by Sort third party
func (ms *MultiSorter) Len() int {
	return len(ms.changes)
}

// Swap ms  called by sort to swap two records by sort parameters
func (ms *MultiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less either swaps the two concurrent values by the sort rules
func (ms *MultiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

// SortBy from the request query get the sort parameters
func SortBy(requestSortParameters string, arr []models.UserParams) error {
	inputSplit := strings.Split(requestSortParameters, ",")
	var (
		orderFunc             = make([]LessFunc, 0, len(inputSplit))
		v         interface{} = arr[0]
	)

	for _, inputSplitItem := range inputSplit {
		inputSplitItemSplit := strings.Split(inputSplitItem, ":")
		IsDesc := slices.Contains(inputSplitItemSplit, "desc")
		userField, err := GetFieldByJSONTag(strings.ToLower(inputSplitItemSplit[0]), v)
		if err != nil {
			return err
		}
		orderFunc = append(orderFunc, GetLessFunc(getType(v), userField.Name, IsDesc))
	}
	OrderedBy(orderFunc...).Sort(arr)
	return nil
}

// getType returns the type of interface as string
func getType(myVar interface{}) string {
	return reflect.TypeOf(myVar).Name()
}

// GetFieldByJSONTag returns the field name as a string from the json value supplied by the request query params
func GetFieldByJSONTag(jsonTagValue string, s interface{}) (reflect.StructField, error) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return reflect.StructField{}, errors.New("incorrect structure")
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("json"), ",")[0]
		if v == jsonTagValue {
			return rt.Field(i), nil
		}
	}
	return reflect.StructField{}, errors.New(" request query sort parameter not found " + jsonTagValue)
}

// GetLessFunc supplies the output function from the
func GetLessFunc(name, field string, direction bool) LessFunc {
	if direction {
		if name == "UserParams" {
			switch field {
			case "Forename":
				return func(c1, c2 *models.UserParams) bool { return c1.Forename > c2.Forename }
			case "Lastname":
				return func(c1, c2 *models.UserParams) bool { return c1.Lastname > c2.Lastname }
			case "Email":
				return func(c1, c2 *models.UserParams) bool { return c1.Email > c2.Email }
			default:
				return func(c1, c2 *models.UserParams) bool { return c1.ID > c2.ID }
			}
		}
	} else {
		if name == "UserParams" {
			switch field {
			case "Forename":
				return func(c1, c2 *models.UserParams) bool { return c1.Forename < c2.Forename }
			case "Lastname":
				return func(c1, c2 *models.UserParams) bool { return c1.Lastname < c2.Lastname }
			case "Email":
				return func(c1, c2 *models.UserParams) bool { return c1.Email < c2.Email }
			default:
				return func(c1, c2 *models.UserParams) bool { return c1.ID < c2.ID }
			}
		}
	}
	return func(c1, c2 *models.UserParams) bool { return c1.ID < c2.ID }
}

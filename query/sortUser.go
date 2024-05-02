package query

import (
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// lessFunc used by multiSorter OrderedBy  used to hold the seq of sort
type lessFunc func(p1, p2 *models.UserParams) bool

// multiSorter
type multiSorter struct {
	changes []models.UserParams
	less    []lessFunc
}

// Sort some description of the function
func (ms *multiSorter) Sort(changes []models.UserParams) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy some descrription of the function
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

// Swap ms  called by sort to swap two records by sort parameters
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less ms
func (ms *multiSorter) Less(i, j int) bool {
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
func SortBy(requestSortParameters string, arr []models.UserParams) {
	//forenameAsc := func(c1, c2 *models.UserParams) bool { return c1.Forename < c2.Forename }
	//forenameDesc := func(c1, c2 *models.UserParams) bool { return c1.Forename > c2.Forename }
	//lastnameAsc := func(c1, c2 *models.UserParams) bool { return c1.Lastname < c2.Lastname }
	//lastnameDesc := func(c1, c2 *models.UserParams) bool { return c1.Lastname > c2.Lastname }
	//emailAsc := func(c1, c2 *models.UserParams) bool { return c1.Email < c2.Email }
	//emailDesc := func(c1, c2 *models.UserParams) bool { return c1.Email > c2.Email }
	//idAsc := func(c1, c2 *models.UserParams) bool { return c1.ID < c2.ID }
	//idDesc := func(c1, c2 *models.UserParams) bool { return c1.ID > c2.ID }

	inputSplit := strings.Split(requestSortParameters, ",")

	var (
		orderFunc  []lessFunc
		v          interface{}
		userparams models.UserParams
	)
	v = userparams

	for _, inputSplitItem := range inputSplit {
		inputSplitItemSplit := strings.Split(inputSplitItem, ":")
		IsDesc := slices.Contains(inputSplitItemSplit, "desc")
		userfield := GetBooleanFieldValueByJsonTag(strings.ToLower(inputSplitItemSplit[0]), v)
		fmt.Println(userfield)
		if strings.ToLower(inputSplitItemSplit[0]) == "forename" {
			if IsDesc {
				orderFunc = append(orderFunc, func(c1, c2 *models.UserParams) bool { return c1.Forename > c2.Forename })
			} else {
				orderFunc = append(orderFunc, func(c1, c2 *models.UserParams) bool { return c1.Forename < c2.Forename })
			}

		}

		OrderedBy(orderFunc...).Sort(arr)

	}
}

func GetBooleanFieldValueByJsonTag(jsonTagValue string, s interface{}) bool {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == jsonTagValue {
			r := reflect.ValueOf(s)
			field := reflect.Indirect(r).FieldByName(f.Name)
			return field.Bool()
		}
	}
	return false
}

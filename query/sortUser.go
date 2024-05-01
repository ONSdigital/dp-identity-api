package query

import (
	"github.com/ONSdigital/dp-identity-api/models"
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

// ms
func (ms *multiSorter) Sort(changes []models.UserParams) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// ms
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
	forenameAsc := func(c1, c2 *models.UserParams) bool { return c1.Forename < c2.Forename }
	forenameDesc := func(c1, c2 *models.UserParams) bool { return c1.Forename > c2.Forename }
	lastnameAsc := func(c1, c2 *models.UserParams) bool { return c1.Lastname < c2.Lastname }
	lastnameDesc := func(c1, c2 *models.UserParams) bool { return c1.Lastname > c2.Lastname }
	emailAsc := func(c1, c2 *models.UserParams) bool { return c1.Email < c2.Email }
	emailDesc := func(c1, c2 *models.UserParams) bool { return c1.Email > c2.Email }
	idAsc := func(c1, c2 *models.UserParams) bool { return c1.ID < c2.ID }
	idDesc := func(c1, c2 *models.UserParams) bool { return c1.ID > c2.ID }

	inputSplit := strings.Split(requestSortParameters, ",")

	var orderFunc []lessFunc
	for _, inputSplitItem := range inputSplit {
		inputSplitItemSplit := strings.Split(inputSplitItem, ":")
		if len(inputSplitItemSplit) > 1 {
			if inputSplitItemSplit[0] == "forename" {
				if inputSplitItemSplit[1] == "asc" {
					orderFunc = append(orderFunc, forenameAsc)
				} else {
					orderFunc = append(orderFunc, forenameDesc)
				}
			} else if inputSplitItemSplit[0] == "lastname" {
				if inputSplitItemSplit[1] == "asc" {
					orderFunc = append(orderFunc, lastnameAsc)
				} else {
					orderFunc = append(orderFunc, lastnameDesc)
				}
			} else if inputSplitItemSplit[0] == "email" {
				if inputSplitItemSplit[1] == "asc" {
					orderFunc = append(orderFunc, emailAsc)
				} else {
					orderFunc = append(orderFunc, emailDesc)
				}
			} else if inputSplitItemSplit[0] == "id" {
				if inputSplitItemSplit[1] == "asc" {
					orderFunc = append(orderFunc, idAsc)
				} else {
					orderFunc = append(orderFunc, idDesc)
				}
			}
		} else {
			if inputSplitItemSplit[0] == "forename" {
				orderFunc = append(orderFunc, forenameAsc)
			} else if inputSplitItemSplit[0] == "lastname" {
				orderFunc = append(orderFunc, lastnameAsc)
			} else if inputSplitItemSplit[0] == "email" {
				orderFunc = append(orderFunc, emailAsc)
			} else {
				orderFunc = append(orderFunc, idAsc)
			}
		}

		OrderedBy(orderFunc...).Sort(arr)

	}
}

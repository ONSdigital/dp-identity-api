package query

import (
	"github.com/ONSdigital/dp-identity-api/models"
	"sort"
	"strings"
)

type lessFunc func(p1, p2 *models.UserParams) bool

type multiSorter struct {
	changes []models.UserParams
	less    []lessFunc
}

func (ms *multiSorter) Sort(changes []models.UserParams) {
	ms.changes = changes
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}
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

func SortBy(paramsSlice string, arr []models.UserParams) {
	forenameAsc := func(c1, c2 *models.UserParams) bool { return c1.Forename < c2.Forename }
	forenameDesc := func(c1, c2 *models.UserParams) bool { return c1.Forename > c2.Forename }
	lastnameAsc := func(c1, c2 *models.UserParams) bool { return c1.Lastname < c2.Lastname }
	lastnameDesc := func(c1, c2 *models.UserParams) bool { return c1.Lastname > c2.Lastname }
	emailAsc := func(c1, c2 *models.UserParams) bool { return c1.Email < c2.Email }
	emailDesc := func(c1, c2 *models.UserParams) bool { return c1.Email > c2.Email }
	idAsc := func(c1, c2 *models.UserParams) bool { return c1.ID < c2.ID }
	idDesc := func(c1, c2 *models.UserParams) bool { return c1.ID > c2.ID }

	p1 := strings.Split(paramsSlice, ",")

	var orderFunc []lessFunc
	for _, p2 := range p1 {
		p3 := strings.Split(p2, ":")
		if len(p3) > 1 {
			if p3[0] == "forename" {
				if p3[1] == "asc" {
					orderFunc = append(orderFunc, forenameAsc)
				} else {
					orderFunc = append(orderFunc, forenameDesc)
				}
			} else if p3[0] == "lastname" {
				if p3[1] == "asc" {
					orderFunc = append(orderFunc, lastnameAsc)
				} else {
					orderFunc = append(orderFunc, lastnameDesc)
				}
			} else if p3[0] == "email" {
				if p3[1] == "asc" {
					orderFunc = append(orderFunc, emailAsc)
				} else {
					orderFunc = append(orderFunc, emailDesc)
				}
			} else if p3[0] == "id" {
				if p3[1] == "asc" {
					orderFunc = append(orderFunc, idAsc)
				} else {
					orderFunc = append(orderFunc, idDesc)
				}
			}
		} else {
			if p3[0] == "forename" {
				orderFunc = append(orderFunc, forenameAsc)
			} else if p3[0] == "lastname" {
				orderFunc = append(orderFunc, lastnameAsc)
			} else if p3[0] == "email" {
				orderFunc = append(orderFunc, emailAsc)
			} else {
				orderFunc = append(orderFunc, idAsc)
			}
		}

		OrderedBy(orderFunc...).Sort(arr)

	}
}

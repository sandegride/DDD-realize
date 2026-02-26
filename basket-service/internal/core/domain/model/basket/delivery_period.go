package basket

import (
	"basket-service/internal/pkg/errs"
	"strings"
)

var (
	Morning = newDeliveryPeriod(1, "morning", 6, 12)
	Midday  = newDeliveryPeriod(2, "midday", 12, 17)
	Evening = newDeliveryPeriod(3, "evening", 17, 24)
	Night   = newDeliveryPeriod(4, "night", 0, 6)

	DeliveryPeriods = []*DeliveryPeriod{Morning, Midday, Evening, Night}
)

type DeliveryPeriod struct {
	id   int
	name string
	from int
	to   int
}

func newDeliveryPeriod(id int, name string, from int, to int) *DeliveryPeriod {
	return &DeliveryPeriod{
		id:   id,
		name: name,
		from: from,
		to:   to,
	}
}

func RestoreDeliveryPeriod(id int, name string, from int, to int) *DeliveryPeriod {
	return &DeliveryPeriod{
		id:   id,
		name: name,
		from: from,
		to:   to,
	}
}

func (dp *DeliveryPeriod) Equal(other *DeliveryPeriod) bool {
	return dp.id == other.id
}

func (dp *DeliveryPeriod) ID() int {
	return dp.id
}

func (dp *DeliveryPeriod) Name() string {
	return dp.name
}

func (dp *DeliveryPeriod) From() int {
	return dp.from
}

func (dp *DeliveryPeriod) To() int {
	return dp.to
}

func GetDeliveryPeriodByName(name string) (*DeliveryPeriod, error) {
	for _, v := range DeliveryPeriods {
		if strings.ToLower(v.name) == strings.ToLower(name) {
			return v, nil
		}
	}
	return nil, errs.NewObjectNotFoundError("Name", nil)
}

func GetDeliveryPeriodByID(id int) (*DeliveryPeriod, error) {
	for _, v := range DeliveryPeriods {
		if v.id == id {
			return v, nil
		}
	}
	return nil, errs.NewObjectNotFoundError("ID", nil)
}

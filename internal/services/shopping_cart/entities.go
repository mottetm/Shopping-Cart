package shopping_cart

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
)

func getHashId() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = "github.com/mottetm/Shopping-Cart/Items#Id"
	hd.MinLength = 10
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	return h
}

var h *hashids.HashID = getHashId()

type ItemId int

func (id ItemId) MarshalJSON() ([]byte, error) {
	s, err := h.Encode([]int{int(id)})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return json.Marshal(s)
}

func (id *ItemId) UnmarshalJSON(d []byte) error {
	ids, err := h.DecodeWithError(string(d))
	if err != nil {
		return errors.WithStack(err)
	}
	*id = ItemId(ids[0])

	return nil
}

func (id ItemId) String() string {
	s, _ := h.Encode([]int{int(id)})
	return s
}

type ItemName string
type ItemQuantity int
type ReservationId string

type Item struct {
	Id            ItemId        `json:"id"`
	Name          ItemName      `json:"name"`
	Quantity      ItemQuantity  `json:"quantity,omitempty"`
	ReservationId ReservationId `json:"reservationId,omitempty"`
}

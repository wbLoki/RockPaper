package pkg

type Hub struct {
	Pools map[string]*Pool
}

func NewHub() *Hub {
	return &Hub{
		Pools: make(map[string]*Pool),
	}
}

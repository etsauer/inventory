package apis

type InventoryItem struct {
	Metadata ResourceMeta  `json:"metadata"`
	Spec     InventorySpec `json:"spec"`
}

type InventorySpec struct {
	FillLevel   int8   `json:"filllevel"`
	Maker       string `json:"maker"`
	Product     string `json:"product"`
	Description string `json:"description"`
}

type ResourceMeta struct {
	Name    string            `json:"name"`
	Barname string            `json:"barname"`
	Labels  map[string]string `json:"labels"`
}

var Items []InventoryItem

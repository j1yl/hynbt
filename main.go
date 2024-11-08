package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Tnze/go-mc/nbt"
)

type ItemData struct {
	Info     map[string]interface{} `json:"info"`
	Enchants map[string]interface{} `json:"enchants"`
}

func main() {
	data := "H4sIAAAAAAAAAFVVTW/iRhgeAkmAZJvcettOpU0VygYZjDHspcoSsqDmYxWy7EbVyhrbAx7hD2SP03DsH9hD1VO17TVSf0YOvfbYe+79C1Xf1wOERhH2vPO8837M8z4uE1IiOVEmhOQ2yIZwczs5stmN0lDmyiQv2aRECjx0PIJ/ObLzLrRjzqbM9nkuT0p94fJTn00S2P23TLZdkcx8NgensyjmRbB+RejDvfmGs5gOHbC9gqVbbxomPNuHhqZpFfItIE5YwCbZrlPVDR2e/LCqa5UMVq2/bNaNWqdeId8BdihjHk6kp9BNrVPTO+sOrR+qxkd4dg6rprE6wdTMmlkhr+CAbiwkXY9YbzUOlIOhHaCHCx6Ng4WvqdVrzdZBhRyB7+soTBN6LCVzpnQ449xdpGHUjIMlvqaZANcAPggl930xgSZmsVi12e6oAGvVdVpGrd2okBfgcc4mwqGnInQVvtGqtdsK2KxDpRXyHFCnPI4cIecKo+sLQKsGpzzPevBwbz9+/gRvHxfL9uOvv+ASbqUDCTzc++q364mAx4yOoAEvsQtdn7NbvraGdgmH+XQ0GKBrp5vaIgmeACcxm0Qh7QNteEwHS3MvdGE5BELAYzQgbbTdcSeVa4f37mY8FtgeOlgZT0XM6XEy446kg0VQsCWSwt2L6Zr7IJgxX4QThC1tZ1x6YJRzDKqhIYpkhlm5naVwe6OVx9CBgsMJJj8iOhoCgUk+nTmczzwocXXAtZdidX4Uu6u2XAvJQvq98P2s3pXviAUzEWcde7LxMAoiYNIIruMbHADO/CRjElAQZ8bN+EllRN8L6fE4qZE6jlLMQqmAdfx9/O2nBZXRiYWuIg1V14+Utdc5iCBoOHp2GXA4Cmw8zPT5LfdrkEsFBulKTDx55PgCmgTx04TTeZTG1PFZklBmC2zu1wDef7gfY70J0rBl6AbY9vCuUCBcms1KVt2KQp6A5D2kF6OOxzAjiOBC8ZlQpDNcQgX1bJhMfieBmaoVmJ2+JO6UPdxjFbg6v7nuD7r05N3Fm97lBR2+v7w6oU/AIilcsICTKpgWifXnSDq4T8j68fMf6//Y0t8/kTLZ62FsGPRY2MDYJE/2YgZlzK10NomZy1H0QATLiWQysewomhKy8WWe7HuRtGYRWCPLQSnFppRJYcKDpEiK3cvz18fXlka23/auTnvd6xL5Ig39yJly10r8SCaomxukPDx++7Y/uOoBdOVUJLvLVwvOI4XLi5sPxf9hl8cSiBVErhgLHpOtcVY4JMdiaUVj60cWr5KDIp4tarIyHoCtXMTPASn2b+C4weVFmezidwC4F3DgX54UxWLyAJzPk4IPEwWvW7DjLC4bltt5suVkcgKLzTzZ9tUowqoAW0k2VWprM8GZUz77qS9FwCS3HCVOClJKlnOq/CElmEFoGiqMCr4zBumwWCYdKrMyXymM8irdLudRuWxzpUkqxjM3EzPLy8RMdWd3jNpjJZn2KFjxdjHAarkjn/RgWXYmkmp7V6I0WNNMGtR+yV+qFGYBl1VIU2j4C7tlOqZu60es0XCOmq4+PmKOrh11bN5qdxqddsOxC6QE7eFAvGAGCf78199//oOU2VJCgN/q/wB/2Fwh2gcAAA\u003d\u003d"
	item := ParseNBTString(data)
	json, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		panic(err)
	}
	f, err := os.Create("example.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Writer.Write(f, json)
	if err != nil {
		panic(err)
	}

}

func ParseNBTString(raw string) ItemData {
	start := time.Now()
	// decode the NBT string
	z, _ := base64.StdEncoding.DecodeString(raw)
	gzreader, _ := gzip.NewReader(bytes.NewReader(z))
	decoder := nbt.NewDecoder(gzreader)
	var val map[string]interface{}
	decoder.Decode(&val)
	extras := val["i"].([]interface{})[0].(map[string]interface{})["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})

	// handle missing enchantments key
	var enchants map[string]interface{}
	if _, ok := extras["enchantments"]; ok {
		enchants = extras["enchantments"].(map[string]interface{})
		delete(extras, "enchantments")
	} else {
		enchants = map[string]interface{}{}
	}

	elapsed := time.Since(start)
	fmt.Printf("time taken: %v\n", elapsed.Microseconds())

	return ItemData{
		Enchants: enchants,
		Info:     extras,
	}
}

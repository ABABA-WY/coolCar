package main

import (
	trippb "coolcar/proto/gen/go"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func main() {
	trip := trippb.Trip{
		Start:        "aaa",
		End:          "bbb",
		DurantionSec: 12,
		FeeCent:      34,
	}
	fmt.Println(&trip)
	b, err := proto.Marshal(&trip) //b二进制流
	if err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", b)
	var trip2 trippb.Trip
	err = proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&trip2)
	b, err = json.Marshal(&trip2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

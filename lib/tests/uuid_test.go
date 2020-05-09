package tests

import (
	"fmt"
	"testing"

	"julive/tools/uuid"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.GenerateV1())
	fmt.Println(uuid.GenerateV4())
	fmt.Println(uuid.Xid())
	fmt.Println(uuid.Sonyflake())
}

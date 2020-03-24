package core

import (
	"fmt"
	"testing"
)

func TestRemove(t *testing.T) {
	slice := []string{"chensq", "zengxf", "yangf", "zhangq", "zhouxd", "chenghl"}
	fmt.Println(Remove(slice, "cdd"))
	fmt.Println(Remove(slice, "chensq"))
	fmt.Println(Remove(slice, "chenghl"))
	fmt.Println(Remove(slice, "yangf"))
}

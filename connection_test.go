package main

import (
	"fmt"
	"testing"
)

func TestRemove(t *testing.T) {
	slice := []string{"chensq", "zengxf", "yangf", "zhangq", "zhouxd", "chenghl"}
	fmt.Println(remove(slice, "cdd"))
	fmt.Println(remove(slice, "chensq"))
	fmt.Println(remove(slice, "chenghl"))
	fmt.Println(remove(slice, "yangf"))
}

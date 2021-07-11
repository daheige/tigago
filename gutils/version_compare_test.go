package gutils

import (
	"log"
	"testing"
)

func TestVersionCompare(t *testing.T) {
	v1 := "v1.2.3"
	v2 := "v1.2.5"
	log.Println("v1 <= v2 : ", VersionCompare(v1, v2, "<="))

	v1 = "1.2.3"
	v2 = "1.2.6"
	log.Println("v1 <= v2 : ", VersionCompare(v1, v2, "<="))

	v1 = "1.2.3.1"
	v2 = "1.2.4"
	log.Println("v1 <= v2 : ", VersionCompare(v1, v2, "<="))

	v1 = "1.2.1"
	v2 = "1.1.3"
	log.Println("v1 >= v2 : ", VersionCompare(v1, v2, ">="))
	log.Println("v1 >= v2 : ", VersionCompare(v1, v2, "ge"))

	v1 = "1.2.2"
	v2 = "1.2.2"
	log.Println("v1 == v2 : ", VersionCompare(v1, v2, "="))

	v1 = "1.2.2"
	v2 = "1.2.1"
	log.Println("v1 > v2 : ", VersionCompare(v1, v2, ">"))
	log.Println("v1 > v2 : ", VersionCompare(v1, v2, "gt"))

	v1 = "1.1.2.1"
	v2 = "1.1.2.2"
	log.Println("v1 < v2 : ", VersionCompare(v1, v2, "<"))
	log.Println("v1 lt v2 : ", VersionCompare(v1, v2, "lt"))

	v1 = "1.1.2"
	v2 = "1.1.1"

	log.Println("v1 != v2 : ", VersionCompare(v1, v2, "!="))
	log.Println("v1 ne v2 : ", VersionCompare(v1, v2, "ne"))
	log.Println("v1 neq v2 : ", VersionCompare(v1, v2, "ne"))
}

func TestFormVersionKeywords(t *testing.T) {
	log.Println(VersionCompare("1.2.3-alpha", "1.2.3RC7", ">="))

	log.Println(VersionCompare("1.2.3-beta", "1.2.3pl", "lt"))

	log.Println(VersionCompare("1.1_dev", "1.2any", "eq"))
}

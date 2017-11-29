// package main
package bchainlibs

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/onrik/gomerkle"
	"strconv"
	"bytes"
	"sort"
)

type Common struct {
	Gender int
	From   string
	To     string
}

type Foo struct {
	Id    string
	Name  string
	Extra Common
}

type Bar struct {
	Id    string
	Name  string
	Extra Common
}





type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }



func main() {

	foo := Foo{Id: "123", Name: "Joe", Extra: Common{Gender: -1, From: "Bob", To: "Alice"}}
	bar := foo
	bar2 := bar
	bar.Id = "2344"
	bar.Extra.Gender = 1
	bar2.Id = "9876"
	bar2.Extra.From = "Tupac"

	fmt.Println(foo)
	fmt.Println(bar)
	fmt.Println(bar2)
	fmt.Printf("%p\n", &foo)
	fmt.Printf("%p\n", &bar)
	fmt.Printf("%p\n", &bar2)


	// -------------------------------------------------------

	puzzle := "This is a puzzle!"

	t1 := time.Now().UnixNano()
	h := sha256.New()
	h.Write([]byte( puzzle ))
	checksum := h.Sum(nil)
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("%x\n", checksum)
	fmt.Printf("%s\n", sha256Hash)

	t2 := time.Now().UnixNano()

	sum := sha256.Sum256([]byte( puzzle ))
	sha256Hash2 := hex.EncodeToString(sum[:])
	fmt.Printf("%x\n", sum)
	fmt.Printf("%s\n", sha256Hash2)

	t3 := time.Now().UnixNano()

	fmt.Printf("=> %d\n", t2-t1)
	fmt.Printf("=> %d\n", t3-t2)

	// -------------------------------------------------------

	data := [][]byte{
		[]byte("Buzz"),
		[]byte("Lenny"),
		[]byte("Squeeze"),
		[]byte("Wheezy"),
		[]byte("Jessie"),
		[]byte("Stretch"),
		[]byte("Buster"),
	}
	tree := gomerkle.NewTree(sha256.New())
	tree.AddData(data...)

	err := tree.Generate()
	if err != nil {
		panic(err)
	}

	// Proof for Jessie
	proof := tree.GetProof(4)
	leaf := tree.GetLeaf(4)
	fmt.Printf(hex.EncodeToString(tree.Root()) + "\n")
	fmt.Printf(string(strconv.AppendBool([]byte("bool:"), tree.VerifyProof(proof, tree.Root(), leaf))))

	data2 := [][]byte{
		[]byte("Buzz"),
		[]byte("Lenny"),
		[]byte("Squeeze"),
		[]byte("Wheezy"),
		[]byte("Jessie"),
		[]byte("Stretch"),
		[]byte("Buster"),
	}
	tree2 := gomerkle.NewTree(sha256.New())
	tree2.AddData(data2...)

	err2 := tree2.Generate()
	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("\n")
	fmt.Printf(hex.EncodeToString(tree.Root()) + "\n")
	fmt.Printf(hex.EncodeToString(tree2.Root()) + "\n")

	if bytes.Equal(tree.Root(), tree2.Root()) {
		fmt.Printf( "Equal Trees\n")
	} else {
		fmt.Printf( "NOT Equal Trees\n")
	}

	// -------------------------------------------------------

	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	fmt.Println(people)
	sort.Sort(ByAge(people))
	fmt.Println(people)
}
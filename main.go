package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
)

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func sigma0(x uint32) uint32 {
	var rotr7 = bits.RotateLeft32(x, -7)
	var rotr18 = bits.RotateLeft32(x, -18)
	var shr = x >> 3
	return rotr7 ^ rotr18 ^ shr
}

func sigma1(x uint32) uint32 {
	var rotr17 = bits.RotateLeft32(x, -17)
	var rotr19 = bits.RotateLeft32(x, -19)
	var shr = x >> 10
	return rotr17 ^ rotr19 ^ shr
}

func usigma0(x uint32) uint32 {
	var rotr2 = bits.RotateLeft32(x, -2)
	var rotr13 = bits.RotateLeft32(x, -13)
	var rotr22 = bits.RotateLeft32(x, -22)
	return rotr2 ^ rotr13 ^ rotr22
}

func usigma1(x uint32) uint32 {
	var rotr6 = bits.RotateLeft32(x, -6)
	var rotr11 = bits.RotateLeft32(x, -11)
	var rotr25 = bits.RotateLeft32(x, -25)
	return rotr6 ^ rotr11 ^ rotr25
}

func ch(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func maj(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

// Returns initializing constants
func constants() []uint32 {
	var primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311}

	r := make([]uint32, 64)
	for i := 0; i < 64; i++ {
		var root = math.Pow(float64(primes[i]), 1.0/3.0)
		var fractional = root - math.Floor(root)

		hex := ""
		for j := 0; j < 8; j++ {
			product := fractional * 16
			carry := int64(math.Floor(product))
			fractional = product - math.Floor(product)
			hex += strconv.FormatInt(carry, 16)
		}
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			panic(err)
		}
		r[i] = uint32(val)
	}
	return r
}

func sha256(str string) string {
	// constants
	consts := constants()

	// padding
	// clear()
	input := []byte(str)

	padd := make([]byte, 64)
	copy(padd, input)
	length := len(input)
	size := uint32(length * 8)
	padd[63] = byte(size)
	padd[length] = 128

	// fmt.Println("Message block:")
	// fmt.Printf("%08b\n", padd)
	// fmt.Println("--------------")

	// message schedule
	schedule := make([]uint32, 64)
	// fmt.Println("Schedule:")

	for i := 0; i < 16; i++ {
		schedule[i] = binary.BigEndian.Uint32(padd[i*4 : i*4+4])
	}

	for i := 16; i < 64; i++ {
		schedule[i] = sigma1(schedule[i-2]) + schedule[i-7] + sigma0(schedule[i-15]) + schedule[i-16]
	}
	// for i := 0; i < 64; i++ {
	// 	if i < 10 {
	// 		fmt.Printf("W%d  %032b\n", i, schedule[i])
	// 	} else {
	// 		fmt.Printf("W%d %032b\n", i, schedule[i])
	// 	}
	// }

	// init a - h
	a := uint32(math.Trunc((math.Sqrt(2) - math.Floor(math.Sqrt(2))) * math.Pow(2, 32)))
	b := uint32(math.Trunc((math.Sqrt(3) - math.Floor(math.Sqrt(3))) * math.Pow(2, 32)))
	c := uint32(math.Trunc((math.Sqrt(5) - math.Floor(math.Sqrt(5))) * math.Pow(2, 32)))
	d := uint32(math.Trunc((math.Sqrt(7) - math.Floor(math.Sqrt(7))) * math.Pow(2, 32)))
	e := uint32(math.Trunc((math.Sqrt(11) - math.Floor(math.Sqrt(11))) * math.Pow(2, 32)))
	f := uint32(math.Trunc((math.Sqrt(13) - math.Floor(math.Sqrt(13))) * math.Pow(2, 32)))
	g := uint32(math.Trunc((math.Sqrt(17) - math.Floor(math.Sqrt(17))) * math.Pow(2, 32)))
	h := uint32(math.Trunc((math.Sqrt(19) - math.Floor(math.Sqrt(19))) * math.Pow(2, 32)))

	a0, b0, c0, d0, e0, f0, g0, h0 := a, b, c, d, e, f, g, h

	// fmt.Printf("a = %d\n", a)
	// fmt.Printf("b = %d\n", b)
	// fmt.Printf("c = %d\n", c)
	// fmt.Printf("d = %d\n", d)
	// fmt.Printf("e = %d\n", e)
	// fmt.Printf("f = %d\n", f)
	// fmt.Printf("g = %d\n", g)
	// fmt.Printf("h = %d\n", h)

	// compression (H0)

	var t1 uint32
	var t2 uint32
	for i := 0; i < 64; i++ {
		// fmt.Println("--------------")

		t1 = usigma1(e) + ch(e, f, g) + h + consts[i] + schedule[i]
		t2 = usigma0(a) + maj(a, b, c)

		b, c, d, e, f, g, h = a, b, c, d, e, f, g
		a = t1 + t2
		e = e + t1

		// fmt.Printf("a = %032b\n", a)
		// fmt.Printf("b = %032b\n", b)
		// fmt.Printf("c = %032b\n", c)
		// fmt.Printf("d = %032b\n", d)
		// fmt.Printf("e = %032b\n", e)
		// fmt.Printf("f = %032b\n", f)
		// fmt.Printf("g = %032b\n", g)
		// fmt.Printf("h = %032b\n", h)

		// reader := bufio.NewReader(os.Stdin)
		// reader.ReadString('\n')
	}

	a = a0 + a
	b = b0 + b
	c = c0 + c
	d = d0 + d
	e = e0 + e
	f = f0 + f
	g = g0 + g
	h = h0 + h

	// fmt.Println("++++++++++++++")
	// fmt.Printf("a = %032b\n", a)
	// fmt.Printf("b = %032b\n", b)
	// fmt.Printf("c = %032b\n", c)
	// fmt.Printf("d = %032b\n", d)
	// fmt.Printf("e = %032b\n", e)
	// fmt.Printf("f = %032b\n", f)
	// fmt.Printf("g = %032b\n", g)
	// fmt.Printf("h = %032b\n", h)

	return strconv.FormatUint(uint64(a), 16) + strconv.FormatUint(uint64(b), 16) + strconv.FormatUint(uint64(c), 16) + strconv.FormatUint(uint64(d), 16) + strconv.FormatUint(uint64(e), 16) + strconv.FormatUint(uint64(f), 16) + strconv.FormatUint(uint64(g), 16) + strconv.FormatUint(uint64(h), 16)
}

func main() {
	// calculates sha256 of a small string(up to 32 characters)
	fmt.Println("sha256 of \"hello world\" is:")
	fmt.Println(sha256("hello world"))
}

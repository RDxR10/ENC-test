package main

import (
        "fmt"
        "math/rand"
        "time"
)


type Z struct {
        s []int
        t   int
}

func NewZ(s []int, t int) *Z {
        return &Z{
                s: s,
                t: t,
        }
}

func (z *Z) Gen() int {
        tB := z.s[z.t]
        x := z.s[0] ^ tB
        for i := 0; i < len(z.s)-1; i++ {
                z.s[i] = z.s[i+1]
        }
        z.s[len(z.s)-1] = x

        return x
}



func (z *Z) Enc(t_ []byte) []byte {
        rand.Seed(time.Now().UnixNano())
        k := make([]byte, len(t_))
        for i := 0; i < len(k); i++ {
                k[i] = byte(rand.Intn(256))
        }
        enc := make([]byte, len(t_))
        for i := 0; i < len(t_); i++ {
                enc[i] = t_[i] ^ k[i] ^ byte(z.Gen())
        }

        return enc
}

func main() {
        s := []int{1, 0, 1, 1, 1, 0, 1, 0}
        t := 4
        z := NewZ(s, t)
        t_ := []byte("ENC-test")
        enc_ := z.Enc(t_)
        fmt.Printf("%v\n", enc_)
}

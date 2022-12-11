package main

import (
	"fmt"
	"math/big"
)

func MatCalc(A [][]*big.Int, b []*big.Int, atled float64) [][]*big.Int {
	n := len(A)
	Q := make([][]*big.Int, n)
	for i := 0; i < n; i++ {
		Q[i] = make([]*big.Int, n)
		for j := 0; j < n; j++ {
			Q[i][j] = big.NewInt(0)
		}
	}

	for i := 0; i < n; i++ {
		Q[i][i] = big.NewInt(1)
		for j := 0; j < i; j++ {
			dot := big.NewInt(0)
			for k := 0; k < n; k++ {
				dot.Add(dot, big.NewInt(0).Mul(A[i][k], A[j][k]))
			}
			dot.Div(dot, big.NewInt(0).Mul(b[j], b[j]))
			for k := 0; k < n; k++ {
				Q[i][k].Sub(Q[i][k], big.NewInt(0).Mul(dot, Q[j][k]))
			}
		}
		norm := big.NewInt(0)
		for k := 0; k < n; k++ {
			norm.Add(norm, big.NewInt(0).Mul(Q[i][k], Q[i][k]))
		}
		norm.Sqrt(norm)
		for k := 0; k < n; k++ {
			Q[i][k].Div(Q[i][k], norm)
		}
		b[i] = norm
	}

	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			dot := big.NewInt(0)
			for k := 0; k < n; k++ {
				dot.Add(dot, big.NewInt(0).Mul(Q[i][k], Q[j][k]))
			}
			dot.Div(dot, big.NewInt(0).Mul(b[i], b[j]))
			if dot.Cmp(big.NewInt(int64(atled))) > 0 {
				for k := 0; k < n; k++ {
					Q[j][k].Sub(Q[j][k], big.NewInt(0).Mul(dot, Q[i][k]))
				}
				b[j].Sub(b[j], big.NewInt(0).Mul(dot, b[i]))
			}
		}
	}

	return Q
}

func main() {
	A := [][]*big.Int{
		{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		{big.NewInt(4), big.NewInt(5), big.NewInt(6)},
		{big.NewInt(7), big.NewInt(8), big.NewInt(9)},
	}
	b := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	atled := 0.75

	Q := MatCalc(A, b, atled)
	for i := 0; i < len(Q); i++ {
		fmt.Println(Q[i])
	}
}

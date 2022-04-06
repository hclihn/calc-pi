package main

import (
	"fmt"
  "math/big"
  "runtime"
  "math"
  "time"
)

// Based on Bailey–Borwein–Plouffe formula (BBP formula)
// Source: https://gist.github.com/linyd2005/81fd00a772ed002cb538 

func GenWorker(p uint) func(id int, result chan *big.Float) {
	B1 := new(big.Float).SetPrec(p).SetInt64(1)
	B2 := new(big.Float).SetPrec(p).SetInt64(2)
	B4 := new(big.Float).SetPrec(p).SetInt64(4)
	B5 := new(big.Float).SetPrec(p).SetInt64(5)
	B6 := new(big.Float).SetPrec(p).SetInt64(6)
	B8 := new(big.Float).SetPrec(p).SetInt64(8)
	B16 := new(big.Float).SetPrec(p).SetInt64(16)

	return func(id int, result chan *big.Float) {
		Bn := new(big.Float).SetPrec(p).SetInt64(int64(id))

		C1 := new(big.Float).SetPrec(p).SetInt64(1)
		for i := 0; i < id; i++ {
			C1.Mul(C1, B16)
		}

		C2 := new(big.Float).SetPrec(p).Mul(B8, Bn)

		T1 := new(big.Float).SetPrec(p).Add(C2, B1)
		T1.Quo(B4, T1)

		T2 := new(big.Float).SetPrec(p).Add(C2, B4)
		T2.Quo(B2, T2)

		T3 := new(big.Float).SetPrec(p).Add(C2, B5)
		T3.Quo(B1, T3)

		T4 := new(big.Float).SetPrec(p).Add(C2, B6)
		T4.Quo(B1, T4)

		R := new(big.Float).SetPrec(p).Sub(T1, T2)
		R.Sub(R, T3).Sub(R, T4).Quo(R, C1)

		result <- R
	}
}

func GenWorkers(p uint, n int) []func(result chan *big.Float) {
	B1 := new(big.Float).SetPrec(p).SetInt64(1)
	B2 := new(big.Float).SetPrec(p).SetInt64(2)
	B4 := new(big.Float).SetPrec(p).SetInt64(4)
	B5 := new(big.Float).SetPrec(p).SetInt64(5)
	B6 := new(big.Float).SetPrec(p).SetInt64(6)
	B8 := new(big.Float).SetPrec(p).SetInt64(8)
	B16 := new(big.Float).SetPrec(p).SetInt64(16)

  workers := make([]func(result chan *big.Float), n)
  C1 := new(big.Float).SetPrec(p).SetInt64(1)
  for i := range workers {
    C1n := new(big.Float).Set(C1)
    Bn := new(big.Float).SetPrec(p).SetInt64(int64(i))
	  workers[i] = func(result chan *big.Float) {
  		C2 := new(big.Float).SetPrec(p).Mul(B8, Bn)
  
  		T1 := new(big.Float).SetPrec(p).Add(C2, B1)
  		T1.Quo(B4, T1)
  
  		T2 := new(big.Float).SetPrec(p).Add(C2, B4)
  		T2.Quo(B2, T2)
  
  		T3 := new(big.Float).SetPrec(p).Add(C2, B5)
  		T3.Quo(B1, T3)
  
  		T4 := new(big.Float).SetPrec(p).Add(C2, B6)
  		T4.Quo(B1, T4)
  
  		R := new(big.Float).SetPrec(p).Sub(T1, T2)
  		R.Sub(R, T3).Sub(R, T4)
      if i != 0 { 
        R.Quo(R, C1n)
      }
  
  		result <- R
  	} 
    C1.Mul(C1, B16)
  }
  return workers
}

func GenWorkerR(p uint) func(id int, result chan *big.Rat) {
	B1 := new(big.Rat).SetInt64(1)
	B2 := new(big.Rat).SetInt64(2)
	B4 := new(big.Rat).SetInt64(4)
	B5 := new(big.Rat).SetInt64(5)
	B6 := new(big.Rat).SetInt64(6)
	B8 := new(big.Rat).SetInt64(8)
	B16 := new(big.Rat).SetInt64(16)

	return func(id int, result chan *big.Rat) {
		Bn := new(big.Rat).SetInt64(int64(id))

		C1 := new(big.Rat).SetInt64(1)
		for i := 0; i < id; i++ {
			C1.Mul(C1, B16)
		}

		C2 := new(big.Rat).Mul(B8, Bn)

		T1 := new(big.Rat).Add(C2, B1)
		T1.Quo(B4, T1)

		T2 := new(big.Rat).Add(C2, B4)
		T2.Quo(B2, T2)

		T3 := new(big.Rat).Add(C2, B5)
		T3.Quo(B1, T3)

		T4 := new(big.Rat).Add(C2, B6)
		T4.Quo(B1, T4)

		R := new(big.Rat).Sub(T1, T2)
		R.Sub(R, T3).Sub(R, T4).Quo(R, C1)

		result <- R
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

  const nBytes = 64
  const nb = 8 * nBytes

  p := uint(nb + 32) // big float precision
  n := nb / 4 + 1 // hex digits
  nd := int(float64(nb) / math.Log2(10) + 1.0) // decimal digits

  start := time.Now()
	result := make(chan *big.Float, n)
	worker := GenWorker(p)

	pi := new(big.Float).SetPrec(p).SetInt64(0)

	for i := 0; i < n; i++ {
		go worker(i, result)
	}

	for i := 0; i < n; i++ {
		pi.Add(pi, <-result)
	}

	dur := time.Since(start)
	fmt.Printf("take %v to calculate %d (%d hex) digits (%d/%d bits) long pi \n", dur, nd, n, nb, p)
	fmt.Printf("%[1]*.[2]*[3]f\n", 1, nd, pi)
  
  xx := new(big.Float).SetMantExp(pi, nb-pi.MantExp(nil))
  x, _ := xx.Int(nil)
  
  fmt.Printf("pi exp: %d; xx exp : %d\n", pi.MantExp(nil), xx.MantExp(nil))
  fmt.Printf("nbits: %d\n", x.BitLen())
  
  buf := make([]byte, nBytes)
  x.FillBytes(buf)
  fmt.Printf("Buf: % 02x \n", buf) 

  // use workers
  start = time.Now()
	result = make(chan *big.Float, n)
	workers := GenWorkers(p, n)

	pi = new(big.Float).SetPrec(p).SetInt64(0)

	for i := 0; i < n; i++ {
		go workers[i](result)
	}

	for i := 0; i < n; i++ {
		pi.Add(pi, <-result)
	}

	dur = time.Since(start)
	fmt.Printf("Workers: take %v to calculate %d (%d hex) digits (%d/%d bits) long pi \n", dur, nd, n, nb, p)
	fmt.Printf("%[1]*.[2]*[3]f\n", 1, nd, pi)
  
  xx = new(big.Float).SetMantExp(pi, nb-pi.MantExp(nil))
  x, _ = xx.Int(nil)
  
  fmt.Printf("pi exp: %d; xx exp : %d\n", pi.MantExp(nil), xx.MantExp(nil))
  fmt.Printf("nbits: %d\n", x.BitLen())
  
  buf = make([]byte, nBytes)
  x.FillBytes(buf)
  fmt.Printf("Buf: % 02x \n", buf) 
  
  // use Rat
  start = time.Now()
  resultR := make(chan *big.Rat, n)
	workerR := GenWorkerR(p)

	piR := new(big.Rat).SetInt64(0)

	for i := 0; i < n; i++ {
		go workerR(i, resultR)
	}

	for i := 0; i < n; i++ {
		piR.Add(piR, <-resultR)
	}

	dur = time.Since(start)
	fmt.Printf("take %v to calculate in big.Rat %d (%d hex) digits (%d/%d bits) long pi \n", dur, nd, n, nb, p)
	fmt.Printf("%s\n", piR.FloatString(nd))

  nbPi := new(big.Int).Quo(piR.Num(), piR.Denom()).BitLen()
  xxR := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(nb-int64(nbPi)), nil))
  xxR.Mul(xxR, piR)
  xR := new(big.Int).Quo(xxR.Num(), xxR.Denom())
  
  fmt.Printf("nbits: %d\n", xR.BitLen())
  
  xR.FillBytes(buf)
  fmt.Printf("Buf: % 02x \n", buf) 
}

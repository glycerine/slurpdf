package slurpdf

import (
	"fmt"
	//"math"
	//"os"
	"reflect"
	"strings"
	"testing"
	"time"
	//cv "github.com/glycerine/goconvey/convey"
)

// fewer dependencies
type smallConvey struct{}

var cv = &smallConvey{}

func (s *smallConvey) AssertTrue(expr bool) {
	if !expr {
		panic(fmt.Sprintf("assertion false at %v", fileLine(2)))
	}
}

func (s *smallConvey) Convey(desc string, t *testing.T, f func()) {
	//fmt.Printf("%v\n", desc)
	f()
}
func (s *smallConvey) ShouldResemble(a, b any) {
	if !reflect.DeepEqual(a, b) {
		panic(fmt.Sprintf("ShouldResemble false at %v", fileLine(2)))
	}
}

func Test001_slurp_in_data(t *testing.T) {

	cv.Convey("read a .csv dataframe from disk into memory using all cores", t, func() {
		//fmt.Printf("read a .csv dataframe from disk into memory using all cores\n")

		fn := "data/test001.csv"
		d := NewSlurpDataFrameNoStrings()
		t0 := time.Now()
		err := d.Slurp(fn)
		panicOn(err)
		vv("slurped in fn '%v' in %v", fn, time.Since(t0))

		nr := d.Nrow // Number of cases in full database
		nc := d.Ncol // Number of columns (X variables and Y target)
		vv("we see nc = %v, nr= %v", nc, nr)

		// illustrate how to use the testing framework
		if nc != 4 {
			panic("expected nc == 4")
		}
		if nr != 5 {
			panic("expected nr == 5")
		}

		// expected
		eh := []string{"x1", "x2", "x3", "y"}
		em := [][]float64{
			[]float64{1.49152459063627, 2.49152289572389, 1.67045378357525, 1},
			[]float64{0.391160302666239, 1.39115885339298, 3.55649293629879, 1},
			[]float64{0.434211270774665, 1.43420995498692, 3.30302332417163, 0},
			[]float64{0.136364617767486, 1.13636348140514, 8.33327222273116, 0},
			[]float64{1.136364617767486, -3.13636348140514, -9.273116, 0},
		}
		_ = em
		cv.ShouldResemble(d.Header, strings.Join(eh, ","))
		cv.ShouldResemble(d.Colnames, eh)
		for i := range em {
			cv.ShouldResemble(d.MatFullRow(i), em[i])
		}

		// ExtractCols
		xi0 := 3
		xi1 := 5
		wcol := []int{1, 3}
		n, nvar, xx, cn := d.ExtractCols(xi0, xi1, wcol)
		cv.AssertTrue(n == 2)
		cv.AssertTrue(nvar == 2)
		cv.ShouldResemble(cn, []string{"x2", "y"})
		cv.ShouldResemble(xx[0], em[3][1])
		cv.ShouldResemble(xx[1], em[3][3])
		cv.ShouldResemble(xx[2], em[4][1])
		cv.ShouldResemble(xx[3], em[4][3])

		// ExtractXXYY
		xi0 = 4
		xi1 = 5  // just the last row
		xj0 := 2 // just the 3rd col
		xj1 := 3
		yj := 3 // and the target
		n, nvar, xx, yy, colnames, targetname := d.ExtractXXYY(xi0, xi1, xj0, xj1, yj)
		cv.AssertTrue(n == 1)
		cv.AssertTrue(nvar == 1)
		cv.ShouldResemble(colnames, []string{"x3"})
		cv.ShouldResemble(targetname, "y")
		cv.ShouldResemble(xx, em[xi0][xj0:xj1])
		cv.ShouldResemble(yy, em[xi0][yj:(yj+1)])

	})

}

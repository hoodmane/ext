package main

import (
    "testing"
    "fmt"
)


func TestBasepExpansion(t *testing.T) {

	tables := []struct {
		n int
		p int
		output []int
	}{
		{8,  3, []int {2, 2}},
		{33, 5, []int {3, 1, 1}},
	}
    for _, table := range tables {
        output := basepExpansion(table.n, table.p, 0)
        if !eqListsQ(output, table.output) {
            t.Errorf("Ran basepExpansion(%v,%v) expected %v got %v", table.n, table.p, table.output, output)
        }
    }
}

func TestDirectBinomial(t *testing.T) {
	tables := []struct {
		n int
		k int
        p int
		output int
	}{
        {21, 2, 23, 210},
        {13, 9, 23, 715},
        {12, 8, 23, 495},
        {13, 8, 23, 1287},
        {14, 8, 23, 3003},
        {14, 9, 23, 2002},
        {15, 5, 23, 3003},
        {15, 8, 23, 6435},  
        {15, 9, 23, 5005},
        {16, 9, 23, 11440},
    }
    for _, table := range tables {
        output := directBinomial(table.n, table.k, table.p)
        if output != table.output  % table.p {
            t.Errorf("Ran directBinomial(%v,%v) expected %v, got %v", table.n, table.k, table.output % table.p, output)
        }
    }    
}

func TestMultinomial2(t *testing.T) {
	tables := []struct {
		l []int
		output int
	}{
        {[]int {1, 2}, 1},
        {[]int {1, 3}, 0},
        {[]int {1, 4}, 1},
        {[]int {2, 4}, 1},
        {[]int {1, 5}, 0},
        {[]int {2, 5}, 1},
        {[]int {2, 6}, 0},
        {[]int {2, 4, 8}, 1},
    }
    for _, table := range tables {
        output := Multinomial2(table.l)
        if output != table.output {
            t.Errorf("Ran Multinomial2(%v) expected %v, got %v", table.l, table.output, output)
        }
    }        
}
    
func TestBinomial2(t *testing.T) {
    tables := []struct {
        n int
        k int
        output int
    }{
        {4, 2, 0},
        {72, 46, 0},
        {82, 66, 1},
        {165, 132, 1},
        {169, 140, 0},
    }
    for _, table := range tables {
        output := Binomial2(table.n, table.k)
        if output != table.output {
            t.Errorf("Ran Binomial2(%v,%v) expected %v, got %v", table.n, table.k, table.output, output)
        }
    }        
}


func TestMultinomialOdd(t *testing.T) {
	tables := []struct {
		l []int
        p int
		output int
	}{
        {[]int {1090, 730}, 3, 1},
        {[]int {108054, 758}, 23, 18},
        {[]int {3, 2}, 7, 3},
    }
    for _, table := range tables {
        output := MultinomialOdd(table.l, table.p)
        if output != table.output {
            t.Errorf("Ran MultinomialOdd(%v, %v) expected %v, got %v", table.l, table.p, table.output, output)
        }
    }        
}
//
func TestBinomialOdd(t *testing.T) {
 
}

func TestXiDegrees(t *testing.T) {
	tables := []struct {
		n int
        p int
		output []int
	}{
        {17,   2, []int{1, 3, 7, 15}},
        {17,   3, []int{1, 4, 13}},
        {400, 17, []int{1, 18, 307}},
    }
    
    for _, table := range tables {
        output := XiDegrees(table.n, table.p)
        if !eqListsQ(output, table.output) {
            t.Errorf("Ran XiDegrees(%v, %v) expected %v, got %v", table.n, table.p, table.output, output)
        }
    }   
    
}

func TestWeightedIntegerVectors(t *testing.T) {
    tables := []struct {
        n int
        l []int
        output [][]int
    }{
        {10, []int{4, 1}, [][]int{[]int{0, 10}, []int{1, 6}, []int{2, 2}}},
        {7, []int {7, 3, 1}, [][]int{[]int{0, 0, 7}, []int{0, 1, 4}, []int{0, 2, 1}, []int{1, 0, 0}}},
        {20, []int{13, 4, 1}, [][]int{[]int{0, 0, 20}, []int{0, 1, 16}, []int{0, 2, 12}, []int{0, 3, 8}, []int{1, 0, 7}, []int{0, 4, 4}, []int{1, 1, 3}, []int{0, 5, 0}}},
    } 
     
    for _, table := range tables {
        call_str := fmt.Sprintf("WeightedIntegerVectors(%v, %v, %v)", table.n, table.l, []int{100,100,100,100})
        gen := WeightedIntegerVectors(table.n, table.l, []int{100,100,100,100})
        error_str := checkGeneratorOfListsOutput(call_str, table.output, gen)
        if error_str != "" {
            t.Error(error_str)
        }
    }    
}


func TestRestrictedPartitions(t *testing.T) {
    tables := []struct {
        n int
        l []int
        output [][]int
    }{
        {8,  []int{1},      [][]int{}},
        {10, []int{2,4,6},  [][]int{[]int{0, 1, 1}}},
        {10, []int{2,2,4,6},[][]int{[]int{0, 0, 1, 1}, []int{1, 1, 0, 1}}},
    }
    for _, table := range tables {
        call_str := fmt.Sprintf("WeightedIntegerVectors(%v, %v, %v)", table.n, table.l, []int{1,1,1,1})
        gen := WeightedIntegerVectors(table.n, table.l, []int{1,1,1,1})
        error_str := checkGeneratorOfListsOutput(call_str, table.output, gen)
        if error_str != "" {
            t.Error(error_str)
        }
    }
}

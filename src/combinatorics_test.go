package main

import (
    "testing"
    "fmt"
)

func testEqLists(a, b []int) bool {

    // If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

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
        if !testEqLists(output, table.output) {
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
        reverse bool
		output []int
	}{
        {17,   2, true,  []int{15, 7, 3, 1}},
        {17,   2, false, []int{1, 3, 7, 15}},
        {17,   3, true,  []int{13, 4, 1}},
        {400, 17, true,  []int{307, 18, 1}},
    }
    
    for _, table := range tables {
        output := XiDegrees(table.n, table.p, table.reverse)
        if !testEqLists(output, table.output) {
            t.Errorf("Ran XiDegrees(%v, %v, %v) expected %v, got %v", table.n, table.p, table.reverse, table.output, output)
        }
    }   
    
}

func CheckGeneratorOfListsOutput(t *testing.T, gen <- chan []int, gen_name string, desired_output [][]int){
    output_strings := make(map[string]int)
    for _, tuple := range desired_output {
        output_strings[fmt.Sprint(tuple)] ++
    }
    for tuple := range gen {
        str_tuple := fmt.Sprint(tuple)
        if output_strings[str_tuple] == 0 {
            t.Errorf("Unexpected result %v appeared in generator %v. Was expecting the multiset %v.", tuple, gen_name, desired_output)
        }
        output_strings[str_tuple] --
    }    
    for k, v := range output_strings {
        if v > 0 {
            t.Errorf("Expected result %v in generator %v but it didn't appear. Was expecting the multiset %v.", k, gen_name, desired_output)
        }
    }
}

func TestWeightedIntegerVectors(t *testing.T) {
    tables := []struct {
        n int
        l []int
        output [][]int
    }{
        {10, []int{1, 4}, [][]int{[]int{10, 0}, []int{6, 1}, []int{2, 2}}},
        {7, []int {1, 3, 7}, [][]int{[]int{7, 0, 0}, []int{4, 1, 0}, []int{1, 2, 0}, []int{0, 0, 1}}},
        {20, []int{1, 4, 13}, [][]int{[]int{20, 0, 0}, []int{16, 1, 0}, []int{12, 2, 0}, []int{8, 3, 0}, []int{7, 0, 1}, []int{4, 4, 0}, []int{3, 1, 1}, []int{0, 5, 0}}},
    } 
     
    for _, table := range tables {
        gen := WeightedIntegerVectors(table.n, table.l)
        gen_name := fmt.Sprintf("WeightedIntegerVectors(%v, %v)", table.n, table.l)
        CheckGeneratorOfListsOutput(t, gen, gen_name, table.output)
    }    
}


func TestRestrictedPartitions(t *testing.T) {
    tables := []struct {
        n int
        l []int
        output [][]int
    }{
        {8,  []int{1},      [][]int{}},
        {10, []int{6,4,2},  [][]int{[]int{6,4}}},
        {10, []int{6,4,2,2},[][]int{[]int{6,4}, []int{6, 2, 2}}},
    }
    for _, table := range tables {
        gen := RestrictedPartitions(table.n, table.l)
        gen_name := fmt.Sprintf("RestrictedPartitions(%v, %v)", table.n, table.l)
        CheckGeneratorOfListsOutput(t, gen, gen_name, table.output)
    }    

}

package main

import (
    "testing"
    "fmt"
)


func TestAdem2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		i int
        j int
		output strToInt
	}{
        {1,1, strToInt{}},
        {1,2, strToInt{"{[] [3]}" : 1}},
        {2,2, strToInt{"{[] [3 1]}" : 1}},
        {4,2, strToInt{"{[] [4 2]}" : 1}}, // Admissible
        {4,4, strToInt{"{[] [6 2]}" : 1, "{[] [7 1]}" : 1}},
        {5,7, strToInt{}},
        {6,7, strToInt{"{[] [13]}" : 1, "{[] [12 1]}" : 1, "{[] [10 3]}" : 1}},
    }
    
    length_tables := []struct {
        i int
        j int
        output_length int
    }{
        {100, 100, 16},
        {200, 200, 23},
    }    
    
    for _, table := range tables {
        call_str := fmt.Sprintf("Adem2(%v, %v)", table.i, table.j)
        output := Adem2(table.i, table.j).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("Adem2(%v, %v)", table.i, table.j)
        output_length := len(Adem2(table.i, table.j).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}


func TestAdemGeneric(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
        p int
		i int
        epsilon int
        j int
		output strToInt
	}{
        {3, 1,0,1, strToInt{"{[0 0] [2]}" : 2}},
        {3, 1,1,1, strToInt{"{[1 0] [2]}" : 1, "{[0 1] [2]}" : 1}},
        {3, 1,0,2, strToInt{}},
        {3, 1,1,2, strToInt{"{[0 1] [3]}" : 1, "{[1 0] [3]}" : 2}},
        {3, 2,0,1, strToInt{}},
        {3, 3,0,1, strToInt{"{[0 0 0] [3 1]}" : 1}},
        {3, 5,0,7, strToInt{"{[0 0 0] [11 1]}" : 1}},
        {3, 25,1,20, strToInt{"{[1 0 0] [38 7]}" : 1, "{[0 1 0] [38 7]}" : 1, "{[0 1 0] [37 8]}" : 1}},
        
        {5, 1, 1, 2, strToInt{"{[0 1] [3]}" : 1, "{[1 0] [3]}" : 2}},
        
        {23, 1,0,1, strToInt{"{[0 0] [2]}" : 2}},
        {23, 1,1,1, strToInt{"{[1 0] [2]}" : 1, "{[0 1] [2]}" : 1}},   
        {23, 2,0,1, strToInt{"{[0 0] [3]}" : 3}},
        {23, 5,0,7, strToInt{"{[0 0] [12]}" : 10}},
    }

    length_tables := []struct {
        p int
        i int
        epsilon int
        j int
        output_length int
    }{
        {3,  1000, 1, 1000, 107},
        {23, 2000, 1, 5000, 9},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("AdemGeneric(%v, %v, %v, %v)", table.p, table.i, table.epsilon, table.j)
        output := AdemGeneric(table.p, table.i, table.epsilon, table.j).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }    

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemGeneric(%v, %v, %v, %v)", table.p, table.i, table.epsilon, table.j)
        output_length := len(AdemGeneric(table.p, table.i, table.epsilon, table.j).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }    
}



func TestMakeMonoAdmissible2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		l []int
		output strToInt
	}{
        {[]int{}, strToInt{"{[] []}" : 1}},
        {[]int{12}, strToInt{"{[] [12]}" : 1}},
        {[]int{2, 1}, strToInt{"{[] [2 1]}" : 1}},
        {[]int{2, 2}, strToInt{"{[] [3 1]}" : 1}},
        {[]int{2, 2, 2}, strToInt{"{[] [5 1]}" : 1}},
    }

    length_tables := []struct {
        l []int
        output_length int
    }{
        
    }

    
    for _, table := range tables {
        output := MakeMonoAdmissible2(table.l).GetCoeffMap()
        call_str := fmt.Sprintf("MakeMonoAdmissible2(%v)", table.l)
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("MakeMonoAdmissible2(%v)", table.l)
        output_length := len(MakeMonoAdmissible2(table.l).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}    

func TestMakeMonoAdmissibleGeneric(t *testing.T){
    type strToInt map[string]int
    Mono := func(od, ev []int) Monomial{
        return Monomial{od, ev}
    }
    tables := []struct {
        p int
		l Monomial
		output strToInt
	}{
        {3, Mono([]int {0}, []int {}), strToInt{"{[0] []}" : 1 }},
        {7, Mono([]int {0, 0, 0}, []int {2, 1}), strToInt{"{[0 0] [3]}" : 3 }},
    }

    length_tables := []struct {
        p int
        l Monomial
        output_length int
    }{
        
    }

    for _, table := range tables {
        call_str := fmt.Sprintf("MakeMonoAdmissibleGeneric(%v, %v)", table.p, table.l)
        output := MakeMonoAdmissibleGeneric(table.p, table.l).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("MakeMonoAdmissibleGeneric(%v, %v)", table.p, table.l)
        output_length := len(MakeMonoAdmissibleGeneric(table.p, table.l).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}    



func TestProduct2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		l []int
        k []int
		output strToInt
	}{

    }

    length_tables := []struct {
        l []int
        k []int
        output_length int
    }{
        
    }

    
    for _, table := range tables {
        call_str := fmt.Sprintf("Product2(%v, %v)", table.l, table.k)
        output := Product2(table.l, table.k).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("Product2(%v, %v)", table.l, table.k)
        output_length := len(Product2(table.l, table.k).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}

func TestProductGeneric(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
        p int
		l Monomial
        k Monomial
		output strToInt
	}{

    }

    length_tables := []struct {
        p int
        l Monomial
        k Monomial
        output_length int
    }{
        
    }

    
    for _, table := range tables {
        call_str := fmt.Sprintf("ProductGeneric(%v, %v, %v)", table.p, table.l, table.k)
        output := ProductGeneric(table.p, table.l, table.k).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("ProductGeneric(%v, %v, %v)", table.p, table.l, table.k)
        output_length := len(ProductGeneric(table.p, table.l, table.k).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}


func TestBasis2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		n int
		output [][]int
	}{
        {0, [][]int{[]int{}}},
        {1, [][]int{[]int{1}}},
        {2, [][]int{[]int{2}}},
        {3, [][]int{[]int{3},[]int{2,1}}},
        {4, [][]int{[]int{4},[]int{3,1}}},
        {7, [][]int{[]int{4,2,1},[]int{5,2}, []int{6,1},[]int{7,}}},
    }

    length_tables := []struct {
        n int
        output_length int
    }{
        {50, 145},
        {100, 1189},
        {150, 5020},
        //{200, 15499},
        //{500, 1010134},
    }

    
    for _, table := range tables {
        call_str := fmt.Sprintf("Basis2(%v)", table.n)
        output := Basis2(table.n)
        checkGeneratorOfListsOutput(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("Basis2(%v)", table.n)
        output := Basis2(table.n)
        checkGeneratorOfListsLength(t, call_str, table.output_length, output)
    }
}
    

func TestBasisGeneric(t *testing.T){
    Mono := func(od, ev []int) Monomial{
        return Monomial{od, ev}
    }
    
    tables := []struct {
		p int
        n int
		output []Monomial
	}{
        {3, 0, []Monomial{Mono([]int{0}, []int{})}},
        {3, 1, []Monomial{Mono([]int{1}, []int{})}},
        {3, 2, []Monomial{}},
        {3, 4, []Monomial{Mono([]int{0, 0}, []int{1})}},
        {3, 5, []Monomial{Mono([]int{1, 0}, []int{1}), Mono([]int{0, 1}, []int{1})}},
        {3, 10,[]Monomial{Mono([]int{1, 1}, []int{2})}},
    }
    
    length_tables := []struct {
		p int
        n int
		output_length int
	}{    
        {3, 400, 395},
        {3, 800, 3905},
        {5, 4000, 1200},
        {23, 88000, 198},
        //{23, 176000, 692},
    }
    
    
    for _, table := range tables {
        call_str := fmt.Sprintf("BasisGeneric(%v, %v)", table.p, table.n)
        output := BasisGeneric(table.p, table.n)
        checkGeneratorOfMonomialsOutput(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("BasisGeneric(%v, %v)", table.p, table.n)
        output := BasisGeneric(table.p, table.n)
        checkGeneratorOfMonomialsLength(t, call_str, table.output_length, output)
    }
}



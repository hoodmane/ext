package main

/*
import (
    "testing"
    "fmt"
)

func TestAdemRelation2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		i int
        j int
		output strToInt
	}{
        {1,1, strToInt{}},
        {1,2, strToInt{"{0 [3]}" : 1}},
        {2,2, strToInt{"{0 [3 1]}" : 1}},
        {4,2, strToInt{"{0 [4 2]}" : 1}}, // Admissible
        {4,4, strToInt{"{0 [6 2]}" : 1, "{0 [7 1]}" : 1}},
        {5,7, strToInt{}},
        {6,7, strToInt{"{0 [13]}" : 1, "{0 [12 1]}" : 1, "{0 [10 3]}" : 1}},
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
        call_str := fmt.Sprintf("AdemRelation2(%v, %v)", table.i, table.j)
        output := AdemRelation2(table.i, table.j).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemRelation2(%v, %v)", table.i, table.j)
        output_length := len(AdemRelation2(table.i, table.j).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}


func TestAdemRelationGeneric(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
        p int
		i int
        epsilon int
        j int
		output strToInt
	}{
        {3, 1,0,1, strToInt{"{0 [2]}" : 2}},
        {3, 1,1,1, strToInt{"{1 [2]}" : 1, "{2 [2]}" : 1}},
        {3, 1,0,2, strToInt{}},
        {3, 1,1,2, strToInt{"{2 [3]}" : 1, "{1 [3]}" : 2}},
        {3, 2,0,1, strToInt{}},
        {3, 3,0,1, strToInt{"{0 [3 1]}" : 1}},
        {3, 5,0,7, strToInt{"{0 [11 1]}" : 1}},
        {3, 25,1,20, strToInt{"{1 [38 7]}" : 1, "{2 [38 7]}" : 1, "{2 [37 8]}" : 1}},
        
        {5, 1, 1, 2, strToInt{"{2 [3]}" : 1, "{1 [3]}" : 2}},
        
        {23, 1,0,1, strToInt{"{0 [2]}" : 2}},
        {23, 1,1,1, strToInt{"{1 [2]}" : 1, "{2 [2]}" : 1}},   
        {23, 2,0,1, strToInt{"{0 [3]}" : 3}},
        {23, 5,0,7, strToInt{"{0 [12]}" : 10}},
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
        call_str := fmt.Sprintf("AdemRelationGeneric(%v, %v, %v, %v)", table.p, table.i, table.epsilon, table.j)
        output := AdemRelationGeneric(table.p, table.i, table.epsilon, table.j).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }    

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemRelationGeneric(%v, %v, %v, %v)", table.p, table.i, table.epsilon, table.j)
        output_length := len(AdemRelationGeneric(table.p, table.i, table.epsilon, table.j).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }    
}



func TestMakeMonoAdmissible2(t *testing.T){
    type strToInt map[string]int
    tables := []struct {
		l []int
		output strToInt
	}{
        {[]int{}, strToInt{"{0 []}" : 1}},
        {[]int{12}, strToInt{"{0 [12]}" : 1}},
        {[]int{2, 1}, strToInt{"{0 [2 1]}" : 1}},
        {[]int{2, 2}, strToInt{"{0 [3 1]}" : 1}},
        {[]int{2, 2, 2}, strToInt{"{0 [5 1]}" : 1}},
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
    int_vec_to_bit_string := func(od []int) uint64 {
        bit_string := uint64(0)
        for i, k := range od {
            if k == 1 {
                bit_string += 1 << uint(i)
            }
        }
        return bit_string
    }
    Mono := func(od, ev []int) Monomial{
        return Monomial{int_vec_to_bit_string(od), ev}
    }
    tables := []struct {
        p int
		l Monomial
		output strToInt
	}{
        {3, Mono([]int {0}, []int {}), strToInt{"{0 []}" : 1 }},
        //{7, Mono([]int {0, 0, 0}, []int {2, 1}), strToInt{"{0 [3]}" : 3 }},
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



func TestAdemProduct2(t *testing.T){
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
        call_str := fmt.Sprintf("AdemProduct2(%v, %v)", table.l, table.k)
        output := AdemProduct2(table.l, table.k).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemProduct2(%v, %v)", table.l, table.k)
        output_length := len(AdemProduct2(table.l, table.k).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}

func TestAdemProductGeneric(t *testing.T){
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
        call_str := fmt.Sprintf("AdemProductGeneric(%v, %v, %v)", table.p, table.l, table.k)
        output := AdemProductGeneric(table.p, table.l, table.k).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemProductGeneric(%v, %v, %v)", table.p, table.l, table.k)
        output_length := len(AdemProductGeneric(table.p, table.l, table.k).GetCoeffMap())
        checkEqLengths(t, call_str, table.output_length, output_length)
    }
}


func TestAdemBasis2(t *testing.T){
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
        //{150, 5020},
        //{200, 15499},
        //{300, 87977},
        //{400, 335566},
        //{500, 1010134},
    }

    
    for _, table := range tables {
        call_str := fmt.Sprintf("AdemBasis2(%v)", table.n)
        output := AdemBasis2(table.n)
        checkGeneratorOfListsOutput(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemBasis2(%v)", table.n)
        output := AdemBasis2(table.n)
        checkGeneratorOfListsLength(t, call_str, table.output_length, output)
    }
}
    

func TestAdemBasisGeneric(t *testing.T){
    int_vec_to_bit_string := func(od []int) uint64 {
        bit_string := uint64(0)
        for i, k := range od {
            if k == 1 {
                bit_string += 1 << uint(i)
            }
        }
        return bit_string
    }
    Mono := func(od, ev []int) Monomial{
        return Monomial{int_vec_to_bit_string(od), ev}
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
        //{5, 4000, 1200},
        //{23, 88000, 198},
        //{23, 176000, 692},
    }
    
    
    for _, table := range tables {
        call_str := fmt.Sprintf("AdemBasisGeneric(%v, %v)", table.p, table.n)
        output := AdemBasisGeneric(table.p, table.n)
        checkGeneratorOfMonomialsOutput(t, call_str, table.output, output)
    }

    for _, table := range length_tables {
        call_str := fmt.Sprintf("AdemBasisGeneric(%v, %v)", table.p, table.n)
        output := AdemBasisGeneric(table.p, table.n)
        checkGeneratorOfMonomialsLength(t, call_str, table.output_length, output)
    }
}

*/

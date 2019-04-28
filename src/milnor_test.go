package main

import (
    "testing"
    "fmt"
)

    
func TestInitializeMilnorMatrix(t *testing.T){
    M := initialize_milnor_matrix([]int{4,5,6},[]int{1,2,3, 4, 5})
    if len(M) != 4 {
        t.Errorf("Length of M should be 4 but is %v", len(M))
    }
    if len(M[0]) != 6 {
        t.Errorf("Depth of M should be 6 but is %v", len(M[0]))
    }
    if !eqListsQ(M[0], []int{0, 1,2,3, 4, 5}){
        t.Errorf("First row of M should be [0 1 2 3 4 5] but is %v", M[0])
    }
    first_col := []int{0,4,5,6}
    for i := 0; i < 4; i++ {
        if M[i][0] != first_col[i] {
            t.Errorf("First column of M should be %v but is [%v %v %v %v]", first_col, M[0][0], M[1][0], M[2][0], M[3][0])
            break
        }
    }
    for i:= 1; i < 4; i++ {
        for j:= 1; j < 6; j++ {
            if M[i][j] != 0 {
                t.Errorf("Entry (%v, %v) of M is nonzero", i, j)
            }
        }
    }
}



func TestStepMilnorMatrix(t *testing.T) {

}    

func TestMilnorMatrices(t *testing.T) {
    tables := []struct {
		p int
        r []int
        s []int
        output [][][]int
	}{
        {2, []int{}, []int{}, [][][]int{[][]int{[]int{0}}}},
        {2, []int{}, []int{5}, [][][]int{[][]int{[]int{0,5}}}},
        {2, []int{5}, []int{}, [][][]int{[][]int{[]int{0}, []int{5}}}},
        {2, []int{1}, []int{1}, [][][]int{[][]int{[]int{0,1}, []int{1,0}}}},
        {2, []int{0, 2}, []int{1}, [][][]int{[][]int{[]int{0,1}, []int{0,0}, []int{2,0}},[][]int{[]int{0,0}, []int{0,0}, []int{0,1}}}},
        {2, []int{0, 4}, []int{1,1},  [][][]int{[][]int{[]int{0, 1, 1}, []int{0, 0, 0}, []int{4, 0, 0}}, [][]int{[]int{0, 0, 1}, []int{0, 0, 0}, []int{2, 1, 0}}, [][]int{[]int{0, 1, 0}, []int{0, 0, 0}, []int{0, 0, 1}}}},
        {3, []int{0, 3}, []int{1}, [][][]int{[][]int{[]int{0, 1}, []int{0, 0}, []int{3, 0}}, [][]int{[]int{0, 0}, []int{0, 0}, []int{0, 1}}}},
        {3, []int{0, 9}, []int{1,1}, [][][]int{[][]int{[]int{0, 1, 1}, []int{0, 0, 0}, []int{9, 0, 0}}, [][]int{[]int{0, 0, 1}, []int{0, 0, 0}, []int{6, 1, 0}}, [][]int{[]int{0, 1, 0}, []int{0, 0, 0}, []int{0, 0, 1}}}},
    }
    
    for _, table := range tables {
        output := milnor_matrices(table.p, table.r, table.s)
        idx := 0
        for matrix := range output {
            if idx >= len(table.output) {
                t.Errorf("Expecting only %v outputs, got %v", len(table.output), idx + 1)
                continue
            }
            if !eqMatrixQ(matrix, table.output[idx]) {
                t.Errorf("Expecting %v as output %v of milnor_matrices(%v, %v, %v) got %v", table.output[idx], idx, table.p, table.r, table.s, matrix)
            }
            idx++
        }
    }
}



func TestRemoveTrailingZeroes(t *testing.T) {
    length_tables := []struct {
        l []int
        output_length int
	}{
        { []int{1,0,0,1,0,0,1}, 7 },
        { []int{1,0,0,1,0,0,0}, 4 },
        { []int{1,0,0,0,0,0,0}, 1 },
        { []int{0,0,0,0,0,0,0}, 0 },
        { []int{}, 0 },
    }
    
    for _, table := range length_tables {
        length := len(remove_trailing_zeroes(table.l))
        if  length != table.output_length {
            t.Errorf("Expecting remove_trailing_zeroes(%v) to have length %v, got length %v", table.l, table.output_length, length)
        }
    }
}

func TestMilnorProductEven(t *testing.T) {
    tables := []struct {
        p int
        r []int
        s []int
        output map[string]int
	}{
        {3, []int{1},   []int{1},   map[string]int{"{0 [2]}" : 2}},
        {3, []int{1},   []int{0, 1},map[string]int{"{0 [1 1]}" : 1}},
        {2, []int{0,2}, []int{1},   map[string]int{"{0 [1 2]}": 1, "{0 [0 0 1]}": 1}},
        {3, []int{0,3}, []int{1},   map[string]int{"{0 [1 3]}": 1, "{0 [0 0 1]}": 1}},
        {3, []int{0,9}, []int{1,1}, map[string]int{"{0 [1 10]}": 1, "{0 [0 7 1]}" : 1, "{0 [1 0 0 1]}": 1}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductEven(%v, %v, %v)", table.p, table.r, table.s)
        output := MilnorProductEven(table.p, table.r, table.s).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }
}

func TestMilnorProductFullQpart(t *testing.T) {
    int_vec_to_bit_string := func(od []int) uint64 {
        bit_string := uint64(0)
        for _, k := range od {
            bit_string += 1 << uint(k)
        }
        return bit_string
    }
    Mono := func(od, ev []int) Monomial{
        return Monomial{int_vec_to_bit_string(od), ev}
    }
    
    mstring := func(odd []int, ev string) string{
        bit_string := int_vec_to_bit_string(odd)
        return fmt.Sprintf("{%v %v}", bit_string, ev)
    }
    
    tables := []struct {
        p int
        m Monomial
        f []int
        output map[string]int
	}{
        {3, Mono([]int{},[]int{1}),[]int{0}, map[string]int{mstring([]int{0}, "[1]") : 1, mstring([]int{1}, "[]"): 1}},
        {5, Mono([]int{0,2},[]int{5}),[]int{1}, map[string]int{mstring([]int{0, 1, 2}, "[5]") : 4}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductFullQpart(%v, %v, %v)", table.p, table.m, table.f)
        output := MilnorProductFullQpart(table.p, table.m, int_vec_to_bit_string(table.f)).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }
}


func TestMilnorProductFull(t *testing.T) {
    int_vec_to_bit_string := func(od []int) uint64 {
        bit_string := uint64(0)
        for _, k := range od {
            bit_string += 1 << uint(k)
        }
        return bit_string
    }
    Mono := func(od, ev []int) Monomial{
        return Monomial{int_vec_to_bit_string(od), ev}
    }
    
    mstring := func(odd []int, ev string) string{
        bit_string := int_vec_to_bit_string(odd)
        return fmt.Sprintf("{%v %v}", bit_string, ev)
    }
    
    tables := []struct {
        p int
        r Monomial
        s Monomial
        output map[string]int
	}{
        {3, Mono([]int{},[]int{1}),Mono([]int{0},[]int{}), map[string]int{mstring([]int{0}, "[1]") : 1, mstring([]int{1}, "[]") : 1}},
        {5, Mono([]int{0,2},[]int{5}),Mono([]int{1},[]int{1}), map[string]int{mstring([]int{0, 1, 2}, "[0 1]") : 4, mstring([]int{0, 1, 2}, "[6]") : 4}},
        {7, Mono([]int{0,2,4},[]int{}),Mono([]int{1,3},[]int{}), map[string]int{mstring([]int{0, 1, 2, 3, 4}, "[]"): 6}},
        {7, Mono([]int{0,2,4},[]int{}),Mono([]int{1,5},[]int{}), map[string]int{mstring([]int{0, 1, 2, 4, 5}, "[]"): 1}},
        {3, Mono([]int{},[]int{6}),Mono([]int{},[]int{2}), map[string]int{mstring([]int{}, "[0 2]") : 1, mstring([]int{}, "[4 1]") : 1, mstring([]int{}, "[8]") : 1}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductFull(%v, %v, %v)", table.p, table.r, table.s)
        output := MilnorProductFull(table.p, table.r, table.s).GetCoeffMap()
        checkEqStrToIntMaps(t, call_str, table.output, output)
    }
}


func TestMilnorProduct2(t *testing.T){

}
            

func TestBasisEven(t *testing.T){
    empty_profile := ProfileList{[]int{},false,false}
    empty_full_profile := FullProfile{empty_profile, empty_profile}
    EvProf := func(profile []int, truncated bool) FullProfile {
        return FullProfile{empty_profile, ProfileList{profile, truncated, true}}
    }
    Alg := func(p int) MilnorAlgebra {
        return MinimalMilnorAlgebra{p,p!=2,empty_full_profile}
    }
    
    
    tables := []struct {
        alg MilnorAlgebra
        n int
        output [][]int
	}{
        {Alg(2), 2, [][]int{[]int{2}}},
        {Alg(2), 3, [][]int{[]int{0, 1}, []int{3}}},
        {Alg(2), 4, [][]int{[]int{1, 1}, []int{4}}},
        {Alg(2), 7, [][]int{[]int{0, 0, 1}, []int{1, 2}, []int{4, 1}, []int{7}}},
        
        {MinimalMilnorAlgebra{2, false,EvProf([]int{2,1}, true)}, 4, [][]int{[]int{1, 1}}},
        {MinimalMilnorAlgebra{2, false,EvProf([]int{}, false)},   4, [][]int{[]int{1, 1}, []int{4}}},
    }
    
    length_tables := []struct {
        alg MilnorAlgebra
        n int
        output_length int
	}{
    
    }    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorBasis2(%v, %v)", table.alg, table.n)
        output := MilnorBasis2(table.alg, table.n)
        checkGeneratorOfListsOutput(t, call_str, table.output, output)
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("MilnorBasis2(%v, %v)", table.alg, table.n)
        output := MilnorBasis2(table.alg, table.n)
        checkGeneratorOfListsLength(t, call_str, table.output_length, output)
    }    
}




func TestBasisGenericQPart(t *testing.T){
    
}


func TestMilnorBasisGeneric(t *testing.T){
    empty_profile := ProfileList{[]int{},false,false}
    empty_full_profile := FullProfile{empty_profile, empty_profile}
    Prof := func(odd_profile , even_profile []int, truncated bool) FullProfile {
        return FullProfile{ProfileList{odd_profile, truncated, true}, ProfileList{even_profile, truncated, true}}
    }
    Alg := func(p int) MilnorAlgebra {
        return MinimalMilnorAlgebra{p,p!=2,empty_full_profile}
    }
    
    int_vec_to_bit_string := func(od []int) uint64 {
        bit_string := uint64(0)
        for _, k := range od {
            bit_string += 1 << uint(k)
        }
        return bit_string
    }
    Mono := func(od, ev []int) Monomial{
        return Monomial{int_vec_to_bit_string(od), ev}
    }


    tables := []struct {
        alg MilnorAlgebra
        n int
        output []Monomial
	}{
        {Alg(3), 1, []Monomial{Mono([]int{0},[]int{})}},
        {Alg(3), 9, []Monomial{Mono([]int{1},[]int{1}), Mono([]int{0},[]int{2})}},
        {Alg(3), 17, []Monomial{Mono([]int{2},[]int{}),Mono([]int{1},[]int{3}),Mono([]int{0},[]int{0,1}),Mono([]int{0},[]int{4})}},
        {Alg(5), 48, []Monomial{Mono([]int{},[]int{0, 1}),Mono([]int{},[]int{6})}},
    }
    
    length_tables := []struct {
        alg MilnorAlgebra
        n int
        output_length int
	}{
        {Alg(3), 100, 13},
        {Alg(7), 200, 0},
        {Alg(7), 240, 3},
        {Alg(7), 240, 3},
        {MinimalMilnorAlgebra{7, true, Prof([]int{}, []int{}, true)}, 240, 0},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorBasisGeneric(%v, %v)", table.alg, table.n)
        output := MilnorBasisGeneric(table.alg, table.n)        
        checkGeneratorOfMonomialsOutput(t, call_str, table.output, output)
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("MilnorBasisGeneric(%v, %v)", table.alg, table.n)
        output := MilnorBasisGeneric(table.alg, table.n)
        checkGeneratorOfMonomialsLength(t, call_str, table.output_length, output)
    }    

}

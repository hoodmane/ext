package main

import (
    "testing"
    "fmt"
    "unsafe"
)

var empty_profile = Profile{[]int{},false,false}
var empty_full_profile = FullProfile{empty_profile, empty_profile}

func EvProf(profile []int, truncated bool) FullProfile {
    return FullProfile{empty_profile, Profile{profile, truncated, true}}
}
    
func Prof(odd_profile , even_profile []int, truncated bool) FullProfile {
    return FullProfile{Profile{odd_profile, truncated, true}, Profile{even_profile, truncated, true}}
}
func Alg(p int) MilnorAlgebra {
    return MilnorAlgebra{p,p!=2,empty_full_profile,""}
}

var A2 = Alg(2)
var A3 = Alg(3)
var A5 = Alg(5)
var A7 = Alg(7)


func int_vec_to_bit_string(p int, od []int) (int, uint64 ){
    bit_string := uint64(0)
    deg := 0
    for _, k := range od {
        deg += tau_degrees[p][k]
        bit_string += 1 << uint(k)
    }
    return deg, bit_string
}
type MBE MilnorBasisElement
func Mono(algebra MilnorAlgebra, od, ev []int) MBE{
    p := algebra.p
    p_deg := 0
    for idx, n := range ev {
        p_deg += n * xi_degrees[p][idx]
    }
    if algebra.generic {
        p_deg *= 2*(p-1)
    }
    q_deg, q_part := int_vec_to_bit_string(p, od)
    return MBE(MilnorBasisElement{q_deg, q_part, p_deg, ev})
}

func convertMBEList(list []MBE) []MilnorBasisElement {
    return *(*[]MilnorBasisElement)(unsafe.Pointer(&list))
}

func mstring(p int, odd []int, ev string) string{
    _, bit_string := int_vec_to_bit_string(p, odd)
    return fmt.Sprintf("{%v %v}", bit_string, ev)
}


func MonomialIndexToIntToStringToInt(algebra MilnorAlgebra, output map[MonomialIndex]int) map[string]int{
    str_to_int := make(map[string]int)
    for key, value := range output {
        b := GetMilnorBasisElementFromIndex(algebra, key)
        str_to_int[b.String()] = value    
    }
    return str_to_int
}

func TestSetup(t *testing.T){
    GenerateMilnorBasis(A2, 100)
    GenerateMilnorBasis(A3, 200)
    GenerateMilnorBasis(A5, 300)
    GenerateMilnorBasis(A7, 300)
}

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


          

func TestBasisEven(t *testing.T){
    tables := []struct {
        alg MilnorAlgebra
        n int
        output [][]int
	}{
        {A2, 2, [][]int{[]int{2}}},
        {A2, 3, [][]int{[]int{0, 1}, []int{3}}},
        {A2, 4, [][]int{[]int{1, 1}, []int{4}}},
        {A2, 7, [][]int{[]int{0, 0, 1}, []int{1, 2}, []int{4, 1}, []int{7}}},
        
        {MilnorAlgebra{2, false,EvProf([]int{2,1}, true),""}, 4, [][]int{[]int{1, 1}}},
        {MilnorAlgebra{2, false,EvProf([]int{}, false),""},   4, [][]int{[]int{1, 1}, []int{4}}},
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
        error_str := checkGeneratorOfListsOutput(call_str, table.output, output)
        if error_str != "" {
            t.Error(error_str)
        }
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("MilnorBasis2(%v, %v)", table.alg, table.n)
        output := MilnorBasis2(table.alg, table.n)
        error_str := checkGeneratorOfListsLength(call_str, table.output_length, output)
        if error_str != "" {
            t.Error(error_str)
        }        
    }    
}




func TestBasisGenericQPart(t *testing.T){
    
}


func TestMilnorBasisGeneric(t *testing.T){
    tables := []struct {
        alg MilnorAlgebra
        n int
        output []MBE
	}{
        {A3, 1, []MBE{Mono(A3, []int{0},[]int{})}},
        {A3, 9, []MBE{Mono(A3, []int{1},[]int{1}), Mono(A3, []int{0},[]int{2})}},
        {A3, 17, []MBE{Mono(A3, []int{2},[]int{}),Mono(A3, []int{1},[]int{3}),Mono(A3, []int{0},[]int{0,1}),Mono(A3, []int{0},[]int{4})}},
        {A5, 48, []MBE{Mono(A5, []int{},[]int{0, 1}),Mono(A5, []int{},[]int{6})}},
    }
    
    length_tables := []struct {
        alg MilnorAlgebra
        n int
        output_length int
	}{
        {A3, 100, 13},
        {A7, 200, 0},
        {A7, 240, 3},
        {A7, 240, 3},
        {MilnorAlgebra{7, true, Prof([]int{}, []int{}, true),""}, 240, 0},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorBasisGeneric(%v, %v)", table.alg.String(), table.n)
        output := MilnorBasisGeneric(table.alg, table.n)        
        error_str := checkGeneratorOfMonomialsOutput(call_str, convertMBEList(table.output), output)
        if error_str != "" {
            t.Error(error_str)
        }
    }
    
    for _, table := range length_tables {
        call_str := fmt.Sprintf("MilnorBasisGeneric(%v, %v)", table.alg, table.n)
        output := MilnorBasisGeneric(table.alg, table.n)
        error_str := checkGeneratorOfMonomialsLength(call_str, table.output_length, output)
        if error_str != "" {
            t.Error(error_str)
        }        
    }    

}



func TestMilnorProductEven(t *testing.T) {
    tables := []struct {
        algebra MilnorAlgebra
        r []int
        s []int
        output map[string]int
	}{
        {A3, []int{1},   []int{1},   map[string]int{"{0 [2]}" : 2}},
        {A3, []int{1},   []int{0, 1},map[string]int{"{0 [1 1]}" : 1}},
        {A2, []int{0,2}, []int{1},   map[string]int{"{0 [1 2]}": 1, "{0 [0 0 1]}": 1}},
        {A3, []int{0,3}, []int{1},   map[string]int{"{0 [1 3]}": 1, "{0 [0 0 1]}": 1}},
        {A3, []int{0,9}, []int{1,1}, map[string]int{"{0 [1 10]}": 1, "{0 [0 7 1]}" : 1, "{0 [1 0 0 1]}": 1}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductEven(%v, %v, %v)", table.algebra, table.r, table.s)
        r_deg := 0
        p := table.algebra.p
        var q int
        if table.algebra.generic {
            q = 2*(p-1)
        } else {
            q = 1
        }
        for idx, power := range table.r {
            r_deg += xi_degrees[p][idx] * power
        }
        r_deg *= q
        s_deg := 0
        for idx, power := range table.s {
            s_deg += xi_degrees[p][idx] * power
        }
        s_deg *= q
        r := MilnorBasisElement{0, 0, r_deg, table.r}
        s := MilnorBasisElement{0, 0, s_deg, table.s}
        output := MilnorProductEven(table.algebra, r, s).GetCoeffMap()
        error_str := checkEqStrToIntMaps(call_str, table.output, MonomialIndexToIntToStringToInt(table.algebra, output))
        if error_str != "" {
            t.Error(error_str)
        }
    }
}

func TestMilnorProductFullQpart(t *testing.T) {    
    tables := []struct {
        algebra MilnorAlgebra
        m MBE
        f []int
        output map[string]int
	}{
        {A3, Mono(A3, []int{},[]int{1}),[]int{0}, map[string]int{mstring(3, []int{0}, "[1]") : 1, mstring(3, []int{1}, "[]"): 1}},
        {A5, Mono(A5, []int{0,2},[]int{5}),[]int{1}, map[string]int{mstring(5, []int{0, 1, 2}, "[5]") : 4}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductFullQpart(%v, %v, %v)", table.algebra, table.m, table.f)
        _, bit_string := int_vec_to_bit_string(table.algebra.p, table.f)
        output := MilnorProductFullQpart(table.algebra, MilnorBasisElement(table.m), bit_string).GetCoeffMap()
        error_str := checkEqStrToIntMaps(call_str, table.output, MonomialIndexToIntToStringToInt(table.algebra, output))
        if error_str != "" {
            t.Error(error_str)
        }
    }
}


func TestMilnorProductFull(t *testing.T) {
    tables := []struct {
        algebra MilnorAlgebra
        r MBE
        s MBE
        output map[string]int
	}{ // Some of these tests involve Q's that are too far out to handle with the current "make a giant table" approach
        {A3, Mono(A3, []int{},[]int{1}),Mono(A3, []int{0},[]int{}), map[string]int{mstring(3, []int{0}, "[1]") : 1, mstring(3, []int{1}, "[]") : 1}},
        {A5, Mono(A5, []int{0,2},[]int{5}),Mono(A5, []int{1},[]int{1}), map[string]int{mstring(5, []int{0, 1, 2}, "[0 1]") : 4, mstring(5, []int{0, 1, 2}, "[6]") : 4}},
        //{A7, Mono(A7, []int{0,2,4},[]int{}),Mono(A7, []int{1,3},[]int{}), map[string]int{mstring(7, []int{0, 1, 2, 3, 4}, "[]"): 6}},
        //{A7, Mono(A7, []int{0,2,4},[]int{}),Mono(A7, []int{1,5},[]int{}), map[string]int{mstring(7, []int{0, 1, 2, 4, 5}, "[]"): 1}},
        {A3, Mono(A3, []int{},[]int{6}),Mono(A3, []int{},[]int{2}), map[string]int{mstring(3, []int{}, "[0 2]") : 1, mstring(3, []int{}, "[4 1]") : 1, mstring(3, []int{}, "[8]") : 1}},
    }
    
    for _, table := range tables {
        call_str := fmt.Sprintf("MilnorProductFull(%v, %v, %v)", table.algebra, table.r, table.s)
        output := MilnorProductFull(table.algebra, MilnorBasisElement(table.r), MilnorBasisElement(table.s)).GetCoeffMap()
        error_str := checkEqStrToIntMaps(call_str, table.output, MonomialIndexToIntToStringToInt(table.algebra, output))
        if error_str != "" {
            t.Error(error_str)
        }
    }
}

/*



func TestMilnorProduct2(t *testing.T){

}
  

/**/

//
//  File: adem.go
//  Author: Hood Chatham
//
//  Defines the basic operations on Adem algebras.
//
//  The implementations here are based on the Sage code, in particular from the two
//  files steenrod_algebra_mult.py and (the basis code from) steenrod_algebra_basis.py
//  I have made significant improvements in code legibility.
//



package main

/*

import (
     "fmt"
//     "github.com/pkg/profile"
//     "sync"
)
func main(){
    for x := range AdemBasisGeneric(3, 5) {
        fmt.Println("out:", x)
    }
}


type AdemAlgebra struct {
    p int
    generic bool
}

type AdemElement struct {
    Vector
    generic bool
}

func NewAdemZeroVector2(size_hint int) AdemElement {
    return AdemElement{NewZeroVector(2, size_hint), false}
}

func NewAdemZeroVectorGeneric(p, size_hint int) AdemElement {
    return AdemElement{NewZeroVector(p, size_hint), true}
}

func NewAdemBasisVector2(even_part []int) AdemElement{
    return AdemElement{NewBasisVector2(even_part), false}
}

func NewAdemBasisVectorGeneric(p int, odd_part uint64, even_part []int) AdemElement{
    return AdemElement{NewBasisVector(p, odd_part, even_part), true}
}

//Return the adem relation Sqa * Sqb when p=2
func AdemRelation2(a, b int) AdemElement {
    if b == 0 {
        return NewAdemBasisVector2([]int {a})
    } 
    if a == 0 {
        return NewAdemBasisVector2([]int {b})
    } 
    if a >= 2*b {
        return NewAdemBasisVector2([]int {a, b})
    }
    result := NewAdemZeroVector2(1 + a/2)
    for j := 0; j < 1 + a/2; j++ {
        if Binomial2(b- j - 1, a - 2*j) == 1 {
            if j == 0 {
                result.Set2([]int {a + b}, 1)
            } else {
                result.Set2([]int {a + b - j, j}, 1)
            }
        }
    }
    return result
}


//Return the generic adem relation for P(A)*P(B) or P(A) * beta * P(B)
func AdemRelationGeneric(p, A, bockstein, B int)  AdemElement {
    if A == 0 {
        return NewAdemBasisVectorGeneric(p, uint64(bockstein), []int{B} )
    }
    if B == 0 {
        return NewAdemBasisVectorGeneric(p, uint64(bockstein) << 1, []int{A} )
    }
    if A >= p*B + bockstein { // Admissible
        return NewAdemBasisVectorGeneric(p, uint64(bockstein) << 1, []int{A, B} )
    }
    
    result := NewAdemZeroVectorGeneric(p, A/p * (1 + bockstein))
    for j := 0; j <= A/p; j++ {
        coeff := BinomialOdd((B-j) * (p-1) - 1 + bockstein, A - p*j, p)
        coeff *= MinusOneToTheN(A + j)
        coeff = coeff % p        
        if coeff != 0 && j == 0 {
            result.Set(uint64(bockstein), []int{A + B}, coeff)
        } else if coeff != 0 && j != 0 {
            result.Set(uint64(bockstein), []int{A + B - j, j}, coeff)
        }
    }
    if bockstein == 1 {
        for j := 0; j <= (A-1)/p; j++ {
            coeff := BinomialOdd((B - j)*(p - 1) - 1, A - p*j - 1, p)
            coeff *= MinusOneToTheN(A + j - 1)
            if coeff != 0 && j == 0 {
                result.Set(uint64(bockstein) << 1, []int{A + B}, coeff)
            } else if coeff != 0 && j != 0 {
                result.Set(uint64(bockstein) << 1, []int{A + B - j, j}, coeff)
            }
        }
    }
    return result
}

//@memoized
//Reduce a monomial into a linear combination of admissible monomials when p = 2
func MakeMonoAdmissible2(mono []int) AdemElement {   
    if len(mono) == 1 {
        return NewAdemBasisVector2(mono)
    }
    first_inadmissible_index := -1
    for j := 0; j < len(mono) - 1; j++ {
        if mono[j] < 2*mono[j+1] {
            first_inadmissible_index = j
            break
        }
    }
    if first_inadmissible_index == -1 { // Didn't find any inadmissible indices
        return NewAdemBasisVector2(mono)
    }
    result := NewAdemZeroVector2(-1)
    j := first_inadmissible_index
    y := AdemRelation2(mono[j], mono[j + 1])
    for _, m := range y.GetBasisVectorMap() {
        x := m.even_part
        new_mono := make([]int, len(mono) + len(x) - 2)
        copy(new_mono, mono[:j])
        copy(new_mono[j:], x)
        copy(new_mono[j+len(x):], mono[j+2:])
        new := MakeMonoAdmissible2(new_mono)
        result.Add(new)
    }
    return result
}

//Reduce a monomial into a linear combination of admissible monomials for the generic Steenrod algebra"""
func MakeMonoAdmissibleGeneric(p int, mono Monomial) AdemElement {
    // check to see if admissible:
    odd_part := mono.odd_part
    even_part := mono.even_part
    first_inadmissible_index := -1
    for j := 0; j < len(even_part) - 1; j++ {
        if even_part[j] < int((odd_part >> uint(j+1)) & 1) + p * even_part[j+1] {
            first_inadmissible_index = j
            break
        }
    }
    if first_inadmissible_index == -1 { // Didn't find any inadmissible indices
        return NewAdemBasisVectorGeneric(p, odd_part, even_part)
    }
    result := NewAdemZeroVectorGeneric(p, -1)
    j := uint(first_inadmissible_index)
    y := AdemRelationGeneric(p, even_part[j],int((odd_part>>(j+1))&1), even_part[j+1])
    coeff_map := y.GetCoeffMap()
    for key, m := range y.GetBasisVectorMap() {
        coeff := coeff_map[key]
        reln_even_part := m.even_part
        reln_odd_part := m.odd_part
        
        if (odd_part>>j)&1     == 1 && reln_odd_part & 1 == 1  || 
           (odd_part>>(j+2))&1 == 1 && (reln_odd_part >> uint(len(reln_even_part))) & 1 == 1 {
            continue
        }
        
        new_even_part := make([]int, len(even_part) + len(reln_even_part) - 2)
        copy(new_even_part, even_part[:j])
        copy(new_even_part[j:], reln_even_part)
        copy(new_even_part[j+uint(len(reln_even_part)):], even_part[j+2:])
        
        new_odd_part := odd_part & ( (1 << j+1) - 1)
        new_odd_part += reln_odd_part << j
        new_odd_part += (odd_part & ^( (1 << j+3) - 1)) >> uint(len(reln_even_part)) - 1
        
        new := MakeMonoAdmissibleGeneric(p, Monomial{new_odd_part, new_even_part})
        result.ScaleAndAdd(new, coeff) 
    }
    return result
}

//Multiply monomials m1 and m2 and write the result in the Adem basis for p = 2.
func AdemProduct2(m1, m2 []int) AdemElement {
    m := make([]int, len(m1) + len(m2))
    copy(m, m1)
    copy(m[len(m1):], m2)
    //fmt.Println(m)
    return MakeMonoAdmissible2(m)
}

//Multiply monomials m1 and m2 and write the result in the Adem basis in the generic case.
func AdemProductGeneric(p int, m1, m2 Monomial) AdemElement {
    if (m1.odd_part >> uint(len(m1.even_part))) & 1 == 1 && m2.odd_part & 1 == 1 {
        return NewAdemZeroVectorGeneric(p, 0)
    }
    odd_part := m1.odd_part
    odd_part += m2.odd_part << uint(len(m1.even_part))
    
    even_part := make([]int, len(m1.even_part) + len(m2.even_part))
    copy(even_part, m1.even_part)
    copy(even_part[len(m1.even_part):], m2.even_part)
    return MakeMonoAdmissibleGeneric(p, Monomial{odd_part, even_part})
}


type dim_last_epsilon struct{
    dim int
    last int
    epsilon int
}
var adem_basis_2_table map[dim_last_epsilon] [][]int = make(map[dim_last_epsilon] [][]int)
var adem_basis_generic_table map[int] map[dim_last_epsilon] []Monomial = make(map[int] map[dim_last_epsilon] []Monomial)

//Get the basis for the n dimensional part of the Steenrod algebra.
//Build basis recursively.  last = last term.
//last is >= bound, and we will append (last,) to the end of
//elements from basis (n - last, bound=2 * last).
//This means that 2 last <= n - last, or 3 last <= n.
//extra_length tells us how many extra elements are planned for insertion onto the end
//so we don't have to do any copies.
func ademBasis2Helper(n, bound, extra_length int) <- chan []int {
    ch := make(chan []int, 50)
    go func() {
        defer close(ch)    
        if n == 0 {
            ch <- make([]int, 0, extra_length)
            return
        }
        result := make([]int, 1, 1 + extra_length)
        result[0] = n
        ch <- result
        for last := bound; last <= n / 3; last ++ {
            ademBasis2HelperLast(ch, n - last, last, extra_length)
        }
    }()
    return ch
}

func ademBasis2HelperLast(ch chan<- []int, n, last, extra_length int) {
    dim_last := dim_last_epsilon{n, last, 0}
    elt_list, ok := adem_basis_2_table[dim_last]
    if ok {
        for _, elt := range elt_list {
            ch <- elt
        }
        return
    }
    
    for vec := range ademBasis2Helper(n , 2 * last, extra_length + 1) {
        vec = append(vec, last)
        elt_list  = append(elt_list, vec)
        ch <- vec
    }
    adem_basis_2_table[dim_last] = elt_list
}

//Get the basis for the n dimensional part of the Steenrod algebra.
func AdemBasis2(n int) <-chan []int {
    return ademBasis2Helper(n, 1, 0)
}

//Get the basis for the n dimensional part of the Steenrod algebra.
func ademBasisGenericHelper(p, n, bound, extra_length int) <-chan Monomial {
    ch := make(chan Monomial, 200)

    go func() {
        defer close(ch)      
        if n == 0 {
            odd_result  := uint64(0)
            even_result := make([]int, 0, extra_length)
            ch <- Monomial{odd_result, even_result}
            return
        }
        if n == 1 {
            odd_result  := uint64(1)
            even_result := make([]int, 0, extra_length)
            ch <- Monomial{odd_result, even_result}
            return
        }
            
        // append P^{last} beta^{epsilon}
        for epsilon := 0; epsilon <= 1; epsilon++ {
            lower_bound := bound + epsilon 
            // Without this lower bound edge case we lose for instance the element (0, 1, 1) in degree 5.
            // I don't have a good explanation for what it means yet.            
            if bound == 1 {
                lower_bound = 1
            }
            for last := lower_bound; last <= n / (2*(p - 1)); last++ {
                ademBasisGenericHelperLast(p, ch, n, epsilon, last, extra_length)
            }
        }
    }()
    return ch
}

func ademBasisGenericHelperLast(p int, ch chan<- Monomial, n, epsilon, last, extra_length int) {
    dim_last := dim_last_epsilon{n, last, epsilon}
    elt_list, ok := adem_basis_generic_table[p][dim_last]
    if ok {
        for _, mono := range elt_list {
            odd_part := mono.odd_part
            even_part := mono.even_part
            if extra_length > 0 {
                output_even_part := make([]int, len(even_part), len(even_part) + extra_length)
                copy(output_even_part, even_part)
                even_part = output_even_part
            }
            ch <- Monomial{odd_part, even_part}
        }
        return
    }
    elt_list = make([]Monomial, 0, 20)
    for vec := range ademBasisGenericHelper(p, n - 2*(p -1)*last - epsilon, p * last, extra_length + 1) {
        odd_part := vec.odd_part
        even_part := vec.even_part
        odd_part += uint64(epsilon) << uint(len(vec.even_part)+1)
        even_part = append(even_part, last)
        mono := Monomial{odd_part, even_part}
        elt_list  = append(elt_list, mono)
        ch <- mono
    }
    adem_basis_generic_table[p][dim_last] = elt_list
}

func AdemBasisGeneric(p, n int) <- chan Monomial {
    _, ok := adem_basis_generic_table[p]
    if !ok {
        basis_table := make(map[dim_last_epsilon] []Monomial)
        adem_basis_generic_table[p] = basis_table
    }
    return ademBasisGenericHelper(p, n, 1, 0)
}

*/

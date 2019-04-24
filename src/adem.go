//
//  File: adem.go
//  Author: Hood Chatham
//
//  Defines the basic operations on Adem algebras. No object oriented code,
//  all functions take basis vectors as inputs and output dictionaries representing
//  Fp-linear combinations of basis vectors. The relevant calls will be wrapped
//  by methods of the AdemAlgebra and AdemElement classes in steenrod.py.
//  They are systemtically extended over vector inputs using the
//  @linearextension decorator.
//
//  The algorithms here are all copied from Sage, in particular from the two
//  files steenrod_algebra_mult.py and (the basis code from) steenrod_algebra_basis.py
//  I have made significant improvements in code legibility and removed all references to Sage code.
//

package main

import (
    "fmt"
)

func main(){
//    a := Adem2(2, 3)
//    fmt.Println("a", a.GetCoeffMap())
//    
//    b := AdemGeneric(3, 2, 1, 4)
//    fmt.Println("b", b.GetCoeffMap())
//    c := MakeMonoAdmissible2([]int{24,12,4,6,1})
//    fmt.Println("c", c.GetCoeffMap())
//    
//    d := MakeMonoAdmissibleGeneric(3, Monomial{[]int{0,1,0},[]int{1,1}})
//    fmt.Println("d", d.GetCoeffMap())
//    
//    e := Product2([]int{2,1}, []int {2, 1})
//    fmt.Println("e", e.GetCoeffMap())
    
    for x := range BasisGeneric(3, 4) {
        fmt.Println(x)
    }
}

type AdemAlgebra struct {
    p int
    generic bool
}

type AdemElement struct {
    Vector
}


//Return the adem relation Sqa * Sqb when p=2
func Adem2(a, b int) AdemElement {
    if b == 0 {
        return AdemElement{NewBasisVector2([]int {a})}
    } 
    if a == 0 {
        return AdemElement{NewBasisVector2([]int {b})}
    } 
    if a >= 2*b {
        return AdemElement{NewBasisVector2([]int {a, b})}
    }
    result := AdemElement{NewZeroVector(2, 1 + a/2)}
    for j := 0; j < 1 + a/2; j++ {
        if Binomial2(b-j-1, a-2*j) == 1 {
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
func AdemGeneric(p, A, bockstein, B int)  AdemElement {
    if A == 0 {
        return AdemElement{NewBasisVector(p, []int{bockstein, 0}, []int{B} )}
    }
    if B == 0 {
        return AdemElement{NewBasisVector(p, []int{0, bockstein}, []int{A} )}
    }
    if A >= p*B + bockstein { // Admissible
        return AdemElement{NewBasisVector(p, []int{0, bockstein, 0}, []int{A, B} )}
    }
    
    result := AdemElement{NewZeroVector(p, A/p * (1 + bockstein))}
    for j := 0; j <= A/p; j++ {
        coeff := BinomialOdd((B-j) * (p-1) - 1 + bockstein, A - p*j, p)
        coeff *= MinusOneToTheN(A + j)
        coeff = coeff % p        
        if coeff != 0 && j == 0 {
            result.Set([]int{bockstein, 0}, []int{A + B}, coeff)
        } else if coeff != 0 && j != 0 {
            result.Set([]int{bockstein, 0, 0}, []int{A + B - j, j}, coeff)
        }
    }
    if bockstein == 1 {
        for j := 0; j <= (A-1)/p; j++ {
            coeff := BinomialOdd((B-j) * (p - 1) - 1, A - p*j - 1, p)
            coeff *= MinusOneToTheN(A + j - 1)
            if coeff != 0 && j == 0 {
                result.Set([]int{0, 1}, []int{A + B}, coeff)
            } else if coeff != 0 && j != 0 {
                result.Set([]int{0, 1, 0}, []int{A + B - j, j}, coeff)
            }
        }
    }
    return result
}

//@memoized
//Reduce a monomial into a linear combination of admissible monomials when p = 2
func MakeMonoAdmissible2(mono []int) AdemElement {   
    if len(mono) == 1 {
        return AdemElement{NewBasisVector2(mono)}
    }
    if len(mono) == 2 {
        return Adem2(mono[0], mono[1])
    }
    first_inadmissible_index := -1
    for j := 0; j < len(mono) - 1; j++ {
        if mono[j] < 2*mono[j+1] {
            first_inadmissible_index = j
            break
        }
    }
    if first_inadmissible_index == -1 { // Didn't find any inadmissible indices
        return AdemElement{NewBasisVector2(mono)}
    }
    result := AdemElement{NewZeroVector(2, -1)}
    j := first_inadmissible_index
    y := Adem2(mono[j], mono[j+1])
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
        if even_part[j] < odd_part[j+1] + p * even_part[j+1] {
            first_inadmissible_index = j
            break
        }
    }
    if first_inadmissible_index == -1 { // Didn't find any inadmissible indices
        return AdemElement{NewBasisVector(p, odd_part, even_part)}
    }
    result := AdemElement{NewZeroVector(p, -1)}
    j := first_inadmissible_index
    y := AdemGeneric(p, even_part[j],odd_part[j+1], even_part[j+1])
    coeff_map := y.GetCoeffMap()
    for key, m := range y.GetBasisVectorMap() {
        coeff := coeff_map[key]
        reln_even_part := m.even_part
        reln_odd_part := m.odd_part
        
        if odd_part[j]   == 1 && reln_odd_part[0] == 1  || 
           odd_part[j+2] == 1 && reln_odd_part[len(reln_odd_part) - 1] == 1 {
            continue
        }
        
        new_even_part := make([]int, len(even_part) + len(reln_even_part) - 2)
        new_odd_part  := make([]int, len(odd_part) + len(reln_odd_part) - 3)
        copy(new_even_part, even_part[:j])
        copy(new_even_part[j:], reln_even_part)
        copy(new_even_part[j+len(reln_even_part):], even_part[j+2:])
        
        copy(new_odd_part, odd_part[:j+1])
        new_odd_part[j] += reln_odd_part[0]
        copy(new_odd_part[j+1:], reln_odd_part[1:])
        new_odd_part[j + len(reln_odd_part) - 1] += odd_part[j + 2]
        copy(new_odd_part[j + len(reln_odd_part):], odd_part[j + 3:])    
        
        new := MakeMonoAdmissibleGeneric(p, Monomial{new_odd_part, new_even_part})
        result.ScaleAndAdd(new, coeff) 
    }
    return result
}

//Multiply monomials m1 and m2 and write the result in the Adem basis for p = 2.
func Product2(m1, m2 []int) AdemElement {
    m := make([]int, len(m1) + len(m2))
    copy(m, m1)
    copy(m[len(m1):], m2)
    //fmt.Println(m)
    return MakeMonoAdmissible2(m)
}

//Multiply monomials m1 and m2 and write the result in the Adem basis in the generic case.
func ProductGeneric(p int, m1, m2 Monomial) AdemElement {
    if m1.odd_part[len(m1.odd_part) - 1] == 1 && m2.odd_part[0] == 1 {
        return AdemElement{NewZeroVector(p, 0)}
    }
    odd_part := make([]int, len(m1.odd_part) + len(m2.odd_part) - 1)
    copy(odd_part, m1.odd_part)
    odd_part[len(m1.odd_part) - 1] += m2.odd_part[0]
    copy(odd_part[len(m1.odd_part):], m2.odd_part[1:])
    
    even_part := make([]int, len(m1.even_part) + len(m2.even_part))
    copy(even_part, m1.even_part)
    copy(even_part[len(m1.even_part):], m2.even_part)
    return MakeMonoAdmissibleGeneric(p, Monomial{odd_part, even_part})
}


//Get the basis for the n dimensional part of the Steenrod algebra.
//Build basis recursively.  last = last term.
//last is >= bound, and we will append (last,) to the end of
//elements from basis (n - last, bound=2 * last).
//This means that 2 last <= n - last, or 3 last <= n.
//extra_length tells us how many extra elements are planned for insertion onto the end
//so we don't have to do any copies.
func basis2Helper(n, bound, extra_length int) <- chan []int {
    ch := make(chan []int)
//    yield := func(result []int) {
//        cpy := make([]int, len(result))
//        copy(cpy, result)
//        ch <- cpy
//    }
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
            for vec := range basis2Helper(n - last, 2 * last, extra_length + 1) {
                vec = append(vec, last)
                ch <- vec
            }
        }
    }()
    return ch
}

//Get the basis for the n dimensional part of the Steenrod algebra.
func Basis2(n int) <-chan []int {
    return basis2Helper(n, 1, 0)
}

//Get the basis for the n dimensional part of the Steenrod algebra.
func basisGenericHelper(p, n, bound, extra_length int) <-chan Monomial {
    ch := make(chan Monomial)
    yield := func(odd_result, even_result []int) {
        ch <- Monomial{odd_result, even_result}
    }
    go func() {
        defer close(ch)      
        if n == 0 {
            odd_result  := make([]int, 1, 1 + extra_length)
            even_result := make([]int, 0, extra_length)
            odd_result[0] = 0
            yield(odd_result, even_result)
        }
        if n == 1 {
            odd_result  := make([]int, 1, 1 + extra_length)
            even_result := make([]int, 0, extra_length)
            odd_result[0] = 1
            yield(odd_result, even_result)
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
                remaining_degree := n - 2*(p -1)*last - epsilon
                for vec := range basisGenericHelper(p, remaining_degree, p * last, extra_length + 1) {
                    odd_part := vec.odd_part
                    even_part := vec.even_part
                    odd_part = append(odd_part, epsilon)
                    even_part = append(even_part, last)
                    yield(odd_part, even_part)
                }
            }
        }
    }()
    return ch
}

func BasisGeneric(p, n int) <- chan Monomial {
    return basisGenericHelper(p, n, 1, 0)

}

package main

import (
    "fmt"
)

type EvenMonomial []int

type Monomial struct {
    odd_part uint64
    even_part []int
}

func (m Monomial) ToString() string{
    return fmt.Sprint(m)
}

type VectorSubtype interface {
    GetPrime() int
    GetCoeffMap() map[string]int
    GetBasisVectorMap() map[string]Monomial
}

type Vector struct {
    p int
    coeff_map map[string]int
    basis_vector_map map[string]Monomial
}

func NewZeroVector(p int, size_hint int) Vector {
    if size_hint < 0 {
        return Vector { p, make(map[string]int, size_hint), make(map[string]Monomial) }    
    }
    return Vector { p, make(map[string]int, size_hint), make(map[string]Monomial, size_hint) }
}

func NewBasisVector(p int, odd_part uint64, even_part []int) Vector{
    m := Monomial{odd_part, even_part}
    coeff_map := make(map[string]int, 1)
    basis_vector_map := make(map[string]Monomial, 1)
    key := m.ToString()
    coeff_map[key] = 1
    basis_vector_map[key] = m
    return Vector{p, coeff_map, basis_vector_map}
}

func NewBasisVector2(m []int) Vector {
    return NewBasisVector(2, 0, m)
}


func (v Vector) GetPrime() int{
    return v.p
}

func (v Vector) GetCoeffMap() map[string]int{
    return v.coeff_map
}

func (v Vector) GetBasisVectorMap() map[string]Monomial{
    return v.basis_vector_map
}

func (v Vector) Get(m Monomial) int{
    return v.coeff_map[m.ToString()]
}

func (v Vector) Set(odd_part uint64, even_part []int, c int){
    m := Monomial{odd_part, even_part}
    key := m.ToString()
    v.coeff_map[key] = ModPositive(c, v.p)
    v.basis_vector_map[key] = m
}

func (v Vector) Set2(even_part []int, c int){
    m := Monomial{0, even_part}
    key := m.ToString()
    v.coeff_map[key] = ModPositive(c, v.p)
    v.basis_vector_map[key] = m
}

func (v Vector) AddBasisVector(m Monomial, c int){
    key := m.ToString()
    v.coeff_map[key] += c
    v.coeff_map[key] = ModPositive(v.coeff_map[key], v.p)
    v.basis_vector_map[key] = m
}

func (v Vector) ScaleAndAdd(w VectorSubtype, c int) {
    v_coeff_map := v.GetCoeffMap()
    v_basis_map := v.GetBasisVectorMap()
    w_coeff_map := w.GetCoeffMap()
    w_basis_map := w.GetBasisVectorMap()
    for key, coeff := range w_coeff_map {
        v_coeff_map[key] += c * coeff
        v_basis_map[key] = w_basis_map[key]
        if v_coeff_map[key] == 0 {
            delete(v_coeff_map, key)
            delete(v_basis_map, key)
        }
    }
}

func (v Vector) Add(w VectorSubtype) {
    v.ScaleAndAdd(w, 1)
}


func (v Vector) Scale(c int) {
    v_coeff_map := v.GetCoeffMap()
    for key := range v_coeff_map {
        v_coeff_map[key] *= c
        v_coeff_map[key] = ModPositive(v_coeff_map[key], v.p)
    }
}

func Sum(vectors... VectorSubtype) Vector {
    p := vectors[0].GetPrime()
    size_estimate := 0
    for _, v := range vectors {
        size_estimate += len(v.GetCoeffMap())
    }
    result := NewZeroVector(p, size_estimate)
    for _, v := range vectors {
        result.Add(v)
    }
    return result
}



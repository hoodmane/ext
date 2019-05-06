package main

import (
    "fmt"
)

type EvenMonomial []int

type Monomial struct {
    odd_part uint64
    even_part []int
}

type MonomialIndex struct {
    degree uint 
    index uint 
}

func (m Monomial) String() string{
    return fmt.Sprintf("{%v %v}", m.odd_part, m.even_part)
}

type VectorSubtype interface {
    GetPrime() int
    GetCoeffMap() map[MonomialIndex]int
}

type Vector struct {
    p int
    coeff_map map[MonomialIndex]int
}

func NewZeroVector(p int, size_hint int) Vector {
    if size_hint < 0 {
        return Vector { p, make(map[MonomialIndex]int, size_hint) }    
    }
    return Vector { p, make(map[MonomialIndex]int, size_hint)}
}

func NewBasisVector(p int, idx MonomialIndex) Vector{
    coeff_map := make(map[MonomialIndex]int, 1)
    coeff_map[idx] = 1
    return Vector{p, coeff_map}
}

func (v Vector) GetPrime() int{
    return v.p
}

func (v Vector) GetCoeffMap() map[MonomialIndex]int{
    return v.coeff_map
}

func (v Vector) Get(m MonomialIndex) int{
    return v.coeff_map[m]
}

func (v Vector) Set(idx MonomialIndex, c int){
    v.coeff_map[idx] = ModPositive(c, v.p)
}


func (v Vector) AddBasisVector(idx MonomialIndex, c int){
    v.coeff_map[idx] += c
    v.coeff_map[idx] = ModPositive(v.coeff_map[idx], v.p)
    if v.coeff_map[idx] == 0 {
        delete(v.coeff_map, idx)
    }    
}

func (v Vector) ScaleAndAdd(w VectorSubtype, c int) {
    w_coeff_map := w.GetCoeffMap()
    for idx, coeff := range w_coeff_map {
        v.coeff_map[idx] += c * coeff
        v.coeff_map[idx] = ModPositive(v.coeff_map[idx], v.p)
        if v.coeff_map[idx] == 0 {
            delete(v.coeff_map, idx)
        }
    }
}

func (v Vector) Add(w VectorSubtype) {
    v.ScaleAndAdd(w, 1)
}


func (v Vector) Scale(c int) {
    for key := range v.coeff_map {
        v.coeff_map[key] *= c
        v.coeff_map[key] = ModPositive(v.coeff_map[key], v.p)
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



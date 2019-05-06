//  File: milnor.go
//  Author: Hood Chatham
//  
//  Defines the basic operations on Milnor algebras. No object oriented code,
//  all functions take basis vectors as inputs and output dictionaries representing
//  Fp-linear combinations of basis vectors. The relevant calls will be wrapped
//  by methods of the MilnorAlgebra and MilnorElement classes in steenrod.py.
//  They are systemtically extended over vector inputs using the
//  @linearextension decorator.
//  
//  The algorithms here are all copied from Sage, in particular from the two
//  files steenrod_algebra_mult.py and (the basis code from) steenrod_algebra_basis.py
//  I have made significant improvements in code legibility and removed all references to Sage code.

package main

import (
    "math"
    "fmt"
)

//func main(){
//    empty_profile := ProfileList{[]int{},false,false}
//    empty_full_profile := FullProfile{empty_profile, empty_profile}
//    Alg := func(p int) MilnorAlgebra {
//        return MilnorAlgebra{p,p!=2,empty_full_profile}
//    }
//    
//    for x := range MilnorBasisGenericQpart(3,(uint64(1) << uint(10))-1, 22){
//        fmt.Println("out:", x)
//    }
//
//    fmt.Println("")
//    fmt.Println("A(3).basis(1)")    
//    for x := range MilnorBasisGeneric(Alg(3), 9){
//        fmt.Println(x)
//    }
//}

type FullProfile struct {
    odd_part Profile
    even_part Profile
}

type Profile struct {
    profile []int
    truncated bool
    restricted bool
}

func (P Profile) getIndex(index int) int{
    if index < len(P.profile) {
        return P.profile[index]
    } 
    if P.truncated {
        return 0
    } 
    return math.MaxInt32
}

func (P Profile) getExponent(p, index int) int{
    if index < len(P.profile) {
        return pow(p, P.profile[index])
    } 
    if P.truncated {
        return 1
    } 
    return math.MaxInt32
}


type MilnorElement struct {
    Vector
}

type MilnorAlgebra struct {
    p int
    generic bool
    profile FullProfile
    name string
}

func (A MilnorAlgebra) String() string {
    if A.name == "" {
        name := fmt.Sprintf("{%v, %v", A.p, A.generic)
        if(!A.profile.even_part.restricted && !A.profile.odd_part.restricted){
            name += "}"
        } else {
            name += fmt.Sprintf("%v}", A.profile)
        }
        A.name = name
    }
    return A.name
}




type Q_part struct {
    bit_string uint64
    degree int
}

type MilnorBasisElement struct {
    q_degree int
    q_part uint64
    p_degree int
    p_part []int
}

func (m MilnorBasisElement) String() string{
    return fmt.Sprintf("{%v %v}", m.q_part, m.p_part)
}


var milnor_basis_table = make(map[string] [][]MilnorBasisElement)
var milnor_basis_name_to_index_map = make(map[string] [] map[string] MonomialIndex)

func MilnorBasisEven(xi_degrees []int, profile_list []int, n int) <-chan []int {
    ch := make(chan []int, 20)
    go func(){
        defer close(ch)
        for exponents := range WeightedIntegerVectors(n, xi_degrees, profile_list) {
            exponents = remove_trailing_zeroes(exponents)
            ch <- exponents
        }
    }()
    return ch
}

var milnor_basis_Q_table = make(map[int] [][]Q_part)
var milnor_basis_Q_table_size = make(map[int] int)
var tau_degrees = make(map[int] []int)
var xi_degrees  = make(map[int] []int)

func generateMilnorBasisQpartTable(p int, n int){
    q := 2*(p-1)
    table, ok := milnor_basis_Q_table[p]
    prev_dim := milnor_basis_Q_table_size[p]
    bit_string_min := uint64(1) << uint(prev_dim)    
    if !ok {
        milnor_basis_Q_table[p] = make([][]Q_part, q) 
        table = milnor_basis_Q_table[p]
        for residue := 0; residue < q; residue ++ {
            table[residue] = make([]Q_part, 0, 10)
        }  
        table[0] = append(table[0], Q_part{0, 0})
    }
    milnor_basis_Q_table_size[p] = len(tau_degrees[p])    
    bit_string_max := uint64(1) 
    bit_string_max <<= uint(len(tau_degrees[p]))
    total := 0
    //The total starts out as tau_degrees[prev_dim] but we update by tau_degrees[prev_dim] - Sum(smaller tau_i's).
    //So initialize total = Sum(smaller tau_i's)
    for i := 0; i < prev_dim; i++ {
        total += tau_degrees[p][i]
    }
    //The residue starts out as 1, but we update by 1 - # of trailing 0's.     
    //On the first pass, # of trailing 0's is prev_dim, so initialize residue = prev_dim 
    residue := prev_dim 
    for bit_string := bit_string_min;  bit_string < bit_string_max; bit_string++ {
         // Iterate over trailing zero bits
         v := (bit_string ^ (bit_string - 1)) >> 1
         c := 0 // c is counting the number of zero bits at the end of bit_string.
         for ; v !=0 ; c++ {
             v >>= 1;
             total -= tau_degrees[p][c]
         }
         total += tau_degrees[p][c]
         // residue is the number of 1's in bit_string mod q. 
         // c bits were unset when we incremented and one new bit was set so we have to add 1 - c.
         residue += 1 - c
         residue = ModPositive(residue, q)
         table[residue] = append(table[residue], Q_part{bit_string, total})
     }
}

// Returns the "Q-part" of the basis in degree q_deg.
// This means return the set of monomials Q(i_1) * ... * Q(i_k) where i_1 < ... < i_k
// and the product is in q_deg. Basically it's just an issue of finding partitions of
// q_deg into parts of size |Q(i_j)|, and then there's a profile condition.
func MilnorBasisGenericQpart(p int, profile uint64, n int) <-chan Q_part {
    ch := make(chan Q_part, 20)
    go func(){
        defer close(ch)
        if NumTausLeqN(n, p) > milnor_basis_Q_table_size[p] {
            generateMilnorBasisQpartTable(p, n)
        }
        tau_monomial_list := milnor_basis_Q_table[p][n % (2*p - 2)]
        for _, tau_mono := range tau_monomial_list {
            if tau_mono.degree > n {
                return
            }
            if tau_mono.bit_string & (^profile) == 0 {
               ch <- tau_mono
            }
        }
    }()
    return ch
}


func GenerateMilnorBasis(algebra MilnorAlgebra, n int) {
    algebra_name := algebra.String()
    p := algebra.p
    tau_degrees[p] = TauDegrees(n, p)
    xi_degrees[p] = XiDegrees(n, p)  
    
    table, ok := milnor_basis_table[algebra_name]
    old_n := len(table) - 1
    
    if !ok {
        milnor_basis_table[algebra_name] = make([][]MilnorBasisElement, n + 1)
        milnor_basis_name_to_index_map[algebra_name] = make([]map[string] MonomialIndex, n + 1)
        table = milnor_basis_table[algebra_name]
    }
    name_table := milnor_basis_name_to_index_map[algebra_name]
    
    if cap(table) < n {
        new_table := make([][]MilnorBasisElement, n + 1)
        copy(new_table, table)
        milnor_basis_table[algebra_name] = new_table
        table = new_table
        
        new_name_table := make([]map[string] MonomialIndex, n + 1)
        copy(new_name_table, name_table)
        milnor_basis_name_to_index_map[algebra_name] = new_name_table
        name_table = new_name_table
    }
    
    for k := old_n + 1; k <= n; k++ {
        table[k] = make([]MilnorBasisElement, 0, 10)
        name_table[k] = make(map[string] MonomialIndex)
    }
    
    if algebra.generic {
        GenerateMilnorBasisGeneric(algebra, old_n, n)
    } else {
        GenerateMilnorBasis2(algebra, old_n, n)
    }
}

func GenerateMilnorBasis2(algebra MilnorAlgebra, old_n, n int){
    algebra_name := algebra.String()
    p := algebra.p
    profile := algebra.profile.even_part
    
    table := milnor_basis_table[algebra_name]
    name_table := milnor_basis_name_to_index_map[algebra_name]
    
    profile_list := make([]int, len(xi_degrees[p]))
    for idx := range xi_degrees[p] {
        profile_list[idx] = profile.getExponent(p, idx) - 1
    }
    for k := old_n + 1; k <= n; k++ {
        for x := range MilnorBasisEven(xi_degrees[p], profile_list, k) {
            m := MilnorBasisElement{0, 0, k, x}
            table[k] = append(table[k], m)
            name_table[k][m.String()] = MonomialIndex{uint(k), uint(len(table[k]) - 1)}
        }
    }
}

func GenerateMilnorBasisGeneric(algebra MilnorAlgebra, old_n, n int){
    algebra_name := algebra.String()
    table := milnor_basis_table[algebra_name]
    name_table := milnor_basis_name_to_index_map[algebra_name]    
    
    for k := old_n + 1; k <= n; k++ {
        for x := range MilnorBasisGeneric(algebra, k) {
            table[k] = append(table[k], x)
            name_table[k][x.String()] = MonomialIndex{uint(k), uint(len(table[k]) - 1)}
        }
    }
}

// Return the even part of the basis in degree n * 2*(p-1).
// In the nongeneric case, this actually just gets the whole degree n basis.
// Note the factor of two difference between 2*(2-1) and 1.
func MilnorBasis2(algebra MilnorAlgebra, n int) <-chan []int {
    profile := algebra.profile.even_part
    p := algebra.p
    profile_list := make([]int, len(xi_degrees[p]))
    for idx := range xi_degrees[p] {
        profile_list[idx] = profile.getExponent(p, idx) - 1
    }
    return MilnorBasisEven(xi_degrees[p], profile_list, n)
}

// Get the basis in degree n for the generic steenrod algebra at the prime p.
// We just put together the "even part" of the basis and the "Q part".
func MilnorBasisGeneric(algebra MilnorAlgebra, n int) <-chan MilnorBasisElement {
    ch := make(chan MilnorBasisElement, 20)
    p := algebra.p
    num_taus := NumTausLeqN(n, p)
    full_profile := algebra.profile
    odd_profile := uint64(0)
    for idx := 0; idx < num_taus; idx ++ {
        if full_profile.odd_part.getIndex(idx) > 0 {
            odd_profile += 1 << uint(idx)
        }
    }    
    even_profile_list := make([]int, len(xi_degrees[p]))
    for idx := range xi_degrees[p] {
        even_profile_list[idx] = full_profile.even_part.getExponent(p, idx) - 1
    }    
    go func(){
        defer close(ch)
        if n == 0 {
            ch <- MilnorBasisElement{0, 0, 0, []int {}}
            return
        }
        // p_deg records the desired degree of the P part of the basis element.
        // Since p-parts are always divisible by 2p-2, we divide by this first.
        // pow(p, -1) returns 1, so min_q_deg is 0 if q divides n evenly.
        for qs_q_deg := range MilnorBasisGenericQpart(p, odd_profile, n){
            q_part := qs_q_deg.bit_string
            q_deg := qs_q_deg.degree
            p_deg := (n - q_deg) / (2*(p-1))
            P_parts := MilnorBasisEven(xi_degrees[p], even_profile_list, p_deg)
            for p_part := range P_parts {
                ch <- MilnorBasisElement{q_deg, q_part, 2*(p-1)*p_deg, p_part}
            }
        }
    }()
    return ch
}

func GetMilnorBasisElementFromIndex(algebra MilnorAlgebra, idx MonomialIndex) MilnorBasisElement{
    return milnor_basis_table[algebra.String()][idx.degree][idx.index]
}

func GetIndexFromMilnorBasisElement(algebra MilnorAlgebra, b MilnorBasisElement) MonomialIndex{
    return milnor_basis_name_to_index_map[algebra.String()][b.q_degree + b.p_degree][b.String()]
}

func NewMilnorZeroVector(p, size_hint int) MilnorElement {
    return MilnorElement{NewZeroVector(p, size_hint)}
}

func NewMilnorBasisVector2(algebra MilnorAlgebra, deg int,  even_part []int) MilnorElement{
    return MilnorElement{NewBasisVector(2, 
        GetIndexFromMilnorBasisElement(algebra, MilnorBasisElement{0, 0, deg, even_part})),
    }
}

func NewMilnorBasisVectorGeneric(algebra MilnorAlgebra, odd_deg int, odd_part uint64, even_deg int, even_part []int) MilnorElement {
    return MilnorElement{NewBasisVector(algebra.p, 
        GetIndexFromMilnorBasisElement(algebra, MilnorBasisElement{odd_deg, odd_part, even_deg, even_part})),
    }
}


func allocate_milnor_matrix(rows, cols int) [][]int {
    M := make([][]int, rows)
    e := make([]int, rows * cols)
    for i := range M {
        M[i] = e[i * cols:(i + 1) * cols]
    }
    return M    
}

// Initializes an len(r)+1 by len(s)+1 matrix
// Puts r along the first column and s along the first row and zeroes everywhere else.
func initialize_milnor_matrix(r, s []int) [][]int {
    rows := len(r) + 1
    cols := len(s) + 1
    M := allocate_milnor_matrix(rows, cols)
    copy(M[0][1:], s)
    for i, v := range r {
        M[i+1][0] = v
    }
    return M
}

func copy_milnor_matrix_starting_in_row(target, source [][]int, row, cols int) {
    copy(target[0][row*cols:cap(target[0])], source[0][row*cols:cap(source[0])])
}


// This seems to move an i x j block of M back to the first row and column.
// To be honest, I don't really know what the point is, but the milnor_matrices
// function was a little long and this seemed like a decent chunk to extract.
// At least it contains all of the steps that modify M so that seems like a good thing.
func step_milnor_matrix(M [][]int, r, s []int, i, j, x int) [][]int {
    rows := len(r) + 1
    cols := len(s) + 1
    N := allocate_milnor_matrix(rows, cols)
    copy(N[0],M[0])
    for row := 1; row < i; row ++ {
        N[row][0] = r[row-1]
        for col := 1; col < cols; col++ {
            N[0][col] += M[row][col]
        }
    }
    copy_milnor_matrix_starting_in_row(N, M, i, cols)
    for col := 1; col < j; col++ {
        N[0][col] += M[i][col]
        N[i][col] = 0
    }
    N[0][j] --
    N[i][j] ++
    N[i][0] = x
    return N
}

// Generator for Milnor matrices. milnor_product_even iterates over this.
// Uses the same algorithm Monks does in his Maple package to iterate through
// the possible matrices: see
// https://monks.scranton.edu/files/software/Steenrod/steen.html
func milnor_matrices(p int, r, s []int) <-chan [][]int {
    ch := make(chan [][]int)
    go func(){
        defer close(ch)
        rows := len(r) + 1
        cols := len(s) + 1
        M := initialize_milnor_matrix(r, s)
        ch <- M
        for found := true; found ; {
            found = false
            for i := 1; !found && i < rows; i++ {
                total := M[i][0]
                for j := 1; j < cols; j++ {
                    column_above_is_empty := true
                    for k := 0; k < i; k++ {
                        if M[k][j] != 0 {
                            column_above_is_empty = false
                            break
                        }
                    }
                    p_to_the_j := pow(p, j)                    
                    if total < p_to_the_j || column_above_is_empty {
                        total += M[i][j] * p_to_the_j
                    } else {
                        M = step_milnor_matrix(M, r, s, i, j, total - p_to_the_j)
                        found = true
                        ch <- M
                        break 
                    }
                }
            }
        }
    }()
    return ch
}

//Remove trailing zeroes from the list l.
func remove_trailing_zeroes(l []int) []int {
    for i := len(l) - 1; i >= 0; i-- {
        if l[i] != 0 {
            return l[:i+1]
        }
    }
    return l[:0]
}

func max(a, b int) int{
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int{
    if a < b {
        return a
    }
    return b
}

// Handles the multiplication in the even subalgebra of the Steenrod algebra P.
// When p = 2, this is isomorphic to the whole Steenrod algebra so this method does everything.
func MilnorProductEven(algebra MilnorAlgebra, r_elt, s_elt MilnorBasisElement) MilnorElement {
    p := algebra.p
    result := NewMilnorZeroVector(p, -1)
    r := r_elt.p_part
    s := s_elt.p_part
    output_degree := (r_elt.p_degree + s_elt.p_degree)
    rows := len(r) + 1
    cols := len(s) + 1
    diags := len(r) + len(s)
    for M := range milnor_matrices(p, r, s) {
        // check diagonals
        coeff := 1
        diagonal_sums := make([]int, diags)
        for n := 1; n <= diags; n++ {
            i_min := max(0, n - cols + 1)
            i_max := min(1 + n, rows)
            nth_diagonal := make([]int, i_max - i_min + 1)
            nth_diagonal_sum := 0
            index := 0
            for i := i_min; i < i_max; i++ {    
                nth_diagonal[index] = M[i][n-i]
                nth_diagonal_sum += nth_diagonal[index]
                index++
            }
            coeff *= Multinomial(nth_diagonal, p)
            coeff = coeff % p
            if coeff == 0 {
                break
            }
            diagonal_sums[n-1] = nth_diagonal_sum
        }
        if coeff != 0 {
            diagonal_sums = remove_trailing_zeroes(diagonal_sums)
            m := MilnorBasisElement{0, 0, output_degree, diagonal_sums}
            idx := GetIndexFromMilnorBasisElement(algebra, m)
            result.AddBasisVector(idx, coeff)
        }
    }
    return result
}

func inListQ(l []int, n int) bool{
    for _, v := range l {
        if v == n {
            return true
        }
    }
    return false
}


// Reduce m1 * f = (Q_e0 Q_e1 ... P(r1, r2, ...)) * (Q_f0 Q_f1 ...) into the form Sum of Q's * P's
// Result is represented as dictionary of pairs of tuples.
func MilnorProductFullQpart(algebra MilnorAlgebra, m1 MilnorBasisElement, f uint64) MilnorElement{
    p := algebra.p
    var q int
    if algebra.generic {
        q = 2*p-2
    } else {
        q = 1
    }
    p_degree := m1.p_degree
    q_degree := m1.q_degree
    result := NewMilnorBasisVectorGeneric(algebra, q_degree, m1.q_part, p_degree, m1.p_part)
    for k := 0; f & ^((1 << uint(k)) - 1) != 0; k++ {
        if f & (1 << uint(k)) == 0 {
            continue
        }
        q_degree += tau_degrees[p][k]
        old_result := result
        result = NewMilnorZeroVector(p, -1)
        p_to_the_k := pow(p, k)
        for idx, coeff := range old_result.GetCoeffMap() {
            mono := GetMilnorBasisElementFromIndex(algebra, idx)
            for i := 0; i < len(mono.p_part) + 1; i++ {
                if mono.q_part & (1 << uint(k+i)) != 0 {
                    continue
                }
                // Make sure mono.p_part[i - 1] is large enough to deduct p^k from it
                if i > 0 && mono.p_part[i - 1] < p_to_the_k {
                    continue 
                }

                q_mono := mono.q_part
                q_degree := mono.q_degree
                p_mono := mono.p_part
                p_degree = mono.p_degree
                if i > 0 {
                    new_p_mono := make([]int, len(p_mono))
                    copy(new_p_mono, p_mono)
                    new_p_mono[i - 1] -= p_to_the_k
                    p_mono = remove_trailing_zeroes(new_p_mono) 
                    p_degree -= q * p_to_the_k * xi_degrees[p][i-1]
                }

                // insert(q_mono, len(q_mono) - qs_gt_k_plus_i, k+i)
                q_mono += 1 << uint(k+i)
                q_degree += tau_degrees[p][k+i]
                
                larger_Qs := 0
                v := q_mono >> uint(k + i + 1)
                for  ; v != 0; v >>= 1 {
                    larger_Qs += int(v & 1);
                }
                coeff *= MinusOneToTheN(larger_Qs)
                
                out_idx := GetIndexFromMilnorBasisElement(algebra, MilnorBasisElement{q_degree, q_mono, p_degree, p_mono})
                result.AddBasisVector(out_idx, coeff)
            }
        }
    }
    return result
}    

// Product of Milnor basis elements defined by m1 and m2 at the prime p.
// 
// INPUT:
// 
// - m1 - pair of tuples (e,r), where e is an increasing tuple of
//   non-negative integers and r is a tuple of non-negative integers
// - m2 - pair of tuples (f,s), same format as m1
// - p -- odd prime number
// 
// OUTPUT:
// 
// Dictionary of terms of the form (tuple: coeff), where 'tuple' is
// a pair of tuples, as for r and s, and 'coeff' is an integer mod p.
// 
// This computes the product of the Milnor basis elements
// $Q_{e_1} Q_{e_2} ... P(r_1, r_2, ...)$ and
// $Q_{f_1} Q_{f_2} ... P(s_1, s_2, ...)$.
func MilnorProductFull(algebra MilnorAlgebra, m1, m2 MilnorBasisElement) MilnorElement {
    p := algebra.p
    s := MilnorBasisElement{0, 0, m2.p_degree, m2.p_part}
    m1_times_f := MilnorProductFullQpart(algebra, m1, m2.q_part)
    // Now for the Milnor matrices.  For each entry '(e,r): coeff' in answer,
    // multiply r with s.  Record coefficient for matrix and multiply by coeff.
    // Store in 'result'.
    if len(m2.p_part) == 0 {
        return m1_times_f
    }
    
    result := NewMilnorZeroVector(p, -1)
    for er_idx, coeff := range m1_times_f.GetCoeffMap() {
        er_mono := GetMilnorBasisElementFromIndex(algebra, er_idx)
        r := MilnorBasisElement{0, 0, er_mono.p_degree, er_mono.p_part}
        prod := MilnorProductEven(algebra, r, s)
        for idx, c := range prod.GetCoeffMap() {
            m := GetMilnorBasisElementFromIndex(algebra, idx)
            m  = MilnorBasisElement{er_mono.q_degree, er_mono.q_part, m.p_degree, m.p_part}
            out_idx := GetIndexFromMilnorBasisElement(algebra, m)
            result.AddBasisVector(out_idx, coeff*c)
        }
    }
    return result
}

// Multiplication of Milnor basis elements in the non generic case.
func MilnorProduct2(algebra MilnorAlgebra, r, s MilnorBasisElement) MilnorElement {
    return MilnorProductEven(algebra, r, s)
}

func MilnorProductGeneric(algebra MilnorAlgebra, r, s MilnorBasisElement) MilnorElement {
    return MilnorProductFull(algebra, r, s)
}


// Multiply r and s in the Milnor algebra determined by algebra.
// Note that since profile functions determine subalgebras, the product
// doesn't need to care about the profile function at all.
func MilnorProduct(algebra MilnorAlgebra, r, s MilnorBasisElement) MilnorElement {
    if algebra.generic {
        return MilnorProductFull(algebra, r, s)
    } else {
        return MilnorProductEven(algebra, r, s)
    }
}

package main

//import "fmt"

func MinusOneToTheN(n int) int {
    return -(n % 2 * 2 - 1) 
}

func ModPositive(n, p int) int {
    return ((n % p) + p) % p
}


func lengthOfBasepExpansion(n, p int) int {
    i := 0
    for ; n > 0; n /= p {
        i++    
    }
    return i
}

func basepExpansion(n, p, padlength int) []int {
    if padlength == 0 {
        padlength = lengthOfBasepExpansion(n, p)
    }
    result := make([]int, padlength)
    i := 0
    for ; n > 0; n /= p {
        result[i] = n % p
        i++
    }
    return result    
}

var binomial_table map[int] [][]int = make(map[int] [][]int)

func directBinomialInitializeTable(p int){
    table_p := make([][]int, p)
    for i := 0; i < p; i++ {
        table_p[i] = make([]int, p)
    }
    for n := 0; n < p; n ++ {
        entry := 1
        table_p[n][0] = entry
        for k := 1; k <= n; k++ {
            entry *= (n + 1 - k)
            entry /= k
            table_p[n][k] = entry % p            
        }
    }
    binomial_table[p] = table_p
}

// "direct" binomial coefficient.
func directBinomial(n, k, p int) int {
    table, table_exists := binomial_table[p]
    if !table_exists {
        directBinomialInitializeTable(p)
        table = binomial_table[p]
    }
    return table[n][k]
}

//Mod 2 multinomial coefficient of the list l
func Multinomial2(l []int) int {
    bit_or := 0
    sum := 0
    for _, v := range l {
        sum += v
        bit_or |= v
        if bit_or < sum {
            return 0
        }
    }
    return 1
}

//Mod 2 binomial coefficient n choose k
func Binomial2(n, k int) int {
    if n < k || k < 0 {
        return 0
    } else {
        if (n-k) & k == 0 {
            return 1
        } else {
            return 0
        }
    }
}

//Mod p multinomial coefficient of l. If p is 2, more efficient to use Multinomial2.
func MultinomialOdd(l []int, p int) int {
    n := 0 
    for _, v := range l {
        n += v
    }
    answer := 1
    n_expansion := basepExpansion(n, p, 0)
    l_expansions := make([][]int, len(l))
    for i, x := range l {
        l_expansions[i] = basepExpansion(x, p, len(n_expansion) )
    }
    for index, _ := range n_expansion {
        multi := 1
        partial_sum := 0
        for _, expansion := range l_expansions {
            partial_sum += expansion[index]
            if partial_sum > n_expansion[index] {
                return 0
            }
            multi *= directBinomial(partial_sum, expansion[index], p)
            multi = multi % p
        }
        answer = (answer * multi) % p
    }
    return answer
}

//Mod p binomial coefficient n choose k. If p is 2, more efficient to use Binomial2.
func BinomialOdd(n, k, p int) int {
    if n < k || k < 0 {
        return 0
    }
    return MultinomialOdd([]int{ n-k, k }, p)
}

//Dispatch to Multinomial2 or MultinomialOdd
func Multinomial(l []int, p int) int {
    if p == 2 {
        return Multinomial2(l)
    } else {
        return MultinomialOdd(l, p)
    }
}

//Dispatch to Binomial2 or BinomialOdd
func Binomial(n, k, p int) int {
    if p == 2 {
        return Binomial2(n, k)
    } else {
        return BinomialOdd(n, k, p)
    }
}

//Reverse a list of integers in place
func ReverseInPlace(l []int) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

//Get the degrees of the xi_i's that occur in dimension at most n
func XiDegrees(n, p int, reverse_output bool) []int {
    if n <= 0 {
        return []int {}
    }
    // First determine length of output list
    N := n*(p-1) + 1
    xi_max := 0
    for ; N > 0; {
        N /= p
        xi_max += 1
    }
    //Now produce it.
    result := make([]int, xi_max - 1)
    entry := 0
    p_to_the_d := 1
    for d := 0; d < xi_max - 1; d++ {
        entry += p_to_the_d
        p_to_the_d *= p    
        result[d] = entry
    }
    if reverse_output {
        ReverseInPlace(result)
    }
    return result
}

//Helper function for WeightedIntegerVectorsGeneral
//If b == 0 return a, otherwise return min(a, b). Second output is the boolean first_output == a
func min_if_b_not_zero_else_a(a, b int) (int, bool) {
    if a <= b || b == 0 {
        return a, true
    }
    return b, false
}

//Compute integral weights so that the dot product l . weights == n. 
//If max_weight is not zero, each weight is limited to be at most max_weight.
func WeightedIntegerVectorsGeneral(n int, l []int, max_weight int) <-chan []int {
    /*
    Iterate over all ``l`` weighted integer vectors with total weight ``n``.

    INPUT:

    - ``n`` -- an integer
    - ``l`` -- the weights in weakly decreasing order

    EXAMPLES::

        sage: from sage.combinat.integer_vector_weighted import iterator_fast
        sage: list(iterator_fast(3, [2,1,1]))
        [[1, 1, 0], [1, 0, 1], [0, 3, 0], [0, 2, 1], [0, 1, 2], [0, 0, 3]]
        sage: list(iterator_fast(2, [2]))
        [[1]]

    Test that :trac:`20491` is fixed::

        sage: type(list(iterator_fast(2, [2]))[0][0])
        <type 'sage.rings.integer.Integer'&gt
    */
    ch := make(chan []int)
    yield := func(result []int) {
        cpy := make([]int, len(result))
        copy(cpy, result)
        ch <- cpy
    }
    go func() {
        defer close(ch)
        if n < 0 {
            return
        }
    
        if len(l) == 0 {
            if n == 0 {
                ch <- []int {}
            }
            return
        }
        
        cur := make([]int, len(l))
        if len(l) == 1 {
            if n % l[0] == 0 {
                ratio := n / l[0]
                _, ratio_leq_max_weight := min_if_b_not_zero_else_a(ratio, max_weight)
                if ratio_leq_max_weight {
                    cur[0] = ratio
                    yield(cur)
                }
            }
            return
        }
        
        k := 0
        cur[0], _ = min_if_b_not_zero_else_a(n / l[k], max_weight)
        cur[0]++ 
        rem := n - cur[0] * l[k] // Amount remaining
        for ; k >= 0 ; {
            cur[k] -= 1
            rem += l[k]
            switch {
                case rem == 0: {
                    yield(cur)
                }
                case cur[k] < 0 || rem < 0: {
                    rem += cur[k] * l[k]
                    cur[k] = 0
                    k --
                }
                case len(l) == k + 2: {
                    if rem % l[k + 1] == 0 {
                        ratio := rem / l[k + 1]
                        _, ratio_leq_max_weight := min_if_b_not_zero_else_a(ratio, max_weight)
                        if ratio_leq_max_weight {
                            cur[k + 1] = ratio
                            yield(cur)
                            cur[k + 1] = 0
                        }
                    }
                }
                default: {
                    k ++
                    cur[k], _ = min_if_b_not_zero_else_a(rem / l[k], max_weight)
                    cur[k] ++ 
                    rem -= cur[k] * l[k]
                }
            }
        }
    }()
    return ch
}

//The no max_weight case of WeightedIntegerVectorsGeneral
func WeightedIntegerVectors(n int, l []int) <- chan []int {
    return WeightedIntegerVectorsGeneral(n, l, 0)
}

//The max_weight == 1 case of WeightedIntegerVectorsGeneral. 
//Returns the list of parts in the partition rather than the weight list.
func RestrictedPartitions(n int, l []int) <- chan []int {
    /*
            restricted_partitions(10, [6,4,2])
            [[6, 4], [6, 2, 2], [4, 4, 2], [4, 2, 2, 2], [2, 2, 2, 2, 2]]
            restricted_partitions(10, [6,4,2,2,2])
            [[6, 4], [6, 2, 2], [4, 4, 2], [4, 2, 2, 2], [2, 2, 2, 2, 2]]
            restricted_partitions(10, [6,4,4,4,2,2,2,2,2,2])
            [[6, 4], [6, 2, 2], [4, 4, 2], [4, 2, 2, 2], [2, 2, 2, 2, 2]]
    */
    yield := make(chan []int)
    go func() {
        defer close(yield)
        for weights := range WeightedIntegerVectorsGeneral(n, l, 1) {
            result := make([]int, 0, len(l))
            for idx, i := range weights {
                if i == 1 {
                    result = append(result, l[idx])
                }
            }
            yield <- result
        }
    }()
    return yield
}


//func main() {
//    for partition := range RestrictedPartitions(10, []int{6,4,2,2}){        
//        fmt.Println("result: ", partition)
//    }   
//}

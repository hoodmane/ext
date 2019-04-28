package main

//import "fmt"

func MinusOneToTheN(n int) int {
    return -(n % 2 * 2 - 1) 
}

func pow(a, b int) int {
    res := 1
    for ; b > 0 ; {
        if b & 1 != 0 {
            res *= a
        }
        a *= a
        b >>= 1;
    }
    return res;
}

func logp(n, p int) int {
    result := 0
    for ; n > 0; {
        n /= p
        result += 1
    }
    return result
}

func p_to_the_n_minus_1_over_p_minus_1(p, n int) int {
    return (pow(p, n) - 1) / (p - 1);
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


//Get the degrees of the xi_i's that occur in dimension at most n
func XiDegrees(n, p int) []int {
    if n <= 0 {
        return []int {}
    }
    // Determine length of output list
    xi_max := logp(n*(p-1) + 1, p)

    result := make([]int, xi_max - 1)
    entry := 0
    p_to_the_d := 1
    for i := 0; i < xi_max - 1; i++ {
        entry += p_to_the_d
        p_to_the_d *= p    
        result[i] = entry
    }
    return result
}

func NumTausLeqN(n, p int) int {
    return logp((n + 1) / 2, p)
}

func TauDegrees(n, p int) []int {
    // Determine length of output list
    tau_max := NumTausLeqN(n, p)
    result := make([]int, tau_max)
    p_to_the_d := 1
    for i := 0; i < tau_max; i++ {
        result[i] = 2 * p_to_the_d - 1
        p_to_the_d *= p    
    }
    return result
}

func threshold_max_weight(a int, index int, max_weights []int) (int, bool) {
    if a <= max_weights[index] {
        return a, true
    } else {
        return max_weights[index], false
    }
}

//Compute integral weights so that the dot product l . weights == n. 
//If max_weight is not zero, each weight is limited to be at most max_weight.
func WeightedIntegerVectors(n int, l []int, max_weights []int) <-chan []int {
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
        rem := n
        k := len(l) - 1 
        for  ;  k >= 0 && l[k] > n; k-- {}
        
        if k == -1 {
            if n == 0 {
                ch <- cur 
            }
            return
        }
        if k == 0 {
            if rem % l[0] == 0 {
                ratio, ratio_leq_max_weight := threshold_max_weight(rem / l[0], 0, max_weights)
                if ratio_leq_max_weight {
                    cur[0] = ratio
                    ch <- cur 
                }
            }
            return
        }
        
        cur[k], _ = threshold_max_weight(rem / l[k], k, max_weights)
        cur[k]++ 
        rem -= cur[k] * l[k]
        
        for ; k < len(cur); {
            cur[k] -= 1
            rem += l[k]
            switch {
                case rem == 0: {
                    result := make([]int, len(cur))
                    copy(result, cur)
                    ch <- result
                }
                case cur[k] < 0 || rem < 0: {
                    rem += cur[k] * l[k]
                    cur[k] = 0
                    k ++
                }
                case k == 1: {
                    if rem % l[0] == 0 {
                        ratio, ratio_leq_max_weight := threshold_max_weight(rem / l[0], 0, max_weights)
                        if ratio_leq_max_weight {
                            cur[0] = ratio
                            result := make([]int, len(cur))
                            copy(result, cur)                            
                            ch <- result
                            cur[0] = 0
                        }
                    }
                }
                default: {
                    k --
                    cur[k], _ = threshold_max_weight(rem / l[k], k, max_weights)
                    cur[k] ++ 
                    rem -= cur[k] * l[k]
                }
            }
        }
    }()
    return ch
}


//func main() {
//    for partition := range RestrictedPartitions(10, []int{6,4,2,2}){        
//        fmt.Println("result: ", partition)
//    }   
//}

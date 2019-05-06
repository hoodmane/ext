package main

import (
    "fmt"
)

func eqMatrixQ(a, b[][]int) bool {
    if len(a) != len(b){
        return false
    }
    for i := 0; i < len(a); i++ {
        if !eqListsQ(a[i], b[i]){
            return false
        }
    }
    return true
}


func eqListsQ(a, b []int) bool {
    // If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func eqStrToIntMapsQ(a, b map [string]int) bool {
    // If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}


func checkEqStrToIntMaps(call_str string, expected_output, output map [string]int) string{
    if !eqStrToIntMapsQ(expected_output, output) {
        return fmt.Sprintf("Ran %v expected %v, got %v", call_str, expected_output, output)
    }
    return ""
}

func checkEqLengths(call_str string, expected_length, output_length int) string{
    if expected_length != output_length {
        return fmt.Sprintf("Ran %v, expected output of length %v, actual length was %v", call_str, expected_length, output_length)
    }
    return ""
}

func checkGeneratorOfListsOutput(call_str string, expected_output [][]int, gen <- chan []int) string{
    output_strings := make(map[string]int)
    for _, tuple := range expected_output {
        output_strings[fmt.Sprint(tuple)] ++
    }
    for tuple := range gen {
        str_tuple := fmt.Sprint(tuple)
        if output_strings[str_tuple] == 0 {
            return fmt.Sprintf("Unexpected result %v appeared in generator %v. Was expecting the multiset %v.", tuple, call_str, expected_output)
        }
        output_strings[str_tuple] --
    }    
    for k, v := range output_strings {
        if v > 0 {
            return fmt.Sprintf("Expected result %v in generator %v but it didn't appear. Was expecting the multiset %v.", k, call_str, expected_output)
        }
    }
    return ""
}

func checkGeneratorOfMonomialsOutput(call_str string, expected_output []MilnorBasisElement, gen <- chan MilnorBasisElement) string{
    output_strings := make(map[string]int)
    for _, tuple := range expected_output {
        output_strings[fmt.Sprint(tuple)] ++
    }
    for tuple := range gen {
        str_tuple := fmt.Sprint(tuple)
        if output_strings[str_tuple] == 0 {
            return fmt.Sprintf("Unexpected result %v appeared in generator %v. Was expecting the multiset %v.", tuple, call_str, expected_output)
        }
        output_strings[str_tuple] --
    }    
    for k, v := range output_strings {
        if v > 0 {
            return fmt.Sprintf("Expected result %v in generator %v but it didn't appear. Was expecting the multiset %v.", k, call_str, expected_output)
        }
    }
    return ""
}


func checkGeneratorOfListsLength(call_str string, expected_length int, gen <- chan []int) string {
    length := 0
    for range gen {
        length ++
    } 
    if expected_length != length {
        return fmt.Sprintf("Expected generator to return %v elements returned %v instead.", expected_length, length)
    }
    return ""
}

func checkGeneratorOfMonomialsLength(call_str string, expected_length int, gen <- chan MilnorBasisElement) string {
    length := 0
    for range gen {
        length ++
    } 
    if expected_length != length {
        return fmt.Sprintf("Expected generator to return %v elements returned %v instead.", expected_length, length)
    }
    return ""
}

package main

import (
    "fmt"
    "testing"
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


func checkEqStrToIntMaps(t *testing.T, call_str string, expected_output, output map [string]int){
    if !eqStrToIntMapsQ(expected_output, output) {
        t.Errorf("Ran %v expected %v, got %v", call_str, expected_output, output)
    }
}

func checkEqLengths(t *testing.T, call_str string, expected_length, output_length int){
    if expected_length != output_length {
        t.Errorf("Ran %v, expected output of length %v, actual length was %v", call_str, expected_length, output_length)
    }
}

func checkGeneratorOfListsOutput(t *testing.T, call_str string, expected_output [][]int, gen <- chan []int){
    output_strings := make(map[string]int)
    for _, tuple := range expected_output {
        output_strings[fmt.Sprint(tuple)] ++
    }
    for tuple := range gen {
        str_tuple := fmt.Sprint(tuple)
        if output_strings[str_tuple] == 0 {
            t.Errorf("Unexpected result %v appeared in generator %v. Was expecting the multiset %v.", tuple, call_str, expected_output)
        }
        output_strings[str_tuple] --
    }    
    for k, v := range output_strings {
        if v > 0 {
            t.Errorf("Expected result %v in generator %v but it didn't appear. Was expecting the multiset %v.", k, call_str, expected_output)
        }
    }
}

func checkGeneratorOfMonomialsOutput(t *testing.T, call_str string, expected_output []Monomial, gen <- chan Monomial){
    output_strings := make(map[string]int)
    for _, tuple := range expected_output {
        output_strings[fmt.Sprint(tuple)] ++
    }
    for tuple := range gen {
        str_tuple := fmt.Sprint(tuple)
        if output_strings[str_tuple] == 0 {
            t.Errorf("Unexpected result %v appeared in generator %v. Was expecting the multiset %v.", tuple, call_str, expected_output)
        }
        output_strings[str_tuple] --
    }    
    for k, v := range output_strings {
        if v > 0 {
            t.Errorf("Expected result %v in generator %v but it didn't appear. Was expecting the multiset %v.", k, call_str, expected_output)
        }
    }
}

func checkGeneratorOfListsLength(t *testing.T, call_str string, expected_length int, gen <- chan []int) {
    length := 0
    for range gen {
        length ++
    } 
    if expected_length != length {
        t.Errorf("Expected generator to return %v elements returned %v instead.", expected_length, length)
    }
}

func checkGeneratorOfMonomialsLength(t *testing.T, call_str string, expected_length int, gen <- chan Monomial) {
    length := 0
    for range gen {
        length ++
    } 
    if expected_length != length {
        t.Errorf("Expected generator to return %v elements returned %v instead.", expected_length, length)
    }
}



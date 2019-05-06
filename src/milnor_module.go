package main

import (
    "os"
    "encoding/json"
    "fmt"
    "io/ioutil"
)

type module_json struct {
    Gens map[string]int `json:"gens"`
    Sq_actions []SqAction `json:"sq_actions"`
    Milnor_actions []MilnorAction `json:"milnor_actions"`
}

type module struct {
    Gens []int
}

type SqAction struct {
    Op int `json:"op"`
    Input string `json:"input"`
    Output []struct{
        Gen string
        Coeff int
    }
} 

type MilnorAction struct {
    Op []int `json:"op"`
    Input string `json:"input"`
    Output []struct{
        Gen string
        Coeff int
    }
} 

func main(){

    jsonFile, err := os.Open("test.json")
    if err != nil {
        fmt.Println(err)
    }
    
    defer jsonFile.Close()
    
    byteValue, _ := ioutil.ReadAll(jsonFile)
    
    var module module_json
    json.Unmarshal(byteValue , &module)
    fmt.Println(module)
    
}

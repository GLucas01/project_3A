package main

import (

    "fmt"
	"time"
	"sync"
)


//var wg_main sync.WaitGroup

       

func mul_mat(mat_1, mat_2 [][]int, result *[100][100]int , mat_row_1,mat_col_1,mat_col_2 int) {
    wg := sync.WaitGroup{}
    for i := 0; i < mat_row_1; i++ {
        for j := 0; j < mat_col_2; j++ {
            wg.Add(1)
            go func(i, j int) {
                defer wg.Done()
                total := 0
                for k := 0; k < mat_col_1; k++ {
                    total += mat_1[i][k] * mat_2[k][j]
                    
                }
                result[i][j]=total
            }(i, j)
        }
    }
    wg.Wait()

}

	
func main() {

        var mat_row_1, mat_col_1 int
        var mat_row_2, mat_col_2 int

        fmt.Print("Enter no of rows of mat_1: ")
        fmt.Scanln( & mat_row_1)

        fmt.Print("Enter no of column of mat_1: ")
        fmt.Scanln( & mat_col_1)
        
        fmt.Print("Enter no of rows of mat_2: ")
        fmt.Scanln( & mat_row_2)

        fmt.Print("Enter no of column of mat_2: ")
        fmt.Scanln( & mat_col_2)
		
		//error:invalid array length mat_row_1
		//var mat_1[mat_row_1][mat_col_1] int
		//var mat_2[mat_row_2][mat_col_2] int
		
		//error:panic: runtime error: makeslice: cap out of range [can't do multidimension]
		//var mat_1 = make([][]int,mat_row_1,mat_col_1)
		//var mat_2 = make([][]int,mat_row_2,mat_col_2)
		if mat_col_1 != mat_row_2 {

			fmt.Println("Error: The matrix cannot be multiplied")

		}else{
		mat_1 := make([][]int, mat_row_1)
		for i := range mat_1 {
			mat_1[i] = make([]int, mat_col_1)
		}
		
		mat_2 := make([][]int, mat_row_2)
		for i := range mat_2 {
			mat_2[i] = make([]int, mat_col_2)
		}
		
	
		var result[100][100] int

        fmt.Println("\nEnter first matrix: \n")

        for m, r := range mat_1 {
			for l := range r {
				fmt.Scan(&mat_1[m][l])
			}
		}

        fmt.Println("\nEnter second matrix: \n")

        for mm, rr := range mat_2 {
			for ll := range rr {
				fmt.Scan(&mat_2[mm][ll])
			}
		}

		start := time.Now().UnixNano() 
		//start := time.Now().UnixNano() / int64(time.Millisecond)
        // Multiplication of two matrix
        
        mul_mat(mat_1, mat_2, &result, mat_row_1, mat_col_1, mat_col_2)

		end := time.Now().UnixNano() 
		//end := time.Now().UnixNano() / int64(time.Millisecond)
		diff := (end - start)
        fmt.Println("Duration(ns):", start,end,diff)
        fmt.Println("\nMatrix Multiplication: \n")

        for i:=0; i < mat_row_1; i++{

            for j:=0; j < mat_col_2; j++{

                fmt.Printf("%d ", result[i][j])

            }
            fmt.Println("\n")

        }
	}
}

package main

import (
	"encoding/gob" // bibliothèque pour l'échange de données entre des codes en go en utilisant le format binaire
	"fmt"
	"net"
	"time"
	"runtime"
)

type Result struct {
	i int
	j int
	res int
}

//function pour faire la calculation de la matrice 
func mul_mat(mat_1, mat_2 [][]int, mat_row_1_start, mat_row_1_end, mat_col_1, mat_col_2 int, results chan Result) {   // prends en arg 2 matrices d'entiers, la commence de la ligne matrice 1(5 par 5), la fin de la ligne matrice 1(5 par 5), le nombre de colonne matrice 1, le nombre de colonne matrice 2 et un channel resultat 
    for i := mat_row_1_start; i < mat_row_1_end; i++ {
        for j := 0; j < mat_col_2; j++ {
             func(i, j int) {
                total := 0
                for k := 0; k < mat_col_1; k++ {
                    total += mat_1[i][k] * mat_2[k][j]
                }
                
                response := Result{i: i, j: j, res: total} //response stocké le resultat dans une structure Result
				results <- response //stocke la reponse dans le channel
            }(i, j)
        }
    }

}

func multiplicationMatrices(mat1 [][]int, mat2 [][]int) [][]int { // prends 2 matrices et renvoies le résultat
	mat_row_1 :=len(mat1)
	mat_col_1 :=len(mat1[0])
	mat_col_2 :=len(mat2[0])
	
	results := make(chan Result) //creer un channel
	
	start := time.Now().UnixNano() 
	worker_number:=8 //nombre de worker
	row_per_worker :=(mat_row_1)/worker_number //nombre de la ligne par un worker
	mod := (mat_row_1)%worker_number //modulo de la ligne
	
	if mod != 0{
			
		for i :=0; i< mat_row_1-mod ; i = i+row_per_worker {
			go mul_mat(mat1, mat2,i,i+row_per_worker, mat_col_1, mat_col_2,results)
		}
		
		new :=mat_row_1-mod
		end :=mat_row_1
		go mul_mat(mat1, mat2,new,end, mat_col_1, mat_col_2,results)
	
	}else{
			
		for i := 0; i < mat_row_1; i = i+row_per_worker {
			go mul_mat(mat1, mat2,i,i+row_per_worker, mat_col_1, mat_col_2,results)	
		}
	}
	
	
	ans := make([][]int, len(mat1)) // créer une matrice résultat de la taille de mat1 avec des zéros partout
	for i := range ans {            // pour chaque ligne de la matrice on crée une nouvelle colonne de la taille de la première colonne de mat2.
		ans[i] = make([]int, len(mat2[0]))
	}
	for a := 1; a <= mat_row_1*mat_col_2; a++ {
			x:= <-results //retirer le resultat depuis le channel

		ans[x.i][x.j] = x.res
/*		for i:=0; i < mat_row_1; i++{
			for j:=0; j < mat_col_2; j++{	
					
					if (x.i==i && x.j==j) {
						ans[i][j] = x.res //mettre le resultat dans la matrice resultat en suivrant l'ordre
					}

			}
		}*/	
	}
	end := time.Now().UnixNano() 
	diff := (end - start)
	fmt.Println("Duration(ns):", start,end,diff)
	return ans
}

func Connection(conn net.Conn) { // utilisation de gob pour décoder les matrices recues par le client tcp.
	//start := time.Now().UnixNano() 
	decoder := gob.NewDecoder(conn)
	var mat1, mat2 [][]int
	err := decoder.Decode(&mat1) // decode la première matrice
	if err != nil {
		fmt.Println(err)
		return
	}
	err = decoder.Decode(&mat2) // decode la seconde matrice
	if err != nil {
		fmt.Println(err)
		return
	}
	result := multiplicationMatrices(mat1, mat2) // multiplie les 2 matrices décodées avec les fonctions définies ci-dessus
	encoder := gob.NewEncoder(conn)              // créer un nouvel encodage pour les données envoyées via  conn
	err = encoder.Encode(result)                 //encode le résultat pour le renvoyer au client tcp.
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Close() //ferme la connexion
	//end := time.Now().UnixNano() 
	//diff := (end - start)
	//fmt.Println("Duration(ns):", start,end,diff)
}

func main() { // main pour écouter le client sur le port 8080
	fmt.Printf("NumCPU; %v\n", runtime.NumCPU())
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	for { //boucle infinie pour accepter les connexions.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go Connection(conn) // pour chaque connexion il appelle handleconnection pour gérer la connexion.
	}
}

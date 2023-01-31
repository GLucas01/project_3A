package main

import (
	"encoding/gob" // bibliothèque pour l'échange de données entre des codes en go en utilisant le format binaire
	"fmt"
	"net"
	"sync"
)

func multiplication(mat1 [][]int, mat2 [][]int, resultat [][]int, i int, j int, wg *sync.WaitGroup) { // prends en arg 3 matrices d'entiers, deux entiers i et j qui sont les index de ligne et colonne et un pointeur
	sum := 0
	for k := 0; k < len(mat1[0]); k++ { //boucle interne avec l'index k qui parcout les colonnes de la première matrice et les lignes de la seconde
		sum = sum + mat1[i][k]*mat2[k][j]
	}
	resultat[i][j] = sum // résultat de la multiplication d'une ligne par une colonne. Représente un élément de la matrice du résultat.
	wg.Done()
}

func multiplicationMatrices(mat1 [][]int, mat2 [][]int) [][]int { // prends 2 matrices et renvoies le résultat
	resultat := make([][]int, len(mat1)) // créer une matrice résultat de la taille de mat1 avec des zéros partout
	for i := range resultat {            // pour chaque ligne de la matrice on crée une nouvelle colonne de la taille de la première colonne de mat2.
		resultat[i] = make([]int, len(mat2[0]))
	}

	var wg sync.WaitGroup            // wg pour synchroniser les goroutines.
	for i := 0; i < len(mat1); i++ { // len(mat1) correspond au nbr de lignes de la matrice   1
		for j := 0; j < len(mat2[0]); j++ { // len(mat2[0]) correspond au nbr de colonnes de la matrice 2
			wg.Add(1)
			go multiplication(mat1, mat2, resultat, i, j, &wg) // multiplication des matrices à l'aide des goroutines. Concurrence sur chaque élément
		} // on peut faire de la concurrence sur chaque ligne en mettant la goroutine après le for i.
	}
	wg.Wait() // attente de synchronisation des goroutines.
	return resultat
}

func Connection(conn net.Conn) { // utilisation de gob pour décoder les matrices recues par le client tcp.
	start := time.Now().UnixNano()
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
	end := time.Now().UnixNano()
	diff := (end - start)
	fmt.Println("Duration(ns):", start, end, diff)
}

func main() { // main pour écouter le client sur le port 8080
	ln, err := net.Listen("tcp", ":8080")
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

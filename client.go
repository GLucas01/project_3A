package main

import (
	"bufio"
	"encoding/gob" //bibliothèque magique
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()

	// lis les matrices depuis les fichiers txt
	fmt.Println("Calcul du produit matriciel par  Shawn & Lucas")
	var fileName1 string
	fmt.Print("Entrez le nom du fichier à envoyer : ")
	fmt.Scan(&fileName1)
	mat1, err := readMatrixFromFile(fileName1)
	if err != nil {
		fmt.Println("Erreur lecture mat1:", err.Error())
		return
	}
	var fileName2 string
	fmt.Print("Entrez le nom du deuxième fichier à envoyer : ")
	fmt.Scan(&fileName2)
	mat2, err := readMatrixFromFile(fileName2)
	if err != nil {
		fmt.Println("Erreur lecture mat2:", err.Error())
		return
	}

	// envoies les matrices au serveur
	enc := gob.NewEncoder(conn)
	err = enc.Encode(mat1)
	if err != nil {
		fmt.Println("Erreur encodage mat1", err.Error())
		return
	}
	err = enc.Encode(mat2)
	if err != nil {
		fmt.Println("Erreur encodage mat2", err.Error())
		return
	}

	// Recoit la matrice résultat du serveur
	dec := gob.NewDecoder(conn)
	var result [][]int
	err = dec.Decode(&result)
	if err != nil {
		fmt.Println("Erreur decodage du résultat", err.Error())
		return
	}
	fmt.Println("Matrice résultat: ")
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i])
	}
	fmt.Println("Calculer le produit de matrices à la main c'est has-been")
	fmt.Println("Merci PFR et Goland")
}

func readMatrixFromFile(fileName string) ([][]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matrix [][]int
	var row []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, s := range strings.Split(scanner.Text(), " ") {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			row = append(row, i)
		}
		matrix = append(matrix, row)
		row = nil
	}
	return matrix, scanner.Err()
}

import random

matrix_size=1000
file_name="Mat.txt"

for j in range(matrix_size) :
	row=""
	for i in range(matrix_size-1) :
		row+=str(random.randint(0, 100))+" "
	
	row+=str(random.randint(0, 100))
	with open(file_name, "a") as file:
		file.write(row+"\n")


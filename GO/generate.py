import random


for j in range(1000) :
	row=""
	for i in range(999) :
		row+=str(random.randint(0, 100))+" "
	
	row+=str(random.randint(0, 100))
	with open("n2.txt", "a") as file:
		file.write(row+"\n")


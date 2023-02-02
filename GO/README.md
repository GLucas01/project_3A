In this go project, we used go routines and channel to optimize the calculation time for the matrix multiplication.<br />
At the end of this project, we found out that the calculation time will not have a significant change due to the powerful processors of the machine until a larger size of matrix is tested. (eg. 1000\*1000 matrix) <br />

--START THE PROGRAM--
1. Download all the source code in a directory
2. Be sure to install ``go`` in your machine
3. Start the server program in terminal to start listen for the connection
```
go run server.go
```
4. Before start the client program, you can write your own matrix in a ```txt``` file. Or you can just create a random matrix using the python code given
```
python3 generate.py
```
Make sure you set the matrix size desired and change the file name
5. Open a new terminal and start the client program
```
go run client.go
```
You can insert the ```txt``` file for the matrix you want to calculate
6. In the server program, it will display the time taken (in nano seconds) for the matrix calculation


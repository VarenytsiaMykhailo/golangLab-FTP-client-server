package main

import (
	"bufio"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var DOWNLOAD_DIR string = "./" //default
func main() {
	fmt.Println("ENTER THE ADDRESS TO CONNECT (example: students.yss.su:21)")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	var addr string = string(line)
	fmt.Println("YOUR ADDR TO CONNECT:", addr)
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	fmt.Println("ENTER THE DOWNLOAD PATH (example: D:/GOPROJECTS/src/ftpLab/client/data for test/) OR \"DF\", TO USE DEFAULT PATH")
	line, _, _ = bufio.NewReader(os.Stdin).ReadLine()//D:/GOPROJECTS/src/ftpLab/client/data for test/
	if string(line) != "DF" {
		DOWNLOAD_DIR = string(line)
	}
	fmt.Println("YOUR DOWNLOAD PATH:", DOWNLOAD_DIR)

	var login string
	fmt.Println("ENTER THE LOGIN TO CONNECT OR \"anonymous\" FOR READ-ONLY MODE")
	fmt.Scan(&login)
	fmt.Println("YOUR LOGIN TO CONNECT:", login)

	var pass string
	fmt.Println("ENTER THE PASSWORD TO CONNECT OR \"anonymous\" FOR READ-ONLY MODE")
	fmt.Scan(&pass)
	fmt.Println("YOUR PASSWORD TO CONNECT:", pass)

	if err = c.Login(login, pass); err != nil { //если передать "anonymous", "anonymous" в качестве логина и пароля, то клиент будет обладать только правами чтения
		fmt.Println(err)
		fmt.Println("ONLY READ MODE")
	}

	/*	if err = c.MakeDir("Mikle"); err != nil {
		fmt.Println(err)
	}*/

	// Do something with the FTP conn
/*	if err = c.ChangeDir("Mikle"); err != nil {
		fmt.Println(err)
	}
	if err = sendFile(c, "./client/data for test/img.jpg", "image.jpg"); err != nil {
		fmt.Println(err)
	}
	if err = sendFile(c, "./client/data for test/img.jpg", "image2"); err != nil {
		fmt.Println(err)
	}
	if err = c.MakeDir("testDir"); err != nil {
		fmt.Println(err)
	}
	if err = filesList(c); err != nil {
		fmt.Println(err)
	}
	if err = c.Delete("image2"); err != nil {
		fmt.Println(err)
	}
	if err = c.RemoveDir("testDir"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("After deleting:")
	if err = filesList(c); err != nil {
		fmt.Println(err)
	}
	if err = getFile(c, "image.jpg", "./client/data for test/downloaded/"); err != nil {
		fmt.Println(err)
	}*/
	err = nil
	fmt.Println("ENTER THE COMMAND:")
	for {
		//bufio.NewReader(os.Stdin).ReadString('\n') //очищаем буфер (нужно, чтобы небыло проблем с выводом в консоль)
		var command string
		fmt.Print("< ")
		fmt.Scan(&command)
		/*line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		slOfLine := strings.Split(string(line), " ")
		command = string(slOfLine[0])*/
		if command == "LIST" { // LIST
			res, err := c.List("")
			if err == nil {
				fmt.Println("> LIST OF CURRENT DIR:")
				for _, v := range res {
					fmt.Println(v.Name)
				}
			}
		} else if command == "CD" { // CD test_dir
			//fmt.Scan(&arg1)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			err = c.ChangeDir(string(line))
		} else if command == "MKDIR" { // MKDIR test_dir
			//fmt.Scan(&arg1)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			err = c.MakeDir(string(line))
		} else if command == "RMDIR" { // RMDIR test_dir
			//fmt.Scan(&arg1)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			err = c.RemoveDir(string(line))
		} else if command == "RMFILE" { // RMFILE image.jpg
			//fmt.Scan(&arg1)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			err = c.Delete(string(line))
		}else if command == "GET" { // GET image.jpg
			//fmt.Scan(&arg1)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			err = getFile(c, string(line), DOWNLOAD_DIR)
		} else if command == "SEND" { // SEND ./client/data for test/img1.jpg
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			slOfLine := strings.Split(string(line), "/")
			err = sendFile(c, string(line), slOfLine[len(slOfLine)-1])
		}else if command == "CDD" { //change donwnload dir
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			DOWNLOAD_DIR = string(line)
			fmt.Println("YOUR DOWNLOAD PATH:", DOWNLOAD_DIR)
		} else if command == "EXIT" {
			break
		} else {
			fmt.Println("> UNKNOWN COMMAND. TRY AGAIN")
		}
		if err != nil {
			fmt.Println("> ERROR:" + err.Error())
			err = nil
		}
	}

	if err := c.Quit(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

}

func sendFile(c *ftp.ServerConn, filePath string, fileNameOnServer string) error {
	data, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "ERROR in sendFile function: os.Open:")
	}
	if err = c.Stor(fileNameOnServer, bufio.NewReader(data)); err != nil {
		return errors.Wrap(err, "ERROR in sendFile function: c.Stor:")
	}
	return nil
}

func getFile(c *ftp.ServerConn, fileNameOnServer string, pathForSavingOfFile string) error {
	r, err := c.Retr(fileNameOnServer) //получаем ответ *ftp.Response с ftp сервера с содержимым файла
	if err != nil {
		return errors.Wrap(err, "ERROR in getFile function: c.Retr:")
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r) //считываем в буфер информацию из ответа от сервера
	if !strings.HasSuffix(pathForSavingOfFile, "/") {
		pathForSavingOfFile += "/"
	}
	pathForSavingOfFile += fileNameOnServer
	file, err := os.Create(pathForSavingOfFile)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "ERROR in getFile function: Cant to create file:")
	}

	if _, err = file.Write(buf); err != nil {
		return errors.Wrap(err, "ERROR in getFile function: file.Write:")
	}
	return nil
}

func filesList(c *ftp.ServerConn) error {
	res, err := c.List("")
	if err != nil {
		return errors.Wrap(err, "ERROR in filesList function:")
	}
	for _, v := range res {
		fmt.Println(v)
	}
	return nil
}

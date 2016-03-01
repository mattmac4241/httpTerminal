package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func EvalCommand(message string) string {
	command := strings.Split(message, " ")
	//add empty string to command to prevent having to check everytime
	//and let the functions perform that
	if len(command) <= 1 {
		//perform twice because most args for a function is two(mv,cp)
		for i := 0; i < 2; i++ {
			command = append(command, " ")
		}
	}
	//evaluate the command, if it catches run the function command
	switch command[0] {
	case "ls":
		return ls()
	case "cat":
		return cat(command[1])
	case "cd":
		return cd(command[1])
	case "rm":
		return rm(command[1])
	case "rmf":
		return removeFolder(command[1])
	case "mkdir":
		return mkdir(command[1])
	case "pwd":
		return pwd()
	case "mv":
		return mv(command[1], command[2])
	case "cp":
		return cp(command[1], command[2])
	case "run":
		//have to join otherwise only first part is used and won't work
		cmd := strings.Join(command[1:], " ")
		return run(cmd)
	case "-h":
		return help()
	default:
		return "Not a valid call, send -h for help"
	}
}

//print the current working directory
func ls() string {
	fileList := ""
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return "Failed to get list of files"
	}
	//loop through all the files in the directory and append them to fileList
	for _, f := range files {
		fileList += fmt.Sprintf("%s\n", f.Name())
	}
	return fileList
}

//return the contents of a file
func cat(fileName string) string {
	file, err := getFile(fileName)
	if err != nil {
		return "Failed to cat file"
	}
	return string(file)
}

//change the directory
func cd(path string) string {
	err := os.Chdir(path)
	if err != nil {
		return "Failed to change path"
	}
	return fmt.Sprintf("Changed path to %s", pwd())
}

//remove the file
func rm(filename string) string {
	if notAllowed(filename) == true {
		return "Not allowed to remove that folder or file "
	}
	err := os.Remove(filename)
	if err != nil {
		return "Failed to remove file or file does not exist"
	}
	return "File removed"
}

//remove a folder
func removeFolder(path string) string {
	//check that folder is not protected
	if notAllowed(path) == true {
		return "Not allowed to remove that directory"
	}
	//check that folder does exist
	_, err := os.Open(path)
	if err != nil {
		return "Not able to remove folder or folder does not exist"
	}
	//remove folder
	err = os.RemoveAll(path)
	if err != nil {
		return "Failed to delete folder"
	}

	return "Deleted folder"
}

//create a directory
func mkdir(folder string) string {
	err := os.Mkdir(folder, 0700)
	if err != nil {
		return "Failed to create directory"
	}
	return "Created the directory"
}

//print the current directory
func pwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "Failed to get directory"
	}
	return wd
}

//move file into another file
func mv(fileFrom, fileTo string) string {
	if notAllowed(fileFrom) == true {
		return "Not allowed to move that file"
	}
	//copy fileFrom into fileTo
	err := copy(fileFrom, fileTo)
	if err != nil {
		return "Failed to move file"
	}
	//then delete fileFrom
	rm(fileFrom)
	return "Successfully moved file"
}

//copy a file into another
func cp(fileFrom, fileTo string) string {
	err := copy(fileFrom, fileTo)
	if err != nil {
		return "Failed to copy file"
	}

	return "Succesfully copied file"

}

//exec commands
func run(command string) string {
	var cmd []byte
	var err error
	parts := strings.Fields(command) //seperate command
	if notValidCommand(parts[0]) == true {
		return "Not allowed to run that command"
	}
	switch {
	case len(parts) == 1:
		cmd, err = exec.Command(parts[0]).Output()
	case len(parts) > 1:
		args := strings.Join(parts[1:], " ")
		cmd, err = exec.Command(parts[0], args).Output()
	default:
		return "No command to run "
	}

	if err != nil {
		fmt.Println(err)
		return "Failed to run"
	}
	return string(cmd)
}

//get the contents of a file
func getFile(file string) ([]byte, error) {
	var content []byte
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("Failed to open file")
	}
	r := bufio.NewReader(f)
	buf := make([]byte, 2048)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		content = append(content, buf[:n]...)
		if n == 0 {
			break
		}
	}
	return content, nil

}

//copy a file
func copy(fileFrom, fileTo string) error {
	//open file and get contents
	file1, err := getFile(fileFrom)
	if err != nil {
		return errors.New("Failed to copy file")
	}
	//create a new file
	file2, err := os.Create(fileTo)
	if err != nil {
		return errors.New("Failed to copy file")
	}
	//write the new file
	_, err = file2.Write(file1)
	if err != nil {
		return errors.New("Failed to copy file")
	}

	return nil
}

//check a file/folder to see if it is allowed to be used
func notAllowed(fileName string) bool {
	notAllowedFiles := []string{"src/server/server.go", "src/commands/commands.go", "src/commands", "src/server"}
	for _, file := range notAllowedFiles {
		if compareFiles(fileName, file) == true {
			return true
		}
	}
	return false
}

//compare files/directory to determine if they are the same
func compareFiles(fileName1, fileName2 string) bool {
	stat1, err1 := os.Stat(fileName1)
	stat2, err2 := os.Stat(fileName2)
	if err1 != nil || err2 != nil {
		return true
	}
	return os.SameFile(stat1, stat2)
}

//check command for eval to prevent already defined commands
func notValidCommand(command string) bool {
	switch command {
	case "cd", "ls", "rm", "mkdir", "pwd", "mv", "cp":
		return true
	default:
		return false
	}
}

//return the help string
func help() string {
	help := fmt.Sprintf("Command\t Action\n" +
		"ls List files\n" +
		"pwd  Print Working Directory\n" +
		"cd [pathname] Change directory\n" +
		"cat [filename] Print contents of file\n" +
		"rm [filename] Remove file\n" +
		"rmf [foldername] Remove folder\n" +
		"mv [filefrom] [fileto] Move file\n" +
		"cp [filefrom] [fileto] Copy file to\n" +
		"run [command] Run command or program")
	return help
}

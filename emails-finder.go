package main

import (
 "bufio"
 "fmt"
 "os"
 "regexp"
 "sync"
)

func main() {
 if len(os.Args) != 3 {
  fmt.Println("Yo, look how to use this: ./email_extractor <input_file> <output_file>")
  return
 }

 inputFile := os.Args[1]
 outputFile := os.Args[2]

 // Открываем файл в котором предположительно содержатся мейлы
 file, err := os.Open(inputFile)
 if err != nil {
  fmt.Printf("No, bro, your file is not correct %s: %v\n", inputFile, err)
  return
 }
 defer file.Close()

 // Создаем результатирующий файл
 outFile, err := os.Create(outputFile)
 if err != nil {
  fmt.Printf("Hell, some troubles with your output file %s: %v\n", outputFile, err)
  return
 }
 defer outFile.Close()

 writer := bufio.NewWriter(outFile)
 defer writer.Flush()

 // Регулярка для мейлов с поддержкой кириллицы
 emailRegex := regexp.MustCompile(`\b[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}\b`)

 linesChan := make(chan string, 100)
 emailsChan := make(chan string, 100)
 var wg sync.WaitGroup

 // Запускаем воркеров для обработки строк
 // todo: воркеркаунт вытащить в аргумент терминала
 workerCount := 5
 for i := 0; i < workerCount; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   for line := range linesChan {
    emails := emailRegex.FindAllString(line, -1)
    for _, email := range emails {
     emailsChan <- email
    }
   }
  }()
 }

 // Читаем файл и суем строки в канал
 go func() {
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
   linesChan <- scanner.Text()
  }
  close(linesChan)
  if err := scanner.Err(); err != nil {
   fmt.Printf("Ohhh, shit happens: %v\n", err)
  }
 }()

 // Закрываем канал emailsChan, когда все горутины отработают
 go func() {
  wg.Wait()
  close(emailsChan)
 }()

 // Записываем email'ы в файл
 for email := range emailsChan {
  _, err := writer.WriteString(email + "\n")
  if err != nil {
   fmt.Printf("Look, man, some trouble with writing to file: %v\n", err)
   return
  }
 }

 fmt.Println("Aee, your pure emails into:", outputFile)
}

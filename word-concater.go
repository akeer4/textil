package main

import (
 "bufio"
 "fmt"
 "os"
 "sync"
)

func main() {
 if len(os.Args) != 5 {
  fmt.Println("Heeey, look how to use it, bro: ./words-concater <file1> <file2> <delimiter> <output_file>")
  return
 }

 file1, err := os.Open(os.Args[1])
 if err != nil {
  fmt.Printf("Blyat, <file1> is not correct!!! %s: %v\n", os.Args[1], err)
  return
 }
 defer file1.Close()

 file2, err := os.Open(os.Args[2])
 if err != nil {
  fmt.Printf("Noo, god, NOO, <file2> is not correct!!! %s: %v\n", os.Args[2], err)
  return
 }
 defer file2.Close()

 delimiter := os.Args[3]

 outputFile, err := os.Create(os.Args[4])
 if err != nil {
  fmt.Printf("Noo, god, NOO, <output_file> is not correct!!! %s: %v\n", os.Args[4], err)
  return
 }
 defer outputFile.Close()

 words1Chan := make(chan string)
 words2Chan := make(chan string)
 resultsChan := make(chan string)

 var wg sync.WaitGroup

 // Чекаем file1 и пиздячим слова в words1Chan
 wg.Add(1)
 go func() {
  defer wg.Done()
  scanner := bufio.NewScanner(file1)
  for scanner.Scan() {
   words1Chan <- scanner.Text()
  }
  close(words1Chan)
  if err := scanner.Err(); err != nil {
   fmt.Printf("Da ebaniy nasos, I can't read the <file1> %s: %v\n", os.Args[1], err)
  }
 }()

 // Чекаем file2 и пиздячим слова в words2Chan
 wg.Add(1)
 go func() {
  defer wg.Done()
  scanner := bufio.NewScanner(file2)
  for scanner.Scan() {
   words2Chan <- scanner.Text()
  }
  close(words2Chan)
  if err := scanner.Err(); err != nil {
   fmt.Printf("Woooops, dude, some trouble with your <file2> %s: %v\n", os.Args[2], err)
  }
 }()

 // Берем все слова из words2Chan и сохраняем в срез
 var words2 []string
 for word2 := range words2Chan {
  words2 = append(words2, word2)
 }

 // Конкатим слова и пишем в resultsChan
 wg.Add(1)
 go func() {
  defer wg.Done()
  for word1 := range words1Chan {
   for _, word2 := range words2 {
    resultsChan <- fmt.Sprintf("%s%s%s\n", word1, delimiter, word2)
   }
  }
  close(resultsChan)
 }()

 // Кидаем результаты в итоговый фалй
 wg.Add(1)
 go func() {
  defer wg.Done()
  writer := bufio.NewWriter(outputFile)
  for result := range resultsChan {
   _, err := writer.WriteString(result)
   if err != nil {
    fmt.Printf("Oh no, problem writing to output file: %v\n", err)
    return
   }
  }
  writer.Flush()
 }()

 // Ждём когда потухнут все горутины
 wg.Wait()

 fmt.Println("Thanks balls, we did it, bro! WE DID IT!!! Here your output file:", os.Args[4])
}

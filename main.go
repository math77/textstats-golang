package main

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "log"
  "time"
  "math/rand"
  "strconv"
  "encoding/json"
)

//Estrutura que representa as estatisticas salvas no arquivo.
type Stats struct {
  FileName string `json:"filename"`
  Caracters string `json:"caracters"`
  Words string `json:"words"`
  FreqLetters map[string]int `json:"letters_frequencys"`
  FreqWords map[string]int `json:"words_frequencys"`
}

func main(){
  rand.Seed(time.Now().UnixNano())

  fptr := flag.String("fpath", "", "file path to read from")
  flag.Parse()

  if *fptr == "" {
    log.Fatal("Utilize -fpath=filepath para passar um arquivo a ser analisado!")
  }

  words, err := scanWords(*fptr)
  if err != nil {
    log.Fatal(err)
  }

  amWords := amountWords(words)
  amCarac := amountCaracters(words)
  freqWords := frequencyWords(words)
  freqLetters := frequencyLetters(words)


  data := &Stats{
    FileName: *fptr,
    Words: strconv.Itoa(amWords),
    Caracters: strconv.Itoa(amCarac),
    FreqLetters: freqLetters, FreqWords: freqWords}

  jsonData, err := json.MarshalIndent(data, "", "")

  if err != nil {
    panic(err)
  }

  writeToFile(jsonData)
}

/*
  Ler o arquivo palavra por palavra.
*/
func scanWords(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)

  scanner.Split(bufio.ScanWords)

  var words []string

  for scanner.Scan() {
    words = append(words, scanner.Text())
  }

  return words, nil
}

/*
 Conta a quantidade de caracteres em um texto.
*/
func amountCaracters(words []string) int {
  count := 0
  for _, word := range words {
    //converte string em slice de rune e retorna o tamanho.
    count += len([]rune(word))
  }
  return count
}

/*
  Conta a quantidade de palavras em um arquivo.
*/
func amountWords(words []string) int {
  return len(words)
}

/*
  Recebe a lista de palavras extraidas do arquivo e retorna um mapa com as palavras
  e a frequência com que aparecem no texto.

  Cria um mapa, iterar sobre o slice de palavras, verifica se já existe no mapa
  uma chave com o nome da palavra atual, caso exista o código incrementa +1 ao valor na chave
  caso contrário cria uma nova chave no mapa com a palavra atual e adiciona 1 ao valor.
*/
func frequencyWords(words []string) map[string]int {
  frequencies := make(map[string]int)

  for _, word := range words {
    if val, ok := frequencies[word]; ok {
      frequencies[word] = val+1
    } else {
      frequencies[word] = 1
    }
  }

  return frequencies
}

/*
  Recebe a lista de palavras extraidas do arquivo e retorna um mapa com as letras
  e a frequência com que aparecem no texto.

  Cria um mapa, iterar sobre o slice de palavras, criar um slice de letras e iterar sobre ele
  verifica se já existe no mapa uma chave com a letra atual, caso exista o código incrementa +1 ao valor na chave
  caso contrário cria uma nova chave no mapa com a letra atual e adiciona 1 ao valor.
*/
func frequencyLetters(words []string) map[string]int {
  frequencies := make(map[string]int)

  for _, word := range words {
    for _, letter := range []rune(word){
      if val, ok := frequencies[string(letter)]; ok {
        frequencies[string(letter)] = val+1
      } else {
        frequencies[string(letter)] = 1
      }
    }
  }
  return frequencies
}

/*
 Gera um número aleatorio para ser usado na função que gera
 uma string aletoria usada como nome do arquivo.
*/
func randomInt(min, max int) int {
  return min + rand.Intn(max-min)
}

/*
 Gera uma string aleatoria para ser usada como nome do arquivo.
*/
func randomString(len int) string {
  bytes := make([]byte, len)
  for index := 0; index < len; index++ {
    bytes[index] = byte(randomInt(65, 90))
  }
  return string(bytes)
}

/*
  Cria um arquivo .json e escreve os resultados das analises nele.
*/

func writeToFile(jsonData []byte) {
  filename := randomString(10)+".json"

  file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0660)
  if err != nil {
    fmt.Println(err)
    file.Close()
    return
  }

  defer file.Close()
  file.Write(jsonData)
  file.Close()
  fmt.Println("Salvo ", file.Name())
}

package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type bufferTools interface {
	Add(data chan int, quantytyData int)
	Get(duration time.Duration)
}

func (cb *ChannelBuffer) Add(data chan int, quantytyData int) {
	i := 0
	var input string

	for {
		fmt.Scanln(&input)
		num, _ := strconv.Atoi(input)
		data <- num
		i++
		log.Println("Добавлено значение:", num)
		if i == quantytyData {
			close(data)
			break
		}
	}
}
func (cb *ChannelBuffer) Get(interval time.Duration) {
	defer cb.wg.Done()

	for v := range cb.data {
		num, err := strconv.Atoi(strconv.Itoa(v))
		if err != nil {
			continue
		}
		if num == 0 {
			log.Printf("Введенное значение '%d' не может быть равно нулю, и будет отброшено", num)
			continue
		}
		if num < 0 {
			log.Printf("Введенное значение '%d' не может быть отрицательным, и будет отброшено", num)
			continue
		}
		if num%3 != 0 {
			log.Printf("Введенное значение '%d' должно быть кратно 3, значение будет отброшено", num)
			continue
		}
		time.Sleep(interval * time.Second)
		log.Println("Получено значение:", v)
	}
}

type ChannelBuffer struct {
	data chan int
	wg   sync.WaitGroup
}

func NewChannelBuffer(size int) *ChannelBuffer {
	return &ChannelBuffer{
		data: make(chan int, size),
		wg:   sync.WaitGroup{},
	}
}

func main() {
	const quantityData int = 4       // колличество исходных данных для ввода
	const bufferSize int = 4         // размер буфера
	const interval time.Duration = 2 // интервал времени чтения в секундах
	const goCount = 2                // колличество горутин

	fmt.Println("Количество данных")
	fmt.Println(quantityData)

	bf := NewChannelBuffer(bufferSize)
	bt := bufferTools(bf)

	fmt.Println("Введите значение:")

	bt.Add(bf.data, quantityData) // добавляем данные в буфер

	bf.wg.Add(goCount)
	for i := 0; i < goCount; i++ {
		go bt.Get(interval) // получаем данные из буфера
	}
	bf.wg.Wait()
	fmt.Println("Все данные считаны")

}

package chanel

import (
	"strings"
	"time"
)

// Agregate chanel is used to remove duplicate in a chanel

type AgregateChanel struct {
	inputChanel   chan string
	outputChanel  chan string
	previousValue string
	timeout       time.Duration
}

func NewAgregateChanel(timeout time.Duration)*AgregateChanel{
	ag :=  &AgregateChanel{
		inputChanel:  make(chan string,10),
		outputChanel: make(chan string,10),
		timeout:      time.Second*timeout,
	}
	go ag.watch()
	return ag
}

func (ag AgregateChanel)Add(data string){
	ag.inputChanel <- data
}

// Get value in a chanel, blocking
func (ag AgregateChanel)Get()string{
	return <- ag.outputChanel
}

func (ag *AgregateChanel) watch() {
	// At each run, wait for message (timeout if necessary)
	for {
		// If previous is empty, no need to used timeout
		if strings.EqualFold("",ag.previousValue){
			ag.readChan()
		}else{
			ag.readChanWithTimeout()
		}
	}
}

func (ag * AgregateChanel)readChan(){
	value := <- ag.inputChanel
	ag.manageValue(value)
}

func (ag * AgregateChanel)readChanWithTimeout(){
	select {
	case value := <-ag.inputChanel:
		ag.manageValue(value)
		break
	case <-time.NewTimer(ag.timeout).C:
		ag.manageTimeout()
	}
}

func (ag *AgregateChanel)manageValue(value string){
	// Compare value with previous. If same, wait next value. If different, send previous if different
	if strings.EqualFold("", ag.previousValue) {
		ag.previousValue = value
	} else {
		if !strings.EqualFold(ag.previousValue, value) {
			ag.outputChanel <- ag.previousValue
			ag.previousValue = value
		}
		// Else wait
	}
}

func (ag * AgregateChanel)manageTimeout(){
	if !strings.EqualFold("",ag.previousValue){
		ag.outputChanel <- ag.previousValue
		ag.previousValue = ""
	}
}
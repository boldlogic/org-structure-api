package models

import (
	"errors"
)

var (
	ErrValidation         = errors.New("некорректный запрос")                  //некорректный запрос
	ErrConflict           = errors.New("запись с таким ключом уже существует") //запись с таким ключом уже существует
	ErrNotFound           = errors.New("запись с таким ключом не существует")  //запись с таким ключом не существует
	ErrBusinessValidation = errors.New("некорректные данные в запросе")        //некорректные данные в запросе
	ErrParentNotFound     = errors.New("parent_id не существует")              //parent_id не существует
	ErrCycle              = errors.New("цикл в дереве подразделений")          //цикл в дереве подразделений
)

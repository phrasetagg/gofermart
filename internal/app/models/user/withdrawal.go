package user

import (
	"encoding/json"
	"time"
)

type Withdrawal struct {
	OrderNumber string    `json:"order"`
	Value       float64   `json:"sum"`
	CreatedAt   time.Time `json:"processed_at"`
}

func (w Withdrawal) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type Alias Withdrawal

	aliasValue := struct {
		Alias
		// переопределяем поле внутри анонимной структуры
		CreatedAt string `json:"processed_at"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		Alias: Alias(w),
		// задаём значение для переопределённого поля
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}

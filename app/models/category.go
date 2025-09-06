package models

// categoriesテーブルと対応する構造体(model)
type Category struct {
	ID   int
	Name string
}

func GetOrCreateCategory(name string) (*Category, error) {
	var category Category
	// データがなければ作成
	tx := DB.FirstOrCreate(&category, Category{Name: name})
	if tx.Error != nil {
		return nil, tx.Error
	}
	// pointerを返す
	return &category, nil
}

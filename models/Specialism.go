package models

type Specialism struct {
	ID int `json:"id"  gorm:"AUTO_INCREMENT"`
	Name string `json:"name"`
	Related []Specialism `json:"specialties" gorm:"many2many:relation_specialties;ForeignKey:id;AssociationForeignKey:related_id"`
}
